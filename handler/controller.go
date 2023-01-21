package handler

import (
	"crypto/rsa"
	"gorm.io/gorm"
)

type Controller struct {
	DB         *gorm.DB
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}
