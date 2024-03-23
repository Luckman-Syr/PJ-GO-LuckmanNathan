package models

import "github.com/asaskevich/govalidator"

type Photo struct {
	GormModel
	Title     string `json:"title" gorm:"not null;" valid:"required~Your title is required"`
	Caption   string `json:"caption" gorm:"not null;"`
	Photo_url string `json:"photo_url" gorm:"not null;" valid:"required~Your photo is required"`
	UserID    uint   `json:"user_id" gorm:"not null;"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"User"`
}

func (p *Photo) BeforeCreate() (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil

	return
}
