package repository_test

import (
	"log"
	"testing"

	"github.com/letenk/todo-list/config"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/repository"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()

	// Drop table after test
	db := config.SetupDB()
	db.Migrator().DropTable(&domain.ActivityGroup{})
}

func createRandomActivityGroup(t *testing.T) {
	t.Parallel()
	db := config.SetupDB()
	activityGroupRepository := repository.NewRepositoryActivityGroup(db)

	activityGroup := domain.ActivityGroup{
		Title: jabufaker.RandomString(20),
		Email: jabufaker.RandomEmail(),
	}

	// Save to db
	newActivityGroup, err := activityGroupRepository.Save(activityGroup)
	if err != nil {
		log.Panic(err)
	}

	// Test pass
	assert.Equal(t, activityGroup.Title, newActivityGroup.Title)
	assert.Equal(t, activityGroup.Email, newActivityGroup.Email)
	assert.NotEmpty(t, newActivityGroup.ID)
	assert.NotEmpty(t, newActivityGroup.CreatedAt)
	assert.NotEmpty(t, newActivityGroup.UpdatedAt)
	assert.Empty(t, newActivityGroup.DeletedAt)
}

func TestCreateActivityGroup(t *testing.T) {
	createRandomActivityGroup(t)
}
