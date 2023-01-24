package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateStateOauthCookie(c *gin.Context) string {
	expiration := time.Now().Add(2 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	return state
}
