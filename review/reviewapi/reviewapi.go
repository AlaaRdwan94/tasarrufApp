package reviewapi

import (
	"context"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/review"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v3"
	"net/http"
	"strconv"
)

// ReviewAPI is the API handler for review related API endpoint
type ReviewAPI struct {
	ReviewUsecase review.Usecase
}

// newReviewRequest represents the review body
type newReviewRequest struct {
	CustomerID uint   `json:"customerID"`
	PartnerID  uint   `json:"partnerID"`
	Stars      int    `json:"stars"`
	Content    string `json:"content"`
}

func (req *newReviewRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.CustomerID, validation.Required),
		validation.Field(&req.PartnerID, validation.Required),
		validation.Field(&req.Stars, validation.Required, validation.Max(5), validation.Min(1)),
	)
}

// CreateReviewAPI returns a new review API instance
func CreateReviewAPI(u review.Usecase) ReviewAPI {
	api := ReviewAPI{
		ReviewUsecase: u,
	}
	return api
}

// CreateReview handles POST /review endpoint
func (h *ReviewAPI) CreateReview(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req newReviewRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	err = req.Validate()
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	newReview := &entities.Review{
		CustomerID: req.CustomerID,
		PartnerID:  req.PartnerID,
		Stars:      req.Stars,
		Content:    req.Content,
	}
	newReview, err = h.ReviewUsecase.Create(ctx, newReview)
	if err != nil {
		entities.SendValidationError(c, "only customer users can create reviews", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "review created successfully",
		"review":  newReview,
	})
	return
}

// UpdateReview handles PUT /review/:id endpoint
func (h *ReviewAPI) UpdateReview(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	reviewID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	var req newReviewRequest
	err = c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	err = req.Validate()
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	review, err := h.ReviewUsecase.GetByID(ctx, uint(reviewID))
	if err != nil {
		entities.SendNotFoundError(c, "Review not found", err)
		return
	}
	review.Content = req.Content
	review.Stars = req.Stars
	review, err = h.ReviewUsecase.Update(ctx, review)
	if err != nil {
		entities.SendValidationError(c, "You are not the owner of this review", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "review updated successfully",
		"review":  review,
	})
	return
}

// DeleteReview handles DELETE /review/:id endpoint
func (h *ReviewAPI) DeleteReview(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	reviewID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	review, err := h.ReviewUsecase.GetByID(ctx, uint(reviewID))
	if err != nil {
		entities.SendNotFoundError(c, "Review not found", err)
		return
	}
	review, err = h.ReviewUsecase.Delete(ctx, review)
	if err != nil {
		entities.SendValidationError(c, "You are not the owner of this review", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "review deleted successfully",
		"review":  review,
	})
	return
}

// GetByID handles GET /review/:id endpoint
func (h *ReviewAPI) GetByID(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	reviewID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	review, err := h.ReviewUsecase.GetByID(ctx, uint(reviewID))
	if err != nil {
		entities.SendNotFoundError(c, "Review not found", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"review": review,
	})
}

// GetMyReviews handles GET /reviews/:id endpoint
func (h *ReviewAPI) GetMyReviews(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	reviews, err := h.ReviewUsecase.GetByCustomerID(ctx, userID)
	if err != nil {
		entities.SendNotFoundError(c, "Reviews not found", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
	})
}

// GetPartnerReviews handles GET /reviews/:id endpoint
func (h *ReviewAPI) GetPartnerReviews(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	partnerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	reviews, err := h.ReviewUsecase.GetByPartnerID(ctx, uint(partnerID))
	if err != nil {
		entities.SendNotFoundError(c, "Reviews not found", err)
		return
	}
	average, err := h.ReviewUsecase.GetAverageRatings(ctx, uint(partnerID))
	if err != nil {
		entities.SendNotFoundError(c, "Reviews not found", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"average": average,
	})
}
