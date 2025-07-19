package services

import (
	"encoding/json"
	"log"
	"sync"

	"recursiveDine/internal/repositories"

	"github.com/gorilla/websocket"
)

type KitchenService struct {
	orderRepo *repositories.OrderRepository
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mutex     sync.RWMutex
}

type KitchenUpdate struct {
	Type    string                `json:"type"`
	OrderID uint                  `json:"order_id"`
	Order   *repositories.Order   `json:"order,omitempty"`
	Status  repositories.OrderStatus `json:"status,omitempty"`
}

func NewKitchenService(orderRepo *repositories.OrderRepository) *KitchenService {
	service := &KitchenService{
		orderRepo: orderRepo,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}

	// Start the broadcast goroutine
	go service.handleMessages()

	return service
}

func (s *KitchenService) AddClient(conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.clients[conn] = true

	// Send current kitchen orders to new client
	orders, err := s.orderRepo.GetKitchenOrders()
	if err != nil {
		log.Printf("Error fetching kitchen orders: %v", err)
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"type":   "initial_orders",
		"orders": orders,
	})
	if err != nil {
		log.Printf("Error marshaling initial orders: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Error sending initial orders: %v", err)
		delete(s.clients, conn)
		conn.Close()
	}
}

func (s *KitchenService) RemoveClient(conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.clients[conn]; ok {
		delete(s.clients, conn)
		conn.Close()
	}
}

func (s *KitchenService) BroadcastOrderUpdate(orderID uint, updateType string) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		log.Printf("Error fetching order for broadcast: %v", err)
		return
	}

	update := KitchenUpdate{
		Type:    updateType,
		OrderID: orderID,
		Order:   order,
		Status:  order.Status,
	}

	data, err := json.Marshal(update)
	if err != nil {
		log.Printf("Error marshaling kitchen update: %v", err)
		return
	}

	select {
	case s.broadcast <- data:
	default:
		log.Printf("Broadcast channel full, dropping message")
	}
}

func (s *KitchenService) handleMessages() {
	for {
		select {
		case message := <-s.broadcast:
			s.mutex.RLock()
			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error writing to client: %v", err)
					delete(s.clients, client)
					client.Close()
				}
			}
			s.mutex.RUnlock()
		}
	}
}

func (s *KitchenService) GetActiveOrders() ([]repositories.Order, error) {
	return s.orderRepo.GetKitchenOrders()
}

func (s *KitchenService) GetClientCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.clients)
}
