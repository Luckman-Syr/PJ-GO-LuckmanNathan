package controllers

import (
	"context"
	"mime/multipart"
	"net/http"
	"project-akhir/config"
	"project-akhir/database"
	"project-akhir/helpers"
	"project-akhir/models"
	"project-akhir/utils"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type resPhoto struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	Photo_url string `json:"photo_url"`
	UserID    uint   `json:"user_id"`
}

func PhotoUpload(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	photo := models.Photo{}

	file, ok := c.Get("file")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}

	filename := "photo_"
	imageUrl, errUplod := utils.UploadToCloudinary(file.(multipart.File),
		filename+c.PostForm("title"))
	if errUplod != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errUplod.Error()})
		return
	}

	photo.Photo_url = imageUrl
	photo.UserID = uint(userData["id"].(float64))
	photo.Title = c.PostForm("title")
	photo.Caption = c.PostForm("caption")

	err := db.Debug().Create(&photo).Error

	if err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	type photoResponse struct {
		resPhoto
		CreatedAt string `json:"created_at"`
	}

	data := photoResponse{
		resPhoto: resPhoto{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			Photo_url: imageUrl,
			UserID:    photo.UserID,
		},
		CreatedAt: photo.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	helpers.Created(c, "Photo created successfully", data)
}

func PhotoGetAll(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	photos := []models.Photo{}

	err := db.Debug().Where("user_id = ?", userData["id"]).Preload("User").Find(&photos).Error

	if err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	type PhotoResponse struct {
		resPhoto
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		User      struct {
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"User"`
	}

	photoResponses := []PhotoResponse{}

	for _, photo := range photos {
		photoResponse := PhotoResponse{
			resPhoto: resPhoto{
				ID:        photo.ID,
				Title:     photo.Title,
				Caption:   photo.Caption,
				Photo_url: photo.Photo_url,
				UserID:    photo.UserID,
			},
			CreatedAt: photo.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: photo.UpdatedAt.Format("2006-01-02 15:04:05"),
			User: struct {
				Email    string `json:"email"`
				Username string `json:"username"`
			}{
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		}
		photoResponses = append(photoResponses, photoResponse)
	}

	helpers.Ok(c, "Photos retrieved successfully", photoResponses)
}

func PhotoUpdate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	id := c.Param("photoId")
	contentType := helpers.GetContentType(c)
	photo := models.Photo{}

	errAuth := db.Debug().Where("id = ? AND user_id = ?", id, userData["id"]).First(&photo).Error

	if errAuth != nil {
		helpers.Unauthorized(c, "You are not authorized to update this photo")
		return
	}

	file, ok := c.Get("file")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}

	filename := "photo_"

	ctx := context.Background()
	cld, errConnect := config.SetupCloudinary()
	if errConnect != nil {
		helpers.BadRequest(c, errConnect.Error())
		return
	}
	cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: filename + photo.Title})

	imageUrl, errUplod := utils.UploadToCloudinary(file.(multipart.File),
		filename+c.PostForm("title"))
	if errUplod != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errUplod.Error()})
		return
	}

	if contentType != appJSON {
		c.ShouldBind(&photo)
	} else {
		c.ShouldBindJSON(&photo)
	}
	photo.Title = c.PostForm("title")
	photo.Caption = c.PostForm("caption")
	photo.Photo_url = imageUrl

	err := db.Debug().Where("id = ? AND user_id = ?", id, userData["id"]).Updates(&photo).Error

	if err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	type photoResponse struct {
		resPhoto
		UpdateAt string `json:"updated_at"`
	}

	data := photoResponse{
		resPhoto: resPhoto{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			Photo_url: photo.Photo_url,
			UserID:    photo.UserID,
		},
		UpdateAt: photo.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helpers.Ok(c, "Photo updated successfully", data)

}

func PhotoDelete(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	id := c.Param("photoId")
	photo := models.Photo{}

	errAuth := db.Debug().Where("id = ? AND user_id = ?", id, userData["id"]).First(&photo).Error

	if errAuth != nil {
		helpers.Unauthorized(c, "You are not authorized to delete this photo")
		return
	}

	filename := "photo_"

	ctx := context.Background()
	cld, errConnect := config.SetupCloudinary()
	if errConnect != nil {
		helpers.BadRequest(c, errConnect.Error())
		return
	}
	cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: filename + photo.Title})

	err := db.Debug().Where("id = ? AND user_id = ?", id, userData["id"]).Delete(&photo).Error

	if err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	helpers.JustMessage(c, "Your photo has been successfully deleted")
}
