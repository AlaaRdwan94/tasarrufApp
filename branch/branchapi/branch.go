package branchapi

import (
	"context"
	"errors"
	"github.com/ahmedaabouzied/tasarruf/branch"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
	"net/http"
	"strconv"
)

// BranchAPI is the handler for branch related API endpoints
type BranchAPI struct {
	BranchUsecase branch.Usecase
}

// newBranchRequest
type newBranchRequest struct {
	Country    string `json:"country"`
	CityID     uint   `json:"cityID"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Mobile     string `json:"mobile"`
	CategoryID uint   `json:"categoryID"`
}

type newCategoryRequest struct {
	EnglishName string `json:"englishName"`
	TurkishName string `json:"turkishName"`
}

type newCityRequest newCategoryRequest

// Validate validates the new user request
func (req *newBranchRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.CityID, validation.Required),
		validation.Field(&req.Country, validation.Required, validation.Length(2, 50), is.LowerCase),
		validation.Field(&req.Phone, validation.Required),
		validation.Field(&req.Mobile, validation.Required),
		validation.Field(&req.Address, validation.Required),
	)
}

// Validate validates the new user request
func (req *newCategoryRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.TurkishName, validation.Required),
		validation.Field(&req.EnglishName, validation.Required),
	)
}

// Validate validates the new user request
func (req *newCityRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.TurkishName, validation.Required),
		validation.Field(&req.EnglishName, validation.Required),
	)
}

// CreateBranchAPI returns a new branch API instance
func CreateBranchAPI(u branch.Usecase) BranchAPI {
	api := BranchAPI{
		BranchUsecase: u,
	}
	return api
}

// CreateBranch handles requests for creating new branches
func (h *BranchAPI) CreateBranch(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req newBranchRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been a trouble sending your information to the server, please try again", err)
		return
	}
	err = req.Validate()
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	newBranch := &entities.Branch{
		Country:    req.Country,
		CityID:     req.CityID,
		Address:    req.Address,
		Phone:      req.Phone,
		Mobile:     req.Mobile,
		CategoryID: req.CategoryID,
	}
	newBranch, err = h.BranchUsecase.Create(ctx, newBranch)
	if err != nil {
		entities.SendValidationError(c, "only partner users can create branches", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "branch created successfully",
		"branch":  newBranch,
	})
	return
}

// GetBranch handles the get by id endpoint
func (h *BranchAPI) GetBranch(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	branchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	branch, err := h.BranchUsecase.GetByID(ctx, uint(branchID))
	if err != nil {
		entities.SendNotFoundError(c, "Sorry, we couldn't find the branch you're looking for", err)
		return
	}
	if branch == nil {
		entities.SendNotFoundError(c, "Sorry, we couldn't find the branch you're looking for", errors.New("branch not found"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"branch": branch,
	})
	return
}

// DeleteBranch handles the delete by id endpoint
func (h *BranchAPI) DeleteBranch(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	branchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	branchToDelete, err := h.BranchUsecase.GetByID(ctx, uint(branchID))
	if err != nil {
		entities.SendNotFoundError(c, "Sorry, we couldn't find the branch you're trying to delete", err)
		return
	}
	if branchToDelete == nil {
		entities.SendNotFoundError(c, "Sorry, we couldn't find the branch you're trying to delete", errors.New("branch not found"))
		return
	}
	deletedBranch, err := h.BranchUsecase.Delete(ctx, branchToDelete)
	if err != nil {
		entities.SendValidationError(c, "You're not the owner of this branch", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Branch deleted successfully",
		"branch":  deletedBranch,
	})
	return
}

// EditBranch handles the PUT by id endpoint
func (h *BranchAPI) EditBranch(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	branchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	var req newBranchRequest
	err = c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been a trouble sending your information to the server, please try again", err)
		return
	}
	err = req.Validate()
	if err != nil {
		entities.SendValidationError(c, err.Error(), err)
		return
	}
	branchToEdit, err := h.BranchUsecase.GetByID(ctx, uint(branchID))
	if err != nil {
		entities.SendNotFoundError(c, "Sorry, we couldn't find the branch you're trying to edit", err)
		return
	}
	if branchToEdit == nil {
		entities.SendNotFoundError(c, "Sorry, we couldn't find the branch you're trying to edit", errors.New("branch not found"))
		return
	}
	branchToEdit.CityID = req.CityID
	branchToEdit.Country = req.Country
	branchToEdit.Address = req.Address
	branchToEdit.Phone = req.Phone
	branchToEdit.Mobile = req.Mobile
	editedBranch, err := h.BranchUsecase.Edit(ctx, branchToEdit)
	if err != nil {
		entities.SendValidationError(c, "You're not the owner of this branch", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "Branch updated successfully",
		"branch":  editedBranch,
	})
	return
}

// GetMyBranches handles GET my-branches enpoint
func (h *BranchAPI) GetMyBranches(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	branches, err := h.BranchUsecase.GetByOwner(ctx, userID)
	if err != nil {
		entities.SendServerError(c, "There has been an error while getting your branches. Please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"branches": branches,
	})
	return
}

// GetBranchesOfOwner handles GET owner-branches endpoint
func (h *BranchAPI) GetBranchesOfOwner(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	ownerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	branches, err := h.BranchUsecase.GetByOwner(ctx, uint(ownerID))
	if err != nil {
		entities.SendServerError(c, "There has been an error while getting your branches. Please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"branches": branches,
	})
	return
}

// GetBranchesByLocation handles GET location-branches endpoint
func (h *BranchAPI) GetBranchesByLocation(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	country := c.Query("country")
	city := c.Query("city")
	cityID, err := strconv.ParseUint(city, 10, 64)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while parsing your data. Please try again", err)
	}
	branches, err := h.BranchUsecase.GetByLocation(ctx, country, uint(cityID))
	if err != nil {
		entities.SendServerError(c, "There has been an error while getting your search results. Please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"branches": branches,
	})
	return
}

// CreateCategory handles POST request to /category endpoint
func (h *BranchAPI) CreateCategory(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req newCategoryRequest
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
	newCategory := &entities.Category{
		TurkishName: req.TurkishName,
		EnglishName: req.EnglishName,
	}
	category, err := h.BranchUsecase.CreateCategory(ctx, newCategory)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while creating category , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  "category created successfully",
		"category": category,
	})
	return
}

// EditCategory handles PUT request to /category/:id endpoint
func (h *BranchAPI) EditCategory(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	var req newCategoryRequest
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
	toEditCategory := &entities.Category{
		TurkishName: req.TurkishName,
		EnglishName: req.EnglishName,
	}
	toEditCategory.ID = uint(categoryID)
	category, err := h.BranchUsecase.EditCategory(ctx, toEditCategory)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while creating category , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  "category edited successfully",
		"category": category,
	})
	return
}

// DeleteCategory handles DELTE /category/:id endpoint
func (h *BranchAPI) DeleteCategory(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	category, err := h.BranchUsecase.DeleteCategory(ctx, uint(categoryID))
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your information , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  "category deleted successfully",
		"category": category,
	})
	return
}

// GetByCategory returns the branches of given category
func (h *BranchAPI) GetByCategory(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	branches, err := h.BranchUsecase.GetByCategory(ctx, uint(categoryID))
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"branches": branches,
	})
	return
}

// GetCategories handles GET /category endpoint
func (h *BranchAPI) GetCategories(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	categories, err := h.BranchUsecase.GetCategories(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your request , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

// CreateCity handles POST /city endpoint
func (h *BranchAPI) CreateCity(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req newCityRequest
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
	newCity := &entities.City{
		TurkishName: req.TurkishName,
		EnglishName: req.EnglishName,
	}
	city, err := h.BranchUsecase.CreateCity(ctx, newCity)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while creating category , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "city created successfully",
		"city":    city,
	})
	return

}

// DeleteCity handles DELTE /city/:id endpoint
func (h *BranchAPI) DeleteCity(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	cityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	city, err := h.BranchUsecase.DeleteCity(ctx, uint(cityID))
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your information , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "city deleted successfully",
		"city":    city,
	})
	return
}

// UpdateCity handles PUT /city/:id endpoint
func (h *BranchAPI) UpdateCity(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	cityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	var req newCityRequest
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
	newCity := &entities.City{
		TurkishName: req.TurkishName,
		EnglishName: req.EnglishName,
	}
	newCity.ID = uint(cityID)
	city, err := h.BranchUsecase.UpdateCity(ctx, newCity)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your information , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "city updated successfully",
		"city":    city,
	})
	return
}

// GetCityByID handles Get /city/:id endpoint
func (h *BranchAPI) GetCityByID(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	cityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while sending your information to the server, please try again", err)
		return
	}
	city, err := h.BranchUsecase.GetCityByID(ctx, uint(cityID))
	if err != nil {
		entities.SendNotFoundError(c, "There has been an error while processing your information , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"city": city,
	})
	return
}

// GetAllCities handles get /city endpoint
func (h *BranchAPI) GetAllCities(c *gin.Context) {
	ctx := context.Background()
	cities, err := h.BranchUsecase.GetAllCities(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while processing your information , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"cities": cities,
	})
	return
}

// Search handles GET /branch endpoint
func (h *BranchAPI) Search(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	city := c.Query("cityID")
	category := c.Query("categoryID")
	name := c.Query("brandName")
	cityID, err := strconv.ParseInt(city, 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing the city ID, please try again", err)
		return
	}
	categoryID, err := strconv.ParseInt(category, 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing the category ID, please try again", err)
		return
	}
	branches, err := h.BranchUsecase.SearchBranches(ctx, uint(cityID), uint(categoryID), name)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while getting your data , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"branches": branches,
	})
}
