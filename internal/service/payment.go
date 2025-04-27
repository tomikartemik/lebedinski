package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"net/http"
	"os"
	"strconv"
	"time"
)

type PaymentService struct {
	repoItem      repository.Item
	repoCart      repository.Cart
	repoPromoCode repository.PromoCode
}

func NewPaymentService(repoItem repository.Item, repoCart repository.Cart, repoPromoCode repository.PromoCode) *PaymentService {
	return &PaymentService{
		repoItem:      repoItem,
		repoCart:      repoCart,
		repoPromoCode: repoPromoCode,
	}
}

func (s *PaymentService) CreatePayment(order model.Order) (*model.PaymentResponse, error) {
	idempotenceKey := uuid.New().String()

	cart, err := s.repoCart.GetCartByID(order.CartID)

	if err != nil {
		return nil, err
	}

	amount := 0.0

	for _, cartItem := range cart.Items {
		item, _ := s.repoItem.GetItemByID(cartItem.ItemID)
		amount = amount + float64(cartItem.Quantity*item.ActualPrice)
	}

	if order.Promocode != "" {
		promoCode, err := s.repoPromoCode.GetPromoCodeByCode(order.Promocode)
		if err != nil {
			fmt.Printf("Warning: Promocode '%s' not found or error fetching: %v\n", order.Promocode, err)
		} else {
			if promoCode.NumberOfUses <= 0 {
				fmt.Printf("Info: Promocode '%s' has no uses left.\n", order.Promocode)
			} else if time.Now().After(promoCode.EndDate) {
				fmt.Printf("Info: Promocode '%s' has expired.\n", order.Promocode)
			} else if amount < promoCode.MinAmount {
				fmt.Printf("Info: Order amount %.2f is less than minimum %.2f for promocode '%s'.\n", amount, promoCode.MinAmount, order.Promocode)
			} else {
				discount := amount * (promoCode.DiscountPercentage / 100.0)
				if promoCode.MaxDiscount > 0 && discount > promoCode.MaxDiscount {
					discount = promoCode.MaxDiscount
				}
				amount -= discount
				fmt.Printf("Info: Applied discount %.2f using promocode '%s'. New amount: %.2f\n", discount, order.Promocode, amount)
			}
		}
	}

	if amount < 15000.0 {
		amount += 350
	}

	paymentRequest := model.PaymentRequest{
		Amount: model.Amount{
			Value:    formatAmount(amount),
			Currency: "RUB",
		},
		Description: strconv.Itoa(order.CartID),
		Capture:     true,
		Confirmation: model.Confirmation{
			Type:      "redirect",
			ReturnURL: "https://lebedinski.shop",
		},
		PaymentMethod: model.PaymentMethod{
			Type: "bank_card",
		},
	}

	requestBody, err := json.Marshal(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.yookassa.ru/v3/payments",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	req.SetBasicAuth(os.Getenv("SHOP_ID"), os.Getenv("SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", idempotenceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var paymentResponse model.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
		return nil, fmt.Errorf("response parsing failed: %w", err)
	}

	return &paymentResponse, nil
}

func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}
