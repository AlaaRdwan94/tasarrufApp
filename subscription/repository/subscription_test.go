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
	db.DropTable(entities.Subscription{})
	db.DropTable(entities.Plan{})
	db.AutoMigrate(entities.Subscription{})
	db.AutoMigrate(entities.Plan{})
	return db, nil
}

func TestCreatePlan(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testPlan := &entities.Plan{}
	tp, err := repo.CreatePlan(context.Background(), testPlan)
	if err != nil {
		t.Error(err)
	}
	if tp == nil {
		t.Fail()
	}
}

func TestGetPlanByID(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testPlan := &entities.Plan{}
	tp, err := repo.CreatePlan(context.Background(), testPlan)
	if err != nil {
		t.Error(err)
		return
	}
	if tp == nil {
		t.Fail()
		return
	}
	returnedPlan, err := repo.GetPlanByID(context.Background(), tp.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if returnedPlan.ID != tp.ID {
		t.Fail()
		return
	}
}

func TestDeletePlan(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testPlan := &entities.Plan{}
	tp, err := repo.CreatePlan(context.Background(), testPlan)
	if err != nil {
		t.Error(err)
	}
	if tp == nil {
		t.Fail()
		return
	}
	deletedPlan, err := repo.DeletePlan(context.Background(), tp)
	if err != nil {
		t.Error(err)
	}
	if deletedPlan == nil {
		t.Fail()
	}
	if deletedPlan.ID != tp.ID {
		t.Fail()
	}
}

func TestGetPlans(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testPlan := &entities.Plan{}
	tp, err := repo.CreatePlan(context.Background(), testPlan)
	if err != nil {
		t.Error(err)
	}
	if tp == nil {
		t.Fail()
		return
	}
	ps, err := repo.GetPlans(context.Background())
	if err != nil {
		t.Error(err)
	}
	if len(ps) < 1 {
		t.Fail()
	}
	if ps[0].ID != tp.ID {
		t.Fail()
	}
}

func TestCreateSubscription(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testSubscription := &entities.Subscription{}
	createdSubscription, err := repo.CreateSubscription(context.Background(), testSubscription)
	if err != nil {
		t.Error(err)
		return
	}
	if createdSubscription.ID == 0 {
		t.Fail()
		return
	}
}

func TestDeleteSubscription(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testSubscription := &entities.Subscription{}
	createdSubscription, err := repo.CreateSubscription(context.Background(), testSubscription)
	if err != nil {
		t.Error(err)
		return
	}
	deletedSubscription, err := repo.DeleteSubscription(context.Background(), createdSubscription)
	if err != nil {
		t.Error(err)
		return
	}
	if deletedSubscription.ID != createdSubscription.ID {
		t.Fail()
		return
	}
}

func TestDecrementOffersCount(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testSubscription := &entities.Subscription{
		RemainingOffers: 3,
	}
	createdSubscription, err := repo.CreateSubscription(context.Background(), testSubscription)
	if err != nil {
		t.Error(err)
		return
	}
	cs := *createdSubscription
	decrementedSubscription, err := repo.DecrementOffersCount(context.Background(), createdSubscription)
	if err != nil {
		t.Error(err)
		return
	}
	if decrementedSubscription == nil {
		t.Fail()
		return
	}
	if decrementedSubscription.RemainingOffers != cs.RemainingOffers {
		t.Fail()
		return
	}
}

func TestGetSubscriptionByID(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testSubscription := &entities.Subscription{
		RemainingOffers: 3,
	}
	createdSubscription, err := repo.CreateSubscription(context.Background(), testSubscription)
	if err != nil {
		t.Error(err)
		return
	}
	returnedSubscription, err := repo.GetSubscriptionByID(context.Background(), createdSubscription.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if returnedSubscription == nil {
		t.Fail()
		return
	}
}

func TestGetSubscriptionByUser(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testSubscription := &entities.Subscription{
		UserID: 1,
	}
	createdSubscription, err := repo.CreateSubscription(context.Background(), testSubscription)
	if err != nil {
		t.Error(err)
		return
	}
	returnedSubscription, err := repo.GetSubscriptionByUser(context.Background(), createdSubscription.UserID)
	if err != nil {
		t.Error(err)
		return
	}
	if returnedSubscription.ID != createdSubscription.ID {
		t.Fail()
		return
	}
	if returnedSubscription == nil {
		t.Fail()
		return
	}
}

func TestExpireSubscription(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testSubscription := &entities.Subscription{
		Expired: false,
	}
	createdSubscription, err := repo.CreateSubscription(context.Background(), testSubscription)
	if err != nil {
		t.Error(err)
		return
	}
	if createdSubscription.IsExpired() {
		t.Fail()
		return
	}
	expiredSubscription, err := repo.ExpireSubscription(context.Background(), createdSubscription)
	if err != nil {
		t.Error(err)
		return
	}
	if expiredSubscription == nil {
		t.Fail()
		return
	}
	if !expiredSubscription.IsExpired() {
		t.Fail()
		return
	}
}

func TestRankPlanUp(t *testing.T) {
	db, err := connectToDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	repo := CreateSubscriptionRepository(db)
	testPlan1 := &entities.Plan{
		Rank: 1,
	}
	testPlan2 := &entities.Plan{
		Rank: 2,
	}
	tp1, err := repo.CreatePlan(context.Background(), testPlan1)
	if err != nil {
		t.Error(err)
	}
	if tp1 == nil {
		t.Fail()
	}
	tp2, err := repo.CreatePlan(context.Background(), testPlan2)
	if err != nil {
		t.Error(err)
	}
	if tp2 == nil {
		t.Fail()
	}
	tp1.Rank = tp1.Rank + 1
	newPlans, err := repo.RankPlanUp(context.Background(), tp1)
	if err != nil {
		t.Error(err)
	}
	if newPlans[0].Rank != 2 {
		t.Fail()
	}
	if newPlans[1].Rank != 1 {
		t.Fail()
	}
}
