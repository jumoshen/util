package tfa

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
	"net/url"
	"time"
)

type TfActorAuthenticator struct {
	Issuer    string
	Algorithm string
	Period    int64
	Digits    int
}

func (tfa TfActorAuthenticator) GenQrCode(account, secret string) ([]byte, error) {
	authUrl := fmt.Sprintf(
		`otpauth://totp/%s?secret=%s&issuer=%s&algorithm=%s&digits=%d&period=%d`,
		url.QueryEscape(account), secret, tfa.Issuer,
		tfa.Algorithm, tfa.Digits, tfa.Period,
	)
	return qrcode.Encode(authUrl, qrcode.Medium, 256)
}

func (tfa *TfActorAuthenticator) TotpString(secret string) string {
	// base32编码秘钥：K
	key := make([]byte, base32.StdEncoding.DecodedLen(len(secret)))
	base32.StdEncoding.Decode(key, []byte(secret))

	// 根据当前时间计算随机数：C
	message := make([]byte, 8)
	binary.BigEndian.PutUint64(message, uint64(time.Now().Unix()/tfa.Period))

	// 使用sha1对K和C做hmac得到20个字节的密串：HMAC-SHA-1(K, C)
	hmacsha1 := hmac.New(sha1.New, key)
	hmacsha1.Write(message)
	hash := hmacsha1.Sum([]byte{})

	// 从20个字节的密串取最后一个字节，取该字节的低四位
	offset := hash[len(hash)-1] & 0xF
	truncatedHash := hash[offset : offset+4]

	// 按照大端方式组成一个整数
	bin := (binary.BigEndian.Uint32(truncatedHash) & 0x7FFFFFFF) % 1000000

	// 将数字转成特定长度字符串，不够的补0
	return fmt.Sprintf(`%0*d`, tfa.Digits, bin)
}

// NewGoogleAuthenticator google双因子认证
func NewGoogleAuthenticator(issuer string) *TfActorAuthenticator {
	return &TfActorAuthenticator{
		Issuer:    issuer,
		Algorithm: "SHA1",
		Period:    30,
		Digits:    6,
	}
}
