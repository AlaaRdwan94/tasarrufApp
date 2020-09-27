package usecase

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/branch"
	branchrepo "github.com/ahmedaabouzied/tasarruf/branch/repository"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/review"
	reviewrepo "github.com/ahmedaabouzied/tasarruf/review/repository"
	"github.com/ahmedaabouzied/tasarruf/user"
	userrepo "github.com/ahmedaabouzied/tasarruf/user/repository"
	"github.com/jinzhu/gorm"
	"testing"
)

var reviewUsecase review.Usecase
var userRepo user.Repository
var reviewRepo review.Repository
var branchRepo branch.Repository

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
	db.DropTable(entities.Review{})
	db.AutoMigrate(entities.Branch{})
	db.AutoMigrate(entities.User{})
	db.AutoMigrate(entities.Review{})
	return db, nil
}

func CreateRepos(db *gorm.DB) {
	userRepo = userrepo.CreateUserRepository(db)
	branchRepo = branchrepo.CreateBranchRepository(db)
	reviewRepo = reviewrepo.CreateReviewRepository(db)
}

func setupTest() {
	db, _ := connectToDB()
	CreateRepos(db)
	reviewUsecase = CreateReviewUsecase(reviewRepo, userRepo, branchRepo)
}

func TestCreate(t *testing.T) {
	t.Run("TestCreateReviewByPartner", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			Email:       "testuser@test.com",
			AccountType: "partner",
		}
		testUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
		}
		testPartner := &entities.User{
			Email:       "testpartner@test.com",
			AccountType: "partner",
		}
		testPartner, err = userRepo.Create(context.Background(), testPartner)
		if err != nil {
			t.Error(err)
		}
		newReview := &entities.Review{
			CustomerID: testUser.ID,
			PartnerID:  testPartner.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, testUser.ID)
		_, err = reviewUsecase.Create(ctx, newReview)
		if err == nil {
			t.Fail()
		}
	})
	t.Run("TestCreateReviewByCustomer", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			Email:       "testuser@test.com",
			AccountType: "user",
		}
		testUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
		}
		testPartner := &entities.User{
			Email:       "testpartner@test.com",
			AccountType: "partner",
		}
		testPartner, err = userRepo.Create(context.Background(), testPartner)
		if err != nil {
			t.Error(err)
		}
		newReview := &entities.Review{
			CustomerID: testUser.ID,
			PartnerID:  testPartner.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, testUser.ID)
		cr, err := reviewUsecase.Create(ctx, newReview)
		if err != nil {
			t.Fail()
		}
		if cr == nil {
			t.Fail()
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("TestUpdateReviewByNotOwner", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			Email:       "testuser@test.com",
			AccountType: "user",
		}
		testUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
			return
		}
		testPartner := &entities.User{
			Email:       "testpartner@test.com",
			AccountType: "partner",
		}
		testPartner, err = userRepo.Create(context.Background(), testPartner)
		if err != nil {
			t.Error(err)
			return
		}
		newReview := &entities.Review{
			CustomerID: testUser.ID,
			PartnerID:  testPartner.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, testUser.ID)
		review, err := reviewUsecase.Create(ctx, newReview)
		if err != nil {
			t.Error(err)
			return
		}
		if review == nil {
			t.Fail()
			return
		}
		review.Content = "Updated"
		ctx = context.WithValue(context.Background(), entities.UserIDKey, testPartner.ID)
		_, err = reviewUsecase.Update(ctx, review)
		if err == nil {
			t.Fail()
			return
		}
	})
	t.Run("TestUpdateReviewByOwner", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			Email:       "testuser@test.com",
			AccountType: "user",
		}
		testUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
			return
		}
		testPartner := &entities.User{
			Email:       "testpartner@test.com",
			AccountType: "partner",
		}
		testPartner, err = userRepo.Create(context.Background(), testPartner)
		if err != nil {
			t.Error(err)
			return
		}
		newReview := &entities.Review{
			CustomerID: testUser.ID,
			PartnerID:  testPartner.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, testUser.ID)
		cr, err := reviewUsecase.Create(ctx, newReview)
		if err != nil {
			t.Fail()
			return
		}
		if cr == nil {
			t.Fail()
			return
		}
		cr.Content = "updated"
		_, err = reviewUsecase.Update(ctx, cr)
		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("TestDeleteReviewByNotOwner", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			Email:       "testuser@test.com",
			AccountType: "user",
		}
		testUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
			return
		}
		testPartner := &entities.User{
			Email:       "testpartner@test.com",
			AccountType: "partner",
		}
		testPartner, err = userRepo.Create(context.Background(), testPartner)
		if err != nil {
			t.Error(err)
			return
		}
		newReview := &entities.Review{
			CustomerID: testUser.ID,
			PartnerID:  testPartner.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, testUser.ID)
		review, err := reviewUsecase.Create(ctx, newReview)
		if err != nil {
			t.Error(err)
			return
		}
		if review == nil {
			t.Fail()
			return
		}
		ctx = context.WithValue(context.Background(), entities.UserIDKey, testPartner.ID)
		_, err = reviewUsecase.Delete(ctx, review)
		if err == nil {
			t.Fail()
			return
		}
	})
	t.Run("TestDeleteReviewByOwner", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			Email:       "testuser@test.com",
			AccountType: "user",
		}
		testUser, err := userRepo.Create(context.Background(), testUser)
		if err != nil {
			t.Error(err)
			return
		}
		testPartner := &entities.User{
			Email:       "testpartner@test.com",
			AccountType: "partner",
		}
		testPartner, err = userRepo.Create(context.Background(), testPartner)
		if err != nil {
			t.Error(err)
			return
		}
		newReview := &entities.Review{
			CustomerID: testUser.ID,
			PartnerID:  testPartner.ID,
		}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, testUser.ID)
		cr, err := reviewUsecase.Create(ctx, newReview)
		if err != nil {
			t.Fail()
			return
		}
		if cr == nil {
			t.Fail()
			return
		}
		_, err = reviewUsecase.Delete(ctx, cr)
		if err != nil {
			t.Error(err)
			return
		}
	})
}
func TestGetByID(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		Email:       "testuser@test.com",
		AccountType: "user",
	}
	testUser, err := userRepo.Create(context.Background(), testUser)
	if err != nil {
		t.Error(err)
		return
	}
	testPartner := &entities.User{
		Email:       "testpartner@test.com",
		AccountType: "partner",
	}
	testPartner, err = userRepo.Create(context.Background(), testPartner)
	if err != nil {
		t.Error(err)
		return
	}
	newReview := &entities.Review{
		CustomerID: testUser.ID,
		PartnerID:  testPartner.ID,
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, testUser.ID)
	cr, err := reviewUsecase.Create(ctx, newReview)
	if err != nil {
		t.Fail()
		return
	}
	returnedReview, err := reviewUsecase.GetByID(ctx, cr.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if cr.ID != returnedReview.ID {
		t.Fail()
		return
	}
}

func TestGetByUserID(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		Email:       "testuser@test.com",
		AccountType: "user",
	}
	testUser, err := userRepo.Create(context.Background(), testUser)
	if err != nil {
		t.Error(err)
		return
	}
	testPartner := &entities.User{
		Email:       "testpartner@test.com",
		AccountType: "partner",
	}
	testPartner, err = userRepo.Create(context.Background(), testPartner)
	if err != nil {
		t.Error(err)
		return
	}
	newReview := &entities.Review{
		CustomerID: testUser.ID,
		PartnerID:  testPartner.ID,
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, testUser.ID)
	_, err = reviewUsecase.Create(ctx, newReview)
	if err != nil {
		t.Fail()
		return
	}
	returnedReviews, err := reviewUsecase.GetByCustomerID(ctx, testUser.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(returnedReviews) < 1 {
		t.Fail()
		return
	}
}

func TestGetByPartnerID(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		Email:       "testuser@test.com",
		AccountType: "user",
	}
	testUser, err := userRepo.Create(context.Background(), testUser)
	if err != nil {
		t.Error(err)
		return
	}
	testPartner := &entities.User{
		Email:       "testpartner@test.com",
		AccountType: "partner",
	}
	testPartner, err = userRepo.Create(context.Background(), testPartner)
	if err != nil {
		t.Error(err)
		return
	}
	newReview := &entities.Review{
		CustomerID: testUser.ID,
		PartnerID:  testPartner.ID,
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, testUser.ID)
	_, err = reviewUsecase.Create(ctx, newReview)
	if err != nil {
		t.Fail()
		return
	}
	returnedReviews, err := reviewUsecase.GetByPartnerID(ctx, testPartner.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(returnedReviews) < 1 {
		t.Fail()
		return
	}
}
