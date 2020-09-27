package subscriptionapi

import (
	"context"
	"github.com/ahmedaabouzied/iyzipay-go/iyzipay"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/subscription"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
	"net/http"
	"strconv"
)

// SubscriptionAPI defines the API handlers for subscription routes
type SubscriptionAPI struct {
	SubscriptionUsecase subscription.Usecase
}

type newPlanRequest struct {
	EnglishName        string  `json:"englishName"`
	EnglishDescription string  `json:"englishDescription"`
	TurkishName        string  `json:"turkishName"`
	TurkishDescription string  `json:"turkishDescription"`
	Price              float64 `json:"price,float64"`
	CountOfOffers      uint    `json:"countOfOffers,uint"`
	Image              string  `json:"image"`
	IsDefault          bool    `json:"isDefault"`
}

type paymentRequest struct {
	CardNumber     string `json:"cardNumber"`
	ExpireMonth    string `json:"expireMonth"`
	ExpireYear     string `json:"expireYear"`
	Cvc            string `json:"cvc"`
	IDNumber       string `json:"idNumber"`
	CardHolderName string `json:"cardHolderName"`
}

// CreateSubscriptionAPI returns a new API instance
func CreateSubscriptionAPI(u subscription.Usecase) SubscriptionAPI {
	api := SubscriptionAPI{
		SubscriptionUsecase: u,
	}
	return api
}

// Validate method for the newPlanRequest body
func (req *newPlanRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.EnglishName, validation.Required, validation.Length(2, 50)),
		validation.Field(&req.EnglishDescription, validation.Required),
		validation.Field(&req.TurkishName, validation.Required),
		validation.Field(&req.TurkishDescription, validation.Required),
		validation.Field(&req.Image, validation.Required, is.URL),
		// validation.Field(&req.Price, validation.Required),
		// validation.Field(&req.CountOfOffers, validation.Required),
	)
}

// Validate method for the paymentRequest body
func (req *paymentRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.CardNumber, validation.Required),
		validation.Field(&req.ExpireMonth, validation.Required),
		validation.Field(&req.ExpireYear, validation.Required),
		validation.Field(&req.Cvc, validation.Required),
		validation.Field(&req.CardHolderName, validation.Required),
		validation.Field(&req.IDNumber, validation.Required),
	)
}

// CreatePlan handles requests for creating new subscription plan
func (h *SubscriptionAPI) CreatePlan(c *gin.Context) {
	ctx := context.Background()
	var currentUserID = c.MustGet("userID").(uint)
	var req newPlanRequest
	ctx = context.WithValue(ctx, entities.UserIDKey, currentUserID)
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	err = req.Validate()
	if err != nil {
		entities.SendServerError(c, err.Error(), err)
		return
	}
	newPlan := &entities.Plan{
		TurkishName:        req.TurkishName,
		TurkishDescription: req.TurkishDescription,
		EnglishName:        req.EnglishName,
		EnglishDescription: req.EnglishDescription,
		Price:              req.Price,
		Image:              req.Image,
		CountOfOffers:      req.CountOfOffers,
		IsDefault:          req.IsDefault,
	}
	newPlan, err = h.SubscriptionUsecase.CreatePlan(ctx, newPlan)
	if err != nil {
		entities.SendValidationError(c, "only admin users can create subscription plans", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "plan created successfully",
		"plan":    newPlan,
	})
}

// GetPlans returns all the available subscription plans
func (h *SubscriptionAPI) GetPlans(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	plans, err := h.SubscriptionUsecase.GetAllPlans(ctx)
	if err != nil {
		entities.SendNotFoundError(c, "Sorry there has been an error while getting plans, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"plans": plans,
	})
}

// RankPlanUp ranks the plan with the given ID up
func (h *SubscriptionAPI) RankPlanUp(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	err = h.SubscriptionUsecase.RankPlanUp(ctx, uint(planID))
	if err != nil {
		if err.Error() == "only admin users can delete subscription plans" {
			entities.SendAuthError(c, "only admin users can delete plans", err)
			return
		}
		entities.SendNotFoundError(c, "there has been an error while deleting plan, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "plan ranked up successfully",
	})

}

// DeletePlan handles DELETE requests to plan endpoint
func (h *SubscriptionAPI) DeletePlan(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	deletedPlan, err := h.SubscriptionUsecase.DeletePlan(ctx, uint(planID))
	if err != nil {
		if err.Error() == "only admin users can delete subscription plans" {
			entities.SendAuthError(c, "only admin users can delete plans", err)
			return
		}
		entities.SendNotFoundError(c, "there has been an error while deleting plan, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "plan deleted successfully",
		"plan":    deletedPlan,
	})

}

// SubscribeToPlan handles POST requests to the subscribe endpoint
func (h *SubscriptionAPI) SubscribeToPlan(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	plan, err := h.SubscriptionUsecase.GetPlanByID(ctx, uint(planID))
	if err != nil {
		entities.SendNotFoundError(c, "plan not found", err)
	}
	if plan.Price > 0 {
		var req paymentRequest
		err = c.BindJSON(&req)
		if err != nil {
			entities.SendValidationError(c, "there has been an error while sending your information to the server , please try again", err)
			return
		}
		err = req.Validate()
		if err != nil {
			entities.SendValidationError(c, err.Error(), err)
			return
		}
		paymentDetails := &entities.PaymentDetails{
			IDNumber: req.IDNumber,
			Card: &iyzipay.PaymentCard{
				CardHolderName: req.CardHolderName,
				CardNumber:     req.CardNumber,
				ExpireYear:     req.ExpireYear,
				ExpireMonth:    req.ExpireMonth,
				Cvc:            req.Cvc,
			},
		}
		subscription, err := h.SubscriptionUsecase.SubscribeToPlan(ctx, uint(planID), paymentDetails)
		if err != nil {
			entities.SendValidationError(c, err.Error(), err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success":      "subscribed suscessfully",
			"subscription": subscription,
		})
		return
	}
	subscription, err := h.SubscriptionUsecase.SubscribeToFreePlan(ctx, uint(planID))
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      "subscribed suscessfully",
		"subscription": subscription,
	})
	return
}

// UpgradePlan handles POST requests to the upgrade endpoint
func (h *SubscriptionAPI) UpgradePlan(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	planID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	var req paymentRequest
	err = c.BindJSON(&req)
	if err != nil {
		entities.SendValidationError(c, "there has been an error while sending your information to the server , please try again", err)
		return
	}
	err = req.Validate()
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	paymentDetails := &entities.PaymentDetails{
		IDNumber: req.IDNumber,
		Card: &iyzipay.PaymentCard{
			CardHolderName: req.CardHolderName,
			CardNumber:     req.CardNumber,
			ExpireYear:     req.ExpireYear,
			ExpireMonth:    req.ExpireMonth,
			Cvc:            req.Cvc,
		},
	}
	subscription, err := h.SubscriptionUsecase.UpgradePlan(ctx, uint(planID), paymentDetails)
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      "upgraded suscessfully",
		"subscription": subscription,
	})
}

// AdminUpgradeUserPlan handles POST requests to the upgrade endpoint
func (h *SubscriptionAPI) AdminUpgradeUserPlan(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	planID, err := strconv.ParseInt(c.Query("planID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	customerID, err := strconv.ParseInt(c.Query("userID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	subscription, err := h.SubscriptionUsecase.AdminUpgradeUserPlan(ctx, uint(customerID), uint(planID))
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      "upgraded suscessfully",
		"subscription": subscription,
	})
}

// RenewPlan handles POST requests to the renew endpoint
func (h *SubscriptionAPI) RenewPlan(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req paymentRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendValidationError(c, "there has been an error while sending your information to the server , please try again", err)
		return
	}
	err = req.Validate()
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	paymentDetails := &entities.PaymentDetails{
		IDNumber: req.IDNumber,
		Card: &iyzipay.PaymentCard{
			CardHolderName: req.CardHolderName,
			CardNumber:     req.CardNumber,
			ExpireYear:     req.ExpireYear,
			ExpireMonth:    req.ExpireMonth,
			Cvc:            req.Cvc,
		},
	}
	subscription, err := h.SubscriptionUsecase.RenewPlan(ctx, paymentDetails)
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      "renewed suscessfully",
		"subscription": subscription,
	})
}

// GetMySubscription handles GET requests to the subscription endpoint
func (h *SubscriptionAPI) GetMySubscription(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	subscription, err := h.SubscriptionUsecase.GetMySubscription(ctx)
	if err != nil {
		entities.SendValidationError(c, "user is already subscribed to a plan", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"subscription": subscription,
	})
}

// GetMySubscriptionWithPartner handled GET request to /subscription/partner/:id enpoint
func (h *SubscriptionAPI) GetMySubscriptionWithPartner(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partnerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	subscription, err := h.SubscriptionUsecase.GetMySubscriptionWithPartner(ctx, uint(partnerID))
	if err != nil {
		entities.SendValidationError(c, "user is already subscribed to a plan", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"subscription": subscription,
	})
}

// UpdatePlan handles PUT /plans/:planID endpoint
func (h *SubscriptionAPI) UpdatePlan(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	planID, err := strconv.ParseInt(c.Param("planID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	var req newPlanRequest
	err = c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "there has been an error while parsing your request", err)
		return
	}
	plan := entities.Plan{
		EnglishName:        req.EnglishName,
		TurkishName:        req.TurkishName,
		EnglishDescription: req.EnglishDescription,
		TurkishDescription: req.TurkishDescription,
		CountOfOffers:      req.CountOfOffers,
		Price:              req.Price,
		Image:              req.Image,
		IsDefault:          req.IsDefault,
	}
	updatedPlan, err := h.SubscriptionUsecase.UpdatePlan(ctx, uint(planID), &plan)
	if err != nil {
		entities.SendValidationError(c, "user is not authorized to update plan", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "plan updated successfully",
		"plan":    updatedPlan,
	})
	return
}

func (h *SubscriptionAPI) CreatePlanCategoryAssociation(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	planID, err := strconv.ParseInt(c.Query("planID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	categoryID, err := strconv.ParseInt(c.Query("categoryID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	err = h.SubscriptionUsecase.CreatePlanCategoryAssociation(ctx, uint(planID), uint(categoryID))
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"successs": "Association created",
	})
}

func (h *SubscriptionAPI) RemovePlanCategoryAssociation(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	planID, err := strconv.ParseInt(c.Query("planID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	categoryID, err := strconv.ParseInt(c.Query("categoryID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	err = h.SubscriptionUsecase.RemovePlanCategoryAssociation(ctx, uint(planID), uint(categoryID))
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"successs": "Association created",
	})
}

func (h *SubscriptionAPI) GetCategoriesOfPlan(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	planID, err := strconv.ParseInt(c.Query("planID"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	categories, err := h.SubscriptionUsecase.GetCategoriesOfPlan(ctx, uint(planID))
	if err != nil {
		entities.SendParsingError(c, "there has been an error parsing your request", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
