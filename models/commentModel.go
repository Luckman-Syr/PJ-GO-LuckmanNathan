package models

import "github.com/asaskevich/govalidator"

type Comment struct {
	GormModel
	UserID  uint   `json:"user_id"`
	User    User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	PhotoID uint   `json:"photo_id"`
	Photo   Photo  `gorm:"foreignKey:PhotoID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photo"`
	Message string `json:"message" gorm:"not null;" valid:"required~Your message is required"`
}

func (c *Comment) BeforeCreate() (err error) {
	_, errCreate := govalidator.ValidateStruct(c)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil

	return
}
