package usecase

import (
	"context"
	"testing"

	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/user"
)

var mockUserRepo user.Repository
var mockUserUsecase user.Usecase

type mockUserRepository struct{}
type mockSubscriptionRepository struct{}
type mockReviewRepository struct{}
type mockBranchRepository struct{}

func SetupTests() {
	mockUserRepo := &mockUserRepository{}
	mockSubscriptionRepo := &mockSubscriptionRepository{}
	mockReviewRepo := &mockReviewRepository{}
	mockBranchRepo := &mockBranchRepository{}
	mockUserUsecase = CreateUserUsecase(mockUserRepo, mockSubscriptionRepo, mockReviewRepo, mockBranchRepo)
}

func (r *mockUserRepository) CreateCustomer(_ context.Context, u *entities.User) (*entities.User, error) {
	return u, nil
}

func (r *mockUserRepository) CreatePartner(_ context.Context, u *entities.User, profile *entities.PartnerProfile) (*entities.User, error) {
	u.PartnerProfile = *profile
	return u, nil
}

func (r *mockUserRepository) UpdateCustomer(_ context.Context, u *entities.User) (*entities.User, error) {
	return u, nil
}

func (r *mockUserRepository) UpdatePartner(_ context.Context, u *entities.User, profile *entities.PartnerProfile) (*entities.User, error) {
	u.PartnerProfile = *profile
	return u, nil
}
func (r *mockUserRepository) GetByID(_ context.Context, ID uint) (*entities.User, error) {
	user := entities.User{}
	user.ID = ID
	return &user, nil
}
func (r *mockUserRepository) GetByEmail(_ context.Context, email string) (*entities.User, error) {
	user := entities.User{}
	user.Email = email
	return &user, nil
}
func (r *mockUserRepository) GetByPhone(_ context.Context, phone string) (*entities.User, error) {
	user := entities.User{}
	user.Mobile = phone
	return &user, nil
}
func (r *mockUserRepository) SoftDelete(_ context.Context, u *entities.User) (*entities.User, error) {
	return u, nil
}
func (r *mockUserRepository) Delete(_ context.Context, u *entities.User) (*entities.User, error) {
	return u, nil
}
func (r *mockUserRepository) UpdatePartnerProfile(_ context.Context, p *entities.PartnerProfile) (*entities.PartnerProfile, error) {
	return p, nil
}
func (r *mockUserRepository) DeletePartnerProfile(_ context.Context, p *entities.PartnerProfile) (*entities.PartnerProfile, error) {
	return p, nil
}
func (r *mockUserRepository) GetPartnerProfileByUserID(_ context.Context, userID uint) (*entities.PartnerProfile, error) {
	profile := entities.PartnerProfile{
		PartnerID: userID,
	}
	return &profile, nil
}
func (r *mockUserRepository) CreatePartnerPhoto(_ context.Context, photo *entities.PartnerPhoto) (*entities.PartnerPhoto, error) {
	return photo, nil
}
func (r *mockUserRepository) DeletePartnerPhoto(_ context.Context, photoID uint) (*entities.PartnerPhoto, error) {
	photo := entities.PartnerPhoto{}
	photo.ID = photoID
	return &photo, nil
}
func (r *mockUserRepository) GetPartnerPhotosByPartnerProfileID(_ context.Context, partnerID uint) ([]entities.PartnerPhoto, error) {
	photo := entities.PartnerPhoto{}
	photo.PartnerProfileID = partnerID
	photos := []entities.PartnerPhoto{photo}
	return photos, nil
}
func (r *mockUserRepository) SaveOTP(_ context.Context, password *entities.OTP) (*entities.OTP, error) {
	return password, nil
}
func (r *mockUserRepository) GetOTPByUser(_ context.Context, userID uint) (*entities.OTP, error) {
	otp := entities.OTP{
		UserID: userID,
	}
	return &otp, nil
}
func (r *mockUserRepository) GetCustomersCount(_ context.Context) (int, error) {
	return 5, nil
}
func (r *mockUserRepository) GetPartnersCount(_ context.Context) (int, error) {
	return 5, nil
}
func (r *mockUserRepository) GetAllParnters(_ context.Context) ([]entities.User, error) {
	partner1 := entities.User{}
	partner2 := entities.User{}
	partners := []entities.User{partner1, partner2}
	return partners, nil
}
func (r *mockUserRepository) GetAllCustomers(_ context.Context) ([]entities.User, error) {
	customer1 := entities.User{}
	customer2 := entities.User{}
	customers := []entities.User{customer1, customer2}
	return customers, nil
}
func (r *mockUserRepository) GetNotApprovedPartners(_ context.Context) ([]entities.User, error) {
	partner1 := entities.User{}
	partner2 := entities.User{}
	partners := []entities.User{partner1, partner2}
	return partners, nil
}
func (r *mockUserRepository) CreateExclusiveRecord(_ context.Context, partnerID uint) error {
	return nil
}
func (r *mockUserRepository) GetExclusiveOffers(_ context.Context) ([]entities.Exclusive, error) {
	exclusive := entities.Exclusive{}
	exclusives := []entities.Exclusive{exclusive}
	return exclusives, nil
}
func (r *mockUserRepository) GetExclusiveOffer(_ context.Context, partnerID uint) (*entities.Exclusive, error) {
	exclusive := entities.Exclusive{PartnerID: partnerID}
	return &exclusive, nil
}
func (r *mockUserRepository) DeleteExclusiveOffer(_ context.Context, exclusive *entities.Exclusive) (*entities.Exclusive, error) {
	return exclusive, nil
}

// CreatePlan creates a new subscription plan
func (r *mockSubscriptionRepository) CreatePlan(_ context.Context, p *entities.Plan) (*entities.Plan, error) {
	return p, nil
}

// GetPlanByID returns the plan with the given ID
func (r *mockSubscriptionRepository) GetPlanByID(_ context.Context, planID uint) (*entities.Plan, error) {
	var plan entities.Plan
	plan.ID = planID
	return &plan, nil
}

// DeletePlan deletes the given subscription plan
func (r *mockSubscriptionRepository) DeletePlan(_ context.Context, p *entities.Plan) (*entities.Plan, error) {
	return p, nil
}

// UpdatePlan saves the given plan
func (r *mockSubscriptionRepository) UpdatePlan(_ context.Context, plan *entities.Plan) (*entities.Plan, error) {
	return plan, nil
}

// GetPlans returns all plan records
func (r *mockSubscriptionRepository) GetPlans(_ context.Context) ([]entities.Plan, error) {
	var plans []entities.Plan
	plans[0] = entities.Plan{}
	return plans, nil
}

// CreateSubscription creates a new subscription record
func (r *mockSubscriptionRepository) CreateSubscription(_ context.Context, s *entities.Subscription) (*entities.Subscription, error) {
	return s, nil
}

// DeleteSubscription soft deletes the subscription record
func (r *mockSubscriptionRepository) DeleteSubscription(_ context.Context, s *entities.Subscription) (*entities.Subscription, error) {
	return s, nil
}

// DecrementOffersCount decrements the count of remaining offers in the given subscription
func (r *mockSubscriptionRepository) DecrementOffersCount(_ context.Context, s *entities.Subscription) (*entities.Subscription, error) {
	return s, nil
}

// GetSubscriptionByID returns the subscription with the given ID
func (r *mockSubscriptionRepository) GetSubscriptionByID(_ context.Context, ID uint) (*entities.Subscription, error) {
	var subscription entities.Subscription
	return &subscription, nil
}

// GetSubscriptionByUser return the active subscription of the given user
func (r *mockSubscriptionRepository) GetSubscriptionByUser(_ context.Context, userID uint) (*entities.Subscription, error) {
	var subscription entities.Subscription
	return &subscription, nil
}

// ExpireSubscription sets the given subscription as expired
func (r *mockSubscriptionRepository) ExpireSubscription(_ context.Context, s *entities.Subscription) (*entities.Subscription, error) {
	s.Expired = true
	return s, nil
}

// Create creates a new review DB record
func (r *mockReviewRepository) Create(ctx context.Context, review *entities.Review) (*entities.Review, error) {
	return review, nil
}

// Update updates the db record of the given review
func (r *mockReviewRepository) Update(ctx context.Context, review *entities.Review) (*entities.Review, error) {
	return review, nil
}

// Delete soft deletes the db record of the given review
func (r *mockReviewRepository) Delete(ctx context.Context, review *entities.Review) (*entities.Review, error) {
	return review, nil
}

// GetByID returns the db record of the review with the given ID
func (r *mockReviewRepository) GetByID(ctx context.Context, ID uint) (*entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var review entities.Review
	review.ID = ID
	cancelFunc()
	return &review, nil
}

// GetByCustomerID returns an array of records of reviews with the given customer ID.
func (r *mockReviewRepository) GetByCustomerID(ctx context.Context, ID uint) ([]entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var reviews []entities.Review
	reviews[0] = entities.Review{}
	cancelFunc()
	return reviews, nil
}

// GetByPartnerID returns an array of records of reviews with the given customer ID.
func (r *mockReviewRepository) GetByPartnerID(ctx context.Context, ID uint) ([]entities.Review, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var reviews []entities.Review
	reviews[0] = entities.Review{}
	cancelFunc()
	return reviews, nil
}

func (r *mockReviewRepository) GetAverageRatings(ctx context.Context, partnerID uint) (float64, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var avr float64
	avr = 4.5
	cancelFunc()
	return avr, nil
}

// Create creates a new branch
func (r *mockBranchRepository) Create(ctx context.Context, branch *entities.Branch) (*entities.Branch, error) {
	return branch, nil
}

// Delete deletes the given branch
func (r *mockBranchRepository) Delete(ctx context.Context, branch *entities.Branch) (*entities.Branch, error) {
	return branch, nil
}

// Edit saves the given branch
func (r *mockBranchRepository) Edit(ctx context.Context, branch *entities.Branch) (*entities.Branch, error) {
	return branch, nil
}

// GetByID returns the branch with the given ID
func (r *mockBranchRepository) GetByID(ctx context.Context, branchID uint) (*entities.Branch, error) {
	var b entities.Branch
	b.ID = branchID
	return &b, nil
}

// GetByOwner returns the list of branches associated with the given owner ID
func (r *mockBranchRepository) GetByOwner(ctx context.Context, ownerID uint) ([]entities.Branch, error) {
	var b []entities.Branch
	b[0] = entities.Branch{}
	owner := entities.User{}
	owner.ID = ownerID
	b[0].Owner = &owner
	return b, nil
}

// GetByLocation returns the list of branches with the given city and country
func (r *mockBranchRepository) GetByLocation(ctx context.Context, country string, cityID uint) ([]entities.Branch, error) {
	var b []entities.Branch
	b[0] = entities.Branch{
		Country: country,
	}
	return b, nil
}

// CreateCategory creates a DB record of the given category
func (r *mockBranchRepository) CreateCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	cancelFunc()
	return category, nil
}

func (r *mockBranchRepository) GetCategories(ctx context.Context) ([]entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var categories []entities.Category
	categories[0] = entities.Category{}
	cancelFunc()
	return categories, nil

}

// DeleteCategory deletes the category with the given ID.
func (r *mockBranchRepository) DeleteCategory(ctx context.Context, categoryID uint) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var category entities.Category
	category.ID = categoryID
	cancelFunc()
	return &category, nil
}

// GetCategoryByID returns the category with the given ID
func (r *mockBranchRepository) GetCategoryByID(ctx context.Context, categoryID uint) (*entities.Category, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var category entities.Category
	cancelFunc()
	return &category, nil
}

// GetBranchesByCategory returns an array of branches associated with the category with the given ID.
func (r *mockBranchRepository) GetBranchesByCategory(ctx context.Context, categoryID uint) ([]entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var branches []entities.Branch
	branches[0] = entities.Branch{}
	cancelFunc()
	return branches, nil
}

// CreateCity creates a new city
func (r *mockBranchRepository) CreateCity(ctx context.Context, city *entities.City) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	cancelFunc()
	return city, nil
}

// DeleteCity deletes the given city
func (r *mockBranchRepository) DeleteCity(ctx context.Context, city *entities.City) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	cancelFunc()
	return city, nil
}

// UpdateCity saves the given city
func (r *mockBranchRepository) UpdateCity(ctx context.Context, city *entities.City) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	cancelFunc()
	return city, nil

}

// GetCityByID gets the city with the given ID
func (r *mockBranchRepository) GetCityByID(ctx context.Context, cityID uint) (*entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var city entities.City
	city.ID = cityID
	cancelFunc()
	return &city, nil
}

// GetCities returns all the citites
func (r *mockBranchRepository) GetCities(ctx context.Context) ([]entities.City, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var cities []entities.City
	cities[0] = entities.City{}
	cancelFunc()
	return cities, nil
}

// SearchBranches searches branches with brand name , category id , city id
func (r *mockBranchRepository) SearchBranches(ctx context.Context, city uint, category uint, name string) ([]entities.Branch, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	var branches []entities.Branch
	branches[0] = entities.Branch{}
	cancelFunc()
	return branches, nil
}

func TestCreateCustomer(t *testing.T) {
	SetupTests()
	testUser := entities.User{
		AccountType: "user",
	}
	_, err := mockUserUsecase.CreateCustomer(context.Background(), &testUser, "123456789")
	if err != nil {
		t.Error(err)
	}
}

func TestCreatePartner(t *testing.T) {
	SetupTests()
	testPartner := entities.User{
		AccountType: "partner",
	}
	testPartnerProfile := entities.PartnerProfile{}
	_, err := mockUserUsecase.CreatePartner(context.Background(), &testPartner, &testPartnerProfile, "1234456")
	if err != nil {
		t.Error(err)
	}
}

func TestCreatePartnerWithCity(t *testing.T) {
	SetupTests()
	testCity := entities.City{}
	testCity.ID = 1
	testPartner := entities.User{
		AccountType: "partner",
		CityID:      1,
		City:        testCity,
		PartnerProfile: entities.PartnerProfile{
			CityID: 1,
			City:   testCity,
		},
	}
	u, err := mockUserUsecase.CreatePartner(context.Background(), &testPartner, &testPartner.PartnerProfile, "12345566")
	if err != nil {
		t.Error(err)
	}
	if u.GetCity().ID != testCity.ID {
		t.Fail()
	}
	if u.PartnerProfile.City.ID != testCity.ID {
		t.Fail()
	}
}
