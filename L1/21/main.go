package main

import "fmt"

// Старая платежная система (Adaptee)
type LegacyPayment struct{}

func (l *LegacyPayment) MakePayment(amount float64, customerID string) string {
	return fmt.Sprintf("Legacy payment processed: $%.2f for customer %s", amount, customerID)
}

// Новый интерфейс платежной системы, который мы хотим использовать (Target)
type ModernPayment interface {
	Pay(amount float64) string
}

// Адаптер для старой системы (Adapter)
type LegacyPaymentAdapter struct {
	legacyPayment *LegacyPayment
	customerID    string
}

func (a *LegacyPaymentAdapter) Pay(amount float64) string {
	return a.legacyPayment.MakePayment(amount, a.customerID)
}

func NewLegacyPaymentAdapter(customerID string) ModernPayment {
	return &LegacyPaymentAdapter{
		legacyPayment: &LegacyPayment{},
		customerID:    customerID,
	}
}

// Новая платежная система (уже реализует ModernPayment)
type NewPayment struct{}

func (n *NewPayment) Pay(amount float64) string {
	return fmt.Sprintf("New payment processed: $%.2f", amount)
}

// Клиентский код, который работает с ModernPayment
func ProcessPayment(payment ModernPayment, amount float64) {
	fmt.Println(payment.Pay(amount))
}

func main() {
	// Используем новую систему напрямую
	newPayment := &NewPayment{}
	ProcessPayment(newPayment, 123.45)

	// Используем старую систему через адаптер
	adapter := NewLegacyPaymentAdapter("new123")
	ProcessPayment(adapter, 1.75)
}
