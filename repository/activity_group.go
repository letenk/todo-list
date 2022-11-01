package repository

import (
	"github.com/letenk/todo-list/models/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Save(activityGroup domain.ActivityGroup) (domain.ActivityGroup, error)
	FindAll() ([]domain.ActivityGroup, error)
	FindOne(id int) (domain.ActivityGroup, error)
	Update(activityGroup domain.ActivityGroup) (domain.ActivityGroup, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepositoryActivityGroup(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]domain.ActivityGroup, error) {
	var activityGroups []domain.ActivityGroup

	err := r.db.Find(&activityGroups).Error
	if err != nil {
		return activityGroups, nil
	}

	return activityGroups, nil
}

func (r *repository) FindOne(id int) (domain.ActivityGroup, error) {
	var activityGroup domain.ActivityGroup

	err := r.db.Where("id = ?", id).Find(&activityGroup).Error
	if err != nil {
		return activityGroup, nil
	}

	return activityGroup, nil
}

func (r *repository) Save(activityGroup domain.ActivityGroup) (domain.ActivityGroup, error) {
	err := r.db.Create(&activityGroup).Error
	if err != nil {
		return activityGroup, err
	}

	return activityGroup, nil
}

func (r *repository) Update(activityGroup domain.ActivityGroup) (domain.ActivityGroup, error) {
	err := r.db.Save(&activityGroup).Error
	if err != nil {
		return activityGroup, err
	}

	return activityGroup, nil
}
