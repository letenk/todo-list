package test

import (
	"sync"
	"testing"
	"time"

	"github.com/letenk/todo-list/helper"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/repository"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/require"
)

func createRandomActivityGroupRepository(t *testing.T) domain.ActivityGroup {
	activityGroupRepository := repository.NewRepositoryActivityGroup(ConnTest)

	activityGroup := domain.ActivityGroup{
		Title: jabufaker.RandomString(20),
		Email: jabufaker.RandomEmail(),
	}

	// Save to db
	newActivityGroup, err := activityGroupRepository.Save(activityGroup)
	helper.ErrLogPanic(err)

	// Test pass
	require.Equal(t, activityGroup.Title, newActivityGroup.Title)
	require.Equal(t, activityGroup.Email, newActivityGroup.Email)
	require.NotEmpty(t, newActivityGroup.ID)
	require.NotEmpty(t, newActivityGroup.CreatedAt)
	require.NotEmpty(t, newActivityGroup.UpdatedAt)
	require.Empty(t, newActivityGroup.DeletedAt)

	return newActivityGroup
}

func TestCreateActivityGroup(t *testing.T) {
	defer DropTable()
	t.Parallel()
	createRandomActivityGroupRepository(t)
}

func TestFindAllActivityGroup(t *testing.T) {
	var mutex sync.Mutex
	defer DropTable()
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomActivityGroupRepository(t)
			mutex.Unlock()
		}()
	}

	t.Parallel()
	activityGroupRepository := repository.NewRepositoryActivityGroup(ConnTest)

	// Find all
	activityGroup, err := activityGroupRepository.FindAll()
	helper.ErrLogPanic(err)

	for _, data := range activityGroup {
		require.NotEmpty(t, data.ID)
		require.NotEmpty(t, data.Title)
		require.NotEmpty(t, data.Email)
		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)
		require.Empty(t, data.DeletedAt)
	}
}

func TestFindOneActivityGroup(t *testing.T) {
	defer DropTable()
	// Create random data
	newActivityGroup := createRandomActivityGroupRepository(t)

	t.Parallel()
	activityGroupRepository := repository.NewRepositoryActivityGroup(ConnTest)

	// Find all
	activityGroup, err := activityGroupRepository.FindOne(newActivityGroup.ID)
	helper.ErrLogPanic(err)

	require.Equal(t, newActivityGroup.ID, activityGroup.ID)
	require.Equal(t, newActivityGroup.Title, activityGroup.Title)
	require.Equal(t, newActivityGroup.Email, activityGroup.Email)
	require.NotEmpty(t, activityGroup.CreatedAt)
	require.NotEmpty(t, activityGroup.UpdatedAt)
	require.Empty(t, activityGroup.DeletedAt)
}

func TestUpdateActivityGroupRepository(t *testing.T) {
	defer DropTable()
	newActivityGroup := createRandomActivityGroupRepository(t)
	t.Parallel()
	activityGroupRepository := repository.NewRepositoryActivityGroup(ConnTest)

	dataUpdate := domain.ActivityGroup{
		ID:        newActivityGroup.ID,
		Title:     jabufaker.RandomString(20),
		Email:     jabufaker.RandomEmail(),
		CreatedAt: newActivityGroup.CreatedAt,
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	// update
	updateActivityGroup, err := activityGroupRepository.Update(dataUpdate)
	helper.ErrLogPanic(err)

	require.Equal(t, newActivityGroup.ID, updateActivityGroup.ID)
	require.Equal(t, newActivityGroup.CreatedAt, updateActivityGroup.CreatedAt)
	// require.Equal(t, newActivityGroup.DeletedAt, updateActivityGroup.DeletedAt)
	require.NotEqual(t, newActivityGroup.Title, updateActivityGroup.Title)
	require.NotEqual(t, newActivityGroup.Email, updateActivityGroup.Email)
}

func TestDeleteActivityGroupRepository(t *testing.T) {
	DropTable()
	newActivityGroup := createRandomActivityGroupRepository(t)
	t.Parallel()

	activityGroupRepository := repository.NewRepositoryActivityGroup(ConnTest)

	ok, err := activityGroupRepository.Delete(newActivityGroup)
	helper.ErrLogPanic(err)
	require.True(t, ok)

	activityGroup, err := activityGroupRepository.FindOne(newActivityGroup.ID)
	helper.ErrLogPanic(err)
	require.Equal(t, 0, int(activityGroup.ID))
}
