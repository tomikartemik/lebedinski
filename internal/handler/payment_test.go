package handler

import (
	"bytes"
	"errors"
	"lebedinski/internal/model"
	"lebedinski/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type paymentTestPaymentService struct {
	payment *model.PaymentResponse
	err     error
}

func (s *paymentTestPaymentService) CreatePayment(model.Order) (*model.PaymentResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *paymentTestPaymentService) GetPayment(string) (*model.PaymentResponse, error) {
	return s.payment, s.err
}

type paymentTestOrderService struct {
	mu              sync.Mutex
	order           model.Order
	claimed         bool
	markCalls       int
	failed          chan string
	confirmationRun chan struct{}
}

func (s *paymentTestOrderService) ProcessOrder(model.Order, string) error { return nil }
func (s *paymentTestOrderService) GetAllOrders() ([]model.Order, error)   { return nil, nil }
func (s *paymentTestOrderService) GetOrderByCartID(int) (model.Order, error) {
	return s.order, nil
}
func (s *paymentTestOrderService) GetOrderByPaymentID(string) (model.Order, error) {
	return s.order, nil
}
func (s *paymentTestOrderService) SendOrderConfirmation(string, string) error {
	if s.confirmationRun != nil {
		s.confirmationRun <- struct{}{}
	}
	return nil
}
func (s *paymentTestOrderService) SendOrderShippedNotification(string) error { return nil }
func (s *paymentTestOrderService) DeleteOrder(int) error                     { return nil }
func (s *paymentTestOrderService) UpdateOrder(model.Order) error             { return nil }
func (s *paymentTestOrderService) ChangeStatus(int, string) error            { return nil }
func (s *paymentTestOrderService) MarkPaymentSucceeded(string) (model.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.markCalls++
	s.order.Status = "Paid"
	s.order.PaymentStatus = "Succeeded"
	return s.order, nil
}
func (s *paymentTestOrderService) ClaimFulfillment(string) (bool, error) {
	return s.claimed, nil
}
func (s *paymentTestOrderService) SetFulfillmentFailed(_ string, message string) error {
	if s.failed != nil {
		s.failed <- message
	}
	return nil
}

type paymentTestCdekService struct {
	err    error
	called chan struct{}
}

func (s *paymentTestCdekService) GetToken() (string, error) { return "", nil }
func (s *paymentTestCdekService) CreateCdekOrder(string) (string, error) {
	if s.called != nil {
		s.called <- struct{}{}
	}
	return "", s.err
}
func (s *paymentTestCdekService) GetPvzList(map[string]string) ([]model.Pvz, error) {
	return nil, nil
}

// Compile-time guards keep the test doubles aligned with the production interfaces.
var _ service.Payment = (*paymentTestPaymentService)(nil)
var _ service.Order = (*paymentTestOrderService)(nil)
var _ service.Cdek = (*paymentTestCdekService)(nil)

func paymentTestHandler(payment *model.PaymentResponse, orderService *paymentTestOrderService, cdekError error) *Handler {
	return &Handler{services: &service.Service{
		Payment: &paymentTestPaymentService{payment: payment},
		Order:   orderService,
		Cdek:    &paymentTestCdekService{err: cdekError, called: make(chan struct{}, 1)},
	}}
}

func paymentWebhookRequest(t *testing.T, h *Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/payment/response", h.HandleWebhook)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/payment/response", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	return response
}

func TestPaidWebhookPersistsPaymentBeforeCdekFailure(t *testing.T) {
	payment := &model.PaymentResponse{
		ID:          "payment-1816",
		Status:      "succeeded",
		Paid:        true,
		Description: "1816",
		Amount:      model.Amount{Value: "5750.00", Currency: "RUB"},
	}
	orders := &paymentTestOrderService{
		order:           model.Order{CartID: 1816, PaymentID: payment.ID, Status: "Not Paid"},
		claimed:         true,
		failed:          make(chan string, 1),
		confirmationRun: make(chan struct{}, 1),
	}
	h := paymentTestHandler(payment, orders, errors.New("CDEK [appropriate_pickup_point_not_found]: invalid point"))

	response := paymentWebhookRequest(t, h, `{"event":"payment.succeeded","object":{"id":"payment-1816","status":"succeeded"}}`)
	if response.Code != http.StatusOK {
		t.Fatalf("webhook returned %d, want 200", response.Code)
	}

	select {
	case message := <-orders.failed:
		if !strings.Contains(message, "Оплата подтверждена") {
			t.Fatalf("warning does not preserve payment truth: %q", message)
		}
	case <-time.After(time.Second):
		t.Fatal("fulfillment failure was not persisted")
	}

	orders.mu.Lock()
	markCalls := orders.markCalls
	status := orders.order.Status
	orders.mu.Unlock()
	if markCalls != 1 || status != "Paid" {
		t.Fatalf("payment was not marked first: calls=%d status=%q", markCalls, status)
	}
	select {
	case <-orders.confirmationRun:
		t.Fatal("confirmation ran after CDEK failure")
	default:
	}
}

func TestWebhookRejectsUnconfirmedProviderPayment(t *testing.T) {
	payment := &model.PaymentResponse{ID: "payment-1816", Status: "canceled", Paid: false, Description: "1816"}
	orders := &paymentTestOrderService{order: model.Order{CartID: 1816, PaymentID: payment.ID}, claimed: true}
	h := paymentTestHandler(payment, orders, nil)

	response := paymentWebhookRequest(t, h, `{"event":"payment.succeeded","object":{"id":"payment-1816","status":"succeeded"}}`)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("webhook returned %d, want 400", response.Code)
	}
	if orders.markCalls != 0 {
		t.Fatal("unconfirmed payment was marked paid")
	}
}

func TestDuplicateWebhookDoesNotRepeatFulfillment(t *testing.T) {
	payment := &model.PaymentResponse{ID: "payment-1816", Status: "succeeded", Paid: true, Description: "1816"}
	orders := &paymentTestOrderService{order: model.Order{CartID: 1816, PaymentID: payment.ID}, claimed: false}
	cdekCalled := make(chan struct{}, 1)
	h := &Handler{services: &service.Service{
		Payment: &paymentTestPaymentService{payment: payment},
		Order:   orders,
		Cdek:    &paymentTestCdekService{called: cdekCalled},
	}}

	response := paymentWebhookRequest(t, h, `{"event":"payment.succeeded","object":{"id":"payment-1816","status":"succeeded"}}`)
	if response.Code != http.StatusOK {
		t.Fatalf("webhook returned %d, want 200", response.Code)
	}
	select {
	case <-cdekCalled:
		t.Fatal("duplicate webhook repeated CDEK fulfillment")
	case <-time.After(50 * time.Millisecond):
	}
}
