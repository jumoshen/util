package tfa

import (
	"encoding/base64"
	"fmt"
	"github.com/jumoshen/util"
	"os"
	"testing"
)

var (
	authClient *TfActorAuthenticator
	secret     string
)

func TestMain(m *testing.M) {
	authClient = NewGoogleAuthenticator("jumoshen")
	os.Exit(m.Run())
}

func TestTfActorAuthenticator_GenQrCode(t *testing.T) {
	secret = util.RandBase32String(32)
	codeData, err := authClient.GenQrCode("jumoshen.cn", secret)
	if err != nil {
		t.Errorf("gen code error:%#v", err)
	}

	tfaCrCode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(codeData)

	fmt.Println(secret, tfaCrCode)
}

func TestTfActorAuthenticator_TotpString(t *testing.T) {
	topString := authClient.TotpString("UGJN4ODXXUW4OC74N332BSUPNED5RAQA")

	fmt.Println(topString)
}
