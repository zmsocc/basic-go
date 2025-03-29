//go:build manual

package wechat

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

// 手动跑的。提前验证代码
func Test_service_manual_VerifyCode(t *testing.T) {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID env variable is not set")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("WECHAT_APP_SECRET env variable is not set")
	}
	svc := NewService(appId, appKey, http.DefaultClient)
	// 这个code要微信扫码之后，复制过来用，会过期
	res, err := svc.VerifyCode(context.Background(), "001bKN100hrxEQ1sNj000ZSPHf2bKN1u", "state")
	require.NoError(t, err)
	t.Log(res)
}
