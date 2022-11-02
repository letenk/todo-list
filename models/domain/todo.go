package domain

import "time"

type Todo struct {
	ID              uint64        `gorm:"primary_key"`
	ActivityGroupID uint64        `gorm:"not null"`
	Title           string        `gorm:"type:varchar(191);not null"`
	IsActive        bool          `gorm:"default:true;not null"`
	Priority        string        `gorm:"type:enum('very-high', 'high', 'medium', 'low', 'very-low');default:'very-high';not null"`
	CreatedAt       *time.Time    `gorm:"autoCreateTime"`
	UpdatedAt       time.Time     `gorm:"autoCreateTime"`
	DeletedAt       *time.Time    `gorm:"default:null"`
	ActivityGroup   ActivityGroup // Relation one to many to activity groups
}
