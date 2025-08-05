package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Manager handles WebSocket connections and subscriptions
type Manager struct {
	url               string
	conn              *websocket.Conn
	mu                sync.RWMutex
	subscriptions     map[string]*Subscription
	reconnectInterval time.Duration
	pingInterval      time.Duration
	pongTimeout       time.Duration
	isConnected       bool
	ctx               context.Context
	cancel            context.CancelFunc
}

// Subscription represents a WebSocket subscription
type Subscription struct {
	ID       string
	Type     string
	Channel  chan json.RawMessage
	Request  interface{}
}

// Message represents a WebSocket message
type Message struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

// NewManager creates a new WebSocket manager
func NewManager(url string) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		url:               url,
		subscriptions:     make(map[string]*Subscription),
		reconnectInterval: 5 * time.Second,
		pingInterval:      30 * time.Second,
		pongTimeout:       10 * time.Second,
		ctx:               ctx,
		cancel:            cancel,
	}
}

// Connect establishes WebSocket connection
func (m *Manager) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(m.url, nil)
	if err != nil {
		return fmt.Errorf("websocket dial error: %w", err)
	}

	m.conn = conn
	m.isConnected = true

	// Set pong handler
	m.conn.SetPongHandler(func(string) error {
		m.conn.SetReadDeadline(time.Now().Add(m.pongTimeout))
		return nil
	})

	// Start read and ping loops
	go m.readLoop()
	go m.pingLoop()

	// Resubscribe to all existing subscriptions
	for _, sub := range m.subscriptions {
		if err := m.sendSubscription(sub); err != nil {
			fmt.Printf("Failed to resubscribe %s: %v\n", sub.ID, err)
		}
	}

	return nil
}

// Subscribe creates a new subscription
func (m *Manager) Subscribe(subType string, params interface{}) (*Subscription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create subscription
	sub := &Subscription{
		ID:      fmt.Sprintf("%s_%d", subType, time.Now().UnixNano()),
		Type:    subType,
		Channel: make(chan json.RawMessage, 100),
		Request: params,
	}

	// Store subscription
	m.subscriptions[sub.ID] = sub

	// Send subscription if connected
	if m.isConnected {
		if err := m.sendSubscription(sub); err != nil {
			delete(m.subscriptions, sub.ID)
			return nil, err
		}
	}

	return sub, nil
}

// Unsubscribe removes a subscription
func (m *Manager) Unsubscribe(sub *Subscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.subscriptions, sub.ID)
	close(sub.Channel)

	if m.isConnected {
		msg := map[string]interface{}{
			"method": "unsubscribe",
			"subscription": map[string]interface{}{
				"type": sub.Type,
			},
		}
		return m.conn.WriteJSON(msg)
	}

	return nil
}

// Close closes the WebSocket connection
func (m *Manager) Close() error {
	m.cancel()
	
	m.mu.Lock()
	defer m.mu.Unlock()

	m.isConnected = false

	// Close all subscription channels
	for _, sub := range m.subscriptions {
		close(sub.Channel)
	}
	m.subscriptions = make(map[string]*Subscription)

	if m.conn != nil {
		return m.conn.Close()
	}

	return nil
}

// IsConnected returns connection status
func (m *Manager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isConnected
}

// sendSubscription sends a subscription request
func (m *Manager) sendSubscription(sub *Subscription) error {
	msg := map[string]interface{}{
		"method":       "subscribe",
		"subscription": sub.Request,
	}

	return m.conn.WriteJSON(msg)
}

// readLoop continuously reads messages from WebSocket
func (m *Manager) readLoop() {
	defer func() {
		m.mu.Lock()
		m.isConnected = false
		m.mu.Unlock()
		m.handleReconnect()
	}()

	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			_, message, err := m.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("WebSocket read error: %v\n", err)
				}
				return
			}

			m.handleMessage(message)
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (m *Manager) handleMessage(data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		// Try to unmarshal as a different format
		var altMsg map[string]json.RawMessage
		if err := json.Unmarshal(data, &altMsg); err != nil {
			fmt.Printf("Failed to unmarshal message: %v\n", err)
			return
		}
		// Handle alternative message format
		m.handleAlternativeMessage(altMsg)
		return
	}

	// Broadcast to relevant subscriptions
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, sub := range m.subscriptions {
		// Match subscription based on channel type
		if m.matchesSubscription(msg.Channel, sub) {
			select {
			case sub.Channel <- msg.Data:
			default:
				// Channel full, drop message
				fmt.Printf("Subscription channel full for %s\n", sub.ID)
			}
		}
	}
}

// handleAlternativeMessage handles messages in alternative formats
func (m *Manager) handleAlternativeMessage(msg map[string]json.RawMessage) {
	// Handle different message formats based on Hyperliquid's WebSocket API
	if channel, ok := msg["channel"]; ok {
		var channelStr string
		if err := json.Unmarshal(channel, &channelStr); err == nil {
			// Broadcast to relevant subscriptions
			m.mu.RLock()
			defer m.mu.RUnlock()

			for _, sub := range m.subscriptions {
				if m.matchesSubscription(channelStr, sub) {
					select {
					case sub.Channel <- msg["data"]:
					default:
						// Channel full
					}
				}
			}
		}
	}
}

// matchesSubscription checks if a channel matches a subscription
func (m *Manager) matchesSubscription(channel string, sub *Subscription) bool {
	// Implement matching logic based on Hyperliquid's channel naming
	return channel == sub.Type || channel == fmt.Sprintf("%s:%s", sub.Type, sub.ID)
}

// pingLoop sends periodic ping messages
func (m *Manager) pingLoop() {
	ticker := time.NewTicker(m.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.mu.Lock()
			if m.isConnected {
				if err := m.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					fmt.Printf("Ping error: %v\n", err)
					m.isConnected = false
					m.mu.Unlock()
					return
				}
			}
			m.mu.Unlock()
		}
	}
}

// handleReconnect attempts to reconnect after disconnection
func (m *Manager) handleReconnect() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-time.After(m.reconnectInterval):
			fmt.Println("Attempting to reconnect...")
			if err := m.Connect(); err != nil {
				fmt.Printf("Reconnection failed: %v\n", err)
				continue
			}
			fmt.Println("Reconnected successfully")
			return
		}
	}
}