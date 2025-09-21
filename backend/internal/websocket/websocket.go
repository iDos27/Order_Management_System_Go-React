package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Struktura wiadomości
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Aktualizacja zamówienia
type OrderUpdateMessage struct {
	OrderId   int    `json:"order_id"`
	NewStatus string `json:"new_status"`
	UpdatedBy string `json:"updated_by"`
}

type Client struct {
	Conn *websocket.Conn
	Send chan Message
}

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

// Tworzenie huba WebSocketów
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

// Uruchomienie huba
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("New client connected")
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Println("Client disconnected")
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastOrderUpdate(orderId int, newStatus, updatedBy string) {
	message := Message{
		Type: "order_update",
		Payload: OrderUpdateMessage{
			OrderId:   orderId,
			NewStatus: newStatus,
			UpdatedBy: updatedBy,
		},
	}
	select {
	case h.Broadcast <- message:
		log.Printf("Broadcasting order update: %+v", message)
	default:
		log.Println("Cant send message, no clients connected")
	}
}

func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}

		client := &Client{
			Conn: conn,
			Send: make(chan Message, 256),
		}

		hub.Register <- client

		go writePump(client, hub)
		go readPump(client, hub)
	}
}

// readPump - odbiera wiadomości od klienta
func readPump(client *Client, hub *Hub) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	for {
		var msg Message
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Błąd WebSocket: %v", err)
			}
			break
		}

		// Tutaj możemy obsłużyć wiadomości od klienta
		log.Printf("Otrzymano wiadomość: %s", msg.Type)
	}
}

// writePump - wysyła wiadomości do klienta
func writePump(client *Client, hub *Hub) {
	defer client.Conn.Close()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteJSON(message); err != nil {
				log.Printf("Błąd wysyłania: %v", err)
				return
			}
		}
	}
}
