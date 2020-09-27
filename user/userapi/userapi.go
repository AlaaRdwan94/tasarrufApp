package userapi

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/user"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// UserAPI defines the api handler for user routes
type UserAPI struct {
	UserUsecase user.Usecase
}

type newUserRequest struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Phone             string `json:"phone"`
	Mobile            string `json:"mobile"`
	Email             string `json:"email"`
	Country           string `json:"country"`
	AccountType       string `json:"accountType"`
	CityID            uint   `json:"cityID"`
	Password          string `json:"password"`
	BrandName         string `json:"brandName"`
	MainBranchAddress string `json:"mainBranchAddress"`
	DateOfBirth       string `json:"dateOfBirth"`
}

type updateMainBranch struct {
	DiscountValue     float64 `json:"discountValue"`
	MainBranchAddress string  `json:"mainBranchAddress"`
	Phone             string  `json:"phone"`
	CityID            uint    `json:"cityID"`
	Country           string  `json:"country"`
	BrandName         string  `json:"brandName"`
	OfferDiscription  string  `json:"offerDiscription"`
	CategoryID        int     `json:"categoryID"`
}

type updateUserRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Country     string `json:"country"`
	CityID      uint   `json:"cityID"`
	DateOfBirth string `json:"dateOfBirth"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
}

type emailLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type phoneLoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type verifyUserRequest struct {
	Code string `json:"code"`
}

type forgetPasswordRequest struct {
	Mobile string `json:"mobile"`
}

type updatePassword struct {
	NewPassword string `json:"newPassword"`
}

// CreateUserAPI creates a new user api instance
func CreateUserAPI(u user.Usecase) UserAPI {
	api := UserAPI{
		UserUsecase: u,
	}
	return api
}

// Validate validates the new user request
func (req *newUserRequest) ValidateUser() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.AccountType, validation.Required, is.LowerCase),
		validation.Field(&req.FirstName, validation.Required, validation.Length(2, 50), is.LowerCase),
		validation.Field(&req.LastName, validation.Required, validation.Length(2, 50), is.LowerCase),
		validation.Field(&req.Country, validation.Required, validation.Length(2, 50), is.LowerCase),
		validation.Field(&req.CityID, validation.Required),
		validation.Field(&req.DateOfBirth, validation.Required),
		validation.Field(&req.Mobile, validation.Required),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.Required, validation.Length(8, 50)),
	)
}

func (req *newUserRequest) ValidatePartner() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.AccountType, validation.Required, is.LowerCase),
		validation.Field(&req.FirstName, validation.Required, validation.Length(2, 50), is.LowerCase),
		validation.Field(&req.LastName, validation.Required, validation.Length(2, 50), is.LowerCase),
		validation.Field(&req.Country, validation.Required, validation.Length(2, 50), is.LowerCase),
		validation.Field(&req.BrandName, validation.Required),
		validation.Field(&req.Mobile, validation.Required),
		validation.Field(&req.Phone, validation.Required),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.MainBranchAddress, validation.Required),
		validation.Field(&req.Password, validation.Required, validation.Length(8, 50)),
	)
}

func (req *updateUserRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.FirstName, validation.Required, is.LowerCase),
		validation.Field(&req.LastName, validation.Required, is.LowerCase),
		validation.Field(&req.Country, validation.Required, is.LowerCase),
		validation.Field(&req.CityID, validation.Required),
		validation.Field(&req.Mobile, validation.Required),
		validation.Field(&req.Email, validation.Required, is.Email),
	)
}
func (req *updateUserRequest) CastUser(original *entities.User) *entities.User {
	original.FirstName = req.FirstName
	original.LastName = req.LastName
	original.Country = req.Country
	original.CityID = req.CityID
	original.Email = req.Email
	return original
}

// Validate handles new user request validations
func (req *newUserRequest) Validate() error {
	switch req.AccountType {
	case "user":
		return req.ValidateUser()
	case "partner":
		return req.ValidatePartner()
	case "admin":
		return req.ValidateUser()
	default:
		return req.ValidateUser()
	}
}

// CreateUser handles the user creation endpoint
func (h *UserAPI) CreateUser(c *gin.Context) {
	ctx := context.Background()
	var req newUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been a trouble sending your information to the server, please try again", err)
		return
	}
	req.Email = strings.ToLower(req.Email)
	err = req.Validate()
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	var dbt time.Time
	if req.AccountType == "user" {
		dbt, err = time.Parse(time.RFC3339, req.DateOfBirth)
		if err != nil {
			entities.SendValidationError(c, "Please make sure that the date of birth is valid", err)
			return
		}
	}
	newUser := entities.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Country:     req.Country,
		CityID:      req.CityID,
		DateOfBirth: dbt,
		AccountType: req.AccountType,
		Mobile:      req.Mobile,
	}
	var user *entities.User
	if newUser.AccountType == "partner" {
		profile := entities.PartnerProfile{
			MainBranchAddress: req.MainBranchAddress,
			Phone:             req.Phone,
			CityID:            req.CityID,
			Country:           req.Country,
			BrandName:         req.BrandName,
			Approved:          false,
		}
		user, err = h.UserUsecase.CreatePartner(ctx, &newUser, &profile, req.Password)
		if err != nil {
			entities.SendValidationError(c, "This email or phone number is already registered, please login instead", err)
			return
		}
	} else {
		user, err = h.UserUsecase.CreateCustomer(ctx, &newUser, req.Password)
		if err != nil {
			entities.SendValidationError(c, "This email or phone number is already registered, please login instead", err)
			return
		}
	}
	token, err := entities.GenerateAuthToken(user.ID)
	if err != nil {
		entities.SendServerError(c, "There has been a server error , please try again", err)
		return
	}
	c.JSON(200, gin.H{
		"message": "User Created Successfully",
		"user":    user,
		"token":   token,
	})
}

// IsEmailRegistered returns a boolean value that indicates if the given email is registered on the platform.
//
// It's meant to be called with the email input from the client side to give live feedback on email availability.
func (h *UserAPI) IsEmailRegistered(c *gin.Context) {
	ctx := context.Background()
	email := c.Query("email")
	isRegistered, err := h.UserUsecase.IsEmailRegistered(ctx, email)
	if err != nil {
		entities.SendServerError(c, "A server error has occured , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isRegistered": isRegistered,
	})
}

// IsPhoneRegistered returns a boolean value that indicates if the given phone is registered on the platform.
//
// It's meant to be called with the phone input from the client side to give live feedback on phone availability.
func (h *UserAPI) IsPhoneRegistered(c *gin.Context) {
	ctx := context.Background()
	phone := c.Query("phone")
	isRegistered, err := h.UserUsecase.IsPhoneRegistered(ctx, phone)
	if err != nil {
		entities.SendServerError(c, "A server error has occured , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isRegistered": isRegistered,
	})
}

// EmailLogin handles login with Email
func (h *UserAPI) EmailLogin(c *gin.Context) {
	ctx := context.Background()
	var req emailLoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error sending your information to the server, please try again", err)
		return
	}
	req.Email = strings.ToLower(req.Email)
	user, token, err := h.UserUsecase.EmailLogin(ctx, req.Email, req.Password)
	if err != nil {
		log.Error(err)
		if err.Error() == "not registered" {
			entities.SendNotFoundError(c, "This email is not registered, please signup instead", err)
			return
		}
		if err.Error() == "Your account has been deactivated please contact us" {
			entities.SendAuthError(c, err.Error(), err)
			return
		}
		entities.SendNotFoundError(c, "Password is incorrect", err)
		return
	}
	c.JSON(200, gin.H{
		"message": "Login Successful",
		"user":    user,
		"token":   token,
	})
}

// PhoneLogin handles login with phone number
func (h *UserAPI) PhoneLogin(c *gin.Context) {
	ctx := context.Background()
	var req phoneLoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been a trouble sending your information to the server, please try again", err)
		return
	}
	user, token, err := h.UserUsecase.PhoneLogin(ctx, req.Phone, req.Password)
	if err != nil {
		if err.Error() == "not registered" {
			entities.SendNotFoundError(c, "This phone number is not registered, please signup instaed", err)
			return
		}
		entities.SendNotFoundError(c, "Password is incorrect", err)
		return
	}
	c.JSON(200, gin.H{
		"message": "Login Successful",
		"user":    user,
		"token":   token,
	})
}

// GetUser handles the get user endpoint of the currently logged in user
func (h *UserAPI) GetUser(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	user, err := h.UserUsecase.GetUser(ctx, userID)
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	c.JSON(200, gin.H{
		"user": user,
	})
}

// DeleteUser handles the delete endpoint of the currently logged in user
func (h *UserAPI) DeleteUser(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	user, err := h.UserUsecase.DeleteUser(ctx, userID)
	if err != nil {
		entities.SendValidationError(c, "You have been logged out , please login", err)
		return
	}
	c.JSON(200, gin.H{
		"success": "user deleted successfully",
		"user":    user,
	})
}

// AdminDeleteUser handles the delete endpoint of the currently logged in user
func (h *UserAPI) AdminDeleteUser(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	user, err := h.UserUsecase.AdminDeleteUser(ctx, uint(id))
	if err != nil {
		entities.SendValidationError(c, "You have been logged out , please login", err)
		return
	}
	c.JSON(200, gin.H{
		"success": "user deleted successfully",
		"user":    user,
	})
}

// DeletePartnerPhoto handles the DELETE /partner-photo end point
func (h *UserAPI) DeletePartnerPhoto(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	photoID, err := strconv.ParseInt(c.Param("photoID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}

	user, err := h.UserUsecase.DeletePartnerPhoto(ctx, uint(photoID))
	if err != nil {
		entities.SendValidationError(c, "You don't have permission to delete this photo", err)
		return
	}
	c.JSON(200, gin.H{
		"success": "photo deleted successfully",
		"user":    user,
	})
}

// UpdateProfileImage handles PUT request to /profile-image
func (h *UserAPI) UpdateProfileImage(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	fileHeader, err := c.FormFile("profileImage")
	if err != nil {
		entities.SendParsingError(c, "There has been an error processing your request, please try again", err)
		return
	}
	user, err := h.UserUsecase.UpdateProfileImage(ctx, fileHeader)
	if err != nil {
		entities.SendParsingError(c, "There has been an error processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Success": "Profile Image Updated Successfully",
		"user":    user,
	})
}

// UpdateTradeLicense handles PUT request to /trade-license
func (h *UserAPI) UpdateTradeLicense(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	fileHeader, err := c.FormFile("tradeLicense")
	if err != nil {
		entities.SendParsingError(c, "There has been an error processing your request, please try again", err)
		return
	}
	user, err := h.UserUsecase.UpdateTradeLicense(ctx, fileHeader)
	if err != nil {
		entities.SendParsingError(c, "There has been an error processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Trade License Updated Successfully",
		"user":    user,
	})
}

// AddPartnerPhoto handles PUT request to /partner-photo
func (h *UserAPI) AddPartnerPhoto(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		entities.SendParsingError(c, "There has been an error processing your request, please try again", err)
		return
	}
	user, err := h.UserUsecase.AddPartnerPhoto(ctx, fileHeader)
	if err != nil {
		entities.SendParsingError(c, "There has been an error processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Partner Photo Added Successfully",
		"user":    user,
	})
}

// ResendVerificationCode handles POST request to /user/resend-verification-code
func (h *UserAPI) ResendVerificationCode(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	err := h.UserUsecase.ResendVerficationCode(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error processing your request , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Verification ",
	})
}

// VerifyUser handles POST request to /user/verify
func (h *UserAPI) VerifyUser(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req verifyUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error parsing your request , please try again", err)
		return
	}
	user, err := h.UserUsecase.VerifyUser(ctx, req.Code)
	if err != nil {
		entities.SendValidationError(c, "There has been an error processing your request , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Verification ",
		"user":    user,
	})
}

// ForgetPassword handles POST request to /forget-password
func (h *UserAPI) ForgetPassword(c *gin.Context) {
	ctx := context.Background()
	var req forgetPasswordRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error parsing your request , please try again", err)
		return
	}
	err = h.UserUsecase.RecoverPassword(ctx, req.Mobile)
	if err != nil {
		entities.SendValidationError(c, "There has been an error processing your request , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "One time password was sent to your mobile phone",
	})
}

// ForgetPassword handles POST request to /forget-password
func (h *UserAPI) AdminForgetPassword(c *gin.Context) {
	ctx := context.Background()
	var req forgetPasswordRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error parsing your request , please try again", err)
		return
	}
	err = h.UserUsecase.AdminRecoverPassword(ctx, req.Mobile)
	if err != nil {
		entities.SendValidationError(c, "There has been an error processing your request , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "One time password was sent to your mobile phone",
	})
}

// UpdatePartnerProfile updates the partner profile
func (h *UserAPI) UpdatePartnerProfile(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req updateMainBranch
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error parsing your request , please try again", err)
		return
	}
	profile := &entities.PartnerProfile{
		DiscountValue:     req.DiscountValue,
		MainBranchAddress: req.MainBranchAddress,
		Phone:             req.Phone,
		CityID:            req.CityID,
		Country:           req.Country,
		BrandName:         req.BrandName,
		OfferDiscription:  req.OfferDiscription,
		CategoryID:        uint(req.CategoryID),
	}
	updatedUser, err := h.UserUsecase.UpdatePartnerProfile(ctx, profile)
	if err != nil {
		entities.SendValidationError(c, "There has been an error processing your request , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Main branch details updated successfully",
		"user":    updatedUser,
	})
}

// UpdatePassword handles POST request to /update-pass endpoint
func (h *UserAPI) UpdatePassword(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req updatePassword
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error parsing your request , please try again", err)
		return
	}
	updatedUser, err := h.UserUsecase.UpdatePassword(ctx, req.NewPassword)
	if err != nil {
		switch errors.Cause(err).Error() {
		case "new password is same as the old password":
			entities.SendValidationError(c, "updated password cannot be the same as the old password", err)
			return
		default:
			entities.SendValidationError(c, "There has been an error while processing your request , please try again", err)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "password updated successfully",
		"user":    updatedUser,
	})
}

func (h *UserAPI) ValidateCustomer(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	customerID, err := strconv.ParseInt(c.Query("customerID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	partnerID, err := strconv.ParseInt(c.Query("partnerID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	user, subscription, err := h.UserUsecase.ValidateCustomerPartnerIntegrity(ctx, uint(customerID), uint(partnerID))
	if err != nil {
		if errors.Cause(err).Error() == "user is not subscribed to any plan" {
			entities.SendValidationError(c, "This user is not subscribed to any plan", err)
			return
		}
		entities.SendValidationError(c, "This offer is not owned by your account", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"subscription": subscription,
	})
}

// UpdateUser handles PUT request to /user
func (h *UserAPI) UpdateUser(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req updateUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your request, please try again", err)
		return
	}
	req.Email = strings.ToLower(req.Email)
	user, err := h.UserUsecase.GetUser(ctx, userID)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	var dbt time.Time
	if user.AccountType == "user" {
		dbt, err = time.Parse(time.RFC3339, req.DateOfBirth)
		if err != nil {
			entities.SendValidationError(c, "Please make sure that the date of birth is valid", err)
			return
		}
	}
	user = req.CastUser(user)
	user.DateOfBirth = dbt
	updated, err := h.UserUsecase.UpdateUser(ctx, user)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "user updated successfully",
		"user":    updated,
	})
}

// GetPartnersCount returns the count of partners accounts
func (h *UserAPI) GetPartnersCount(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	count, err := h.UserUsecase.GetPartnersCount(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
	return
}

// GetCustomersCount returns the count of partners accounts
func (h *UserAPI) GetCustomersCount(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	count, err := h.UserUsecase.GetCustomersCount(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
	return
}

// GetAllCustomers handles GET /admin/customers
func (h *UserAPI) GetAllCustomers(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	customers, err := h.UserUsecase.GetAllCustomers(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": customers,
	})
}

// GetAllPartners handles GET /admin/partners
func (h *UserAPI) GetAllPartners(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partners, err := h.UserUsecase.GetAllParnters(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": partners,
	})
}

// GetNotApproved handles GET /partners/not-approved
func (h *UserAPI) GetNotApproved(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partners, err := h.UserUsecase.GetNotApprovedPartners(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": partners,
	})
}

// ApprovePartner handles POST /approve/:id
func (h *UserAPI) ApprovePartner(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partnerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	partner, err := h.UserUsecase.ApprovePartner(ctx, uint(partnerID))
	if partner == nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
	}
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "user approved successfully",
		"user":    partner,
	})
	return
}

// SetPartnerAsExclusive handles POST /exclusive/:id endpoint
func (h *UserAPI) SetPartnerAsExclusive(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partnerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	partner, err := h.UserUsecase.SetPartnerAsExclusive(ctx, uint(partnerID))
	if partner == nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Parter is now a member of exclusive group",
		"user":    partner,
	})
	return
}

// RemovePartnerAsExclusive handles DELETE /exclusive/:id endpoint
func (h *UserAPI) RemovePartnerAsExclusive(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partnerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	partner, err := h.UserUsecase.RemovePartnerAsExclusive(ctx, uint(partnerID))
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Parter has been removed from exclusive group",
		"user":    partner,
	})
	return
}

// GetExclusivePartners handles GET /exclusive endpoint
func (h *UserAPI) GetExclusivePartners(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partners, err := h.UserUsecase.GetExclusivePartners(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": partners,
	})
	return
}

// GetCustomerByID handles GET /admin/customer/:id endpoint
func (h *UserAPI) GetCustomerByID(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	customer, err := h.UserUsecase.GetCustomerByID(ctx, uint(customerID))
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": customer,
	})
	return
}

// GetPartnerByID handles GET /admin/partner/:id endpoint
func (h *UserAPI) GetPartnerByID(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partnerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	partner, err := h.UserUsecase.GetPartnerByID(ctx, uint(partnerID))
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": partner,
	})
	return
}

// Share handles POST /share endpoint
func (h *UserAPI) Share(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	err := h.UserUsecase.Share(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "shared successfully",
	})
}

func (h *UserAPI) ToggleIsSharable(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partnerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	err = h.UserUsecase.ToggleIsSharable(ctx, uint(partnerID))
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Toggled is sharable successfully",
	})
}

// ToggleActive handles POST /admin/activate-user/:id
func (h *UserAPI) ToggleActive(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	user, err := h.UserUsecase.ToggleActiveProperty(ctx, uint(id))
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Toggled is sharable successfully",
		"user":    user,
	})
	return
}

// SearchUsers handles /admin/users/?q=
func (h *UserAPI) SearchUsers(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	searchTerm := c.Query("q")
	var result []entities.IUser
	result, err := h.UserUsecase.SearchUsers(ctx, searchTerm)
	if err != nil {
		entities.SendParsingError(c, "error searching user, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": result,
	})
}
