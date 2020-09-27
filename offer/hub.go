package offer

import (
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/gorilla/websocket"
)

// Hub interface represents an interface for a hub of WS clients
type Hub interface {
	AddUser(user *entities.User, conn *websocket.Conn)
	RemoveUser(ID uint)
	GetUser(ID uint) (*entities.User, error)
	HasUser(ID uint) bool
	SendOfferToUser(ID uint, offer *entities.Offer) error
}
