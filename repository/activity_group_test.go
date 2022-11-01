package repository_test

import (
	"sync"
	"testing"
	"time"

	"github.com/letenk/todo-list/config"
	"github.com/letenk/todo-list/helper"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/repository"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var mutex sync.Mutex
var conn *gorm.DB

func TestMain(m *testing.M) {
	db := config.SetupDB()
	conn = db
	m.Run()
}

func dropTable() {
	// Drop table after test
	conn.Raw("delete from activity_groups")
}

func createRandomActivityGroup(t *testing.T) domain.ActivityGroup {
	activityGroupRepository := repository.NewRepositoryActivityGroup(conn)

	activityGroup := domain.ActivityGroup{
		Title: jabufaker.RandomString(20),
		Email: jabufaker.RandomEmail(),
	}

	// Save to db
	newActivityGroup, err := activityGroupRepository.Save(activityGroup)
	helper.ErrLogPanic(err)

	// Test pass
	assert.Equal(t, activityGroup.Title, newActivityGroup.Title)
	assert.Equal(t, activityGroup.Email, newActivityGroup.Email)
	assert.NotEmpty(t, newActivityGroup.ID)
	assert.NotEmpty(t, newActivityGroup.CreatedAt)
	assert.NotEmpty(t, newActivityGroup.UpdatedAt)
	assert.Empty(t, newActivityGroup.DeletedAt)

	return newActivityGroup
}

func TestCreateActivityGroup(t *testing.T) {
	defer dropTable()
	t.Parallel()
	createRandomActivityGroup(t)
}

func TestFindAllActivityGroup(t *testing.T) {
	defer dropTable()
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomActivityGroup(t)
			mutex.Unlock()
		}()
	}

	t.Parallel()
	activityGroupRepository := repository.NewRepositoryActivityGroup(conn)

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
	defer dropTable()
	// Create random data
	newActivityGroup := createRandomActivityGroup(t)

	t.Parallel()
	activityGroupRepository := repository.NewRepositoryActivityGroup(conn)

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

func TestUpdateActivityGroup(t *testing.T) {
	defer dropTable()
	newActivityGroup := createRandomActivityGroup(t)
	t.Parallel()
	activityGroupRepository := repository.NewRepositoryActivityGroup(conn)

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

func TestDeleteActivityGroup(t *testing.T) {
	dropTable()
	newActivityGroup := createRandomActivityGroup(t)
	t.Parallel()

	activityGroupRepository := repository.NewRepositoryActivityGroup(conn)

	ok, err := activityGroupRepository.Delete(newActivityGroup)
	helper.ErrLogPanic(err)
	assert.True(t, ok)

	activityGroup, err := activityGroupRepository.FindOne(newActivityGroup.ID)
	helper.ErrLogPanic(err)
	assert.Equal(t, 0, activityGroup.ID)
}
