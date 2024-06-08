package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AccessAudience  = "access"
	RefreshAudience = "refresh"
)

var (
	jwtInstance *Jwt
	once        sync.Once
)

type Jwt struct {
	privateKey      *ecdsa.PrivateKey
	publicKey       *ecdsa.PublicKey
	tokenDuration   time.Duration
	refreshDuration time.Duration
	refreshBuffer   time.Duration
}

func New(
	privateKey *ecdsa.PrivateKey, tokenDuration, refreshDuration time.Duration,
	refreshBuffer time.Duration,
) (*Jwt, error) {
	var err error
	once.Do(
		func() {
			j := &Jwt{
				privateKey:      privateKey,
				publicKey:       &privateKey.PublicKey,
				tokenDuration:   tokenDuration,
				refreshDuration: refreshDuration,
				refreshBuffer:   refreshBuffer,
			}
			if err = j.checkJWT(); err == nil {
				jwtInstance = j
			}
		},
	)
	return jwtInstance, err
}

func Update(privateKey *ecdsa.PrivateKey, tokenDuration, refreshDuration, refreshBuffer time.Duration) error {
	var err error
	if privateKey == nil {
		privateKey = jwtInstance.privateKey
	}
	jwtInstance = &Jwt{
		privateKey:      privateKey,
		publicKey:       &privateKey.PublicKey,
		tokenDuration:   tokenDuration,
		refreshDuration: refreshDuration,
		refreshBuffer:   refreshBuffer,
	}
	if err = jwtInstance.checkJWT(); err != nil {
		return err
	}
	return nil
}

func (j *Jwt) UpdateKey() error {
	// XXX: can be chosen from elliptic.P224(), elliptic.P256(), elliptic.P384(), elliptic.P521()
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	j.privateKey = privateKey
	j.publicKey = &privateKey.PublicKey
	return nil
}

func (j *Jwt) checkJWT() error {
	if j.privateKey == nil {
		return fmt.Errorf("private key is nil")
	}
	if j.tokenDuration == 0 {
		return fmt.Errorf("token duration is 0")
	}
	if j.refreshDuration == 0 {
		return fmt.Errorf("refresh duration is 0")
	}
	if j.refreshBuffer == 0 {
		return fmt.Errorf("refresh buffer is 0")
	}
	if j.tokenDuration > j.refreshDuration || j.refreshDuration < j.refreshBuffer || j.tokenDuration < j.refreshBuffer {
		return fmt.Errorf("invalid token, refresh or buffer duration")
	}

	return nil
}

func (j *Jwt) GenerateAccessToken(subject string) (string, error) {
	if subject == "" {
		return "", fmt.Errorf("subject is empty") // TODO: CHANGE ERROR TYPE
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodES256, &jwt.StandardClaims{
			Subject:   subject,
			Audience:  AccessAudience,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(j.tokenDuration).Unix(),
			NotBefore: time.Now().Unix(),
		},
	)

	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *Jwt) GenerateRefreshToken(subject string) (string, error) {
	if subject == "" {
		return "", fmt.Errorf("subject is empty") // TODO: CHANGE ERROR TYPE
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodES256, &jwt.StandardClaims{
			Subject:   subject,
			Audience:  RefreshAudience,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(j.refreshDuration).Unix(),
			NotBefore: time.Now().Unix(),
		},
	)

	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *Jwt) RefreshToken(token string) (string, error) {
	claims, err := j.ExtractClaims(token)
	if err != nil {
		return "", err
	}

	// Check if token is expired
	if time.Unix(int64(claims["exp"].(float64)), 0).Sub(time.Now()) > j.refreshBuffer {
		return "", fmt.Errorf(
			"token is not expired yet: %v", time.Unix(int64(claims["exp"].(float64)), 0).Sub(time.Now()),
		)
	}

	// Generate new token
	newToken, err := j.GenerateAccessToken(claims["sub"].(string))
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (j *Jwt) ExtractClaims(token string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}

	t, err := jwt.ParseWithClaims(
		token, claims, func(token *jwt.Token) (interface{}, error) {
			return j.publicKey, nil
		},
	)
	if err != nil || !t.Valid {
		return nil, err
	}

	return claims, nil
}

func (j *Jwt) VerifyAccessToken(token string) (string, error) {
	claims, err := j.ExtractClaims(token)
	if err != nil {
		return "", err
	}

	if claims["aud"] != AccessAudience {
		return "", fmt.Errorf("invalid audience")
	}
	return claims["sub"].(string), nil
}

func (j *Jwt) VerifyRefreshToken(token string) (string, error) {
	claims, err := j.ExtractClaims(token)
	if err != nil {
		return "", err
	}

	if claims["aud"] != RefreshAudience {
		return "", fmt.Errorf("invalid audience")
	}
	return claims["sub"].(string), nil
}
