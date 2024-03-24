package controllers

import (
	"net/http"
	"project-akhir/database"
	"project-akhir/helpers"
	"project-akhir/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

type resUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	user := models.User{}

	if contentType != appJSON {
		c.ShouldBind(&user)
	} else {
		c.ShouldBindJSON(&user)
	}

	err := db.Debug().Create(&user).Error

	if err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	data := resUser{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Age:      user.Age,
	}
	helpers.Created(c, "User created successfully", data)

}

func UserLogin(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	user := models.User{}

	if contentType != appJSON {
		c.ShouldBind(&user)
	} else {
		c.ShouldBindJSON(&user)
	}

	password := user.Password

	db.Debug().Where("email = ?", user.Email).First(&user)

	comparePass := helpers.CheckPassword(user.Password, password)

	if !comparePass {
		helpers.BadRequest(c, "Invalid email or password")
		return
	}

	token := helpers.GenerateToken(user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func UserUpdate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	id := c.Param("userId")
	_, _, _ = db, contentType, userData
	user := models.User{}

	err := db.Debug().Where("id = ?", id).First(&user).Error

	if user.Email != userData["email"] {
		helpers.Unauthorized(c, "You are not authorized to update this user")
		return
	}
	// gasd
	// if user.ID != userData["id"] {
	if err != nil {
		helpers.NotFound(c, "User not found")
		return
	}

	if contentType != appJSON {
		c.ShouldBind(&user)
	} else {
		c.ShouldBindJSON(&user)
	}

	errUpdate := db.Debug().Model(&user).Where("id = ?", id).Updates(&user).Error

	if errUpdate != nil {
		helpers.BadRequest(c, errUpdate.Error())
		return
	}

	type UserResponse struct {
		resUser
		UpdatedAt string `json:"updated_at"`
	}

	data := UserResponse{
		resUser: resUser{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Age:      user.Age,
		},
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	helpers.Ok(c, "User updated successfully", data)
}

func UserDelete(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	_, _ = db, userData
	user := models.User{}

	err := db.Where("id = ? ", userData["id"]).First(&user).Error

	if user.ID != userData["id"] {
		helpers.Unauthorized(c, "You are not authorized to delete this user")
		return
	}

	if err != nil {
		helpers.NotFound(c, "User not found")
		return
	}

	errDelete := db.Debug().Delete(&user).Error

	if errDelete != nil {
		helpers.BadRequest(c, errDelete.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been sucessfully deleted",
	})

	helpers.JustMessage(c, "User deleted successfully")
}
