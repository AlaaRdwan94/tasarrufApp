package main

import (
	branchapi "github.com/ahmedaabouzied/tasarruf/branch/branchapi"
	_branchrepo "github.com/ahmedaabouzied/tasarruf/branch/repository"
	_branchusecase "github.com/ahmedaabouzied/tasarruf/branch/usecase"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/offer/hub"
	offerapi "github.com/ahmedaabouzied/tasarruf/offer/offerapi"
	_offerrepo "github.com/ahmedaabouzied/tasarruf/offer/repository"
	_offerusecase "github.com/ahmedaabouzied/tasarruf/offer/usecase"
	_reviewrepo "github.com/ahmedaabouzied/tasarruf/review/repository"
	reviewapi "github.com/ahmedaabouzied/tasarruf/review/reviewapi"
	_reviewusecase "github.com/ahmedaabouzied/tasarruf/review/usecase"
	_subscriptionrepo "github.com/ahmedaabouzied/tasarruf/subscription/repository"
	subscriptionapi "github.com/ahmedaabouzied/tasarruf/subscription/subscriptionapi"
	_subscriptionusecase "github.com/ahmedaabouzied/tasarruf/subscription/usecase"
	_supportrepo "github.com/ahmedaabouzied/tasarruf/support/repository"
	supportapi "github.com/ahmedaabouzied/tasarruf/support/supportapi"
	_supportusecase "github.com/ahmedaabouzied/tasarruf/support/usecase"
	_userrepo "github.com/ahmedaabouzied/tasarruf/user/repository"
	_userusecase "github.com/ahmedaabouzied/tasarruf/user/usecase"
	userapi "github.com/ahmedaabouzied/tasarruf/user/userapi"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// InitializeRoutes defines server routes
func InitializeRoutes(db *gorm.DB, router *gin.Engine) {
	hub := hub.CreateUserHub()
	userRepo := _userrepo.CreateUserRepository(db)
	branchRepo := _branchrepo.CreateBranchRepository(db)
	subscriptionRepo := _subscriptionrepo.CreateSubscriptionRepository(db)
	offerRepo := _offerrepo.CreateOfferRepository(db)
	reviewRepo := _reviewrepo.CreateReviewRepository(db)
	supportRepo := _supportrepo.CreateSupportRepository(db)
	userUsecase := _userusecase.CreateUserUsecase(userRepo, subscriptionRepo, reviewRepo, branchRepo, offerRepo)
	branchUsecase := _branchusecase.CreateBranchUsecase(branchRepo, userRepo, subscriptionRepo)
	subscriptionUsecase := _subscriptionusecase.CreateSubscriptionUsecase(subscriptionRepo, userRepo, branchRepo, offerRepo)
	offerUsecase := _offerusecase.CreateOfferUsecase(offerRepo, hub, userRepo, branchRepo, subscriptionRepo)
	reviewUsecase := _reviewusecase.CreateReviewUsecase(reviewRepo, userRepo, branchRepo)
	supportUsecase := _supportusecase.CreateSupportUsecase(supportRepo, userRepo)
	userHandler := userapi.CreateUserAPI(userUsecase)
	branchHandler := branchapi.CreateBranchAPI(branchUsecase)
	subscriptionHandler := subscriptionapi.CreateSubscriptionAPI(subscriptionUsecase)
	offerHandler := offerapi.CreateOfferHandler(offerUsecase, hub, branchUsecase, userUsecase)
	reviewHandler := reviewapi.CreateReviewAPI(reviewUsecase)
	supportHandler := supportapi.CreateSupportAPI(supportUsecase)
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowWebSockets = true
	config.AllowHeaders = []string{"Token", "Content-Type"}
	config.AllowCredentials = true
	config.AllowBrowserExtensions = true
	router.Use(cors.New(config))
	publicRoutes := router.Group("/api/v1/public")
	{
		publicRoutes.POST("/user", userHandler.CreateUser)
		publicRoutes.POST("/email-login", userHandler.EmailLogin)
		publicRoutes.POST("/phone-login", userHandler.PhoneLogin)
		publicRoutes.GET("/is-email-registered", userHandler.IsPhoneRegistered)
		publicRoutes.GET("/is-phone-registered", userHandler.IsEmailRegistered)
		publicRoutes.POST("/forget-password", userHandler.ForgetPassword)
		publicRoutes.POST("/admin/forget-password", userHandler.AdminForgetPassword)
		publicRoutes.GET("/cities", branchHandler.GetAllCities)
		publicRoutes.GET("support-info", supportHandler.GetSupportInfo)
	}
	router.GET("/api/v1/connect", offerHandler.Connect)
	authorizedRoutes := router.Group("/api/v1")
	authorizedRoutes.Use(authUser())
	{
		authorizedRoutes.GET("my-branches", branchHandler.GetMyBranches)
		authorizedRoutes.GET("branches-by-owner/:id", branchHandler.GetBranchesOfOwner)
		authorizedRoutes.GET("branches-by-location", branchHandler.GetBranchesByLocation)
		authorizedRoutes.GET("branches-by-category/:id", branchHandler.GetByCategory)
		userRoutes := authorizedRoutes.Group("/user")
		{
			userRoutes.GET("", userHandler.GetUser)
			userRoutes.PUT("", userHandler.UpdateUser)
			userRoutes.DELETE("", userHandler.DeleteUser)
			userRoutes.PUT("/profile-image", userHandler.UpdateProfileImage)
			userRoutes.PUT("/trade-license", userHandler.UpdateTradeLicense)
			userRoutes.PUT("/partner-photo", userHandler.AddPartnerPhoto)
			userRoutes.DELETE("/partner-photo/:photoID", userHandler.DeletePartnerPhoto)
			userRoutes.POST("/verify", userHandler.VerifyUser)
			userRoutes.POST("/resend-verification-code", userHandler.ResendVerificationCode)
			userRoutes.PUT("main-branch", userHandler.UpdatePartnerProfile)
			userRoutes.POST("/update-pass", userHandler.UpdatePassword)
			userRoutes.GET("/validate-offer", userHandler.ValidateCustomer)
			userRoutes.POST("/share", userHandler.Share)
		}
		branchRoutes := authorizedRoutes.Group("/branch")
		{
			branchRoutes.GET("", branchHandler.Search)
			branchRoutes.POST("", branchHandler.CreateBranch)
			branchRoutes.GET("/:id", branchHandler.GetBranch)
			branchRoutes.DELETE("/:id", branchHandler.DeleteBranch)
			branchRoutes.PUT("/:id", branchHandler.EditBranch)
		}
		planHanRoutes := authorizedRoutes.Group("/plans")
		{
			planHanRoutes.POST("", subscriptionHandler.CreatePlan)
			planHanRoutes.GET("", subscriptionHandler.GetPlans)
			planHanRoutes.DELETE("/:id", subscriptionHandler.DeletePlan)
			planHanRoutes.PUT("/:planID", subscriptionHandler.UpdatePlan)
			planHanRoutes.POST("/:planID", subscriptionHandler.RankPlanUp)
		}
		subscriptionRoutes := authorizedRoutes.Group("/subscription")
		{
			subscriptionRoutes.GET("", subscriptionHandler.GetMySubscription)
			subscriptionRoutes.GET("/partner/:id", subscriptionHandler.GetMySubscriptionWithPartner)
			subscriptionRoutes.POST("/subscribe/:id", subscriptionHandler.SubscribeToPlan)
			subscriptionRoutes.POST("/renew", subscriptionHandler.RenewPlan)
			subscriptionRoutes.POST("/upgrade/:id", subscriptionHandler.UpgradePlan)
		}
		offersRoutes := authorizedRoutes.Group("/offer")
		{
			offersRoutes.POST("", offerHandler.ConsumeOffer)
			offersRoutes.POST("/history", offerHandler.GetMyOffersHistory)
			offersRoutes.POST("/history/mail", offerHandler.SendOffersStaticMail)
			offersRoutes.GET("/:id", offerHandler.GetOffer)
		}
		reviewRoutes := authorizedRoutes.Group("/review")
		{
			reviewRoutes.POST("", reviewHandler.CreateReview)
			reviewRoutes.GET("/:id", reviewHandler.GetByID)
			reviewRoutes.DELETE("/:id", reviewHandler.DeleteReview)
			reviewRoutes.PUT("/:id", reviewHandler.UpdateReview)
		}
		reviewsRoutes := authorizedRoutes.Group("/reviews")
		{
			reviewsRoutes.GET("", reviewHandler.GetMyReviews)
			reviewsRoutes.GET("/:id", reviewHandler.GetPartnerReviews)
		}
		categoryRoutes := authorizedRoutes.Group("/category")
		{
			categoryRoutes.POST("", branchHandler.CreateCategory)
			categoryRoutes.GET("", branchHandler.GetCategories)
			categoryRoutes.DELETE("/:id", branchHandler.DeleteCategory)
			categoryRoutes.PUT("/:id", branchHandler.EditCategory)
		}
		cityRoutes := authorizedRoutes.Group("/city")
		{
			cityRoutes.GET("/:id", branchHandler.GetCityByID)
			cityRoutes.POST("", branchHandler.CreateCity)
			cityRoutes.DELETE("/:id", branchHandler.DeleteCity)
			cityRoutes.PUT("/:id", branchHandler.UpdateCity)
		}
		supportRoutes := authorizedRoutes.Group("/support")
		{
			supportRoutes.POST("", supportHandler.CreateSupportRecord)
			supportRoutes.PUT("", supportHandler.UpdateSupportRecord)
		}
		exclusiveRoutes := authorizedRoutes.Group("/exclusive")
		{
			exclusiveRoutes.GET("", userHandler.GetExclusivePartners)
			exclusiveRoutes.POST("/:id", userHandler.SetPartnerAsExclusive)
			exclusiveRoutes.DELETE("/:id", userHandler.RemovePartnerAsExclusive)
		}
		adminRoutes := authorizedRoutes.Group("/admin")
		{
			adminRoutes.GET("/count/customers", userHandler.GetCustomersCount)
			adminRoutes.GET("/count/partners", userHandler.GetPartnersCount)
			adminRoutes.GET("/count/offers", offerHandler.GetOffersCount)
			adminRoutes.GET("/customers", userHandler.GetAllCustomers)
			adminRoutes.GET("/customer/:id", userHandler.GetCustomerByID)
			adminRoutes.GET("/partners", userHandler.GetAllPartners)
			adminRoutes.GET("/users", userHandler.SearchUsers)
			adminRoutes.GET("/partner/:id", userHandler.GetPartnerByID)
			adminRoutes.GET("/offers", offerHandler.GetAllOffers)
			adminRoutes.GET("/offers/customer/:id", offerHandler.GetOffersOfCustomer)
			adminRoutes.GET("/partners/not-approved", userHandler.GetNotApproved)
			adminRoutes.POST("/approve/:id", userHandler.ApprovePartner)
			adminRoutes.DELETE("/user/:userID", userHandler.AdminDeleteUser)
			adminRoutes.POST("/is-sharable/:id", userHandler.ToggleIsSharable)
			adminRoutes.POST("/upgrade-plan", subscriptionHandler.AdminUpgradeUserPlan)
			adminRoutes.POST("/associate-plan-category", subscriptionHandler.CreatePlanCategoryAssociation)
			adminRoutes.DELETE("/associate-plan-category", subscriptionHandler.RemovePlanCategoryAssociation)
			adminRoutes.GET("/categories", subscriptionHandler.GetCategoriesOfPlan)
			adminRoutes.POST("/activate-user/:id", userHandler.ToggleActive)
		}
	}
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
}
func authUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		if len(token) == 0 {
			err := errors.New("empty token")
			entities.SendAuthError(c, "you are not authorized to access this page, please login first", err)
			return
		}
		id, err := entities.ParseToken(token)
		if err != nil {
			entities.SendAuthError(c, "You are not authorized to access this page, please login first", err)
			return
		}
		c.Set("userID", id)
		c.Next()
	}
}
