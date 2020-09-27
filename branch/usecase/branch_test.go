package usecase

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/branch"
	branchrepo "github.com/ahmedaabouzied/tasarruf/branch/repository"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/user"
	userrepo "github.com/ahmedaabouzied/tasarruf/user/repository"
	"github.com/jinzhu/gorm"
	"testing"
)

var branchUsecase branch.Usecase
var userRepo user.Repository

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
	db.DropTable(entities.User{})
	db.AutoMigrate(entities.Branch{})
	db.AutoMigrate(entities.User{})
	return db, nil
}

func CreateRepos(db *gorm.DB) (user.Repository, branch.Repository) {
	userRepo = userrepo.CreateUserRepository(db)
	branchRepo := branchrepo.CreateBranchRepository(db)
	return userRepo, branchRepo
}

func setupTest() {
	db, _ := connectToDB()
	userRepo, branchRepo := CreateRepos(db)
	branchUsecase = CreateBranchUsecase(branchRepo, userRepo)
}

func TestCreate(t *testing.T) {
	t.Run("TestCreateBranchWithNonPartnerUser", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			AccountType: "user",
		}
		createdUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
		}
		newBranch := &entities.Branch{
			OwnerID: createdUser.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		_, err = branchUsecase.Create(ctx, newBranch)
		if err == nil {
			t.Fail()
		}

	})
	t.Run("TestCreateBranchWithPartnerUser", func(t *testing.T) {
		setupTest()
		setupTest()
		testUser := &entities.User{
			AccountType: "partner",
		}
		createdUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
		}
		newBranch := &entities.Branch{
			OwnerID: createdUser.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		createdBranch, err := branchUsecase.Create(ctx, newBranch)
		if err != nil {
			t.Error(err)
		}
		returnedBranch, err := branchUsecase.GetByID(context.Background(), createdBranch.ID)
		if err != nil {
			t.Error(err)
		}
		if returnedBranch.ID != createdBranch.ID {
			t.Fail()
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("TestDeleteWithOwnerUser", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			AccountType: "partner",
		}
		createdUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
		}
		newBranch := &entities.Branch{
			OwnerID: createdUser.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		createdBranch, err := branchUsecase.Create(ctx, newBranch)
		if err != nil {
			t.Error(err)
		}
		deletedBranch, err := branchUsecase.Delete(ctx, createdBranch)
		if err != nil {
			t.Error(err)
		}
		returnedBranch, err := branchUsecase.GetByID(context.Background(), deletedBranch.ID)
		if err != nil {
			t.Error(err)
		}
		if returnedBranch != nil {
			t.Fail()
		}
	})
	t.Run("TestDeleteWithNonOwnerUser", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			AccountType: "partner",
		}
		createdUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
		}
		newBranch := &entities.Branch{
			OwnerID: createdUser.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		createdBranch, err := branchUsecase.Create(ctx, newBranch)
		if err != nil {
			t.Error(err)
		}
		ctx = context.WithValue(context.Background(), entities.UserIDKey, uint(2))
		_, err = branchUsecase.Delete(ctx, createdBranch)
		if err == nil {
			t.Fail()
		}
	})
}

func TestEdit(t *testing.T) {
	t.Run("TestEditWithTheOwnerUser", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			AccountType: "partner",
		}
		createdUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
		}
		newBranch := &entities.Branch{
			OwnerID: createdUser.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		createdBranch, err := branchUsecase.Create(ctx, newBranch)
		if err != nil {
			t.Error(err)
		}
		createdBranch.City = "edited"
		ctx = context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		editedBranch, err := branchUsecase.Edit(ctx, createdBranch)
		if err != nil {
			t.Error(err)
		}
		returnedBranch, err := branchUsecase.GetByID(context.Background(), editedBranch.ID)
		if err != nil {
			t.Error(err)
		}
		if returnedBranch.City != createdBranch.City {
			t.Fail()
		}
	})
	t.Run("TestEditWithNonOwnerUser", func(t *testing.T) {
		setupTest()

		testUser := &entities.User{
			AccountType: "partner",
		}
		createdUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
		}
		newBranch := &entities.Branch{
			OwnerID: createdUser.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		createdBranch, err := branchUsecase.Create(ctx, newBranch)
		if err != nil {
			t.Error(err)
		}
		createdBranch.City = "edited"
		ctx = context.WithValue(context.Background(), entities.UserIDKey, uint(2))
		_, err = branchUsecase.Edit(ctx, createdBranch)
		if err == nil {
			t.Fail()
		}
	})
}

func TestGetByID(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		AccountType: "partner",
	}
	createdUser, err := userRepo.Create(context.Background(), testUser)
	if err != nil {
		t.Error(err)
	}
	newBranch := &entities.Branch{
		OwnerID: createdUser.ID,
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
	createdBranch, err := branchUsecase.Create(ctx, newBranch)
	if err != nil {
		t.Error(err)
	}
	returnedBranch, err := branchUsecase.GetByID(context.Background(), createdBranch.ID)
	if err != nil {
		t.Error(err)
	}
	if returnedBranch.ID != createdBranch.ID {
		t.Fail()
	}
}

func TestGetByOwner(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		AccountType: "partner",
	}
	createdUser, err := userRepo.Create(context.Background(), testUser)
	if err != nil {
		t.Error(err)
	}
	newBranch := &entities.Branch{
		OwnerID: createdUser.ID,
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
	createdBranch, err := branchUsecase.Create(ctx, newBranch)
	if err != nil {
		t.Error(err)
	}
	returnedBranch, err := branchUsecase.GetByOwner(context.Background(), newBranch.OwnerID)
	if err != nil {
		t.Error(err)
	}
	if returnedBranch[0].ID != createdBranch.ID {
		t.Fail()
	}
}

func TestGetByLocation(t *testing.T) {
	setupTest()

	testUser := &entities.User{
		AccountType: "partner",
	}
	createdUser, err := userRepo.Create(context.Background(), testUser)
	if err != nil {
		t.Error(err)
	}
	newBranch := &entities.Branch{
		Country: "testCountry",
		City:    "testCity",
		OwnerID: createdUser.ID,
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
	createdBranch, err := branchUsecase.Create(ctx, newBranch)
	if err != nil {
		t.Error(err)
	}
	returnedBranch, err := branchUsecase.GetByLocation(context.Background(), newBranch.Country, newBranch.City)
	if err != nil {
		t.Error(err)
	}
	if returnedBranch[0].ID != createdBranch.ID {
		t.Fail()
	}
}
