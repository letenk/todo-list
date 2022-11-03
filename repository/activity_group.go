package repository

import (
	"github.com/letenk/todo-list/models/domain"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	Save(Activity domain.Activity) (domain.Activity, error)
	FindAll() ([]domain.Activity, error)
	FindOne(id uint64) (domain.Activity, error)
	Update(Activity domain.Activity) (domain.Activity, error)
	Delete(Activity domain.Activity) (bool, error)
}

type activityRepository struct {
	db *gorm.DB
}

func NewRepositoryActivity(db *gorm.DB) *activityRepository {
	return &activityRepository{db}
}

func (r *activityRepository) FindAll() ([]domain.Activity, error) {
	var Activitys []domain.Activity

	err := r.db.Find(&Activitys).Error
	if err != nil {
		return Activitys, nil
	}

	return Activitys, nil
}

func (r *activityRepository) FindOne(id uint64) (domain.Activity, error) {
	var Activity domain.Activity

	err := r.db.Where("id = ?", id).Find(&Activity).Error
	if err != nil {
		return Activity, nil
	}

	return Activity, nil
}

func (r *activityRepository) Save(Activity domain.Activity) (domain.Activity, error) {
	err := r.db.Create(&Activity).Error
	if err != nil {
		return Activity, err
	}

	return Activity, nil
}

func (r *activityRepository) Update(Activity domain.Activity) (domain.Activity, error) {
	err := r.db.Save(&Activity).Error
	if err != nil {
		return Activity, err
	}

	return Activity, nil
}

func (r *activityRepository) Delete(Activity domain.Activity) (bool, error) {
	err := r.db.Delete(&Activity).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
