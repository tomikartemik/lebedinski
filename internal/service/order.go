package service

import (
	"fmt"
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"net/smtp"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OrderService struct {
	repoItem      repository.Item
	repoOrder     repository.Order
	repoSize      repository.Size
	repoCart      repository.Cart
	repoPromoCode repository.PromoCode
}

func NewOrderService(repoItem repository.Item, repoOrder repository.Order, repoSize repository.Size, repoCart repository.Cart, repoPromoCode repository.PromoCode) *OrderService {
	return &OrderService{
		repoItem:      repoItem,
		repoOrder:     repoOrder,
		repoSize:      repoSize,
		repoCart:      repoCart,
		repoPromoCode: repoPromoCode,
	}
}

func (s *OrderService) ProcessOrder(order model.Order, paymentID string) error {
	order.PaymentID = paymentID
	order.Status = "Not Paid"
	order.DateTime = time.Now()

	cartItems, err := s.repoOrder.GetCartItemsByCartID(order.CartID)

	if err != nil {
		return fmt.Errorf("не удалось получить товары для CartID %d: %w", order.CartID, err)
	}

	if len(cartItems) == 0 {
		return fmt.Errorf("корзина с ID %d пуста или не найдена", order.CartID)
	}

	err = s.repoOrder.SaveOrder(order)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetAllOrders() ([]model.Order, error) {
	orders, err := s.repoOrder.GetAllOrders()
	if err != nil {
		return nil, err
	}

	// Сортировка по CartID в обратном порядке (по убыванию)
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CartID > orders[j].CartID
	})

	return orders, nil
}

func (s *OrderService) GetOrderByCartID(id int) (model.Order, error) {
	order, err := s.repoOrder.GetOrderByCartID(id)

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *OrderService) SendOrderConfirmation(cartIDStr, total string) error {
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		return err
	}

	order, err := s.repoOrder.GetOrderByCartID(cartID)
	if err != nil {
		return err
	}

	cart, err := s.repoCart.GetCartByID(cartID)
	if err != nil {
		return err
	}

	cartItems, err := s.repoOrder.GetCartItemsByCartID(order.CartID)
	for _, cartItem := range cartItems {
		item, err := s.repoItem.GetItemByID(cartItem.ItemID)
		if err != nil {
			return err
		}
		if item.CustomTailoring == false {
			err := s.repoSize.DecreaseStock(cartItem.ItemID, cartItem.Size, cartItem.Quantity)
			if err != nil {
				return fmt.Errorf("не удалось списать остаток для ItemID %d, Size %s: %w", cartItem.ItemID, cartItem.Size, err)
			}
		}
	}

	if order.Promocode != "" {
		promoCode, err := s.repoPromoCode.GetPromoCodeByCode(order.Promocode)
		if err == nil {
			promoCode.NumberOfUses--
			err = s.repoPromoCode.UpdatePromoCode(promoCode)
			if err != nil {
				fmt.Printf("Error: Failed to update promocode '%s' uses: %v\n", order.Promocode, err)
			}
		}
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	var itemsHTML strings.Builder

	for _, cartItem := range cart.Items {
		item, err := s.repoItem.GetItemByID(cartItem.ItemID)
		if err != nil {
			return fmt.Errorf("не удалось получить информацию о товаре ID %d: %v", cartItem.ItemID, err)
		}

		itemTotal := cartItem.Quantity * item.ActualPrice

		itemsHTML.WriteString(fmt.Sprintf(`
        <tr>
            <td style="padding: 10px; border-bottom: 1px solid #eee; word-break: break-word;">%s<br><span style="color: #777;">%s</span></td>
            <td style="padding: 10px; border-bottom: 1px solid #eee; text-align: center;">%d</td>
            <td style="padding: 10px; border-bottom: 1px solid #eee; text-align: right;">%d руб.</td>
            <td style="padding: 10px; border-bottom: 1px solid #eee; text-align: right;">%d руб.</td>
        </tr>`,
			item.Name,
			cartItem.Size,
			cartItem.Quantity,
			item.ActualPrice,
			itemTotal,
		))
	}

	discountHTML := ""
	if order.Promocode != "" {
		discountHTML = fmt.Sprintf(`
            <tr>
                <td colspan="3" style="padding: 10px; text-align: right; font-weight: bold;">Промокод "%s"</td>
                <td style="padding: 10px; text-align: right;">-%%s руб.</td>
            </tr>`, order.Promocode)
	}

	header := fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: Подтверждение заказа #%d\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n", order.Email, smtpUser, order.CartID)

	body := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <style>
                body { font-family: Arial, sans-serif; color: #333; line-height: 1.6; margin: 0; padding: 0; }
                .container { max-width: 100%%; width: 600px; margin: 0 auto; padding: 0; }
                .header { background-color: #000; color: white; padding: 30px 20px; text-align: center; }
                .brand { font-size: 24px; font-weight: bold; letter-spacing: 2px; margin-bottom: 10px; }
                .content { padding: 20px; background-color: #f9f9f9; }
                table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
                th { background-color: #f2f2f2; padding: 12px 10px; text-align: left; }
                td { padding: 10px; border-bottom: 1px solid #eee; }
                .total { font-weight: bold; font-size: 18px; }
                .footer { margin-top: 30px; font-size: 12px; color: #777; text-align: center; }
                .info-block { background: white; padding: 20px; margin-bottom: 20px; border-radius: 5px; box-shadow: 0 2px 5px rgba(0,0,0,0.05); }
                .info-title { font-weight: bold; margin-bottom: 10px; font-size: 16px; color: #000; }
                
                @media only screen and (max-width: 480px) {
                    table { width: 100%%; display: block; overflow-x: auto; }
                    .container { width: 100%% !important; }
                    .header h2 { font-size: 20px; }
                    .info-block { padding: 15px; }
                    .brand { font-size: 20px; }
                }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <div class="brand">LEBEDINSKI</div>
                    <h2>Ваш заказ #%d подтверждён</h2>
                </div>
                
                <div class="content">
                    <p style="font-size: 16px;">%s, благодарим Вас за покупку в LEBEDINSKI <br/>
                    Ваш заказ успешно оформлен и будет отправлен в указанные на сайте сроки<br/><br/>
                    Как только ваш заказ будет отправлен, пришлем вам сообщение об этом.<br/></p>
                    
                    <div class="info-block">
                        <div class="info-title">Детали заказа</div>
                        <p><strong>Номер заказа:</strong> #%d</p>
                        <p><strong>Код отслеживания:</strong> %s</p>
                        <p><strong>Статус:</strong> %s</p>
                    </div>
                    
                    <div class="info-block">
                        <div class="info-title">Состав заказа</div>
                        <table>
                            <thead>
                                <tr>
                                    <th style="width: 50%%;">Товар</th>
                                    <th style="text-align: center;">Кол-во</th>
                                    <th style="text-align: right;">Цена</th>
                                    <th style="text-align: right;">Сумма</th>
                                </tr>
                            </thead>
                            <tbody>
                                %s
                                %s
                            </tbody>
                            <tfoot>
                                <tr class="total">
                                    <td colspan="3" style="text-align: right;">Итого (с учетом скидок и доставки):</td>
                                    <td style="text-align: right;">%s руб.</td>
                                </tr>
                            </tfoot>
                        </table>
                    </div>
                    
                    <div class="info-block">
                        <div class="info-title">Доставка</div>
                        <p><strong>Пункт выдачи:</strong> %s</p>
                        <p><strong>Телефон:</strong> %s</p>
                        %s
                    </div>
                    
                    <div class="footer">
                        <p>Если у вас есть вопросы, напишите в телеграме @Lebedinski_help.</p>
                        <p>&copy; %d Lebedinski.shop</p>
                    </div>
                </div>
            </div>
        </body>
        </html>
    `,
		order.CartID,
		order.FullName,
		order.CartID,
		order.CdekOrderUUID,
		order.Status,
		itemsHTML.String(),
		discountHTML,
		total,
		order.PointCode,
		order.Phone,
		func() string {
			if order.AdditionalInfo != "" {
				return fmt.Sprintf("<p><strong>Дополнительно:</strong> %s</p>", order.AdditionalInfo)
			}
			return ""
		}(),
		time.Now().Year(),
	)

	msg := []byte(header + body)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{order.Email}, msg)
	if err != nil {
		return fmt.Errorf("ошибка при отправке email: %v", err)
	}
	return nil
}

func (s *OrderService) SendOrderShippedNotification(cartIDStr string) error {
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		return err
	}

	order, err := s.repoOrder.GetOrderByCartID(cartID)
	if err != nil {
		return err
	}

	order.Status = "Sent"
	err = s.repoOrder.UpdateOrder(order)
	if err != nil {
		return err
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	trackingURL := fmt.Sprintf("https://www.cdek.ru/ru/tracking/?order_id=%s", order.CdekOrderUUID)

	header := fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: Ваш заказ #%d отправлен\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n", order.Email, smtpUser, order.CartID)

	body := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <style>
                body { font-family: Arial, sans-serif; color: #333; line-height: 1.6; margin: 0; padding: 0; }
                .container { max-width: 600px; margin: 0 auto; padding: 20px; }
                .header { background-color: #000; color: white; padding: 30px 20px; text-align: center; }
                .brand { font-size: 24px; font-weight: bold; letter-spacing: 2px; margin-bottom: 10px; }
                .content { padding: 20px; }
                .tracking-number { font-size: 18px; font-weight: bold; margin: 20px 0; }
                .tracking-btn { 
                    display: inline-block; 
                    padding: 12px 24px; 
                    background-color: #000; 
                    color: white; 
                    text-decoration: none; 
                    border-radius: 4px; 
                    font-weight: bold;
                    margin: 15px 0;
                }
                .footer { margin-top: 30px; font-size: 12px; color: #777; text-align: center; }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <div class="brand">LEBEDINSKI</div>
                    <h2>Ваш заказ #%d отправлен</h2>
                </div>
                
                <div class="content">
                    <p>%s, ваш заказ был передан в службу доставки СДЭК.</p>
                    
                    <div class="tracking-number">Номер для отслеживания: %s</div>
                    
                    <a href="%s" class="tracking-btn">Отследить заказ</a>
                    
                    <p>Если возникнут вопросы, пишите в телеграм: @Lebedinski_help</p>
                </div>
                
                <div class="footer">
                    <p>&copy; %d Lebedinski.shop</p>
                </div>
            </div>
        </body>
        </html>
    `,
		order.CartID,
		order.FullName,
		order.CdekOrderUUID,
		trackingURL,
		time.Now().Year(),
	)

	msg := []byte(header + body)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{order.Email}, msg)
	if err != nil {
		return fmt.Errorf("ошибка при отправке email: %v", err)
	}
	return nil
}

func (s *OrderService) DeleteOrder(cartID int) error {
	return s.repoOrder.DeleteOrder(cartID)
}

func (s *OrderService) UpdateOrder(order model.Order) error {
	return s.repoOrder.UpdateOrder(order)
}
