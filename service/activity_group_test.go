package service_test

import (
	"fmt"
	"testing"

	"github.com/letenk/todo-list/config"
	"github.com/letenk/todo-list/helper"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/repository"
	"github.com/letenk/todo-list/service"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var conn *gorm.DB

func dropTable() {
	// Drop table after test
	conn.Raw("delete from activity_groups")
}

func TestMain(m *testing.M) {
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

	// Insert
	newActivityGroup, err := service.Insert(data)
	helper.ErrLogPanic(err)
	fmt.Println(newActivityGroup)
	// Test pass
	assert.Equal(t, data.Title, newActivityGroup.Title)
	assert.Equal(t, data.Email, newActivityGroup.Email)
	assert.NotEmpty(t, newActivityGroup.ID)
	assert.NotEmpty(t, newActivityGroup.CreatedAt)
	assert.NotEmpty(t, newActivityGroup.UpdatedAt)
	assert.Empty(t, newActivityGroup.DeletedAt)

	return newActivityGroup
}

func TestCreateActivityGroup(t *testing.T) {
	defer dropTable()
	t.Parallel()
	createRandomActivityGroupService(t)
}
