package structs

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string
	Password string `gorm: "type:text;"`
}
