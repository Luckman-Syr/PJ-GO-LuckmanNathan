package models

import "github.com/asaskevich/govalidator"

type SocialMedia struct {
	GormModel
	Name           string `json:"name" gorm:"not null;" valid:"required~Your name is required"`
	SocialMediaUrl string `json:"social_media_url" gorm:"not null;" valid:"required~Your social media url is required"`
	UserID         uint
	User           User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
}

func (s *SocialMedia) BeforeCreate() (err error) {
	_, errCreate := govalidator.ValidateStruct(s)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil

	return
}
