package controllers

import (
	"fmt"
	"project-akhir/database"
	"project-akhir/helpers"
	"project-akhir/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type resSocialMedia struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	SocialMediaUrl string `json:"social_media_url"`
	UserID         uint   `json:"user_id"`
}

func SocialMediaCreate(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	db := database.GetDB()
	socialMedia := models.SocialMedia{}
	contentType := helpers.GetContentType(c)

	if contentType != appJSON {
		c.ShouldBind(&socialMedia)
	} else {
		c.ShouldBindJSON(&socialMedia)
	}

	socialMedia.UserID = uint(userData["id"].(float64))

	fmt.Println(socialMedia)

	err := db.Create(&socialMedia).Error

	if err != nil {
		helpers.BadRequest(c, "Failed to create social media")
		return
	}

	type SocialMediaResponse struct {
		resSocialMedia
		CreatedAt string `json:"created_at"`
	}

	data := SocialMediaResponse{
		resSocialMedia: resSocialMedia{
			ID:             socialMedia.ID,
			Name:           socialMedia.Name,
			SocialMediaUrl: socialMedia.SocialMediaUrl,
			UserID:         socialMedia.UserID,
		},
		CreatedAt: socialMedia.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	helpers.Created(c, "Social media created successfully", data)

}

func SocialMediaList(c *gin.Context) {
	db := database.GetDB()
	socialMedia := []models.SocialMedia{}
	userData := c.MustGet("userData").(jwt.MapClaims)

	err := db.Where("user_id = ?", uint(userData["id"].(float64))).Preload("User").Find(&socialMedia).Error

	if err != nil {
		helpers.NotFound(c, "Social media not found")
		return
	}

	type resUser struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}

	type SocialMediaResponse struct {
		resSocialMedia
		CreatedAt string  `json:"created_at"`
		UpdateAt  string  `json:"update_at"`
		User      resUser `json:"User"`
	}

	var data []SocialMediaResponse

	for _, v := range socialMedia {
		data = append(data, SocialMediaResponse{
			resSocialMedia: resSocialMedia{
				ID:             v.ID,
				Name:           v.Name,
				SocialMediaUrl: v.SocialMediaUrl,
				UserID:         v.UserID,
			},
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateAt:  v.UpdatedAt.Format("2006-01-02 15:04:05"),
			User: resUser{
				ID:       v.User.ID,
				Username: v.User.Username,
			},
		})
	}

	helpers.Ok(c, "Social media found", data)
}

func SocialMediaUpdate(c *gin.Context) {
	db := database.GetDB()
	socialMedia := models.SocialMedia{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	id := c.Param("socialMediaId")

	err := db.Where("user_id = ?", uint(userData["id"].(float64))).First(&socialMedia, id).Error

	if err != nil {
		helpers.Unauthorized(c, "Social media not found or you don't have access to update this social media")
		return
	}

	if contentType != appJSON {
		c.ShouldBind(&socialMedia)
	} else {
		c.ShouldBindJSON(&socialMedia)
	}

	err = db.Model(&socialMedia).Updates(&socialMedia).Error

	if err != nil {
		helpers.BadRequest(c, "Failed to update social media")
		return
	}

	type SocialMediaResponse struct {
		resSocialMedia
		UpdateAt string `json:"update_at"`
	}

	data := SocialMediaResponse{
		resSocialMedia: resSocialMedia{
			ID:             socialMedia.ID,
			Name:           socialMedia.Name,
			SocialMediaUrl: socialMedia.SocialMediaUrl,
			UserID:         socialMedia.UserID,
		},
		UpdateAt: socialMedia.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helpers.Ok(c, "Social media updated", data)
}

func SocialMediaDelete(c *gin.Context) {
	db := database.GetDB()
	socialMedia := models.SocialMedia{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	id := c.Param("socialMediaId")

	err := db.Where("id = ? AND user_id = ?", id, uint(userData["id"].(float64))).First(&socialMedia).Error

	if err != nil {
		helpers.Unauthorized(c, "Social media not found or you don't have access to delete this social media")
		return
	}

	err = db.Delete(&socialMedia).Error

	if err != nil {
		helpers.BadRequest(c, "Failed to delete social media")
		return
	}

	helpers.JustMessage(c, "Social media deleted")
}
