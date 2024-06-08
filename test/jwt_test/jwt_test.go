package jwt_test

import (
	"os"
	"testing"
	"time"

	"fiber-admin/test/common"
	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
)

const (
	sub            = "test"
	invalidToken   = "invalid token"
	invalidSubject = ""
)

func TestMain(m *testing.M) {
	if err := common.Setup(); err != nil {
		panic(err)
	}
	code := m.Run()
	if err := common.Teardown(); err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestJwtGenerateAccessToken(t *testing.T) {

	var (
		injector = wire.GetInjector()
		j        = injector.Jwt
		err      error
	)
	accessToken, err := j.GenerateAccessToken(sub)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	t.Logf("access token: %s", accessToken)

	accessToken, err = j.GenerateAccessToken(invalidSubject)
	assert.Error(t, err)
	assert.Empty(t, accessToken)
}

func TestJwtGenerateRefreshToken(t *testing.T) {
	var (
		injector = wire.GetInjector()
		j        = injector.Jwt
		err      error
	)
	refreshToken, err := j.GenerateRefreshToken(sub)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)
	t.Logf("refresh token: %s", refreshToken)

	invalidSubject := ""
	refreshToken, err = j.GenerateRefreshToken(invalidSubject)
	assert.Error(t, err)
	assert.Empty(t, refreshToken)
}

func TestJwtVerifyToken(t *testing.T) {
	var (
		injector = wire.GetInjector()
		j        = injector.Jwt
		err      error
	)
	accessToken, err := j.GenerateAccessToken(sub)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	t.Logf("access token: %s", accessToken)

	token, err := j.VerifyAccessToken(accessToken)
	assert.NoError(t, err)
	assert.Equal(t, sub, token)
	t.Logf("token: %s", token)

	invalidToken := "Invalid token"
	token, err = j.VerifyAccessToken(invalidToken)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestJwtRefreshToken(t *testing.T) {
	var (
		injector = wire.GetInjector()
		j        = injector.Jwt
		err      error
	)
	refreshToken, err := j.GenerateRefreshToken(sub)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)
	t.Logf("refresh token: %s", refreshToken)

	time.Sleep(injector.Config.JWTConfig.TokenDuration - injector.Config.JWTConfig.RefreshDuration + 1*time.Second) // wait for token to expire
	newAccessToken, err := j.RefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
	t.Logf("new access token: %s", newAccessToken)

	invalidToken := "Invalid token"
	newAccessToken, err = j.RefreshToken(invalidToken)
	assert.Error(t, err)
	assert.Empty(t, newAccessToken)
}

func TestJwtExtractClaims(t *testing.T) {
	var (
		injector = wire.GetInjector()
		j        = injector.Jwt
		err      error
	)
	accessToken, err := j.GenerateAccessToken(sub)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	t.Logf("access token: %s", accessToken)

	claims, err := j.ExtractClaims(accessToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
	t.Logf("claims: %v", claims)

	claims, err = j.ExtractClaims(invalidToken)
	assert.Error(t, err)
	assert.Empty(t, claims)
}
