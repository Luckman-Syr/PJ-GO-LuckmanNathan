package models

import (
	"project-akhir/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Username string `gorm:"not null;uniqueIndex" json:"username"  form:"username" valid:"required~Your username is required" `
	Email    string `gorm:"not null;uniqueIndex" json:"email"  form:"email" valid:"email~Please enter a valid email address, required~Your email is required"`
	Password string `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required, minstringlength(6)~Your password should be at least 6 characters long"`
	Age      int    `gorm:"not null" json:"age" form:"age" valid:"range(9|100)~You must be at least 9 years old"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HashPassword(u.Password)
	err = nil

	return
}
