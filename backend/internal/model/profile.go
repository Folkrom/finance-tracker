package model

type Profile struct {
	Base
	Currency string `gorm:"type:varchar(3);not null;default:'MXN'" json:"currency"`
	Language string `gorm:"type:varchar(5);not null;default:'en'" json:"language"`
}

func (Profile) TableName() string {
	return "profiles"
}
