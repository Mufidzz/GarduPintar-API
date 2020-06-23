package structs

import "github.com/jinzhu/gorm"

type Log struct {
	gorm.Model
	Phase1A   float32 `gorm:"type:double(11,1)"`
	Phase2A   float32 `gorm:"type:double(11,1)"`
	Phase3A   float32 `gorm:"type:double(11,1)"`
	AllPhaseV int
	TrTemp    int
}
