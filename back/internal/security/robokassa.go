package security

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Robokassa struct {
	MerchantLogin string
	Password1     string
	Password2     string
	TestMode      bool
}

func NewRobokassa(merchantLogin, password1, password2 string, testMode bool) *Robokassa {
	return &Robokassa{
		MerchantLogin: merchantLogin,
		Password1:     password1,
		Password2:     password2,
		TestMode:      testMode,
	}
}

// GeneratePaymentURL создает URL для оплаты через Robokassa
func (r *Robokassa) GeneratePaymentURL(outSum string, invID int64, description string, userID string) string {
	// Генерируем подпись: MD5(MerchantLogin:OutSum:InvId:Password1:Shp_userId=value)
	signStr := fmt.Sprintf("%s:%s:%d:%s:Shp_userId=%s",
		r.MerchantLogin, outSum, invID, r.Password1, userID)
	signature := fmt.Sprintf("%x", md5.Sum([]byte(signStr)))

	baseURL := "https://auth.robokassa.ru/Merchant/Index.aspx"

	params := url.Values{}
	params.Set("MerchantLogin", r.MerchantLogin)
	params.Set("OutSum", outSum)
	params.Set("InvId", strconv.FormatInt(invID, 10))
	params.Set("Description", description)
	params.Set("SignatureValue", signature)
	params.Set("Shp_userId", userID)

	if r.TestMode {
		params.Set("IsTest", "1")
	}

	return baseURL + "?" + params.Encode()
}

// VerifyResultSignature проверяет подпись от ResultURL (когда платеж проведен)
func (r *Robokassa) VerifyResultSignature(outSum string, invID int64, signatureValue string, userID string) bool {
	// Подпись от Result: MD5(OutSum:InvId:Password2:Shp_userId=value)
	signStr := fmt.Sprintf("%s:%d:%s:Shp_userId=%s",
		outSum, invID, r.Password2, userID)
	expectedSignature := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(signStr))))
	receivedSignature := strings.ToUpper(signatureValue)

	fmt.Printf("[Robokassa] Result signature check:\n")
	fmt.Printf("  SignString: %s\n", signStr)
	fmt.Printf("  Expected:   %s\n", expectedSignature)
	fmt.Printf("  Received:   %s\n", receivedSignature)
	fmt.Printf("  Match:      %v\n", expectedSignature == receivedSignature)

	return expectedSignature == receivedSignature
}

// VerifySuccessSignature проверяет подпись от SuccessURL (страница успеха)
func (r *Robokassa) VerifySuccessSignature(outSum string, invID int64, signatureValue string, userID string) bool {
	// Подпись от Success: MD5(MerchantLogin:OutSum:InvId:Password1:Shp_userId)
	signStr := fmt.Sprintf("%s:%s:%d:%s:Shp_userId=%s",
		r.MerchantLogin, outSum, invID, r.Password1, userID)
	expectedSignature := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(signStr))))
	receivedSignature := strings.ToUpper(signatureValue)

	return expectedSignature == receivedSignature
}
