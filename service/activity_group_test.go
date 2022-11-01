package service_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/letenk/todo-list/config"
	"github.com/letenk/todo-list/helper"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/repository"
	"github.com/letenk/todo-list/service"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var mutex sync.Mutex
var conn *gorm.DB

func dropTable() {
	// Drop table after test
	conn.Raw("delete from activity_groups")
}

func TestMain(m *testing.M) {
	// Set env
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "root")
	os.Setenv("MYSQL_HOST", "127.0.0.1")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DBNAME", "todo4")

	db := config.SetupDB()
	conn = db
	m.Run()
}

func createRandomActivityGroupService(t *testing.T) domain.ActivityGroup {
	repository := repository.NewRepositoryActivityGroup(conn)
	service := service.NewServiceActivityGroup(repository)

	data := web.ActivityGroupRequest{
		Title: jabufaker.RandomString(20),
		Email: jabufaker.RandomEmail(),
	}

	// Create
	newActivityGroup, err := service.Create(data)
	helper.ErrLogPanic(err)

	// Test pass
	assert.Equal(t, data.Title, newActivityGroup.Title)
	assert.Equal(t, data.Email, newActivityGroup.Email)
	assert.NotEmpty(t, newActivityGroup.ID)
	assert.NotEmpty(t, newActivityGroup.CreatedAt)
	assert.NotEmpty(t, newActivityGroup.UpdatedAt)
	assert.Nil(t, newActivityGroup.DeletedAt)

	return newActivityGroup
}

func TestCreateActivityGroupServices(t *testing.T) {
	defer dropTable()
	t.Parallel()
	createRandomActivityGroupService(t)
}

func TestGetAllServices(t *testing.T) {
	defer dropTable()
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomActivityGroupService(t)
			mutex.Unlock()
		}()
	}

	t.Parallel()

	repository := repository.NewRepositoryActivityGroup(conn)
	service := service.NewServiceActivityGroup(repository)

	// Get activity groups
	activityGroups, err := service.GetAll()
	helper.ErrLogPanic(err)

	for _, data := range activityGroups {
		require.NotEmpty(t, data.ID)
		require.NotEmpty(t, data.Title)
		require.NotEmpty(t, data.Email)
		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)
		require.Nil(t, data.DeletedAt)
	}

}

func TestGetOneService(t *testing.T) {
	defer dropTable()
	// Create random data
	newActivityGroup := createRandomActivityGroupService(t)

	t.Parallel()
	repository := repository.NewRepositoryActivityGroup(conn)
	service := service.NewServiceActivityGroup(repository)

	// Find all
	activityGroup, err := service.GetOne(newActivityGroup.ID)
	helper.ErrLogPanic(err)

	require.Equal(t, newActivityGroup.ID, activityGroup.ID)
	require.Equal(t, newActivityGroup.Title, activityGroup.Title)
	require.Equal(t, newActivityGroup.Email, activityGroup.Email)
	require.NotEmpty(t, activityGroup.CreatedAt)
	require.NotEmpty(t, activityGroup.UpdatedAt)
	require.Nil(t, activityGroup.DeletedAt)
}

func TestUpdateActivityGroupService(t *testing.T) {
	defer dropTable()
	// Create random data
	newActivityGroup := createRandomActivityGroupService(t)

	t.Parallel()
	repository := repository.NewRepositoryActivityGroup(conn)
	service := service.NewServiceActivityGroup(repository)

	dataUpdated := web.ActivityGroupUpdateRequest{
		Title: jabufaker.RandomString(20),
	}

	t.Run("Update success", func(t *testing.T) {

		updatedActivityGroup, err := service.Update(newActivityGroup.ID, dataUpdated)
		helper.ErrLogPanic(err)

		require.Equal(t, newActivityGroup.ID, updatedActivityGroup.ID)
		require.Equal(t, newActivityGroup.Email, updatedActivityGroup.Email)

		require.NotEqual(t, newActivityGroup.Title, updatedActivityGroup.Title)
		require.NotEqual(t, newActivityGroup.UpdatedAt, updatedActivityGroup.UpdatedAt)

		require.NotEmpty(t, updatedActivityGroup.CreatedAt)
		require.Nil(t, updatedActivityGroup.DeletedAt)

	})

	t.Run("Update failed activity group not found", func(t *testing.T) {
		_, err := service.Update(7329323, dataUpdated)
		require.Error(t, err)

		message := fmt.Sprintf("Activity with ID %d Not Found", 7329323)
		require.Equal(t, message, err.Error())

	})
}

func TestDeleteActivityGroupService(t *testing.T) {
	defer dropTable()
	// Create random data
	newActivityGroup := createRandomActivityGroupService(t)

	t.Parallel()
	repository := repository.NewRepositoryActivityGroup(conn)
	service := service.NewServiceActivityGroup(repository)

	t.Run("Delete success", func(t *testing.T) {

		ok, err := service.Delete(newActivityGroup.ID)
		helper.ErrLogPanic(err)

		require.True(t, ok)

	})

	t.Run("Delete failed activity group not found", func(t *testing.T) {
		ok, err := service.Delete(7329323)
		require.Error(t, err)
		require.False(t, ok)

		message := fmt.Sprintf("Activity with ID %d Not Found", 7329323)
		require.Equal(t, message, err.Error())

	})
}
