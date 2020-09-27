package userapi

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/gin-gonic/gin"
)

var mockUserAPI UserAPI
var router *gin.Engine

type mockUserUsecase struct{}
type successUserResp struct {
	Success string        `json:"success"`
	User    entities.User `json:"user"`
}

// CreateCustomer creates a new user
func (c *mockUserUsecase) CreateCustomer(_ context.Context, u *entities.User, newPassword string) (*entities.User, error) {
	return u, nil
}

// CreatePartner creates a new partner user
func (c *mockUserUsecase) CreatePartner(_ context.Context, u *entities.User, profile *entities.PartnerProfile, newPassword string) (*entities.User, error) {
	return u, nil
}

// IsEmailRegistered returns true if the given email is associated with a user
func (c *mockUserUsecase) IsEmailRegistered(_ context.Context, email string) (bool, error) {
	return true, nil
}

// IsPhoneRegistered returns true if the given phone is associated with a user
func (c *mockUserUsecase) IsPhoneRegistered(_ context.Context, phone string) (bool, error) {
	return true, nil
}

// EmailLogin performs login with Email
func (c *mockUserUsecase) EmailLogin(_ context.Context, email string, password string) (*entities.User, string, error) {
	u := entities.User{}
	u.Email = email
	return &u, "token", nil
}

// PhoneLogin performs login with phone
func (c *mockUserUsecase) PhoneLogin(_ context.Context, phone string, password string) (*entities.User, string, error) {
	u := entities.User{}
	u.Mobile = phone
	return &u, "token", nil
}

// GetUser retruns the user with the given ID
func (c *mockUserUsecase) GetUser(_ context.Context, ID uint) (*entities.User, error) {
	u := entities.User{}
	u.ID = ID
	return &u, nil
}

// DeleteUser retruns the user with the given ID
func (c *mockUserUsecase) DeleteUser(_ context.Context, ID uint) (*entities.User, error) {
	u := entities.User{}
	u.ID = ID
	return &u, nil
}

func (c *mockUserUsecase) UpdateProfileImage(ctx context.Context, _ *multipart.FileHeader) (*entities.User, error) {
	userID := ctx.Value(entities.UserIDKey).(uint)
	u := entities.User{}
	u.ID = userID
	return &u, nil
}

// UpdateTradeLicense uploads the given file to to filestore then updates the partnerProfile DB record
func (c *mockUserUsecase) UpdateTradeLicense(ctx context.Context, fileHeader *multipart.FileHeader) (entities.IPartner, error) {
	userID := ctx.Value(entities.UserIDKey).(uint)
	u := entities.Partner{}
	u.ID = userID
	return &u, nil
}

// AddPartnerPhoto uploads the given file to to filestore then creates the photo DB record
func (c *mockUserUsecase) AddPartnerPhoto(ctx context.Context, fileHeader *multipart.FileHeader) (*entities.User, error) {
	userID := ctx.Value(entities.UserIDKey).(uint)
	u := entities.User{}
	u.ID = userID
	return &u, nil
}

// DeletePartnerPhoto deletes the partner photo with the given ID
func (c *mockUserUsecase) DeletePartnerPhoto(ctx context.Context, photoID uint) (*entities.User, error) {
	userID := ctx.Value(entities.UserIDKey).(uint)
	u := entities.User{}
	u.ID = userID
	return &u, nil
}

// RecoverPassword recovers password
func (c *mockUserUsecase) RecoverPassword(ctx context.Context, mobile string) error {
	_ = ctx.Value(entities.UserIDKey).(uint)
	return nil
}

// ResendVerficationCode sends a 4-digit verification code to the user via SMS.
// Then saves a hash value for the generated code in the user DB record
func (c *mockUserUsecase) ResendVerficationCode(ctx context.Context) error {
	_ = ctx.Value(entities.UserIDKey).(uint)
	return nil
}

// VerifyUser sets the user as verified
func (c *mockUserUsecase) VerifyUser(ctx context.Context, code string) (*entities.User, error) {
	userID := ctx.Value(entities.UserIDKey).(uint)
	u := entities.User{}
	u.ID = userID
	return &u, nil
}

// UpdatePartnerProfile updates the partner profile associated with a user provided in the _
func (c *mockUserUsecase) UpdatePartnerProfile(ctx context.Context, profile *entities.PartnerProfile) (*entities.User, error) {
	userID := ctx.Value(entities.UserIDKey).(uint)
	u := entities.User{}
	u.ID = userID
	u.PartnerProfile = *profile
	return &u, nil
}

// UpdatePassword updates the user hashed password field
func (c *mockUserUsecase) UpdatePassword(ctx context.Context, newPassword string) (*entities.User, error) {
	userID := ctx.Value(entities.UserIDKey).(uint)
	u := entities.User{}
	u.ID = userID
	return &u, nil
}

// ValidateCustomerPartnerIntegrity validates the partner customer integrity
func (c *mockUserUsecase) ValidateCustomerPartnerIntegrity(ctx context.Context, customerID uint, partnerID uint) (*entities.User, *entities.Subscription, error) {
	_ = ctx.Value(entities.UserIDKey).(uint)
	u := entities.User{}
	u.ID = customerID
	sub := entities.Subscription{}
	sub.UserID = customerID
	return &u, &sub, nil
}

// UpdateUser saves the updates the fields of the given user
func (c *mockUserUsecase) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	userID := ctx.Value(entities.UserIDKey).(uint)
	u := entities.User{}
	u.ID = userID
	return &u, nil
}

// GetCustomersCount returns the count of customer users
func (c *mockUserUsecase) GetCustomersCount(ctx context.Context) (int, error) {
	return 5, nil
}

// GetPartnersCount returns the count of customer users
func (c *mockUserUsecase) GetPartnersCount(_ context.Context) (int, error) {
	return 5, nil
}

// GetAllCustomers returns all customer users
func (c *mockUserUsecase) GetAllCustomers(_ context.Context) ([]entities.User, error) {
	customers := make([]entities.User, 1)
	customers[0] = entities.User{}
	return customers, nil
}

// GetAllParnters returns all the partner users
func (c *mockUserUsecase) GetAllParnters(_ context.Context) ([]entities.User, error) {
	customers := make([]entities.User, 1)
	customers[0] = entities.User{}
	return customers, nil
}

// GetNotApprovedPartners returns a list of all not approved partner users
func (c *mockUserUsecase) GetNotApprovedPartners(_ context.Context) ([]entities.User, error) {
	customers := make([]entities.User, 1)
	customers[0] = entities.User{}
	return customers, nil
}

// ApprovePartner sets the partner as approved
func (c *mockUserUsecase) ApprovePartner(_ context.Context, partnerID uint) (*entities.User, error) {
	u := entities.User{}
	u.ID = partnerID
	return &u, nil
}

// SetPartnerAsExclusive sets the partner with the given ID as exclusive
func (c *mockUserUsecase) SetPartnerAsExclusive(_ context.Context, partnerID uint) (*entities.Partner, error) {
	u := entities.Partner{}
	u.ID = partnerID
	return &u, nil
}

// RemovePartnerAsExclusive removes the partner with the given ID from exclusive group
func (c *mockUserUsecase) RemovePartnerAsExclusive(_ context.Context, partnerID uint) (*entities.Partner, error) {
	u := entities.Partner{}
	u.ID = partnerID
	return &u, nil
}

// GetExclusivePartners returns the partners set as exclusive
func (c *mockUserUsecase) GetExclusivePartners(_ context.Context) ([]entities.Partner, error) {
	partners := make([]entities.Partner, 1)
	partners[0] = entities.Partner{}
	return partners, nil
}

// GetCustomerByID returns the customer with the given ID
func (c *mockUserUsecase) GetCustomerByID(_ context.Context, ID uint) (*entities.Customer, error) {
	u := entities.Customer{}
	u.ID = ID
	return &u, nil
}

// GetPartnerByID returns the partner with the given ID
func (c *mockUserUsecase) GetPartnerByID(_ context.Context, ID uint) (*entities.Partner, error) {
	u := entities.Partner{}
	u.ID = ID
	return &u, nil
}

func SetupTests() {
	mockUserAPI = CreateUserAPI(&mockUserUsecase{})
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
}

func TestCreatePartner(t *testing.T) {
	SetupTests()
	w := httptest.NewRecorder()
	partner := newUserRequest{
		FirstName:         "john",
		LastName:          "doe",
		Phone:             "123456789",
		Mobile:            "123456789",
		Email:             "johndoe@test.com",
		Country:           "turket",
		AccountType:       "partner",
		CityID:            1,
		Password:          "123456789",
		BrandName:         "Brandy",
		MainBranchAddress: "Address test",
		DateOfBirth:       "2006-02-02T15:04:05Z",
	}
	body, err := json.Marshal(partner)
	if err != nil {
		t.Error(err)
	}
	_ = router.POST("/", mockUserAPI.CreateUser)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fail()
	}
}

func TestCreatePartnerWithCity(t *testing.T) {
	SetupTests()
	w := httptest.NewRecorder()
	partner := newUserRequest{
		FirstName:         "john",
		LastName:          "doe",
		Phone:             "123456789",
		Mobile:            "123456789",
		Email:             "johndoe@test.com",
		Country:           "turket",
		AccountType:       "partner",
		CityID:            1,
		Password:          "123456789",
		BrandName:         "Brandy",
		MainBranchAddress: "Address test",
		DateOfBirth:       "2006-02-02T15:04:05Z",
	}
	body, err := json.Marshal(partner)
	if err != nil {
		t.Error(err)
	}
	_ = router.POST("/", mockUserAPI.CreateUser)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fail()
	}
	var resp successUserResp
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fail()
	}
	if resp.User.CityID != 1 {
		t.Fail()
	}
}
