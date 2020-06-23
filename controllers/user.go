package controllers

import (
	"../structs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func (idb *InDB) GetUser(c *gin.Context) {
	var (
		user    structs.User
		payload gin.H
		status  int
	)

	id := c.Param("id")
	result := idb.DB.
		Where("id = ?", id).
		First(&user)

	if result.Error != nil {
		if result.RecordNotFound() {
			status = http.StatusOK
			payload = gin.H{
				"message": "Record Not Found",
				"data":    nil,
			}
		} else {
			status = http.StatusBadGateway
			payload = gin.H{
				"messages": "Request Had Some Error",
				"error":    result.Error,
			}
		}
	} else {
		status = http.StatusOK
		payload = gin.H{
			"message": "Success",
			"data":    user,
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) GetUsers(c *gin.Context) {
	var (
		users   []structs.User
		payload gin.H
		status  int
	)

	result := idb.DB.
		Set("gorm:auto_preload", true).
		Find(&users)

	if result.Error != nil {
		if result.RecordNotFound() {
			status = http.StatusOK
			payload = gin.H{
				"message": "Record Not Found",
				"data":    nil,
			}
		} else {
			status = http.StatusInternalServerError
			payload = gin.H{
				"message": "Request Had Some Error",
				"error":   result.Error,
			}
		}
	} else {
		status = http.StatusOK
		payload = gin.H{
			"message": "Success",
			"data":    users,
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) CreateUser(c *gin.Context) {
	var (
		user    structs.User
		payload gin.H
		status  int
	)

	if err := c.BindJSON(&user); err != nil {
		status = http.StatusInternalServerError
		payload = gin.H{
			"message":   "Request Had Some Error",
			"reference": "JSON Binding Failed",
			"error":     err,
		}
	} else if user.Username != "" {
		r, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(r)

		result := idb.DB.
			Set("gorm:auto_preload", true).
			Create(&user)

		if result.Error != nil {
			if result.RecordNotFound() {
				status = http.StatusOK
				payload = gin.H{
					"message": "Record Not Found",
					"data":    nil,
				}
			} else {
				status = http.StatusInternalServerError
				payload = gin.H{
					"message": "Request Had Some Error",
					"error":   result.Error,
				}
			}
		} else {
			status = http.StatusOK
			payload = gin.H{
				"message": "Success",
				"data":    user,
			}
		}
	} else {
		status = http.StatusBadRequest
		payload = gin.H{
			"message": "Request Rejected",
			"error":   0,
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) UpdateUser(c *gin.Context) {
	var (
		user    structs.User
		payload gin.H
		status  int
	)

	id := c.Param("id")

	if err := idb.DB.
		Set("gorm:auto_preload", true).
		Where("id = ?", id).First(&user).Error; err != nil {
		status = http.StatusBadRequest
		payload = gin.H{
			"message": "Request Rejected",
			"error":   0,
		}
	} else {
		if err := c.BindJSON(&user); err != nil {
			status = http.StatusInternalServerError
			payload = gin.H{
				"message":   "Request Had Some Error",
				"reference": "JSON Binding Failed",
				"error":     err,
			}
		} else if user.Username != "" {
			r, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(r)

			result := idb.DB.
				Set("gorm:auto_preload", true).
				Save(&user)

			if result.Error != nil {
				status = http.StatusInternalServerError
				payload = gin.H{
					"message": "Request Had Some Error",
					"error":   result.Error,
				}
			} else {
				status = http.StatusOK
				payload = gin.H{
					"message": "Success",
					"data":    user,
				}
			}
		} else {
			status = http.StatusBadRequest
			payload = gin.H{
				"message": "Request Rejected",
				"error":   0,
			}
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) DeleteUser(c *gin.Context) {
	var (
		user    structs.User
		payload gin.H
		status  int
	)

	id := c.Param("id")
	if idb.DB.
		Set("gorm:auto_preload", true).
		Where("id = ?", id).First(&user).RecordNotFound() {
		status = http.StatusOK
		payload = gin.H{
			"message": "Record Not Found",
			"data":    nil,
		}
	} else {
		result := idb.DB.
			Set("gorm:auto_preload", true).
			Delete(&user)

		if result.Error != nil {
			status = http.StatusInternalServerError
			payload = gin.H{
				"message": "Request Had Some Error",
				"error":   result.Error,
			}
		} else {
			status = http.StatusOK
			payload = gin.H{
				"message": "Success",
				"data":    user,
			}
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) AuthorizeUser(c *gin.Context) {
	var (
		user    structs.User
		payload gin.H
		result  *gorm.DB
		status  int
	)

	credential := c.Param("credential")
	password := c.Param("password")

	authType := 1
	if strings.Contains(credential, "@") {
		authType = 2
	}

	if authType == 1 {
		result = idb.DB.
			Set("gorm:auto_preload", true).
			Where("username = ?", credential).
			First(&user)
	} else {
		result = idb.DB.
			Set("gorm:auto_preload", true).
			Where("email = ?", credential).
			First(&user)
	}

	if result.Error != nil {
		if result.RecordNotFound() {
			status = http.StatusNoContent
			payload = gin.H{
				"message": "Record Not Found",
				"data":    nil,
			}
		} else {
			status = http.StatusInternalServerError
			payload = gin.H{
				"message": "Request Had Some Error",
				"error":   result.Error,
			}
		}
	} else {
		passwordCompareResult := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if passwordCompareResult != nil {
			status = 401
			payload = gin.H{
				"message": "Not Authorized",
				"error":   passwordCompareResult,
			}
		} else {
			status = http.StatusOK
			payload = gin.H{
				"message": "Success",
				"data":    user,
			}
		}
	}

	c.JSON(status, payload)
}
