package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// --- Domain Models ---

type GameState struct {
	HP        int      `json:"hp"`
	Inventory []string `json:"inventory"`
}

type Message struct {
	Sender string `json:"sender"` // "user" or "ai"
	Text   string `json:"text"`
}

type Session struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	State    GameState `json:"state"`
	Messages []Message `json:"messages"`
}

// --- Global State (In-Memory) ---

var (
	sessions = make(map[string]*Session)
	mu       sync.Mutex
)

// --- AI Engine (Mock) ---

type AIResponse struct {
	Text         string
	HPDelta      int
	InventoryAdd []string
}

func ProcessPlayerMessage(msg string, state GameState) AIResponse {
	msg = strings.ToLower(strings.TrimSpace(msg))
	resp := AIResponse{HPDelta: 0, InventoryAdd: []string{}}

	switch {
	case strings.Contains(msg, "look"):
		resp.Text = "You look around. It's dark, but you spot a shiny sword on the ground."
	case strings.Contains(msg, "take") || strings.Contains(msg, "grab"):
		resp.Text = "You grabbed the item!"
		resp.InventoryAdd = append(resp.InventoryAdd, "Sword")
	case strings.Contains(msg, "attack") || strings.Contains(msg, "fight"):
		resp.Text = "You swing wildly at the darkness. You trip and hurt yourself."
		resp.HPDelta = -10
	case strings.Contains(msg, "heal") || strings.Contains(msg, "drink"):
		resp.Text = "You drink a strange potion you found in your pocket. You feel better!"
		resp.HPDelta = 20
	default:
		responses := []string{
			"The shadows seem to whisper back.",
			"You hear a distant echoing footstep.",
			"A cold wind blows through the cavern.",
		}
		resp.Text = responses[rand.Intn(len(responses))]
	}
	return resp
}

// --- WebSocket Handling ---

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ClientMessage struct {
	Action    string `json:"action"` // "create_session" or "chat"
	SessionID string `json:"session_id"`
	Text      string `json:"text"`
}

type ServerUpdate struct {
	Sessions map[string]*Session `json:"sessions"`
}

func sendUpdate(ws *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	update := ServerUpdate{Sessions: sessions}
	if err := ws.WriteJSON(update); err != nil {
		log.Println("write error:", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer ws.Close()

	// Send initial state
	sendUpdate(ws)

	for {
		var req ClientMessage
		if err := ws.ReadJSON(&req); err != nil {
			break
		}

		mu.Lock()
		switch req.Action {
		case "create_session":
			if len(sessions) < 3 {
				id := fmt.Sprintf("sess_%d", time.Now().UnixNano())
				sessions[id] = &Session{
					ID:   id,
					Name: fmt.Sprintf("Adventure %d", len(sessions)+1),
					State: GameState{
						HP:        100,
						Inventory: []string{"Torch"},
					},
					Messages: []Message{
						{Sender: "ai", Text: "Welcome to the dungeon. What do you do?"},
					},
				}
			}
		case "chat":
			if sess, exists := sessions[req.SessionID]; exists && req.Text != "" {
				sess.Messages = append(sess.Messages, Message{Sender: "user", Text: req.Text})

				aiResp := ProcessPlayerMessage(req.Text, sess.State)

				sess.State.HP += aiResp.HPDelta
				if sess.State.HP > 100 {
					sess.State.HP = 100
				}
				if sess.State.HP < 0 {
					sess.State.HP = 0
				}
				sess.State.Inventory = append(sess.State.Inventory, aiResp.InventoryAdd...)

				sess.Messages = append(sess.Messages, Message{Sender: "ai", Text: aiResp.Text})
			}
		}
		mu.Unlock()

		sendUpdate(ws)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) //nolint:staticcheck
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
