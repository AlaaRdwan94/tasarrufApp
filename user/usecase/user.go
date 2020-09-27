package usecase

import (
	"context"
	"fmt"
	"github.com/ahmedaabouzied/tasarruf/branch"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/filestore"
	"github.com/ahmedaabouzied/tasarruf/offer"
	"github.com/ahmedaabouzied/tasarruf/review"
	"github.com/ahmedaabouzied/tasarruf/subscription"
	"github.com/ahmedaabouzied/tasarruf/user"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// MB is one megabyte
const MB = 1 << 20

// UserUsecase implements the User Usecase interface
type UserUsecase struct {
	UserRepository         user.Repository
	SubscriptionRepository subscription.Repository
	ReviewRepository       review.Repository
	BranchRepository       branch.Repository
	OfferRepo              offer.Repository
}

// CreateUserUsecase returns an instance of the user usecase interface
func CreateUserUsecase(uRepo user.Repository, sRepo subscription.Repository, reviewRepo review.Repository, branchRepo branch.Repository, offerRepo offer.Repository) user.Usecase {
	u := UserUsecase{
		UserRepository:         uRepo,
		SubscriptionRepository: sRepo,
		ReviewRepository:       reviewRepo,
		BranchRepository:       branchRepo,
		OfferRepo:              offerRepo,
	}
	return &u
}

// CreateCustomer creates a new user
func (c *UserUsecase) CreateCustomer(ctx context.Context, u *entities.User, newPassword string) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	profileImageURL := os.Getenv("DEFAULT_PROFILE_IMAGE_URL")
	err := u.SetHashedPassword(newPassword, u)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	err = u.SetProfileImageURL(profileImageURL, u)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	verificationCode, code, err := generateVerficationCode()
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "repository error")
	}
	err = u.SetHashedVerificationCode(verificationCode, u)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	newUser, err := c.UserRepository.CreateCustomer(ctx, u)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "repository error")
	}
	city, err := c.BranchRepository.GetCityByID(ctx, newUser.CityID)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	err = sendVerfificationMessage(newUser.Mobile, code)
	if err != nil {
		log.Error(err)
	}
	newUser.City = *city
	cancelFunc()
	return newUser, nil
}

// CreatePartner creates a new partner user
func (c *UserUsecase) CreatePartner(ctx context.Context, u *entities.User, profile *entities.PartnerProfile, newPassword string) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	profileImageURL := os.Getenv("DEFAULT_PROFILE_IMAGE_URL")
	err := u.SetHashedPassword(newPassword, u)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	u.ProfileImageURL = profileImageURL
	verificationSlice, code, err := generateVerficationCode()
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error finding user with the given ID")
	}
	// u.HashedVerificationCode = verificationSlice
	err = u.SetHashedVerificationCode(verificationSlice, u)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error setting verification code")
	}
	newUser, err := c.UserRepository.CreatePartner(ctx, u, profile)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "repository error")
	}
	err = sendVerfificationMessage(newUser.Mobile, code)
	if err != nil {
		log.Error(err)
	}
	cancelFunc()
	return newUser, nil
}

// IsEmailRegistered returns true if the given email is associated with a user
func (c *UserUsecase) IsEmailRegistered(ctx context.Context, email string) (bool, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	_, err := c.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		cancelFunc()
		return false, nil
	}
	cancelFunc()
	return true, nil

}

// IsPhoneRegistered returns true if the given phone is associated with a user
func (c *UserUsecase) IsPhoneRegistered(ctx context.Context, phone string) (bool, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	_, err := c.UserRepository.GetByPhone(ctx, phone)
	if err != nil {
		cancelFunc()
		return false, nil
	}
	cancelFunc()
	return true, nil
}

// EmailLogin performs login with Email
func (c *UserUsecase) EmailLogin(ctx context.Context, email string, password string) (*entities.User, string, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	user, err := c.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		cancelFunc()
		return nil, "", errors.New("not registered")
	}
	if !user.Active {
		cancelFunc()
		return nil, "", errors.New("Your account has been deactivated please contact us")
	}
	if !user.IsPartner() {
		city, err := c.BranchRepository.GetCityByID(ctx, user.CityID)
		if err != nil {
			cancelFunc()
			return nil, "", err
		}
		user.City = *city
	}
	otp, err := c.UserRepository.GetOTPByUser(ctx, user.ID)
	if err != nil {
		cancelFunc()
		return nil, "", errors.Wrap(err, "OTP error")
	}
	user.SetOTP(otp)
	err = user.ValidatePassword(password)
	if err != nil {
		cancelFunc()
		return nil, "", errors.Wrap(err, "Wrong Password")
	}
	token, err := entities.GenerateAuthToken(user.ID)
	if err != nil {
		cancelFunc()
		return nil, "", errors.Wrap(err, "error generating auth token")
	}
	if user.CityID != 0 {
		city, err := c.BranchRepository.GetCityByID(ctx, user.CityID)
		if err != nil {
			cancelFunc()
			return nil, "", err
		}
		user.City = *city
		if user.IsPartner() {
			user.PartnerProfile.City = *city
		}
	}
	cancelFunc()
	return user, token, nil
}

// PhoneLogin performs login with phone
func (c *UserUsecase) PhoneLogin(ctx context.Context, phone string, password string) (*entities.User, string, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	user, err := c.UserRepository.GetByPhone(ctx, phone)
	if err != nil {
		cancelFunc()
		return nil, "", errors.New("not registered")
	}
	if !user.IsPartner() {
		city, err := c.BranchRepository.GetCityByID(ctx, user.CityID)
		if err != nil {
			cancelFunc()
			return nil, "", err
		}
		user.City = *city
	}
	truePass := user.VerifyPassword(password)
	if !truePass {
		cancelFunc()
		return nil, "", errors.New("invalid password")
	}
	token, err := entities.GenerateAuthToken(user.ID)
	if err != nil {
		cancelFunc()
		return nil, "", errors.Wrap(err, "error generating auth token")
	}
	cancelFunc()
	return user, token, nil
}

// GetUser retruns the user with the given ID
func (c *UserUsecase) GetUser(ctx context.Context, ID uint) (*entities.User, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	user, err := c.UserRepository.GetByID(ctx, ID)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	city, err := c.BranchRepository.GetCityByID(ctx, user.CityID)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	user.City = *city
	if user.IsPartner() {
		user.PartnerProfile.City = *city
	}
	if user.AccountType == "user" {
		_, err := c.SubscriptionRepository.GetSubscriptionByUser(ctx, user.ID)
		if err != nil {
			cancelFunc()
			return nil, err
		}
	}
	cancelFunc()
	return user, nil
}

// DeleteUser retruns the user with the given ID
func (c *UserUsecase) DeleteUser(ctx context.Context, ID uint) (*entities.User, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	user, err := c.UserRepository.GetByID(ctx, ID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while deleting user")
	}
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	if currentUserID != user.ID {
		cancelFunc()
		return nil, errors.Wrap(err, "authorization error : current user is not allowed to delete the user with the given id")
	}
	user, err = c.UserRepository.SoftDelete(ctx, user)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	if !user.IsPartner() {
		city, err := c.BranchRepository.GetCityByID(ctx, user.CityID)
		if err != nil {
			cancelFunc()
			return nil, err
		}
		user.City = *city
	}
	cancelFunc()
	return user, nil
}

func (c *UserUsecase) UpdateProfileImage(ctx context.Context, fileHeader *multipart.FileHeader) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting current user")
	}
	if !currentUser.IsPartner() {
		city, err := c.BranchRepository.GetCityByID(ctx, currentUser.CityID)
		if err != nil {
			cancelFunc()
			return nil, err
		}
		currentUser.City = *city
	}
	// open file from file header
	file, err := fileHeader.Open()
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, errors.Wrap(err, "error processing file header")
	}
	// create a byte array of the size of the file
	fileContent := make([]byte, fileHeader.Size)
	// read the file into the byte array
	_, err = file.Read(fileContent)
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, err
	}
	defer file.Close()
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, err
	}
	u, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, err
	}
	oldProfileImageKey := u.ProfileImageKey
	// get content type of the file.
	// read the first 512 bytes of the file byte array
	// check those bytes for the content type
	contentType := http.DetectContentType(fileContent[0:513])
	if contentType == "image/png" || contentType == "image/jpeg" {
		fileName := "profileImage_" + time.Now().String() + fileHeader.Filename
		file := &filestore.File{
			FileName:    fileName,
			ContentType: contentType,
			Size:        fileHeader.Size,
			Body:        fileContent,
		}
		out, err := file.UploadToS3()
		if err != nil {
			cancelFunc()
			log.Error(err)
			return nil, err
		}
		u.SetProfileImageKey(file.FileName, currentUser)
		u.SetProfileImageURL(out.Location, currentUser)
		nu, err := c.UserRepository.UpdateCustomer(ctx, u)
		if err != nil {
			cancelFunc()
			log.Error(err)
			return nil, err
		}
		err = filestore.DeleteFromS3(oldProfileImageKey)
		if err != nil {
			log.Error(err)
		}
		cancelFunc()
		return nu, nil
	}
	cancelFunc()
	return nil, errors.New("file is not an image")
}

// UpdateTradeLicense uploads the given file to to filestore then updates the partnerProfile DB record
func (c *UserUsecase) UpdateTradeLicense(ctx context.Context, fileHeader *multipart.FileHeader) (*entities.Partner, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.getPartnerByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting partner")
	}
	// open file from file header
	file, err := fileHeader.Open()
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, errors.Wrap(err, "error processing file header")
	}
	// create a byte array of the size of the file
	fileContent := make([]byte, fileHeader.Size)
	// read the file into the byte array
	_, err = file.Read(fileContent)
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, errors.Wrap(err, "error processing file")
	}
	defer file.Close()
	oldTradeLicence := currentUser.GetPartnerProfile().LicenseKey
	// get content type of the file.
	// read the first 512 bytes of the file byte array
	// check those bytes for the content type
	contentType := http.DetectContentType(fileContent[0:513])
	if contentType == "image/png" || contentType == "image/jpeg" || contentType == "pdf" {
		fileName := "/trade_licenses_" + time.Now().String() + fileHeader.Filename
		file := &filestore.File{
			FileName:    fileName,
			ContentType: contentType,
			Size:        fileHeader.Size,
			Body:        fileContent,
		}
		out, err := file.UploadToS3()
		if err != nil {
			cancelFunc()
			log.Error(err)
			return nil, errors.Wrap(err, "error uploading file")
		}
		err = currentUser.SetLicenseKey(currentUser, fileName)
		if err != nil {
			cancelFunc()
			return nil, errors.Wrap(err, "error uploading file")
		}
		err = currentUser.SetLicenseURL(currentUser, out.Location)
		if err != nil {
			cancelFunc()
			return nil, err
		}
		_, err = c.UserRepository.UpdatePartnerProfile(ctx, currentUser.GetPartnerProfile())
		if err != nil {
			cancelFunc()
			log.Error(err)
			return nil, errors.Wrap(err, "error uploading file")
		}
		err = filestore.DeleteFromS3(oldTradeLicence)
		if err != nil {
			err := errors.Wrap(err, "error deleting old image")
			log.Error(err)
		}
		cancelFunc()
		return currentUser, nil
	}
	cancelFunc()
	return nil, errors.New("file type is not supported")
}

// AddPartnerPhoto uploads the given file to to filestore then creates the photo DB record
func (c *UserUsecase) AddPartnerPhoto(ctx context.Context, fileHeader *multipart.FileHeader) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	// open file from file header
	file, err := fileHeader.Open()
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, errors.Wrap(err, "error processing file header")
	}
	// create a byte array of the size of the file
	fileContent := make([]byte, fileHeader.Size)
	// read the file into the byte array
	_, err = file.Read(fileContent)
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, err
	}
	defer file.Close()
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, err
	}
	profile, err := c.UserRepository.GetPartnerProfileByUserID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		log.Error(err)
		return nil, err
	}
	// get content type of the file.
	// read the first 512 bytes of the file byte array
	// check those bytes for the content type
	contentType := http.DetectContentType(fileContent[0:513])
	if contentType == "image/png" || contentType == "image/jpeg" {
		fileName := "partner_photo" + time.Now().String() + fileHeader.Filename
		file := &filestore.File{
			FileName:    fileName,
			ContentType: contentType,
			Size:        fileHeader.Size,
			Body:        fileContent,
		}
		out, err := file.UploadToS3()
		if err != nil {
			cancelFunc()
			log.Error(err)
			return nil, err
		}
		photo := &entities.PartnerPhoto{
			PartnerProfileID: profile.ID,
			PhotoURL:         out.Location,
			PhotoKey:         fileName,
		}
		_, err = c.UserRepository.CreatePartnerPhoto(ctx, photo)
		if err != nil {
			cancelFunc()
			log.Error(err)
			return nil, err
		}
		nu, err := c.UserRepository.GetByID(ctx, currentUserID)
		if err != nil {
			cancelFunc()
			log.Error(err)
			return nil, err
		}
		cancelFunc()
		return nu, nil
	}
	cancelFunc()
	return nil, errors.New("file is not an image")
}

// DeletePartnerPhoto deletes the partner photo with the given ID
func (c *UserUsecase) DeletePartnerPhoto(ctx context.Context, photoID uint) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	profile, err := c.UserRepository.GetPartnerProfileByUserID(ctx, currentUserID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	if profile.PartnerID != currentUserID {
		err := errors.New("not authrorized to delete photo")
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	photo, err := c.UserRepository.DeletePartnerPhoto(ctx, photoID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	err = filestore.DeleteFromS3(photo.PhotoKey)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	ru, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return ru, nil
}

// RecoverPassword recovers password
func (c *UserUsecase) AdminRecoverPassword(ctx context.Context, mobile string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	user, err := c.UserRepository.GetByPhone(ctx, mobile)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error finding user with the given mobile")
	}
	if user == nil {
		cancelFunc()
		return errors.New("mobile phone not registered")
	}
	if !user.IsAdmin() {
		cancelFunc()
		return errors.New("mobile phone not registered")
	}
	otp, password, err := generateOTP(user.ID)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error generating OTP")
	}
	otp, err = c.UserRepository.SaveOTP(ctx, otp)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error saving OTP")
	}
	err = sendOTPMessage(user.Mobile, password)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error sending recovery SMS")
	}
	cancelFunc()
	return nil
}

// RecoverPassword recovers password
func (c *UserUsecase) RecoverPassword(ctx context.Context, mobile string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	user, err := c.UserRepository.GetByPhone(ctx, mobile)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error finding user with the given mobile")
	}
	if user == nil {
		cancelFunc()
		return errors.New("mobile phone not registered")
	}
	otp, password, err := generateOTP(user.ID)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error generating OTP")
	}
	otp, err = c.UserRepository.SaveOTP(ctx, otp)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error saving OTP")
	}
	err = sendOTPMessage(user.Mobile, password)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error sending recovery SMS")
	}
	cancelFunc()
	return nil
}

// ResendVerficationCode sends a 4-digit verification code to the user via SMS.
// Then saves a hash value for the generated code in the user DB record
func (c *UserUsecase) ResendVerficationCode(ctx context.Context) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	user, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error finding user with the given mobile")
	}
	verificationSlice, code, err := generateVerficationCode()
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error finding user with the given ID")
	}
	user.HashedVerificationCode = verificationSlice
	if user.AccountType == "partner" {
		user, err = c.UserRepository.UpdatePartner(ctx, user, &user.PartnerProfile)
		if err != nil {
			cancelFunc()
			return errors.Wrap(err, "error updating verification code hash value")
		}
	} else {
		user, err = c.UserRepository.UpdateCustomer(ctx, user)
		if err != nil {
			cancelFunc()
			return errors.Wrap(err, "error updating verification code hash value")
		}
	}
	err = sendVerfificationMessage(user.Mobile, code)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error sending verification SMS")
	}
	cancelFunc()
	return nil
}

// VerifyUser sets the user as verified
func (c *UserUsecase) VerifyUser(ctx context.Context, code string) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	user, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		err := errors.Wrap(err, "error getting current user")
		cancelFunc()
		return nil, err
	}
	trueCode := user.CheckVerificationCode(code)
	if !trueCode {
		err := errors.Wrap(err, "validation code doesn't match")
		cancelFunc()
		return nil, err
	}
	user.Verified = true
	if user.AccountType == "partner" {
		user, err = c.UserRepository.UpdatePartner(ctx, user, &user.PartnerProfile)
		if err != nil {
			cancelFunc()
			return nil, errors.Wrap(err, "error setting user as verified")
		}
	} else {
		user, err = c.UserRepository.UpdateCustomer(ctx, user)
		if err != nil {
			cancelFunc()
			return nil, errors.Wrap(err, "error setting user as verified")
		}
	}
	if !user.IsPartner() {
		city, err := c.BranchRepository.GetCityByID(ctx, user.CityID)
		if err != nil {
			cancelFunc()
			return nil, err
		}
		user.City = *city
	}
	cancelFunc()
	return user, nil
}

// UpdatePartnerProfile updates the partner profile associated with a user provided in the ctx
func (c *UserUsecase) UpdatePartnerProfile(ctx context.Context, profile *entities.PartnerProfile) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	user, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		err := errors.Wrap(err, "error getting current user")
		cancelFunc()
		return nil, err
	}
	user.PartnerProfile.DiscountValue = profile.DiscountValue
	user.PartnerProfile.MainBranchAddress = profile.MainBranchAddress
	user.PartnerProfile.Phone = profile.Phone
	user.PartnerProfile.CityID = profile.CityID
	user.CityID = profile.CityID
	user.PartnerProfile.Country = profile.Country
	user.PartnerProfile.BrandName = profile.BrandName
	user.PartnerProfile.OfferDiscription = profile.OfferDiscription
	user.PartnerProfile.CategoryID = profile.CategoryID
	_, err = c.UserRepository.UpdatePartner(ctx, user, &user.PartnerProfile)
	if err != nil {
		err := errors.Wrap(err, "error updating partner profile")
		cancelFunc()
		return nil, err
	}
	_, err = c.UserRepository.UpdatePartnerProfile(ctx, &user.PartnerProfile)
	if err != nil {
		err := errors.Wrap(err, "error updating partner profile")
		cancelFunc()
		return nil, err
	}
	city, err := c.BranchRepository.GetCityByID(ctx, user.CityID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting city")
	}
	branches, err := c.BranchRepository.GetByOwner(ctx, user.ID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting branches")
	}
	for _, branch := range branches {
		branch.CategoryID = user.PartnerProfile.CategoryID
		_, err := c.BranchRepository.Edit(ctx, &branch)
		if err != nil {
			cancelFunc()
			return nil, errors.Wrap(err, "error updating branch")
		}
	}
	user.City = *city
	user.PartnerProfile.City = *city
	cancelFunc()
	return user, nil
}

func (c *UserUsecase) ToggleIsSharable(ctx context.Context, partnerID uint) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error getting current user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("not authorized to set is sharable")
		cancelFunc()
		return err
	}
	partner, err := c.GetPartnerByID(ctx, partnerID)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error getting partner")
	}
	partner.PartnerProfile.IsSharable = !partner.PartnerProfile.IsSharable
	log.Info(partner.PartnerProfile.IsSharable)
	_, err = c.UserRepository.UpdatePartnerProfile(ctx, partner.PartnerProfile)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error updating partner profile")
	}
	cancelFunc()
	return nil
}

// UpdatePassword updates the user hashed password field
func (c *UserUsecase) UpdatePassword(ctx context.Context, newPassword string) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	user, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		err := errors.Wrap(err, "error getting current user")
		cancelFunc()
		return nil, err
	}
	err = user.SetHashedPassword(newPassword, user)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, err
	}
	var updatedUser *entities.User
	if user.IsPartner() {
		updatedUser, err = c.UserRepository.UpdatePartner(ctx, user, &user.PartnerProfile)
		if err != nil {
			err := errors.Wrap(err, "error saving password")
			cancelFunc()
			return nil, err
		}
	} else {
		updatedUser, err = c.UserRepository.UpdateCustomer(ctx, user)
		if err != nil {
			err := errors.Wrap(err, "error saving password")
			cancelFunc()
			return nil, err
		}
	}
	cancelFunc()
	return updatedUser, nil
}

func generateOTP(userID uint) (*entities.OTP, string, error) {
	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 8)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	pass := string(b)
	hpass, err := entities.EncryptPassword(pass)
	if err != nil {
		return nil, "", errors.Wrap(err, "error generating hash string")
	}
	expireTime := time.Now().Add(60 * time.Minute)
	otp := &entities.OTP{
		HashedPassword: hpass,
		UserID:         userID,
		ExpiryDate:     expireTime,
	}
	return otp, pass, nil
}

func generateVerficationCode() ([]byte, string, error) {
	letters := []rune("0123456789")
	b := make([]rune, 4)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	pass := string(b)
	hpass, err := entities.EncryptPassword(pass)
	if err != nil {
		return nil, "", errors.Wrap(err, "error generating hash string")
	}
	return hpass, pass, nil
}

func sendOTPMessage(mobile string, otp string) error {
	accountSid := os.Getenv("TWILLIO_SID")
	authToken := os.Getenv("TWILLIO_AUTH_TOKEN")
	twillioNumber := os.Getenv("TWILLIO_NUMBER")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	msgData := url.Values{}
	msgData.Set("To", mobile)
	msgData.Set("From", twillioNumber)
	msgData.Set("Body", fmt.Sprintf("Use this password to login to your account\n TASARRUF hesabınıza giriş yapmak için bu şifreyi kullanabilirsiniz\n %s", otp))
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making request to twilio API")
	}
	switch resp.StatusCode {
	case 201:
		return nil
	default:
		return errors.New(resp.Status)
	}
}

func sendVerfificationMessage(mobile string, code string) error {
	log.Info("OPT : ", code)
	accountSid := os.Getenv("TWILLIO_SID")
	authToken := os.Getenv("TWILLIO_AUTH_TOKEN")
	twillioNumber := os.Getenv("TWILLIO_NUMBER")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	msgData := url.Values{}
	msgData.Set("To", mobile)
	msgData.Set("From", twillioNumber)
	msgData.Set("Body", fmt.Sprintf("Your Tasarruf verification code : %s. \n TASARRUF üyelik doğrulama kodunuz: %s.\n", code, code))
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	log.Info(req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making request to twilio API")
	}
	log.Info(resp)
	switch resp.StatusCode {
	case 201:
		return nil
	default:
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch : reading %s : %v\n", resp.Status, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
		return errors.New(resp.Status)
	}
}

func sendYouGotVerifiedMessage(mobile string) error {
	accountSid := os.Getenv("TWILLIO_SID")
	authToken := os.Getenv("TWILLIO_AUTH_TOKEN")
	twillioNumber := os.Getenv("TWILLIO_NUMBER")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	msgData := url.Values{}
	msgData.Set("To", mobile)
	msgData.Set("From", twillioNumber)
	msgData.Set("Body", fmt.Sprintf("Your tasarruf account got verified.\n TASARRUF hesabınız doğrulandı.\n"))
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	log.Info(req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making request to twilio API")
	}
	log.Info(resp)
	switch resp.StatusCode {
	case 201:
		return nil
	default:
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch : reading %s : %v\n", resp.Status, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
		return errors.New(resp.Status)
	}

}

// ValidateCustomerPartnerIntegrity validates the partner customer integrity
func (c *UserUsecase) ValidateCustomerPartnerIntegrity(ctx context.Context, customerID uint, partnerID uint) (*entities.User, *entities.Subscription, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	if currentUserID != partnerID {
		err := errors.New("current user ID doesn't match with the given partner ID")
		cancelFunc()
		log.Error(err)
		return nil, nil, err
	}
	user, err := c.GetCustomerByID(ctx, customerID)
	if err != nil {
		err := errors.Wrap(err, "error getting customer")
		cancelFunc()
		log.Error(err)
		return nil, nil, err
	}
	subscription, err := c.SubscriptionRepository.GetSubscriptionByUser(ctx, customerID)
	if err != nil {
		err := errors.Wrap(err, "error getting subscription")
		cancelFunc()
		log.Error(err)
		return nil, nil, err
	}
	if subscription == nil {
		err := errors.New("user is not subscribed to any plan")
		cancelFunc()
		log.Error(err)
		return nil, nil, err
	}
	if subscription.HasExpirPassed() {
		err := errors.New("subscription has expired, please renew or upgrade the subscription")
		cancelFunc()
		log.Error(err)
		return nil, nil, err
	}
	partner, err := c.getPartnerByID(ctx, partnerID)
	if err != nil {
		return nil, nil, err
	}
	remainingOffers, err := c.SubscriptionRepository.GetCountOfOffersWithPartner(ctx, partner, subscription)
	if err != nil {
		err := errors.New("subscription has expired, please renew or upgrade the subscription")
		cancelFunc()
		log.Error(err)
		return nil, nil, err
	}
	subscription.RemainingOffers = remainingOffers.CountOfOffers
	plan, err := c.SubscriptionRepository.GetPlanByID(ctx, subscription.PlanID)
	if err != nil {
		err := errors.Wrap(err, "error getting plan")
		cancelFunc()
		log.Error(err)
		return nil, nil, err

	}
	subscription.Plan = *plan
	cancelFunc()
	return &user.User, subscription, nil
}

// UpdateUser saves the updates the fields of the given user
func (c *UserUsecase) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	if user.AccountType == "partner" {
		updatedUser, err := c.UserRepository.UpdatePartner(ctx, user, &user.PartnerProfile)
		if err != nil {
			log.Error(err)
			cancelFunc()
			return nil, errors.Wrap(err, "repository error while updating partner")
		}
		cancelFunc()
		return updatedUser, nil
	}
	updatedUser, err := c.UserRepository.UpdateCustomer(ctx, user)
	if err != nil {
		log.Error(err)
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while updating customer")
	}
	log.Info("Updated user city", updatedUser.CityID)
	if updatedUser.IsPartner() {
		city, err := c.BranchRepository.GetCityByID(ctx, updatedUser.CityID)
		if err != nil {
			cancelFunc()
			return nil, err
		}
		updatedUser.City = *city
	}
	cancelFunc()
	return updatedUser, nil
}

// GetCustomersCount returns the count of customer users
func (c *UserUsecase) GetCustomersCount(ctx context.Context) (int, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return 0, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return 0, err
	}
	count, err := c.UserRepository.GetCustomersCount(ctx)
	if err != nil {
		cancelFunc()
		return 0, errors.Wrap(err, "error getting count of customers")
	}
	cancelFunc()
	return count, nil
}

// GetPartnersCount returns the count of customer users
func (c *UserUsecase) GetPartnersCount(ctx context.Context) (int, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return 0, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return 0, err
	}
	count, err := c.UserRepository.GetPartnersCount(ctx)
	if err != nil {
		cancelFunc()
		return 0, errors.Wrap(err, "error getting count of customers")
	}
	cancelFunc()
	return count, nil
}

// GetAllCustomers returns all customer users
func (c *UserUsecase) GetAllCustomers(ctx context.Context) ([]entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return nil, err
	}
	customers, err := c.UserRepository.GetAllCustomers(ctx)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting count of customers")
	}
	cancelFunc()
	return customers, nil
}

// GetAllParnters returns all the partner users
func (c *UserUsecase) GetAllParnters(ctx context.Context) ([]entities.Partner, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return nil, err
	}
	partners, err := c.UserRepository.GetAllParnters(ctx)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting count of customers")
	}
	partnerProfiles := make([]entities.Partner, len(partners))
	for i, partner := range partners {
		partner, err := c.getPartnerByID(ctx, partner.ID)
		if err != nil {
			cancelFunc()
			return nil, errors.Wrap(err, "error getting partner")
		}
		partnerProfiles[i] = *partner
	}
	cancelFunc()
	return partnerProfiles, nil
}

// GetNotApprovedPartners returns a list of all not approved partner users
func (c *UserUsecase) GetNotApprovedPartners(ctx context.Context) ([]entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return nil, err
	}
	partners, err := c.UserRepository.GetNotApprovedPartners(ctx)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting count of customers")
	}
	for i, partner := range partners {
		profile, err := c.UserRepository.GetPartnerProfileByUserID(ctx, partner.ID)
		if err != nil {
			cancelFunc()
			return nil, errors.Wrap(err, "error getting profile of partner")
		}
		partners[i].PartnerProfile = *profile
	}
	cancelFunc()
	return partners, nil
}

// ApprovePartner sets the partner as approved
func (c *UserUsecase) ApprovePartner(ctx context.Context, partnerID uint) (*entities.User, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return nil, err
	}
	partner, err := c.UserRepository.GetByID(ctx, partnerID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting count of customers")
	}
	partner.PartnerProfile.Approved = true
	partnerProfile, err := c.UserRepository.UpdatePartnerProfile(ctx, &partner.PartnerProfile)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting count of customers")
	}
	SEND_APPROVAL_MSG_FEATUE := false
	if SEND_APPROVAL_MSG_FEATUE {
		err = sendYouGotVerifiedMessage(partner.Mobile)
		if err != nil {
			log.Error(err)
		}
	}
	partner.PartnerProfile = *partnerProfile
	cancelFunc()
	return partner, nil
}

// SetPartnerAsExclusive sets the partner with the given ID as exclusive
func (c *UserUsecase) SetPartnerAsExclusive(ctx context.Context, partnerID uint) (*entities.Partner, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return nil, err
	}
	partner, err := c.getPartnerByID(ctx, partnerID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting partner")
	}
	if partner == nil {
		cancelFunc()
		return nil, errors.New("partner not found")
	}
	err = c.UserRepository.CreateExclusiveRecord(ctx, partner.ID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error setting partner as exclusive")
	}
	cancelFunc()
	return partner, nil
}

// RemovePartnerAsExclusive removes the partner with the given ID from exclusive group
func (c *UserUsecase) RemovePartnerAsExclusive(ctx context.Context, partnerID uint) (*entities.Partner, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	if !currentUser.IsAdmin() {
		err := errors.New("only admins are authorized to perform this task")
		cancelFunc()
		return nil, err
	}
	partner, err := c.getPartnerByID(ctx, partnerID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting partner")
	}
	if partner == nil {
		cancelFunc()
		return nil, errors.New("partner not found")
	}
	record, err := c.UserRepository.GetExclusiveOffer(ctx, partner.ID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting exclusive record")
	}
	if record == nil {
		cancelFunc()
		return nil, errors.New("exclusive partner record not found")
	}
	_, err = c.UserRepository.DeleteExclusiveOffer(ctx, record)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error removing partner from exclusive group")
	}
	cancelFunc()
	return partner, nil
}

func (c *UserUsecase) AdminDeleteUser(ctx context.Context, userID uint) (*entities.User, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	toDeleteUser, err := c.UserRepository.GetByID(ctx, userID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while deleting user")
	}
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while deleting user")
	}
	if !currentUser.IsAdmin() {
		err = errors.New("not authorized to delete user")
		cancelFunc()
		return nil, errors.Wrap(err, "repository error while deleting user")
	}
	deletedUser, err := c.UserRepository.Delete(ctx, toDeleteUser)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	if !deletedUser.IsPartner() {
		city, err := c.BranchRepository.GetCityByID(ctx, deletedUser.CityID)
		if err != nil {
			cancelFunc()
			return nil, err
		}
		deletedUser.City = *city
	}
	cancelFunc()
	return deletedUser, nil
}

// GetExclusivePartners returns the partners set as exclusive
func (c *UserUsecase) GetExclusivePartners(ctx context.Context) ([]entities.Partner, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	exclusiveRecords, err := c.UserRepository.GetExclusiveOffers(ctx)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting exclusive partner records")
	}
	partners := make([]entities.Partner, 0)
	for _, record := range exclusiveRecords {
		partner, err := c.getPartnerByID(ctx, record.PartnerID)
		if err != nil {
			cancelFunc()
			log.Info(err)
		} else {
			partners = append(partners, *partner)
		}
	}
	cancelFunc()
	return partners, nil
}

// GetCustomerByID returns the customer with the given ID
func (c *UserUsecase) GetCustomerByID(ctx context.Context, ID uint) (*entities.Customer, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	customer, err := c.getCustomerByID(ctx, ID)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	cancelFunc()
	return customer, nil
}

// GetPartnerByID returns the partner with the given ID
func (c *UserUsecase) GetPartnerByID(ctx context.Context, ID uint) (*entities.Partner, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	partner, err := c.getPartnerByID(ctx, ID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting partner")
	}
	cancelFunc()
	return partner, nil
}

// Share creates a share record for this partner with the currently logged in customer
func (c *UserUsecase) Share(ctx context.Context) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return errors.Wrap(err, "error getting user")
	}
	if currentUser.IsPartner() || currentUser.IsAdmin() {
		cancelFunc()
		return errors.New("only customers are allowed to share partners")
	}
	sharables, err := c.UserRepository.GetSharablePartners(ctx)
	if err != nil {
		cancelFunc()
		return err
	}
	sub, err := c.SubscriptionRepository.GetSubscriptionByUser(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return err
	}
	for _, sharable := range sharables {
		partner, err := c.GetPartnerByID(ctx, sharable.ID)
		if err != nil {
			log.Error(err)
		}
		_, err = c.SubscriptionRepository.GetCountOfOffersWithPartner(ctx, partner, sub)
		if err != nil {
			log.Error(err)
		}
	}
	share := entities.Share{
		CustomerID: currentUser.ID,
	}
	err = c.UserRepository.CreateShare(ctx, &share)
	if err != nil {
		if err.Error() == "customer already shared this partner" {
			cancelFunc()
			return nil
		}
		cancelFunc()
		return errors.Wrap(err, "error creating share database record")
	}
	for _, sharable := range sharables {
		partner, err := c.GetPartnerByID(ctx, sharable.ID)
		if err != nil {
			log.Error(err)
		}
		currentRemainingOffers, err := c.SubscriptionRepository.GetCountOfOffersWithPartner(ctx, partner, sub)
		if err != nil {
			log.Error(err)
		}
		err = c.SubscriptionRepository.SetCountOfOffersWithPartner(ctx, partner, sub, currentRemainingOffers.CountOfOffers+1)
		if err != nil {
			log.Error(err)
		}
	}
	cancelFunc()
	return nil
}

// SearchUsers searches the users for the given search term.
func (c *UserUsecase) SearchUsers(ctx context.Context, searchTerm string) ([]entities.IUser, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	searchTerm = strings.ToLower(searchTerm)
	resultUsers, err := c.UserRepository.SearchUsers(ctx, searchTerm)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error searching users")
	}
	var result []entities.IUser
	for _, user := range resultUsers {
		if user.IsPartner() {
			partner, err := c.getPartnerByID(ctx, user.ID)
			if err != nil {
				log.Error("error getting partner", err)
			} else {
				result = append(result, partner)
			}
		} else {
			customer, err := c.GetCustomerByID(ctx, user.ID)
			if err != nil {
				log.Error("error getting customer", err)
			}
			result = append(result, customer)
		}
	}
	cancelFunc()
	return result, nil
}

// ToggleActiveProperty toggles the active property of a user.
func (c *UserUsecase) ToggleActiveProperty(ctx context.Context, userID uint) (entities.IUser, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	currentUserID := ctx.Value(entities.UserIDKey).(uint)
	log.Info(currentUserID, userID)
	currentUser, err := c.UserRepository.GetByID(ctx, currentUserID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	user, err := c.UserRepository.GetByID(ctx, userID)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error getting user")
	}
	err = user.ToggleActive(currentUser)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	user, err = c.UserRepository.UpdateCustomer(ctx, user)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrap(err, "error updating user")
	}
	cancelFunc()
	return user, nil
}

func (c *UserUsecase) getUserByID(ctx context.Context, ID uint) (entities.IUser, error) {
	var user entities.IUser
	user, err := c.UserRepository.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	switch user.GetAccountType() {
	case "user":
		return c.GetCustomerByID(ctx, ID)
	case "partner":
		return c.GetPartnerByID(ctx, ID)
	case "admin":
		return user, nil
	default:
		return user, nil
	}
}

func (c *UserUsecase) getCustomerByID(ctx context.Context, ID uint) (*entities.Customer, error) {
	user, err := c.UserRepository.GetByID(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "user not found")
	}
	otp, err := c.UserRepository.GetOTPByUser(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "OTP not found")
	}
	subscription, err := c.SubscriptionRepository.GetSubscriptionByUser(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "subscription not found")
	}
	reviews, err := c.ReviewRepository.GetByCustomerID(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "reviews not found")
	}
	city, err := c.BranchRepository.GetCityByID(ctx, user.CityID)
	if err != nil {
		return nil, errors.Wrap(err, "city not found")
	}
	customer := entities.Customer{
		User:         *user,
		OTP:          otp,
		Subscription: subscription,
		Reviews:      reviews,
	}
	customer.City = *city
	return &customer, nil
}

func (c *UserUsecase) getPartnerByID(ctx context.Context, ID uint) (*entities.Partner, error) {
	user, err := c.UserRepository.GetByID(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting partner")
	}
	otp, err := c.UserRepository.GetOTPByUser(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting partner OTP records")
	}
	partnerProfile, err := c.UserRepository.GetPartnerProfileByUserID(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting partner profile")
	}
	reviews, err := c.ReviewRepository.GetByPartnerID(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting partner reviews")
	}
	partner := &entities.Partner{
		User: *user,
	}
	partner.PartnerProfile = partnerProfile
	partner.Reviews = reviews
	partner.OTP = otp
	city, err := c.BranchRepository.GetCityByID(ctx, partnerProfile.CityID)
	if err != nil {
		log.Error(errors.Wrap(err, "error getting partner city"))
		partner.PartnerProfile.City = entities.City{}
		partner.User.City = entities.City{}
	} else {
		partner.PartnerProfile.City = *city
		partner.User.City = *city
	}
	return partner, nil
}
