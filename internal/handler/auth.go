package handler

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	adminSessionCookie = "__Host-lebedinski_admin"
	adminSessionTTL    = 12 * time.Hour
	loginWindow        = 15 * time.Minute
	maxLoginFailures   = 5
)

type adminAuth struct {
	username     string
	passwordHash []byte

	mu       sync.Mutex
	sessions map[string]time.Time
	attempts map[string]loginAttempts
}

type loginAttempts struct {
	windowStarted time.Time
	failures      int
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func newAdminAuthFromEnv() (*adminAuth, error) {
	username := os.Getenv("ADMIN_USERNAME")
	passwordHashBase64 := os.Getenv("ADMIN_PASSWORD_HASH_B64")
	if username == "" || passwordHashBase64 == "" {
		return nil, fmt.Errorf("ADMIN_USERNAME and ADMIN_PASSWORD_HASH_B64 must be configured")
	}
	passwordHash, err := base64.StdEncoding.DecodeString(passwordHashBase64)
	if err != nil {
		return nil, fmt.Errorf("ADMIN_PASSWORD_HASH_B64 is not valid base64: %w", err)
	}
	if _, err := bcrypt.Cost(passwordHash); err != nil {
		return nil, fmt.Errorf("ADMIN_PASSWORD_HASH_B64 does not contain a valid bcrypt hash: %w", err)
	}

	return &adminAuth{
		username:     username,
		passwordHash: passwordHash,
		sessions:     make(map[string]time.Time),
		attempts:     make(map[string]loginAttempts),
	}, nil
}

func (h *Handler) login(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	if h.auth.loginBlocked(c.ClientIP(), time.Now()) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "too many login attempts; try again later"})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 8<<10)
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid login request"})
		return
	}

	usernameMatches := subtle.ConstantTimeCompare([]byte(request.Username), []byte(h.auth.username)) == 1
	passwordMatches := bcrypt.CompareHashAndPassword(h.auth.passwordHash, []byte(request.Password)) == nil
	if !usernameMatches || !passwordMatches {
		h.auth.recordLoginFailure(c.ClientIP(), time.Now())
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	token, err := newSessionToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to create session"})
		return
	}
	expires := time.Now().Add(adminSessionTTL)
	h.auth.storeSession(token, expires, c.ClientIP())
	setSessionCookie(c, token, expires, int(adminSessionTTL.Seconds()))
	c.Status(http.StatusNoContent)
}

func (h *Handler) logout(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	if token, err := c.Cookie(adminSessionCookie); err == nil {
		h.auth.deleteSession(token)
	}
	setSessionCookie(c, "", time.Unix(1, 0), -1)
	c.Status(http.StatusNoContent)
}

func (h *Handler) sessionStatus(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{"authenticated": true})
}

func (h *Handler) requireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(adminSessionCookie)
		if err != nil || !h.auth.validSession(token, time.Now()) {
			c.Header("Cache-Control", "no-store")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authentication required"})
			return
		}
		c.Next()
	}
}

func newSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func setSessionCookie(c *gin.Context, value string, expires time.Time, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     adminSessionCookie,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (a *adminAuth) storeSession(token string, expires time.Time, clientIP string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.removeExpiredLocked(time.Now())
	a.sessions[token] = expires
	delete(a.attempts, clientIP)
}

func (a *adminAuth) validSession(token string, now time.Time) bool {
	if token == "" {
		return false
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	expires, exists := a.sessions[token]
	if !exists || !expires.After(now) {
		delete(a.sessions, token)
		return false
	}
	return true
}

func (a *adminAuth) deleteSession(token string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.sessions, token)
}

func (a *adminAuth) loginBlocked(clientIP string, now time.Time) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.removeExpiredLocked(now)
	attempt, exists := a.attempts[clientIP]
	if !exists || now.Sub(attempt.windowStarted) >= loginWindow {
		return false
	}
	return attempt.failures >= maxLoginFailures
}

func (a *adminAuth) recordLoginFailure(clientIP string, now time.Time) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.removeExpiredLocked(now)
	attempt, exists := a.attempts[clientIP]
	if !exists || now.Sub(attempt.windowStarted) >= loginWindow {
		a.attempts[clientIP] = loginAttempts{windowStarted: now, failures: 1}
		return
	}
	attempt.failures++
	a.attempts[clientIP] = attempt
}

func (a *adminAuth) removeExpiredLocked(now time.Time) {
	for token, expires := range a.sessions {
		if !expires.After(now) {
			delete(a.sessions, token)
		}
	}
	for clientIP, attempt := range a.attempts {
		if now.Sub(attempt.windowStarted) >= loginWindow {
			delete(a.attempts, clientIP)
		}
	}
}
