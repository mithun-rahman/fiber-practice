package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;"`
	Username string    `json:"username" validate:"required,min=3,max=15"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password,omitempty" validate:"digit,uppercase,lowercase,special,required,min=3"`
	Address  string    `json:"address"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Users struct {
	Users []User `json:"users"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}
