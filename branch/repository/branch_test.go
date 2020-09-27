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
	db.DropTable(entities.Branch{})
	db.DropTable(entities.City{})
	db.DropTable(entities.Category{})
	db.AutoMigrate(entities.Branch{})
	db.AutoMigrate(entities.City{})
	db.AutoMigrate(entities.Category{})
	db.Create(&entities.City{})
	return db, nil
}

func TestCreate(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	testBranch := &entities.Branch{}
	testBranch.City.ID = 1
	repo := CreateBranchRepository(db)
	tb, err := repo.Create(context.Background(), testBranch)
	if err != nil {
		t.Error(err)
		return
	}
	rb, err := repo.GetByID(context.Background(), tb.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if rb.ID != tb.ID {
		t.Fail()
		return
	}
}

func TestDelete(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateBranchRepository(db)
	testBranch := &entities.Branch{}
	tb, err := repo.Create(context.Background(), testBranch)
	if err != nil {
		t.Error(err)
	}
	deletedBranch, err := repo.Delete(context.Background(), tb)
	if err != nil {
		t.Error(err)
	}
	returnedBranch, err := repo.GetByID(context.Background(), deletedBranch.ID)
	if err != nil {
		t.Fail()
	}
	if returnedBranch != nil {
		t.Fail()
	}
}

func TestEdit(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateBranchRepository(db)
	testBranch := &entities.Branch{}
	testBranch.City.ID = 1
	tb, err := repo.Create(context.Background(), testBranch)
	if err != nil {
		t.Error(err)
	}
	tb.City = entities.City{}
	tb.City.ID = 1
	editedBranch, err := repo.Edit(context.Background(), tb)
	if err != nil {
		t.Error(err)
	}
	if editedBranch.City.ID != 1 {
		t.Fail()
	}
}

func TestGetByID(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateBranchRepository(db)
	testBranch := &entities.Branch{}
	testBranch.City.ID = 1
	tb, err := repo.Create(context.Background(), testBranch)
	if err != nil {
		t.Error(err)
	}
	rb, err := repo.GetByID(context.Background(), tb.ID)
	if err != nil {
		t.Error(err)
	}
	if rb.ID != tb.ID {
		t.Fail()
	}
}

func TestGetByOwner(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateBranchRepository(db)
	testBranch := &entities.Branch{
		OwnerID: 1,
	}
	testBranch.City.ID = 1
	tb, err := repo.Create(context.Background(), testBranch)
	if err != nil {
		t.Error(err)
	}
	ownerBranches, err := repo.GetByOwner(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}
	if ownerBranches[0].ID != tb.ID {
		t.Fail()
	}
}

func TestGetByLocation(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer db.Close()
	repo := CreateBranchRepository(db)
	testBranch := &entities.Branch{
		City:    entities.City{},
		Country: "country1",
	}
	testBranch.City.ID = 1
	tb, err := repo.Create(context.Background(), testBranch)
	if err != nil {
		t.Error(err)
	}
	branches, err := repo.GetByLocation(context.Background(), testBranch.Country, testBranch.City.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if branches[0].ID != tb.ID {
		t.Fail()
		return
	}
}

func TestCreateCategory(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateBranchRepository(db)
	testCategory, err := repo.CreateCategory(context.Background(), &entities.Category{})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if testCategory == nil {
		t.Fail()
	}
}

func TestEditCategory(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateBranchRepository(db)
	testCategory, err := repo.CreateCategory(context.Background(), &entities.Category{})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if testCategory == nil {
		t.Fail()
	}
	testCategory.TurkishName = "edited"
	_, err = repo.EditCategory(context.Background(), testCategory)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	returnedCategory, err := repo.GetCategoryByID(context.Background(), testCategory.ID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if returnedCategory == nil {
		t.Fail()
	}
	if returnedCategory.TurkishName != testCategory.TurkishName {
		t.Fail()
	}
}
