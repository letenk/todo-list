package service

import (
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/repository"
)

type ActivityGroupService interface {
	Insert(req web.ActivityGroupRequest) (domain.ActivityGroup, error)
	// GetActivityGroups() ([]domain.ActivityGroup, error)
	// GetOne(id int) (domain.ActivityGroup, error)
	// Update(activityGroup domain.ActivityGroup) (domain.ActivityGroup, error)
	// Delete(activityGroup domain.ActivityGroup) (bool, error)
}

type activityGroupService struct {
	repository repository.ActivityGroupRepository
}

func NewServiceActivityGroup(repository repository.ActivityGroupRepository) *activityGroupService {
	return &activityGroupService{repository}
}

func (s *activityGroupService) Insert(req web.ActivityGroupRequest) (domain.ActivityGroup, error) {
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
