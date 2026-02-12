package handlers

import (
	"errors"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"unicorn-auth/internal/config"
	"unicorn-auth/internal/models"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthHandler struct {
	cfg      config.Config
	sec      *security.Security
	users    *repo.UserRepo
	sessions *repo.SessionRepo
}

func NewAuthHandler(cfg config.Config, sec *security.Security, users *repo.UserRepo, sessions *repo.SessionRepo) *AuthHandler {
	return &AuthHandler{cfg: cfg, sec: sec, users: users, sessions: sessions}
}

type registerReq struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"` // user/company
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if !bindStrict(c, &req, 32<<10) {
		return
	}
	req.Login = strings.TrimSpace(req.Login)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.Type = strings.TrimSpace(req.Type)

	if req.Login == "" || len(req.Login) > 64 || req.Password == "" || len(req.Password) < 8 || req.DisplayName == "" || len(req.DisplayName) > 64 {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "bad_request"})
		return
	}
	if req.Type != string(models.UserTypeUser) && req.Type != string(models.UserTypeCompany) {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "bad_request"})
		return
	}

	if u, _ := h.users.FindByLoginNorm(c.Request.Context(), strings.ToLower(req.Login)); u != nil {
		c.JSON(http.StatusConflict, gin.H{"ok": false, "error": "login_taken"})
		return
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "server_error"})
		return
	}

	u := &models.User{
		UserID:       ulid.Make().String(),
		Login:        req.Login,
		DisplayName:  req.DisplayName,
		Type:         models.UserType(req.Type),
		PasswordHash: hash,
	}
	if err := h.users.Create(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "server_error"})
		return
	}

	access, refreshCookie, err := h.issueTokens(c, u.UserID, string(u.Type), []string{"pwd"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "server_error"})
		return
	}
	http.SetCookie(c.Writer, refreshCookie)
	c.JSON(http.StatusOK, gin.H{"ok": true, "accessToken": access})
}

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if !bindStrict(c, &req, 16<<10) {
		return
	}
	req.Login = strings.TrimSpace(req.Login)

	u, err := h.users.FindByLoginNorm(c.Request.Context(), strings.ToLower(req.Login))
	if err != nil || u == nil || u.Status.Deleted || u.Status.Blocked {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "invalid_credentials"})
		return
	}
	ok, _ := security.VerifyPassword(req.Password, u.PasswordHash)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "invalid_credentials"})
		return
	}

	if u.MFA.TOTP.Enabled {
		mfaTok, err := h.sec.Tokens.NewMFAToken(u.UserID, string(u.Type), 10*time.Minute)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}
		c.JSON(200, gin.H{"ok": true, "mfaRequired": true, "mfaToken": mfaTok})
		return
	}

	access, refreshCookie, err := h.issueTokens(c, u.UserID, string(u.Type), []string{"pwd"})
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}
	http.SetCookie(c.Writer, refreshCookie)
	c.JSON(200, gin.H{"ok": true, "accessToken": access})
}

type verifyTotpReq struct {
	MFAToken string `json:"mfaToken"`
	Code     string `json:"code"`
}

func (h *AuthHandler) VerifyTOTP(c *gin.Context) {
	var req verifyTotpReq
	if !bindStrict(c, &req, 8<<10) {
		return
	}
	req.Code = strings.TrimSpace(req.Code)

	claims, err := h.sec.Tokens.ParseMFA(strings.TrimSpace(req.MFAToken))
	if err != nil {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	u, err := h.users.FindByUserID(c.Request.Context(), claims.UserID)
	if err != nil || u == nil || u.Status.Deleted || u.Status.Blocked || !u.MFA.TOTP.Enabled {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	secret, err := h.decryptSecret(u.MFA.TOTP.SecretEncB64)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}
	if !verifyTOTPOnce(u, secret, req.Code) {
		c.JSON(401, gin.H{"ok": false, "error": "invalid_totp"})
		return
	}
	_ = h.users.UpdateByUserID(c.Request.Context(), u.UserID, bson.M{"mfa.totp.lastStep": currentTotpStep()})

	access, refreshCookie, err := h.issueTokens(c, u.UserID, string(u.Type), []string{"pwd", "totp"})
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}
	http.SetCookie(c.Writer, refreshCookie)
	c.JSON(200, gin.H{"ok": true, "accessToken": access})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid := c.GetString("userId")
	u, err := h.users.FindByUserID(c.Request.Context(), uid)
	if err != nil || u == nil || u.Status.Deleted || u.Status.Blocked {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "displayName": u.DisplayName, "type": u.Type})
}

type changePwReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req changePwReq
	if !bindStrict(c, &req, 16<<10) {
		return
	}
	if len(req.NewPassword) < 8 {
		c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
		return
	}
	uid := c.GetString("userId")
	u, err := h.users.FindByUserID(c.Request.Context(), uid)
	if err != nil || u == nil || u.Status.Deleted || u.Status.Blocked {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	ok, _ := security.VerifyPassword(req.OldPassword, u.PasswordHash)
	if !ok {
		c.JSON(401, gin.H{"ok": false, "error": "invalid_credentials"})
		return
	}
	hash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}
	if err := h.users.UpdateByUserID(c.Request.Context(), uid, bson.M{"passwordHash": hash}); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func (h *AuthHandler) TotpEnroll(c *gin.Context) {
	uid := c.GetString("userId")
	u, err := h.users.FindByUserID(c.Request.Context(), uid)
	if err != nil || u == nil || u.Status.Deleted || u.Status.Blocked {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	if u.MFA.TOTP.Enabled {
		c.JSON(409, gin.H{"ok": false, "error": "already_enabled"})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Unicorn",
		AccountName: u.Login,
	})
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}

	enc, err := h.sec.Crypt.EncryptToB64([]byte(key.Secret()))
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}
	exp := time.Now().UTC().Add(10 * time.Minute).Unix()
	if err := h.users.UpdateByUserID(c.Request.Context(), uid, bson.M{
		"mfa.totp.pendingEncB64":  enc,
		"mfa.totp.pendingExpires": exp,
	}); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}

	c.JSON(200, gin.H{"ok": true, "secret": key.Secret(), "otpauth": key.URL(), "expiresIn": int64(600)})
}

type totpEnableReq struct {
	Code string `json:"code"`
}

func (h *AuthHandler) TotpEnable(c *gin.Context) {
	var req totpEnableReq
	if !bindStrict(c, &req, 8<<10) {
		return
	}
	req.Code = strings.TrimSpace(req.Code)

	uid := c.GetString("userId")
	u, err := h.users.FindByUserID(c.Request.Context(), uid)
	if err != nil || u == nil || u.Status.Deleted || u.Status.Blocked {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	if u.MFA.TOTP.Enabled {
		c.JSON(409, gin.H{"ok": false, "error": "already_enabled"})
		return
	}
	if u.MFA.TOTP.PendingEncB64 == "" || u.MFA.TOTP.PendingExpires < time.Now().UTC().Unix() {
		c.JSON(400, gin.H{"ok": false, "error": "no_pending_totp"})
		return
	}
	secret, err := h.decryptSecret(u.MFA.TOTP.PendingEncB64)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}
	if !totp.Validate(req.Code, secret) {
		c.JSON(401, gin.H{"ok": false, "error": "invalid_totp"})
		return
	}

	if err := h.users.UpdateByUserID(c.Request.Context(), uid, bson.M{
		"mfa.totp.enabled":        true,
		"mfa.totp.secretEncB64":   u.MFA.TOTP.PendingEncB64,
		"mfa.totp.pendingEncB64":  "",
		"mfa.totp.pendingExpires": int64(0),
		"mfa.totp.lastStep":       currentTotpStep(),
	}); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	rc, err := c.Request.Cookie("refresh")
	if err != nil || rc.Value == "" {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	sessionID, token, ok := splitRefresh(rc.Value)
	if !ok {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	s, err := h.sessions.FindActiveByID(c.Request.Context(), sessionID)
	if err != nil || s == nil {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	if time.Now().UTC().After(s.ExpiresAt) {
		_ = h.sessions.Revoke(c.Request.Context(), sessionID)
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}
	if !secureEqual(hashToken(token), s.RefreshHash) {
		_ = h.sessions.Revoke(c.Request.Context(), sessionID)
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}

	u, err := h.users.FindByUserID(c.Request.Context(), s.UserID)
	if err != nil || u == nil || u.Status.Deleted || u.Status.Blocked {
		c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
		return
	}

	access, err := h.sec.Tokens.NewAccessToken(u.UserID, string(u.Type), []string{"pwd"}, h.cfg.AccessTTL)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}

	newTok, newHash := newRefreshToken()
	newVal := joinRefresh(sessionID, newTok)
	if err := h.sessions.Rotate(c.Request.Context(), sessionID, newHash, time.Now().UTC().Add(h.cfg.RefreshTTL)); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "server_error"})
		return
	}
	http.SetCookie(c.Writer, h.refreshCookie(newVal))
	c.JSON(200, gin.H{"ok": true, "accessToken": access})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	rc, _ := c.Request.Cookie("refresh")
	if rc != nil && rc.Value != "" {
		if sid, _, ok := splitRefresh(rc.Value); ok {
			_ = h.sessions.Revoke(c.Request.Context(), sid)
		}
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: sameSite(h.cfg.CookieSameSite),
		Domain:   h.cfg.CookieDomain,
	})
	c.JSON(200, gin.H{"ok": true})
}

// helpers

func (h *AuthHandler) issueTokens(c *gin.Context, userID, typ string, amr []string) (string, *http.Cookie, error) {
	access, err := h.sec.Tokens.NewAccessToken(userID, typ, amr, h.cfg.AccessTTL)
	if err != nil {
		return "", nil, err
	}
	rtok, rhash := newRefreshToken()
	s, err := h.sessions.Create(c.Request.Context(), userID, rhash, h.cfg.RefreshTTL)
	if err != nil {
		return "", nil, err
	}
	val := joinRefresh(s.SessionID, rtok)
	return access, h.refreshCookie(val), nil
}

func (h *AuthHandler) refreshCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     "refresh",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: sameSite(h.cfg.CookieSameSite),
		Domain:   h.cfg.CookieDomain,
		MaxAge:   int(h.cfg.RefreshTTL.Seconds()),
	}
}

func sameSite(s string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteStrictMode
	}
}

func bindStrict(c *gin.Context, dst any, maxBytes int64) bool {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "bad_request"})
		return false
	}
	return true
}

func newRefreshToken() (token string, hash string) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	token = base64.RawURLEncoding.EncodeToString(b)
	hash = hashToken(token)
	return
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func joinRefresh(sessionID, token string) string {
	return sessionID + "." + token
}

func splitRefresh(v string) (sessionID, token string, ok bool) {
	parts := strings.Split(v, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func secureEqual(a, b string) bool {
	// constant-time compare
	if len(a) != len(b) {
		return false
	}
	var out byte
	for i := 0; i < len(a); i++ {
		out |= a[i] ^ b[i]
	}
	return out == 0
}

func (h *AuthHandler) decryptSecret(encB64 string) (string, error) {
	plain, err := h.sec.Crypt.DecryptFromB64(encB64)
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(string(plain))
	if s == "" {
		return "", errors.New("empty secret")
	}
	return s, nil
}

func currentTotpStep() int64 {
	// default TOTP step 30s
	return time.Now().UTC().Unix() / 30
}

func verifyTOTPOnce(u *models.User, secret, code string) bool {
	if !totp.Validate(code, secret) {
		return false
	}
	step := currentTotpStep()
	// prevent reuse in same time-step
	if u.MFA.TOTP.LastStep == step {
		return false
	}
	return true
}
