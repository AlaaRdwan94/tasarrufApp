package repository

import (
	"context"
	"testing"
	"time"

	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/jinzhu/gorm"
)

func connectToDB() (*gorm.DB, error) {
	conf := entities.DBConfig{
		Port:     5432,
		Host:     "localhost",
		User:     "tasarruf",
		Password: "password",
		DBName:   "tasarruftestdb",
	}
	db, err := entities.ConnectToDB(&conf)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateUserRepository(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	repo := CreateUserRepository(db)
	t.Log(repo)
}

func TestCreateCustomer(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Fail()
	}
	db.DropTable(&entities.User{})
	db.AutoMigrate(&entities.User{})
	defer db.Close()
	repo := CreateUserRepository(db)
	tu := entities.User{
		Email:     "testmail@gmail.com",
		FirstName: "tester",
		LastName:  "tester",
	}
	_, err = repo.CreateCustomer(context.Background(), &tu)
	if err != nil {
		t.Fail()
	}
}

func TestUpdateUser(t *testing.T) {
	// Connect to db
	db, err := connectToDB()
	if err != nil {
		t.Fail()
	}
	db.DropTable(&entities.User{})
	db.AutoMigrate(&entities.User{})
	defer db.Close()
	repo := CreateUserRepository(db)

	// Create a user
	tu := entities.User{
		Email:     "testmail@gmail.com",
		FirstName: "tester",
		LastName:  "tester",
	}
	_, err = repo.CreateCustomer(context.Background(), &tu)
	if err != nil {
		t.Fail()
	}
}

func TestCreateShare(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Fail()
	}
	db.DropTable(&entities.Share{})
	db.AutoMigrate(&entities.Share{})
	defer db.Close()
	repo := CreateUserRepository(db)

	s := entities.Share{
		CustomerID: 1,
	}
	err = repo.CreateShare(context.Background(), &s)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateDuplicateShare(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Fail()
	}
	db.DropTable(&entities.Share{})
	db.AutoMigrate(&entities.Share{})
	defer db.Close()
	repo := CreateUserRepository(db)

	s := entities.Share{
		CustomerID: 1,
	}
	err = repo.CreateShare(context.Background(), &s)
	if err != nil {
		t.Error(err)
	}
	err = repo.CreateShare(context.Background(), &s)
	if err == nil {
		t.Fail()
	}
}

func TestGetSharesByCustomer(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Fail()
	}
	db.DropTable(&entities.Share{})
	db.AutoMigrate(&entities.Share{})
	defer db.Close()
	repo := CreateUserRepository(db)
	s := entities.Share{
		CustomerID: 1,
	}
	err = repo.CreateShare(context.Background(), &s)
	if err != nil {
		t.Error(err)
	}
	c, err := repo.GetSharesByCustomer(context.Background(), s.CustomerID, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
	if err != nil {
		t.Error(err)
	}
	if c != 1 {
		t.Fail()
	}
}

func TestSearchUsers(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Fail()
		return
	}
	defer db.Close()
	repo := CreateUserRepository(db)
	user1 := &entities.User{
		Email:  "testemail1@test.com",
		Mobile: "1234561",
	}
	user2 := &entities.User{
		Email:  "testemail2@test.com",
		Mobile: "1234562",
	}
	user1, err = repo.CreateCustomer(context.Background(), user1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	user2, err = repo.CreateCustomer(context.Background(), user2)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	r1, err := repo.SearchUsers(context.Background(), user1.Email)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	if len(r1) != 1 {
		t.Fail()
		return
	}
	if r1[0].ID != user1.ID {
		t.Fail()
		return
	}
	r2, err := repo.SearchUsers(context.Background(), user2.Mobile)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if len(r2) != 1 {
		t.Fail()
		return
	}
	if r2[0].ID != user2.ID {
		t.Fail()
		return
	}
	r3, err := repo.SearchUsers(context.Background(), "testemail")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	if len(r3) != 2 {
		t.Fail()
		return
	}
	r4, err := repo.SearchUsers(context.Background(), "123456")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	if len(r4) != 2 {
		t.Fail()
		return
	}
}
