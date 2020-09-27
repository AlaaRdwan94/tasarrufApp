package usecase

import (
	"context"
	"testing"

	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/subscription"
	subscriptionrepo "github.com/ahmedaabouzied/tasarruf/subscription/repository"
	"github.com/ahmedaabouzied/tasarruf/user"
	userrepo "github.com/ahmedaabouzied/tasarruf/user/repository"
	"github.com/jinzhu/gorm"
)

var subscriptionUsecase subscription.Usecase
var userRepo user.Repository
var subscriptionRepo subscription.Repository

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
	db.DropTable(entities.Subscription{})
	db.DropTable(entities.Plan{})
	db.AutoMigrate(entities.Branch{})
	db.AutoMigrate(entities.User{})
	db.AutoMigrate(entities.Subscription{})
	db.AutoMigrate(entities.Plan{})
	return db, nil
}

func CreateRepos(db *gorm.DB) (user.Repository, subscription.Repository) {
	userRepo = userrepo.CreateUserRepository(db)
	subscriptionRepo = subscriptionrepo.CreateSubscriptionRepository(db)
	return userRepo, subscriptionRepo
}

func setupTest() {
	db, _ := connectToDB()
	userRepo, subscriptionRepo := CreateRepos(db)
	subscriptionUsecase = CreateSubscriptionUsecase(subscriptionRepo, userRepo)
}

func TestCreatePlan(t *testing.T) {
	t.Run("TestCreatePlanByAdminUser", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			AccountType: "admin",
		}
		createdUser, err := userRepo.CreateCustomer(context.Background(), testUser)
		if err != nil {
			t.Error(err)
			return
		}
		newPlan := &entities.Plan{}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		cp, err := subscriptionUsecase.CreatePlan(ctx, newPlan)
		if err != nil {
			t.Fail()
			return
		}
		if cp == nil {
			t.Fail()
			return
		}
	})

	t.Run("TestCreatePlanByNotAdminUser", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			AccountType: "user",
		}
		createdUser, err := userRepo.CreateCustomer(context.Background(), testUser)
		if err != nil {
			t.Error(err)
			return
		}
		newPlan := &entities.Plan{}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		cp, err := subscriptionUsecase.CreatePlan(ctx, newPlan)
		if err == nil {
			t.Fail()
			return
		}
		if cp != nil {
			t.Fail()
			return
		}
	})
}

func TestDeletePlan(t *testing.T) {
	t.Run("TestDeletePlanByAdminUser", func(t *testing.T) {
		setupTest()
		testUser := &entities.User{
			AccountType: "admin",
		}
		createdUser, err := userRepo.CreateCustomer(context.Background(), testUser)
		if err != nil {
			t.Error(err)
			return
		}
		newPlan := &entities.Plan{}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		cp, err := subscriptionUsecase.CreatePlan(ctx, newPlan)
		if err != nil {
			t.Fail()
			return
		}
		if cp == nil {
			t.Fail()
			return
		}
		deletedPlan, err := subscriptionUsecase.DeletePlan(ctx, cp.ID)
		if err != nil {
			t.Error(err)
			return
		}
		if deletedPlan == nil {
			t.Fail()
			return
		}
	})

	t.Run("TestDeletePlanByNonAdminUser", func(t *testing.T) {
		setupTest()
		testAdminUser := &entities.User{
			Email:       "testadminuser@test.com",
			AccountType: "admin",
			Mobile:      "123456788",
		}
		testUser := &entities.User{
			Email:       "testuser@test.com",
			AccountType: "user",
			Mobile:      "123456789",
		}
		createdUser, err := userRepo.CreateCustomer(context.Background(), testUser)
		if err != nil {
			t.Error(err)
			return
		}
		createdAdminUser, err := userRepo.CreateCustomer(context.Background(), testAdminUser)
		if err != nil {
			t.Error(err)
			return
		}
		newPlan := &entities.Plan{}
		ctx := context.WithValue(context.Background(), entities.UserIDKey, createdAdminUser.ID)
		cp, err := subscriptionUsecase.CreatePlan(ctx, newPlan)
		if err != nil {
			t.Fail()
			return
		}
		if cp == nil {
			t.Fail()
			return
		}
		ctx = context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
		deletedPlan, err := subscriptionUsecase.DeletePlan(ctx, cp.ID)
		if err == nil {
			t.Fail()
			return
		}
		if deletedPlan != nil {
			t.Fail()
			return
		}
	})
}

func TestGetAllPlans(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		AccountType: "admin",
	}
	createdUser, err := userRepo.CreateCustomer(context.Background(), testUser)
	if err != nil {
		t.Error(err)
		return
	}
	newPlan := &entities.Plan{}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
	cp, err := subscriptionUsecase.CreatePlan(ctx, newPlan)
	if err != nil {
		t.Fail()
		return
	}
	if cp == nil {
		t.Fail()
		return
	}
	plans, err := subscriptionUsecase.GetAllPlans(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if len(plans) < 1 {
		t.Fail()
		return
	}
	if plans[0].ID != cp.ID {
		t.Fail()
		return
	}
}

func TestSubscribeToPlan(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		AccountType: "admin",
	}
	createdUser, err := userRepo.CreateCustomer(context.Background(), testUser)
	if err != nil {
		t.Error(err)
		return
	}
	plan := &entities.Plan{}
	plan, err = subscriptionRepo.CreatePlan(context.Background(), plan)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
	subscription, err := subscriptionUsecase.SubscribeToPlan(ctx, plan.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if subscription == nil {
		t.Fail()
		return
	}
}

func TestUpgradePlan(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		AccountType: "admin",
	}
	createdUser, err := userRepo.CreateCustomer(context.Background(), testUser)
	if err != nil {
		t.Error(err)
		return
	}
	plan := &entities.Plan{
		CountOfOffers: 5,
	}
	plan, err = subscriptionRepo.CreatePlan(context.Background(), plan)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
	subscription, err := subscriptionUsecase.SubscribeToPlan(ctx, plan.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if subscription == nil {
		t.Fail()
		return
	}
	plan2 := &entities.Plan{
		CountOfOffers: 10,
	}
	plan2, err = subscriptionRepo.CreatePlan(context.Background(), plan2)
	if err != nil {
		t.Error(err)
		return
	}
	subscription, err = subscriptionUsecase.UpgradePlan(ctx, plan2.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if subscription == nil {
		t.Fail()
		return
	}
	if subscription.RemainingOffers != plan.CountOfOffers+plan2.CountOfOffers {
		t.Fail()
		return
	}
}

func TestRenewPlan(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		AccountType: "admin",
	}
	createdUser, err := userRepo.Create(context.Background(), testUser)
	if err != nil {
		t.Error(err)
		return
	}
	plan := &entities.Plan{
		CountOfOffers: 5,
	}
	plan, err = subscriptionRepo.CreatePlan(context.Background(), plan)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
	subscription, err := subscriptionUsecase.SubscribeToPlan(ctx, plan.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if subscription == nil {
		t.Fail()
		return
	}
	subscription, err = subscriptionUsecase.RenewPlan(ctx)
	if err != nil {
		t.Fail()
		return
	}
	if subscription == nil {
		t.Fail()
		return
	}
	if subscription.RemainingOffers != 2*plan.CountOfOffers {
		t.Fail()
		return
	}
}

func TestGetMySubscription(t *testing.T) {
	setupTest()
	testUser := &entities.User{
		AccountType: "admin",
	}
	createdUser, err := userRepo.Create(context.Background(), testUser)
	if err != nil {
		t.Error(err)
		return
	}
	plan := &entities.Plan{
		CountOfOffers: 5,
	}
	plan, err = subscriptionRepo.CreatePlan(context.Background(), plan)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.WithValue(context.Background(), entities.UserIDKey, createdUser.ID)
	subscription, err := subscriptionUsecase.SubscribeToPlan(ctx, plan.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if subscription == nil {
		t.Fail()
		return
	}
	returnedSubscription, err := subscriptionUsecase.GetMySubscription(ctx)
	if err != nil {
		t.Fail()
		return
	}
	if returnedSubscription == nil {
		t.Fail()
		return
	}
	if returnedSubscription.ID != subscription.ID {
		t.Fail()
		return
	}
}
