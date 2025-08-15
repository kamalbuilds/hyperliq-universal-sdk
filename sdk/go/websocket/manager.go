package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hyperliquid-labs/hyperliquid-go-sdk/types"
)

type Manager struct {
	url            string
	conn           *websocket.Conn
	mu             sync.RWMutex
	subscriptions  map[string]*Subscription
	handlers       map[string]MessageHandler
	reconnectDelay time.Duration
	maxReconnect   int
	pingInterval   time.Duration
	pongTimeout    time.Duration
	stopCh         chan struct{}
	isConnected    bool
	reconnectCount int
	messageQueue   chan []byte
}

type MessageHandler func(data json.RawMessage) error

type Subscription struct {
	ID       string
	Type     string
	Callback MessageHandler
	Request  types.WSSubscription
}

func NewManager(url string) *Manager {
	return &Manager{
		url:            url,
		subscriptions:  make(map[string]*Subscription),
		handlers:       make(map[string]MessageHandler),
		reconnectDelay: 5 * time.Second,
		maxReconnect:   10,
		pingInterval:   30 * time.Second,
		pongTimeout:    10 * time.Second,
		stopCh:         make(chan struct{}),
		messageQueue:   make(chan []byte, 1000),
	}
}

func (m *Manager) Connect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isConnected {
		return nil
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, m.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	m.conn = conn
	m.isConnected = true
	m.reconnectCount = 0

	go m.readLoop()
	go m.pingLoop()
	go m.processMessages()

	if err := m.resubscribeAll(); err != nil {
		return fmt.Errorf("failed to resubscribe: %w", err)
	}

	return nil
}

func (m *Manager) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isConnected {
		return nil
	}

	close(m.stopCh)
	m.isConnected = false

	if m.conn != nil {
		m.conn.Close()
	}

	return nil
}

func (m *Manager) Subscribe(sub types.WSSubscription, handler MessageHandler) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isConnected {
		return "", fmt.Errorf("not connected")
	}

	subID := m.generateSubscriptionID(sub)
	
	subscription := &Subscription{
		ID:       subID,
		Type:     sub.Type,
		Callback: handler,
		Request:  sub,
	}

	m.subscriptions[subID] = subscription
	m.handlers[sub.Type] = handler

	req := types.WSRequest{
		Method:       "subscribe",
		Subscription: sub,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	if err := m.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return "", fmt.Errorf("failed to send subscription: %w", err)
	}

	return subID, nil
}

func (m *Manager) Unsubscribe(subID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sub, exists := m.subscriptions[subID]
	if !exists {
		return fmt.Errorf("subscription not found: %s", subID)
	}

	req := types.WSRequest{
		Method:       "unsubscribe",
		Subscription: sub.Request,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	if err := m.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send unsubscribe: %w", err)
	}

	delete(m.subscriptions, subID)
	delete(m.handlers, sub.Type)

	return nil
}

func (m *Manager) readLoop() {
	defer func() {
		m.mu.Lock()
		m.isConnected = false
		m.mu.Unlock()
		m.handleReconnect()
	}()

	for {
		select {
		case <-m.stopCh:
			return
		default:
			messageType, data, err := m.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			if messageType == websocket.TextMessage {
				m.messageQueue <- data
			}
		}
	}
}

func (m *Manager) processMessages() {
	for {
		select {
		case <-m.stopCh:
			return
		case data := <-m.messageQueue:
			var msg types.WSMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			m.handleMessage(msg)
		}
	}
}

func (m *Manager) handleMessage(msg types.WSMessage) {
	m.mu.RLock()
	handler, exists := m.handlers[msg.Channel]
	m.mu.RUnlock()

	if !exists {
		if msg.Channel != "pong" {
			log.Printf("No handler for channel: %s", msg.Channel)
		}
		return
	}

	if err := handler(msg.Data); err != nil {
		log.Printf("Handler error for channel %s: %v", msg.Channel, err)
	}
}

func (m *Manager) pingLoop() {
	ticker := time.NewTicker(m.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.mu.Lock()
			if m.isConnected && m.conn != nil {
				ping := map[string]string{"method": "ping"}
				if err := m.conn.WriteJSON(ping); err != nil {
					log.Printf("Failed to send ping: %v", err)
					m.isConnected = false
					m.mu.Unlock()
					return
				}
			}
			m.mu.Unlock()
		}
	}
}

func (m *Manager) handleReconnect() {
	for m.reconnectCount < m.maxReconnect {
		select {
		case <-m.stopCh:
			return
		case <-time.After(m.reconnectDelay):
			m.reconnectCount++
			log.Printf("Attempting reconnect %d/%d", m.reconnectCount, m.maxReconnect)
			
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err := m.Connect(ctx)
			cancel()
			
			if err == nil {
				log.Println("Reconnected successfully")
				return
			}
			
			log.Printf("Reconnect failed: %v", err)
		}
	}
	
	log.Printf("Max reconnect attempts reached")
}

func (m *Manager) resubscribeAll() error {
	for _, sub := range m.subscriptions {
		req := types.WSRequest{
			Method:       "subscribe",
			Subscription: sub.Request,
		}

		data, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("failed to marshal resubscription: %w", err)
		}

		if err := m.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return fmt.Errorf("failed to resubscribe: %w", err)
		}
	}

	return nil
}

func (m *Manager) generateSubscriptionID(sub types.WSSubscription) string {
	switch sub.Type {
	case "allMids":
		return "allMids"
	case "l2Book":
		return fmt.Sprintf("l2Book:%s", sub.Coin)
	case "trades":
		return fmt.Sprintf("trades:%s", sub.Coin)
	case "candle":
		return fmt.Sprintf("candle:%s:%s", sub.Coin, sub.Interval)
	case "userEvents", "orderUpdates":
		return fmt.Sprintf("%s:%s", sub.Type, sub.User)
	default:
		return fmt.Sprintf("%s:%s", sub.Type, sub.Coin)
	}
}

func (m *Manager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isConnected
}

func (m *Manager) GetStats() types.WSStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return types.WSStats{
		Connected:    m.isConnected,
		Reconnects:   int64(m.reconnectCount),
		Subscriptions: len(m.subscriptions),
	}
}