package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	dadataBaseURL       = "https://suggestions.dadata.ru/suggestions/api/4_1/rs"
	dadataMaxBodyBytes  = 8 << 10
	dadataMaxReplyBytes = 2 << 20
	dadataRateWindow    = time.Minute
	dadataRateLimit     = 60
)

type dadataProxy struct {
	token  string
	client *http.Client

	mu       sync.Mutex
	requests map[string]dadataRequestWindow
}

type dadataRequestWindow struct {
	started time.Time
	count   int
}

type dadataQuery struct {
	Query string `json:"query" binding:"required"`
	Count int    `json:"count,omitempty"`
}

func newDadataProxyFromEnv() (*dadataProxy, error) {
	token := strings.TrimSpace(os.Getenv("DADATA_API_TOKEN"))
	if token == "" {
		return nil, fmt.Errorf("DADATA_API_TOKEN must be configured")
	}
	return &dadataProxy{
		token:    token,
		client:   &http.Client{Timeout: 8 * time.Second},
		requests: make(map[string]dadataRequestWindow),
	}, nil
}

func (h *Handler) dadataSuggestAddress(c *gin.Context) {
	query, ok := h.validDadataQuery(c, true)
	if !ok {
		return
	}
	h.proxyDadata(c, "/suggest/address", query)
}

func (h *Handler) dadataFindDelivery(c *gin.Context) {
	query, ok := h.validDadataQuery(c, false)
	if !ok {
		return
	}
	h.proxyDadata(c, "/findById/delivery", gin.H{"query": query.Query})
}

func (h *Handler) validDadataQuery(c *gin.Context, withCount bool) (dadataQuery, bool) {
	if !h.dadata.allow(c.ClientIP(), time.Now()) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "address lookup rate limit exceeded"})
		return dadataQuery{}, false
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, dadataMaxBodyBytes)
	var query dadataQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid address lookup request"})
		return dadataQuery{}, false
	}
	query.Query = strings.TrimSpace(query.Query)
	if len(query.Query) < 2 || len(query.Query) > 512 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "address query must contain between 2 and 512 characters"})
		return dadataQuery{}, false
	}
	if withCount {
		if query.Count == 0 {
			query.Count = 5
		}
		if query.Count < 1 || query.Count > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "suggestion count must be between 1 and 10"})
			return dadataQuery{}, false
		}
	}
	return query, true
}

func (h *Handler) proxyDadata(c *gin.Context, endpoint string, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to prepare address lookup"})
		return
	}
	request, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, dadataBaseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to prepare address lookup"})
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Token "+h.dadata.token)

	response, err := h.dadata.client.Do(request)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "address lookup service unavailable"})
		return
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, dadataMaxReplyBytes))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "invalid address lookup response"})
		return
	}
	c.Data(response.StatusCode, "application/json; charset=utf-8", responseBody)
}

func (d *dadataProxy) allow(clientIP string, now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	for ip, window := range d.requests {
		if now.Sub(window.started) >= dadataRateWindow {
			delete(d.requests, ip)
		}
	}
	window, exists := d.requests[clientIP]
	if !exists {
		d.requests[clientIP] = dadataRequestWindow{started: now, count: 1}
		return true
	}
	if window.count >= dadataRateLimit {
		return false
	}
	window.count++
	d.requests[clientIP] = window
	return true
}
