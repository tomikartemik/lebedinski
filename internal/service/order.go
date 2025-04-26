package service

import (
	"fmt"
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

type OrderService struct {
	repoItem  repository.Item
	repoOrder repository.Order
	repoSize  repository.Size
	repoCart  repository.Cart
}

func NewOrderService(repoItem repository.Item, repoOrder repository.Order, repoSize repository.Size, repoCart repository.Cart) *OrderService {
	return &OrderService{
		repoItem:  repoItem,
		repoOrder: repoOrder,
		repoSize:  repoSize,
		repoCart:  repoCart,
	}
}

func (s *OrderService) ProcessOrder(order model.Order, paymentID string) error {
	order.PaymentID = paymentID
	order.Status = "Paid"

	cartItems, err := s.repoOrder.GetCartItemsByCartID(order.CartID)

	if err != nil {
		return fmt.Errorf("не удалось получить товары для CartID %d: %w", order.CartID, err)
	}

	if len(cartItems) == 0 {
		return fmt.Errorf("корзина с ID %d пуста или не найдена", order.CartID)
	}

	for _, item := range cartItems {
		err := s.repoSize.DecreaseStock(item.ItemID, item.Size, item.Quantity)
		if err != nil {
			return fmt.Errorf("не удалось списать остаток для ItemID %d, Size %s: %w", item.ItemID, item.Size, err)
		}
	}

	err = s.repoOrder.SaveOrder(order)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetAllOrders() ([]model.Order, error) {
	return s.repoOrder.GetAllOrders()
}

func (s *OrderService) GetOrderByCartID(id int) (model.Order, error) {
	order, err := s.repoOrder.GetOrderByCartID(id)

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *OrderService) SendOrderConfirmation(cartIDStr string) error {
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

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Подготовка данных о товарах
	var itemsHTML strings.Builder
	total := 0

	for _, cartItem := range cart.Items {
		// Получаем полную информацию о товаре
		item, err := s.repoItem.GetItemByID(cartItem.ItemID)
		if err != nil {
			return fmt.Errorf("не удалось получить информацию о товаре ID %d: %v", cartItem.ItemID, err)
		}

		itemTotal := cartItem.Quantity * item.ActualPrice
		total += itemTotal

		itemsHTML.WriteString(fmt.Sprintf(`
        <tr>
            <td style="padding: 10px; border-bottom: 1px solid #eee;">%s</td>
            <td style="padding: 10px; border-bottom: 1px solid #eee;">%s</td>
            <td style="padding: 10px; border-bottom: 1px solid #eee; text-align: center;">%d</td>
            <td style="padding: 10px; border-bottom: 1px solid #eee; text-align: right;">%d руб.</td>
            <td style="padding: 10px; border-bottom: 1px solid #eee; text-align: right;">%d руб.</td>
        </tr>`,
			item.Name,         // Название товара из Item
			cartItem.Size,     // Размер из CartItem
			cartItem.Quantity, // Количество из CartItem
			item.ActualPrice,  // Цена из Item
			itemTotal,         // Итоговая сумма за позицию
		))
	}
	// Если есть промокод, добавляем скидку
	discountHTML := ""
	if order.Promocode != "" {
		discountHTML = fmt.Sprintf(`
            <tr>
                <td colspan="4" style="padding: 10px; text-align: right; font-weight: bold;">Промокод "%s"</td>
                <td style="padding: 10px; text-align: right;">-%%s руб.</td>
            </tr>`, order.Promocode)
	}

	// MIME-заголовки для HTML-письма
	header := fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: Подтверждение заказа #%d\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n", order.Email, smtpUser, order.CartID)

	// HTML-тело письма
	body := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <style>
                body { font-family: Arial, sans-serif; color: #333; line-height: 1.6; }
                .container { max-width: 800px; margin: 0 auto; padding: 20px; }
                .header { background-color: #000; color: white; padding: 20px; text-align: center; }
                .content { padding: 20px; background-color: #f9f9f9; }
                table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
                th { background-color: #f2f2f2; padding: 10px; text-align: left; }
                .total { font-weight: bold; font-size: 18px; }
                .footer { margin-top: 30px; font-size: 12px; color: #777; text-align: center; }
                .info-block { background: white; padding: 15px; margin-bottom: 20px; border-radius: 5px; }
                .info-title { font-weight: bold; margin-bottom: 5px; }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h2>Ваш заказ #%d подтверждён</h2>
                </div>
                
                <div class="content">
                    <p>Уважаемый(ая) %s, спасибо за ваш заказ!</p>
                    
                    <div class="info-block">
                        <div class="info-title">Информация о заказе:</div>
                        <p>Номер заказа: #%d</p>
                        <p>Код отслеживания СДЭК: %s</p>
                        <p>Статус заказа: %s</p>
                    </div>
                    
                    <table>
                        <thead>
                            <tr>
                                <th>Товар</th>
                                <th>Размер</th>
                                <th>Кол-во</th>
                                <th>Цена</th>
                                <th>Сумма</th>
                            </tr>
                        </thead>
                        <tbody>
                            %s
                            %s
                        </tbody>
                        <tfoot>
                            <tr class="total">
                                <td colspan="4" style="text-align: right;">Итого:</td>
                                <td style="text-align: right;">%d руб.</td>
                            </tr>
                        </tfoot>
                    </table>
                    
                    <div class="info-block">
                        <div class="info-title">Информация о доставке:</div>
                        <p>Пункт выдачи: %s</p>
                        <p>Телефон для связи: %s</p>
                        %s
                    </div>
                    
                    <p>Мы свяжемся с вами для уточнения деталей. Вы можете отслеживать статус заказа в личном кабинете.</p>
                    
                    <div class="footer">
                        <p>Если у вас есть вопросы, пожалуйста, ответьте на это письмо.</p>
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
				return fmt.Sprintf("<p>Дополнительная информация: %s</p>", order.AdditionalInfo)
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
