package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaim struct {
	MapClaim jwt.MapClaims
	jwt.StandardClaims
}

// Decoding JWT to get payload, not verifying JWT
func Decode(JWTToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(JWTToken, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if err.Error() == jwt.ErrInvalidKeyType.Error() {
		return token, nil
	}
	return nil, err
}

// Generate HS256 JWT token
func GenerateHS256JWT(payload map[string]interface{}, exprireAt time.Time) (string, error) {
	mapClaims := jwt.MapClaims{}
	for key, val := range payload {
		mapClaims[key] = val
	}

	claims := &JWTClaim{
		MapClaim: mapClaims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exprireAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return signedToken, err
}

// Verify JWT func
func VerifyJWT(tokenString string) bool {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return false
	}

	return token.Valid
}
