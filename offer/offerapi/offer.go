package offerapi

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/ahmedaabouzied/tasarruf/branch"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/ahmedaabouzied/tasarruf/offer"
	"github.com/ahmedaabouzied/tasarruf/user"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Handler handles API requests for offers
type Handler struct {
	offersUsecase offer.Usecase
	userUsecase   user.Usecase
	branchUsecase branch.Usecase
	hub           offer.Hub
}

type consumeOfferRequest struct {
	Amount     float64 `json:"amount"`
	CustomerID uint    `json:"customerID"`
	PartnerID  uint    `json:"partnerID"`
}

type dateFilters struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// CreateOfferHandler returns a new API handler for offer endpoints
func CreateOfferHandler(o offer.Usecase, h offer.Hub, b branch.Usecase, u user.Usecase) *Handler {
	handler := Handler{
		offersUsecase: o,
		hub:           h,
		branchUsecase: b,
		userUsecase:   u,
	}
	return &handler
}

// pongwait is the time the server awaits for a pong message
var pongWait = 60 * time.Second

// pingPeriod is the interval for sending a ping message
var pingPeriod = (pongWait * 6) / 10

// wsupgrader upgrades HTTP/HTTPS connection to WS connection
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// AuthMessage the first message expected from the client upon connection
type AuthMessage struct {
	Token string `json:"token"`
}

// Connect handles WS connection on connect endpoint
func (h *Handler) Connect(c *gin.Context) {
	log.Info("conncting")
	h.ConnectWebSocket(c.Writer, c.Request)
}

// ConsumeOffer handles POST request to consume offer endpoint
func (h *Handler) ConsumeOffer(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var req consumeOfferRequest
	err := c.BindJSON(&req)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while processing your request , please try again", err)
		return
	}
	offer, err := h.offersUsecase.ConsumeOffer(ctx, req.CustomerID, req.PartnerID, req.Amount)
	if err != nil {
		entities.SendValidationError(c, errors.Cause(err).Error(), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"receipt": offer,
	})
}

// GetMyOffersHistory handles GET request to offers endpoint
func (h *Handler) GetMyOffersHistory(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var startDate time.Time
	var endDate time.Time
	var err error
	var req dateFilters
	err = c.BindJSON(&req)
	if err != nil {
		entities.SendValidationError(c, "Please select a valid start date", err)
		return
	}
	if req.StartDate != "" && req.EndDate != "" {
		log.Info("START DATE = ", req.StartDate)
		log.Info("END DATE = ", req.EndDate)
		startDate, err = time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			entities.SendValidationError(c, "Please select a valid start date", err)
			return
		}
		endDate, err = time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			entities.SendValidationError(c, "Please select a valid end date", err)
			return
		}
	}
	offers, err := h.offersUsecase.GetMyOffersHistory(ctx, startDate, endDate)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while getting your offers history , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"offers": offers,
	})
}

// GetOffersOfCustomer handles GET /offers/customer/:ID
func (h *Handler) GetOffersOfCustomer(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	offers, err := h.offersUsecase.GetByCustomer(ctx, uint(customerID))
	if err != nil {
		entities.SendParsingError(c, "error getting offers, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"offers": offers,
	})
}

// SendOffersStaticMail handles GET request to offers endpoint
func (h *Handler) SendOffersStaticMail(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	var startDate time.Time
	var endDate time.Time
	var err error
	var req dateFilters
	err = c.BindJSON(&req)
	if err != nil {
		entities.SendValidationError(c, "Please select a valid start date", err)
		return
	}
	if req.StartDate != "" && req.EndDate != "" {
		log.Info("START DATE = ", req.StartDate)
		log.Info("END DATE = ", req.EndDate)
		startDate, err = time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			entities.SendValidationError(c, "Please select a valid start date", err)
			return
		}
		endDate, err = time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			entities.SendValidationError(c, "Please select a valid end date", err)
			return
		}
	}
	err = h.offersUsecase.SendOffersStaticMail(ctx, startDate, endDate)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while getting your offers history , please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "An email has been sent to your email account with your offers history",
	})
}

// GetOffer handles GET request to offers/:id endpoint
func (h *Handler) GetOffer(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	offerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		entities.SendParsingError(c, "There has been an error while parsing your information , please try again", err)
		return
	}
	offer, err := h.offersUsecase.GetOffer(ctx, uint(offerID))
	if err != nil {
		entities.SendValidationError(c, "There has been an error while getting information from the server, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"offer": offer,
	})
}

// ConnectWebSocket connects the user into the web socket hub
func (h *Handler) ConnectWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Info("connecting WS")
	wsupgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var authMessage AuthMessage
	err = conn.ReadJSON(&authMessage)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte("invalid token"))
		log.Error(err.Error())
		conn.Close()
		return
	}
	if authMessage.Token == "" {
		log.Error(authMessage)
		conn.WriteMessage(websocket.CloseMessage, []byte("empty token"))
		conn.Close()
		return
	}
	ID, err := entities.ParseToken(authMessage.Token)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte("invalid token"))
		log.Error(err.Error())
		conn.Close()
		return
	}
	user, err := h.userUsecase.GetUser(context.Background(), ID)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte("invalid token"))
		log.Error(err.Error())
		conn.Close()
		return
	}
	h.hub.AddUser(user, conn)
	conn.WriteMessage(websocket.TextMessage, []byte("success: authenticated successfully"))
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg AuthMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Error(err)
			h.hub.RemoveUser(user.ID)
			return
		}
	}
}

// GetOffersCount handles GET /admin/count/offers
func (h *Handler) GetOffersCount(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	count, err := h.offersUsecase.GetOffersCount(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while getting information from the server, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
	return
}

// GetAllOffers handles GET /admin/all/offers
func (h *Handler) GetAllOffers(c *gin.Context) {
	ctx := context.Background()
	userID := c.MustGet("userID").(uint)
	ctx = context.WithValue(ctx, entities.UserIDKey, userID)
	offers, err := h.offersUsecase.GetAllOffers(ctx)
	if err != nil {
		entities.SendValidationError(c, "There has been an error while getting information from the server, please try again", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"offers": offers,
	})
	return
}
