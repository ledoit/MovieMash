package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const voterCookie = "mm_voter"

func ensureVoter(c *gin.Context) string {
	if value, err := c.Cookie(voterCookie); err == nil && len(value) >= 16 {
		return value
	}

	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	key := hex.EncodeToString(buf)

	secure := os.Getenv("COOKIE_SECURE") == "1"
	sameSite := http.SameSiteLaxMode
	if strings.EqualFold(os.Getenv("COOKIE_SAMESITE"), "none") {
		sameSite = http.SameSiteNoneMode
		secure = true
	}
	c.SetSameSite(sameSite)
	c.SetCookie(voterCookie, key, 400*24*3600, "/", "", secure, true)
	return key
}
