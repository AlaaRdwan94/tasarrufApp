package hub

import (
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// Hub interface represents an interface for a hub of WS clients
type Hub interface {
	AddUser(user *entities.User, conn *websocket.Conn)
	RemoveUser(ID uint)
	GetUser(ID uint) (*entities.User, error)
	HasUser(ID uint) bool
	SendOfferToUser(ID uint, offer *entities.Offer) error
}

// usersHub is an implementation of the Hub interface
type usersHub struct {
	lock        *sync.Mutex
	addUser     chan *userClient
	removeUser  chan uint
	users       map[uint]*userClient
	sendMessage chan *message
}

type message struct {
	Client  *userClient     `json:"client,omitempty"`
	Offer   *entities.Offer `json:"offer,omitempty"`
	Message string          `json:"message,omitempty"`
}

// userClient is a wrapper arrount a user object
type userClient struct {
	user *entities.User
	conn *websocket.Conn
}

// CreateUserHub returns a new instance of the Hub interface
func CreateUserHub() Hub {
	h := usersHub{
		lock:        &sync.Mutex{},
		addUser:     make(chan *userClient),
		removeUser:  make(chan uint),
		users:       make(map[uint]*userClient),
		sendMessage: make(chan *message),
	}
	go h.Run()
	return &h
}

func (h *usersHub) Run() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case addRequest := <-h.addUser:
			h.handleAddRequest(addRequest)
		case removeRequest := <-h.removeUser:
			h.handleRemoveRequest(removeRequest)
		case sendMessageRequest := <-h.sendMessage:
			h.handleSendMessage(sendMessageRequest)
		case <-ticker.C:
			h.handleTicker()
		}
	}
}

func (h *usersHub) handleAddRequest(client *userClient) {
	_, alreadyAdded := h.users[client.user.ID]
	if !alreadyAdded {
		h.lock.Lock()
		h.users[client.user.ID] = client
		h.lock.Unlock()
		log.Info("Added user with ID :", client.user.ID, " to Hub")
	}
}

func (h *usersHub) handleRemoveRequest(ID uint) {
	_, alreadyAdded := h.users[ID]
	if alreadyAdded {
		h.lock.Lock()
		delete(h.users, ID)
		h.lock.Unlock()
		log.Error("disconnected user : ", ID)
	}
}

func (h *usersHub) handleTicker() {
	for _, client := range h.users {
		err := client.conn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			h.lock.Lock()
			delete(h.users, client.user.ID)
			h.lock.Unlock()
		}
	}
}

func (h *usersHub) handleSendMessage(m *message) {
	_, userConnected := h.users[m.Client.user.ID]
	if userConnected {
		err := m.Client.conn.WriteJSON(m)
		if err != nil {
			log.Error("error sending message to user : ", m.Client.user.ID, " : ", err)
		}
	}
}

func (h *usersHub) AddUser(user *entities.User, conn *websocket.Conn) {
	client := userClient{
		user: user,
		conn: conn,
	}
	h.addUser <- &client
}

func (h *usersHub) RemoveUser(ID uint) {
	h.removeUser <- ID
}

func (h *usersHub) GetUser(ID uint) (*entities.User, error) {
	client, ok := h.users[ID]
	if !ok {
		err := errors.New("client not found")
		log.Error(ID, " : ", err)
		return nil, err
	}
	return client.user, nil
}

func (h *usersHub) HasUser(ID uint) bool {
	_, ok := h.users[ID]
	return ok
}

func (h *usersHub) SendOfferToUser(ID uint, offer *entities.Offer) error {
	client, ok := h.users[ID]
	if !ok {
		err := errors.New("client not found")
		log.Error(ID, " : ", err)
		return err
	}
	message := message{
		Client: client,
		Offer:  offer,
	}
	log.Info(message)
	h.sendMessage <- &message
	return nil
}
