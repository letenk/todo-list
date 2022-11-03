package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/repository"
)

type ActivityService interface {
	Create(req web.ActivityRequest) (domain.Activity, error)
	GetAll() ([]domain.Activity, error)
	GetOne(id uint64) (domain.Activity, error)
	Update(id uint64, req web.ActivityUpdateRequest) (domain.Activity, error)
	Delete(id uint64) (bool, error)
}

type activityService struct {
	repository repository.ActivityRepository
}

func NewServiceActivity(repository repository.ActivityRepository) *activityService {
	return &activityService{repository}
}

func (s *activityService) GetAll() ([]domain.Activity, error) {
	// Find all
	Activitys, err := s.repository.FindAll()

	if err != nil {
		return Activitys, err
	}

	return Activitys, nil
}

func (s *activityService) GetOne(id uint64) (domain.Activity, error) {
	// Find one
	Activity, err := s.repository.FindOne(id)

	if err != nil {
		return Activity, err
	}

	return Activity, nil
}

func (s *activityService) Create(req web.ActivityRequest) (domain.Activity, error) {
	Activity := domain.Activity{
		Title: req.Title,
		Email: req.Email,
	}

	// Save
	newActivity, err := s.repository.Save(Activity)
	if err != nil {
		return newActivity, err
	}

	return newActivity, nil
}

func (s *activityService) Update(id uint64, req web.ActivityUpdateRequest) (domain.Activity, error) {
	// Find one
	Activity, err := s.repository.FindOne(id)
	// If activity group not found
	if Activity.ID == 0 {
		message := fmt.Sprintf("Activity with ID %d Not Found", id)
		return Activity, errors.New(message)
	}

	if err != nil {
		return Activity, err
	}

	// Change field title to req update title
	Activity.Title = req.Title
	// Change time field updatUpdatedAted
	Activity.UpdatedAt = time.Now()

	// Update
	updatedActivity, err := s.repository.Update(Activity)
	if err != nil {
		return Activity, err
	}

	return updatedActivity, nil
}

func (s *activityService) Delete(id uint64) (bool, error) {
	// Find one
	Activity, err := s.repository.FindOne(id)
	// If activity group not found
	if Activity.ID == 0 {
		message := fmt.Sprintf("Activity with ID %d Not Found", id)
		return false, errors.New(message)
	}

	if err != nil {
		return false, err
	}

	ok, err := s.repository.Delete(Activity)
	if err != nil {
		return false, err
	}

	return ok, nil
}
