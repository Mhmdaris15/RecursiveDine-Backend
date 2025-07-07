package controllers

import (
	"log"
	"net/http"

	"recursiveDine/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type KitchenController struct {
	kitchenService *services.KitchenService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

func NewKitchenController(kitchenService *services.KitchenService) *KitchenController {
	return &KitchenController{
		kitchenService: kitchenService,
	}
}

// @Summary WebSocket for kitchen updates
// @Description Establish WebSocket connection for real-time kitchen updates
// @Tags kitchen
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param token query string true "JWT token for authentication"
// @Success 101 {string} string "Switching Protocols"
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /kitchen/updates [get]
func (ctrl *KitchenController) HandleWebSocket(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Add client to kitchen service
	ctrl.kitchenService.AddClient(conn)
	defer ctrl.kitchenService.RemoveClient(conn)

	// Handle incoming messages
	for {
		// Read message from client
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// For now, we don't handle incoming messages from clients
		// In a more advanced implementation, you might handle client commands
		// like requesting specific order details, updating order status, etc.
	}
}

// @Summary Get active kitchen orders
// @Description Get all orders currently in the kitchen (confirmed, preparing)
// @Tags kitchen
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} repositories.Order
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kitchen/orders [get]
func (ctrl *KitchenController) GetActiveOrders(c *gin.Context) {
	orders, err := ctrl.kitchenService.GetActiveOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":       orders,
		"client_count": ctrl.kitchenService.GetClientCount(),
	})
}

// @Summary Broadcast order update
// @Description Broadcast order status update to all connected kitchen clients
// @Tags kitchen
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Order update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /kitchen/broadcast [post]
func (ctrl *KitchenController) BroadcastUpdate(c *gin.Context) {
	var req struct {
		OrderID    uint   `json:"order_id" binding:"required"`
		UpdateType string `json:"update_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate update type
	validTypes := []string{"status_update", "new_order", "order_ready", "order_cancelled"}
	valid := false
	for _, validType := range validTypes {
		if req.UpdateType == validType {
			valid = true
			break
		}
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update type"})
		return
	}

	// Broadcast the update
	ctrl.kitchenService.BroadcastOrderUpdate(req.OrderID, req.UpdateType)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Update broadcasted successfully",
		"client_count": ctrl.kitchenService.GetClientCount(),
	})
}
