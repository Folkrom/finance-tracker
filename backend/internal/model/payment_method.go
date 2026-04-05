package model

type PaymentMethodType string

const (
	PaymentMethodCash          PaymentMethodType = "cash"
	PaymentMethodDebitCard     PaymentMethodType = "debit_card"
	PaymentMethodCreditCard    PaymentMethodType = "credit_card"
	PaymentMethodDigitalWallet PaymentMethodType = "digital_wallet"
	PaymentMethodCrypto        PaymentMethodType = "crypto"
)

type PaymentMethod struct {
	Base
	Name    string            `gorm:"type:varchar(255);not null" json:"name"`
	Type    PaymentMethodType `gorm:"type:varchar(50);not null" json:"type"`
	Details *string           `gorm:"type:varchar(255)" json:"details,omitempty"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}
