package repository_test

import (
	"sync"
	"testing"

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

	// Drop table after test
	db.Migrator().DropTable(&domain.ActivityGroup{})
}

func createRandomActivityGroup(t *testing.T) {
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
}

func TestCreateActivityGroup(t *testing.T) {
	t.Parallel()
	createRandomActivityGroup(t)
}

func TestFindAllActivityGroup(t *testing.T) {
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
