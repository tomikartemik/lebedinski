package handler

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func configuredTestHandler(t *testing.T) *Handler {
	t.Helper()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct horse battery staple"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADMIN_USERNAME", "admin")
	t.Setenv("ADMIN_PASSWORD_HASH_B64", base64.StdEncoding.EncodeToString(passwordHash))
	t.Setenv("ALLOWED_ORIGINS", "https://admin.lebedinski.shop")
	t.Setenv("DADATA_API_TOKEN", "test-token")

	handler, err := NewHandler(nil)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func TestAdminSessionLifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := configuredTestHandler(t)
	router := gin.New()
	router.POST("/auth/login", handler.login)
	router.POST("/auth/logout", handler.logout)
	router.GET("/protected", handler.requireAdmin(), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	badLogin := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"admin","password":"wrong"}`))
	badLogin.Header.Set("Content-Type", "application/json")
	badResponse := httptest.NewRecorder()
	router.ServeHTTP(badResponse, badLogin)
	if badResponse.Code != http.StatusUnauthorized {
		t.Fatalf("wrong password returned %d, want %d", badResponse.Code, http.StatusUnauthorized)
	}

	login := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"admin","password":"correct horse battery staple"}`))
	login.Header.Set("Content-Type", "application/json")
	loginResponse := httptest.NewRecorder()
	router.ServeHTTP(loginResponse, login)
	if loginResponse.Code != http.StatusNoContent {
		t.Fatalf("login returned %d, want %d: %s", loginResponse.Code, http.StatusNoContent, loginResponse.Body.String())
	}
	cookies := loginResponse.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("login returned %d cookies, want 1", len(cookies))
	}
	sessionCookie := cookies[0]
	if sessionCookie.Name != adminSessionCookie || !sessionCookie.HttpOnly || !sessionCookie.Secure || sessionCookie.SameSite != http.SameSiteStrictMode {
		t.Fatalf("session cookie does not have the required security attributes: %#v", sessionCookie)
	}

	protected := httptest.NewRequest(http.MethodGet, "/protected", nil)
	protected.AddCookie(sessionCookie)
	protectedResponse := httptest.NewRecorder()
	router.ServeHTTP(protectedResponse, protected)
	if protectedResponse.Code != http.StatusNoContent {
		t.Fatalf("authenticated request returned %d, want %d", protectedResponse.Code, http.StatusNoContent)
	}

	logout := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	logout.AddCookie(sessionCookie)
	logoutResponse := httptest.NewRecorder()
	router.ServeHTTP(logoutResponse, logout)
	if logoutResponse.Code != http.StatusNoContent {
		t.Fatalf("logout returned %d, want %d", logoutResponse.Code, http.StatusNoContent)
	}

	afterLogout := httptest.NewRequest(http.MethodGet, "/protected", nil)
	afterLogout.AddCookie(sessionCookie)
	afterLogoutResponse := httptest.NewRecorder()
	router.ServeHTTP(afterLogoutResponse, afterLogout)
	if afterLogoutResponse.Code != http.StatusUnauthorized {
		t.Fatalf("logged-out session returned %d, want %d", afterLogoutResponse.Code, http.StatusUnauthorized)
	}
}

func TestSensitiveRoutesRequireAdminSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := configuredTestHandler(t).InitRoutes()

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/order/all"},
		{http.MethodGet, "/promocode/all"},
		{http.MethodPost, "/item/new"},
		{http.MethodDelete, "/order?cart_id=1"},
		{http.MethodPost, "/banner/upload"},
	}
	for _, test := range tests {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(test.method, test.path, nil)
		router.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized {
			t.Errorf("%s %s returned %d, want %d", test.method, test.path, response.Code, http.StatusUnauthorized)
		}
	}
}

func TestWildcardOriginIsRejected(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "*")
	if _, err := configuredOrigins(); err == nil {
		t.Fatal("configuredOrigins accepted a wildcard origin")
	}
}
