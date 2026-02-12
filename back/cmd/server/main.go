package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"unicorn-auth/internal/cleanup"
	"unicorn-auth/internal/config"
	"unicorn-auth/internal/db"
	"unicorn-auth/internal/http/router"
	adminmod "unicorn-auth/internal/modules/admin"
	appmod "unicorn-auth/internal/modules/applications"
	chatmod "unicorn-auth/internal/modules/chat"
	companymod "unicorn-auth/internal/modules/company"
	profilemod "unicorn-auth/internal/modules/profile"
	resumemod "unicorn-auth/internal/modules/resumes"
	submod "unicorn-auth/internal/modules/subscription"
	vacmod "unicorn-auth/internal/modules/vacancies"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.MustLoad()

	sec, err := security.NewSecurity(cfg.JWTHS256Secret, cfg.TotpEncKeyB64)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	d := db.Connect(ctx, cfg.MongoURI, cfg.MongoDB)
	defer func() { _ = d.Close(context.Background()) }()
	d.EnsureIndexes(ctx)

	users := repo.NewUserRepo(d)
	sessions := repo.NewSessionRepo(d)
	profiles := repo.NewProfileRepo(d)
	vac := repo.NewVacancyRepo(d)
	resumes := repo.NewResumeRepo(d)
	apps := repo.NewApplicationRepo(d)
	chatRepo := repo.NewChatRepo(d)
	admins := repo.NewAdminRepo(d)
	subs := repo.NewSubscriptionRepo(d)

	bootstrapAdmin(ctx, admins)

	robokassa := security.NewRobokassa(
		cfg.RobokassaMerchantLogin,
		cfg.RobokassaPassword1,
		cfg.RobokassaPassword2,
		cfg.RobokassaTestMode,
	)

	r := router.New(cfg, sec, users, sessions, resumes, vac)

	// Register modules
	profilemod.Register(r, sec, users, profiles)
	companymod.Register(r, profiles)
	vacmod.Register(r, sec, users, vac)
	resumemod.Register(r, sec, users, resumes, apps)
	appmod.Register(r, sec, users, vac, resumes, apps)
	chatmod.Register(r, sec, users, apps, chatRepo, vac, profiles)
	adminmod.Register(r, sec, admins, users)

	// Subscription module
	subCfg := submod.Config{
		Price:            cfg.SubscriptionPrice,
		DurationDays:     cfg.SubscriptionDuration,
		RobokassaEnabled: cfg.RobokassaMerchantLogin != "" && cfg.RobokassaPassword1 != "",
	}
	submod.Register(r, subCfg, sec, robokassa, users, subs, vac, resumes)

	// Запускаем фоновую очистку старых скрытых записей (каждые 24 часа)
	cleaner := cleanup.NewCleaner(apps)
	go cleaner.Start(context.Background(), 24*time.Hour)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func bootstrapAdmin(ctx context.Context, admins *repo.AdminRepo) {
	login := strings.TrimSpace(os.Getenv("ADMIN_BOOTSTRAP_LOGIN"))
	pass := strings.TrimSpace(os.Getenv("ADMIN_BOOTSTRAP_PASSWORD"))
	if login == "" || pass == "" {
		return
	}
	exists, _ := admins.FindByLoginNorm(ctx, strings.ToLower(login))
	if exists != nil {
		return
	}
	hash, err := security.HashPassword(pass)
	if err != nil {
		log.Printf("admin bootstrap: %v", err)
		return
	}
	_, err = admins.Create(ctx, login, hash)
	if err != nil {
		log.Printf("admin bootstrap: %v", err)
		return
	}
	log.Printf("admin bootstrap created: %s", login)
}
