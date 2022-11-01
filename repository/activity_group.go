package repository

import (
	"github.com/letenk/todo-list/models/domain"
	"gorm.io/gorm"
)

type ActivityGroupRepository interface {
	Save(activityGroup domain.ActivityGroup) (domain.ActivityGroup, error)
	FindAll() ([]domain.ActivityGroup, error)
	FindOne(id int) (domain.ActivityGroup, error)
	Update(activityGroup domain.ActivityGroup) (domain.ActivityGroup, error)
	Delete(activityGroup domain.ActivityGroup) (bool, error)
}

type activityGroupRepository struct {
	db *gorm.DB
}

func NewRepositoryActivityGroup(db *gorm.DB) *activityGroupRepository {
	return &activityGroupRepository{db}
}

func (r *activityGroupRepository) FindAll() ([]domain.ActivityGroup, error) {
	var activityGroups []domain.ActivityGroup

	err := r.db.Find(&activityGroups).Error
	if err != nil {
		return activityGroups, nil
	}

	return activityGroups, nil
}

func (r *activityGroupRepository) FindOne(id int) (domain.ActivityGroup, error) {
	var activityGroup domain.ActivityGroup

	err := r.db.Where("id = ?", id).Find(&activityGroup).Error
	if err != nil {
		return activityGroup, nil
	}

	return activityGroup, nil
}

func (r *activityGroupRepository) Save(activityGroup domain.ActivityGroup) (domain.ActivityGroup, error) {
	err := r.db.Create(&activityGroup).Error
	if err != nil {
		return activityGroup, err
	}

	return activityGroup, nil
}

func (r *activityGroupRepository) Update(activityGroup domain.ActivityGroup) (domain.ActivityGroup, error) {
	err := r.db.Save(&activityGroup).Error
	if err != nil {
		return activityGroup, err
	}

	return activityGroup, nil
}

func (r *activityGroupRepository) Delete(activityGroup domain.ActivityGroup) (bool, error) {
	err := r.db.Delete(&activityGroup).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
