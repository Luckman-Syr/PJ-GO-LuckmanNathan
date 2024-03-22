package controllers

import (
	"project-akhir/database"
	"project-akhir/helpers"
	"project-akhir/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type resComment struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
	PhotoID uint   `json:"photo_id"`
	UserID  uint   `json:"user_id"`
}

func CommentCreate(c *gin.Context) {
	comment := models.Comment{}
	photo := models.Photo{}
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	userData := c.MustGet("userData").(jwt.MapClaims)

	if contentType != appJSON {
		c.ShouldBind(&comment)
	} else {
		c.ShouldBindJSON(&comment)
	}

	errFind := db.First(&photo, comment.PhotoID).Error
	if errFind != nil {
		helpers.NotFound(c, "Photo not found")
		return
	}

	comment.UserID = uint(userData["id"].(float64))
	err := db.Create(&comment).Error

	if err != nil {
		helpers.BadRequest(c, "Failed to create comment")
		return
	}

	type CommentResponse struct {
		resComment
		CreatedAt string `json:"created_at"`
	}

	data := CommentResponse{
		resComment: resComment{
			ID:      comment.ID,
			Message: comment.Message,
			PhotoID: comment.PhotoID,
			UserID:  comment.UserID,
		},
		CreatedAt: comment.CreatedAt.String(),
	}

	helpers.Created(c, "Comment created", data)
}

func CommentList(c *gin.Context) {
	comments := []models.Comment{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	db := database.GetDB()
	_ = userData

	err := db.Preload("User").Preload("Photo").Find(&comments).Error

	if err != nil {
		helpers.NotFound(c, "Comments not found")
		return
	}

	type UserResponse struct {
		ID       uint   `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	type PhotoResponse struct {
		ID       uint   `json:"id"`
		Title    string `json:"title"`
		Caption  string `json:"caption"`
		PhotoUrl string `json:"photo_url"`
		UserID   uint   `json:"user_id"`
	}

	type CommentResponse struct {
		resComment
		UpdatedAt string `json:"update_at"`
		CreatedAt string `json:"created_at"`
		User      UserResponse
		Photo     PhotoResponse
	}

	var data []CommentResponse

	for _, comment := range comments {
		data = append(data, CommentResponse{
			resComment: resComment{
				ID:      comment.ID,
				Message: comment.Message,
				PhotoID: comment.PhotoID,
				UserID:  comment.UserID,
			},
			UpdatedAt: comment.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
			User: UserResponse{
				ID:       comment.User.ID,
				Email:    comment.User.Email,
				Username: comment.User.Username,
			},
			Photo: PhotoResponse{
				ID:       comment.Photo.ID,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoUrl: comment.Photo.Photo_url,
				UserID:   comment.Photo.UserID,
			},
		})
	}

	helpers.Ok(c, "Comments list", data)
}

func CommentUpdate(c *gin.Context) {
	comment := models.Comment{}
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	id := c.Param("commentId")
	_ = userData

	errCheck := db.Where("id = ? AND user_id = ?", id, userData["id"]).First(&comment).Error

	if errCheck != nil {
		helpers.Unauthorized(c, "You are not authorized to update this comment")
		return
	}

	if contentType != appJSON {
		c.ShouldBind(&comment)
	} else {
		c.ShouldBindJSON(&comment)
	}

	err := db.Model(&comment).Where("id = ?", c.Param("commentId")).Updates(&comment).Error

	if err != nil {
		helpers.BadRequest(c, "Failed to update comment")
		return
	}

	db.Where("id = ?", id).First(&comment)

	type CommentResponse struct {
		resComment
		UpdatedAt string `json:"update_at"`
	}

	data := CommentResponse{
		resComment: resComment{
			ID:      comment.ID,
			Message: comment.Message,
			PhotoID: comment.PhotoID,
			UserID:  comment.UserID,
		},
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helpers.Ok(c, "Comment updated", data)
}

func CommentDelete(c *gin.Context) {
	comment := models.Comment{}
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	_ = userData

	err := db.Where("id = ? AND user_id = ?", c.Param("commentId"), userData["id"]).Find(&comment).Error

	if err != nil {
		helpers.Unauthorized(c, "You are not authorized to delete this comment")
		return
	}

	errDelete := db.Delete(&comment, c.Param("commentId")).Error

	if errDelete != nil {
		helpers.BadRequest(c, "Failed to delete comment")
		return
	}

	helpers.JustMessage(c, "Comment deleted")
}
