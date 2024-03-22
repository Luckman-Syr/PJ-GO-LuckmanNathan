package config

import (
	"github.com/cloudinary/cloudinary-go/v2"
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	cldSecret := "EFNh0ORtSgSqnxfEGQUC-euuqHc"
	cldName := "dwfsqkchc"
	cldKey := "673597216942313"

	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
	if err != nil {
		return nil, err
	}

	return cld, nil
}
