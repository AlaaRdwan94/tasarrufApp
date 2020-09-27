package supportapi

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/support"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
)

// SupportAPI defines the api handler for support routes
type SupportAPI struct {
	SupportUsecase support.Usecase
}

type supportRequest struct {
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

// CreateSupportAPI creates a new support API instance
func CreateSupportAPI(u support.Usecase) SupportAPI {
	api := SupportAPI{
		SupportUsecase: u,
	}
	return api
}

// Validate validates the new user request
func (req *supportRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Mobile, validation.Required),
		validation.Field(&req.Email, validation.Required, is.Email),
	)
}

// CreateSupportRecord handles POST /support endpoint
func (h *SupportAPI) CreateSupportRecord(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req supportRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information, please try again", err)
		return
	}
	err = req.Validate()
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	info := entities.SupportInfo{
		Email:  req.Email,
		Mobile: req.Mobile,
	}
	createdInfo, err := h.SupportUsecase.Create(ctx, &info)
	if err != nil {
		entities.SendValidationError(c, "You are not authrorized to creates support info records", err)
		return
	}
	c.JSON(200, gin.H{
		"messages": "Support Info Created Successfully",
		"info":     createdInfo,
	})
	return
}

// UpdateSupportRecord handles PUT /support endpoint
func (h *SupportAPI) UpdateSupportRecord(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req supportRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information, please try again", err)
		return
	}
	err = req.Validate()
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	info := entities.SupportInfo{
		Email:  req.Email,
		Mobile: req.Mobile,
	}
	updatedInfo, err := h.SupportUsecase.Update(ctx, &info)
	if err != nil {
		entities.SendValidationError(c, "You are not authrorized to update support info records", err)
		return
	}
	c.JSON(200, gin.H{
		"messages": "Support Info Updated Successfully",
		"info":     updatedInfo,
	})
	return
}

// GetSupportInfo handles GET /support endpoint
func (h *SupportAPI) GetSupportInfo(c *gin.Context) {
	ctx := context.Background()
	info, err := h.SupportUsecase.GetSupportInfo(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while getting support contact information, please try again", err)
		return
	}
	c.JSON(200, gin.H{
		"info": info,
	})
	return
}
