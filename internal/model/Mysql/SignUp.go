package Mysql

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type SignUp struct {
	Base
	TeamName  string `gorm:"type:varchar(90);not null;" json:"teamName"`
	IsHDU     bool   `gorm:"type:bool;not null;" json:"isHDU"`
	School    string `gorm:"type:varchar(90);not null;" json:"school"`
	Member1ID string `gorm:"type:varchar(90);not null;" json:"member1ID"`
	Member2ID string `gorm:"type:varchar(90);not null;" json:"member2ID"`
	Member3ID string `gorm:"type:varchar(90);not null;" json:"member3ID"`
}

type Member struct {
	Base
	Phone          string `gorm:"type:varchar(30);not null;" json:"phone"`
	QQ             string `gorm:"type:varchar(30);not null;" json:"qq"`
	Name           string `gorm:"type:varchar(30);not null;" json:"name"`
	IDNumber       string `gorm:"type:varchar(20);not null;" json:"idNumber"`
	BankCardNumber string `gorm:"type:varchar(30);" json:"bankCardNumber"`
	BankName       string `gorm:"type:varchar(50);" json:"bankName"`
	HDUID          string `gorm:"type:varchar(10);" json:"hduId"`
	Role           string `gorm:"type:varchar(30);not null;" json:"isLeader"`
}

func (c *Member) BeforeCreate(tx *gorm.DB) (err error) { //使用钩子函数
	c.ID = uuid.NewV4().String()
	return
}

func (c *SignUp) BeforeCreate(tx *gorm.DB) (err error) { //使用钩子函数
	c.ID = uuid.NewV4().String()
	return
}
