package controllers

import (
	"../structs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (idb *InDB) GetLastLog(c *gin.Context) {
	var (
		log     structs.Log
		payload gin.H
		status  int
	)

	result := idb.DB.Last(&log)

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
			"data":    log,
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) GetLogsWithSpecificDate(c *gin.Context) {
	var (
		logs    []structs.Log
		payload gin.H
		status  int
	)

	date := c.Param("date")
	result := idb.DB.
		Where("DATE(created_at) = ?", string(date)).Order("id DESC").Find(&logs)

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
			"data":    logs,
		}
	}

	c.JSON(status, payload)
}
func (idb *InDB) GetLogDateLists(c *gin.Context) {
	var (
		logs    []structs.Log
		payload gin.H
		status  int
	)

	result := idb.DB.Select("DISTINCT DATE(created_at) as created_at").Find(&logs)

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
			"data":    logs,
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) GetLogs(c *gin.Context) {
	var (
		logs    []structs.Log
		payload gin.H
		status  int
	)

	result := idb.DB.
		Set("gorm:auto_preload", true).
		Find(&logs)

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
			"data":    logs,
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) CreateLog(c *gin.Context) {
	var (
		log     structs.Log
		payload gin.H
		status  int
	)

	if err := c.BindJSON(&log); err != nil {
		status = http.StatusInternalServerError
		payload = gin.H{
			"message":   "Request Had Some Error",
			"reference": "JSON Binding Failed",
			"error":     err,
		}
	} else {
		result := idb.DB.
			Set("gorm:auto_preload", true).
			Create(&log)

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
				"data":    log,
			}
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) UpdateLog(c *gin.Context) {
	var (
		log     structs.Log
		payload gin.H
		status  int
	)

	id := c.Param("id")

	if err := idb.DB.
		Set("gorm:auto_preload", true).
		Where("id = ?", id).First(&log).Error; err != nil {
		status = http.StatusBadRequest
		payload = gin.H{
			"message": "Request Rejected",
			"error":   0,
		}
	} else {
		if err := c.BindJSON(&log); err != nil {
			status = http.StatusInternalServerError
			payload = gin.H{
				"message":   "Request Had Some Error",
				"reference": "JSON Binding Failed",
				"error":     err,
			}
		} else {
			result := idb.DB.
				Set("gorm:auto_preload", true).
				Save(&log)

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
					"data":    log,
				}
			}
		}
	}

	c.JSON(status, payload)
}

func (idb *InDB) DeleteLog(c *gin.Context) {
	var (
		log     structs.Log
		payload gin.H
		status  int
	)

	id := c.Param("id")
	if idb.DB.
		Set("gorm:auto_preload", true).
		Where("id = ?", id).First(&log).RecordNotFound() {
		status = http.StatusOK
		payload = gin.H{
			"message": "Record Not Found",
			"data":    nil,
		}
	} else {
		result := idb.DB.
			Set("gorm:auto_preload", true).
			Delete(&log)

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
				"data":    log,
			}
		}
	}

	c.JSON(status, payload)
}
