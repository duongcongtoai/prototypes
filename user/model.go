package user

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(100);unique_index"`
	Password string `gorm:"type:vachar(255)"`
}

func (user User) TableName() string {
	return "user"
}
