package repository

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/jinzhu/gorm"
	"testing"
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
	db.DropTable(entities.Review{})
	db.AutoMigrate(entities.Review{})
	return db, nil
}

func TestCreateReview(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateReviewRepository(db)
	testReview := &entities.Review{}
	tr, err := repo.Create(context.Background(), testReview)
	if err != nil {
		t.Error(err)
	}
	if tr == nil {
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateReviewRepository(db)
	testReview := &entities.Review{}
	tr, err := repo.Create(context.Background(), testReview)
	if err != nil {
		t.Error(err)
	}
	if tr == nil {
		t.Fail()
	}
	ur, err := repo.Update(context.Background(), testReview)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if ur == nil {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateReviewRepository(db)
	testReview := &entities.Review{}
	tr, err := repo.Create(context.Background(), testReview)
	if err != nil {
		t.Error(err)
	}
	if tr == nil {
		t.Fail()
	}
	ur, err := repo.Delete(context.Background(), testReview)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if ur == nil {
		t.Fail()
	}
}

func TestGetByID(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateReviewRepository(db)
	tr := &entities.Review{}
	tr, err = repo.Create(context.Background(), tr)
	if err != nil {
		t.Error(err)
	}
	if tr == nil {
		t.Fail()
	}
	rr, err := repo.GetByID(context.Background(), tr.ID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if rr == nil {
		t.Fail()
	}
	if rr.ID != tr.ID {
		t.Fail()
	}
}

func TestGetNullRecordByID(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateReviewRepository(db)
	rr, err := repo.GetByID(context.Background(), 1)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if rr != nil {
		t.Fail()
	}
}

func TestGetByUserID(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateReviewRepository(db)
	tr := &entities.Review{
		CustomerID: 1,
	}
	tr, err = repo.Create(context.Background(), tr)
	if err != nil {
		t.Error(err)
	}
	if tr == nil {
		t.Fail()
	}
	rr, err := repo.GetByCustomerID(context.Background(), tr.CustomerID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if rr == nil {
		t.Fail()
	}
	if len(rr) < 1 {
		t.Fail()
	}
}

func TestGetByPartnerID(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateReviewRepository(db)
	tr := &entities.Review{
		PartnerID: 1,
	}
	tr, err = repo.Create(context.Background(), tr)
	if err != nil {
		t.Error(err)
	}
	if tr == nil {
		t.Fail()
	}
	rr, err := repo.GetByPartnerID(context.Background(), tr.PartnerID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if rr == nil {
		t.Fail()
	}
	if len(rr) < 1 {
		t.Fail()
	}
}
