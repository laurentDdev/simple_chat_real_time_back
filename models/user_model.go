package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Pseudo   string `json:"pseudo"`
	Email    string `gorm:"unique" json:"email"`
	Password string
}
