package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/repository"
)

type ActivityGroupService interface {
	Create(req web.ActivityGroupRequest) (domain.ActivityGroup, error)
	GetAll() ([]domain.ActivityGroup, error)
	GetOne(id int) (domain.ActivityGroup, error)
	Update(id int, req web.ActivityGroupUpdateRequest) (domain.ActivityGroup, error)
	Delete(id int) (bool, error)
}

type activityGroupService struct {
	repository repository.ActivityGroupRepository
}

func NewServiceActivityGroup(repository repository.ActivityGroupRepository) *activityGroupService {
	return &activityGroupService{repository}
}

func (s *activityGroupService) GetAll() ([]domain.ActivityGroup, error) {
	// Find all
	activityGroups, err := s.repository.FindAll()

	if err != nil {
		return activityGroups, err
	}

	return activityGroups, nil
}

func (s *activityGroupService) GetOne(id int) (domain.ActivityGroup, error) {
	// Find one
	activityGroup, err := s.repository.FindOne(id)

	if err != nil {
		return activityGroup, err
	}

	return activityGroup, nil
}

func (s *activityGroupService) Create(req web.ActivityGroupRequest) (domain.ActivityGroup, error) {
	activityGroup := domain.ActivityGroup{
		Title: req.Title,
		Email: req.Email,
	}

	// Save
	newActivityGroup, err := s.repository.Save(activityGroup)
	if err != nil {
		return newActivityGroup, err
	}

	return newActivityGroup, nil
}

func (s *activityGroupService) Update(id int, req web.ActivityGroupUpdateRequest) (domain.ActivityGroup, error) {
	// Find one
	activityGroup, err := s.repository.FindOne(id)
	// If activity group not found
	if activityGroup.ID == 0 {
		message := fmt.Sprintf("Activity with ID %d Not Found", id)
		return activityGroup, errors.New(message)
	}

	if err != nil {
		return activityGroup, err
	}

	// Change field title to req update title
	activityGroup.Title = req.Title
	// Change time field updatUpdatedAted
	activityGroup.UpdatedAt = time.Now()

	// Update
	updatedActivityGroup, err := s.repository.Update(activityGroup)
	if err != nil {
		return activityGroup, err
	}

	return updatedActivityGroup, nil
}

func (s *activityGroupService) Delete(id int) (bool, error) {
	// Find one
	activityGroup, err := s.repository.FindOne(id)
	// If activity group not found
	if activityGroup.ID == 0 {
		message := fmt.Sprintf("Activity with ID %d Not Found", id)
		return false, errors.New(message)
	}

	if err != nil {
		return false, err
	}

	ok, err := s.repository.Delete(activityGroup)
	if err != nil {
		return false, err
	}

	return ok, nil
}
