package middleware

import "testing"

func TestGenAuthTokenByOpenIdAndUserId(t *testing.T) {
	token, _ := GenAuthTokenByOpenIdAndUserId(1)
	t.Logf("%v", token)
}
