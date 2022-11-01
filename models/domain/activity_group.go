package domain

import "time"

type ActivityGroup struct {
	ID        int        `gorm:"primary_key"`
	Email     string     `gorm:"type:varchar(191);not null;unique"`
	Title     string     `gorm:"type:varchar(191);not null"`
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoCreateTime"`
	DeletedAt *time.Time `gorm:"default:null"`
}
