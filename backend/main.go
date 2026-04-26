package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// --- Domain Models ---

type GameState struct {
	Inventory []string `json:"inventory"`
	Effects   []string `json:"effects"`
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

// --- AI Engine ---

// AIResponse is the structure the LLM must return as JSON.
type AIResponse struct {
	Text            string   `json:"text"`
	InventoryAdd    []string `json:"inventory_add"`
	InventoryRemove []string `json:"inventory_remove"`
	EffectsAdd      []string `json:"effects_add"`
	EffectsRemove   []string `json:"effects_remove"`
}

// openAIChatRequest mirrors the /v1/chat/completions request body.
type openAIChatRequest struct {
	Model    string              `json:"model"`
	Messages []openAIChatMessage `json:"messages"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIChatResponse holds the fields we need from the API response.
type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

const defaultSystemPrompt = `You are the Game Master of a text-based adventure game.
Respond to the player's action with a short narrative (1–3 sentences).
You MUST reply ONLY with a valid JSON object (no markdown, no extra text) matching this schema:
{
  "text": "<narrative response>",
  "inventory_add": ["<item>"],
  "inventory_remove": ["<item>"],
  "effects_add": ["<buff or debuff, e.g. drunk, slowed, poisoned>"],
  "effects_remove": ["<effect to remove>"]
}
All list fields default to empty arrays [] when unused.
Keep the game atmospheric and fun.`

func ProcessPlayerMessage(playerMsg string, state GameState, history []Message) AIResponse {
	apiKey := os.Getenv("AI_API_KEY")
	if apiKey == "" {
		log.Println("AI_API_KEY not set – using fallback response")
		return AIResponse{Text: "The shadows whisper... (set AI_API_KEY to enable the AI Game Master)"}
	}

	baseURL := os.Getenv("AI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}

	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "google/gemma-4-26b-a4b-it:free"
	}

	systemPrompt := os.Getenv("SYSTEM_PROMPT")
	if systemPrompt == "" {
		systemPrompt = defaultSystemPrompt
	}

	// Build conversation context (last 10 messages for brevity).
	stateJSON, _ := json.Marshal(state)
	contextMsg := fmt.Sprintf("Current player state: %s", string(stateJSON))

	msgs := []openAIChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "system", Content: contextMsg},
	}

	start := 0
	if len(history) > 10 {
		start = len(history) - 10
	}
	for _, m := range history[start:] {
		role := "user"
		if m.Sender == "ai" {
			role = "assistant"
		}
		msgs = append(msgs, openAIChatMessage{Role: role, Content: m.Text})
	}
	msgs = append(msgs, openAIChatMessage{Role: "user", Content: playerMsg})

	reqBody := openAIChatRequest{
		Model:    model,
		Messages: msgs,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("marshal error:", err)
		return AIResponse{Text: "The magic fails. (internal error)"}
	}

	endpoint := strings.TrimRight(baseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequest("POST", endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Println("request build error:", err)
		return AIResponse{Text: "The magic fails. (internal error)"}
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("HTTP-Referer", "http://localhost:8080")
	httpReq.Header.Set("X-Title", "AIGameMVP")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Println("OpenAI request error:", err)
		return AIResponse{Text: "The GM is lost in thought… (API error)"}
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read error:", err)
		return AIResponse{Text: "The GM is lost in thought… (read error)"}
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned non-200 status %d: %s", resp.StatusCode, string(raw))
		return AIResponse{Text: fmt.Sprintf("The GM is unavailable (HTTP %d). Check logs for details.", resp.StatusCode)}
	}

	var apiResp openAIChatResponse
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		log.Println("unmarshal error:", err, "body:", string(raw))
		return AIResponse{Text: "The GM is lost in thought… (parse error)"}
	}

	if apiResp.Error != nil {
		log.Println("OpenAI API error:", apiResp.Error.Message)
		return AIResponse{Text: "The GM is unavailable: " + apiResp.Error.Message}
	}

	if len(apiResp.Choices) == 0 {
		log.Println("OpenAI returned no choices")
		return AIResponse{Text: "Silence fills the room."}
	}

	content := strings.TrimSpace(apiResp.Choices[0].Message.Content)
	var aiResp AIResponse
	if err := json.Unmarshal([]byte(content), &aiResp); err != nil {
		log.Println("AI response parse error:", err, "content:", content)
		// Fall back: treat the raw content as plain text.
		return AIResponse{Text: content}
	}

	return aiResp
}

// removeFromSlice removes all occurrences of items in toRemove from src (case-insensitive).
func removeFromSlice(src, toRemove []string) []string {
	if len(toRemove) == 0 {
		return src
	}
	removeSet := make(map[string]struct{}, len(toRemove))
	for _, v := range toRemove {
		removeSet[strings.ToLower(v)] = struct{}{}
	}
	result := make([]string, 0, len(src))
	for _, v := range src {
		if _, skip := removeSet[strings.ToLower(v)]; !skip {
			result = append(result, v)
		}
	}
	return result
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
						Inventory: []string{"Torch"},
						Effects:   []string{},
					},
					Messages: []Message{
						{Sender: "ai", Text: "Welcome to the dungeon. What do you do?"},
					},
				}
			}
		case "chat":
			if sess, exists := sessions[req.SessionID]; exists && req.Text != "" {
				sess.Messages = append(sess.Messages, Message{Sender: "user", Text: req.Text})

				// Copy state and history before releasing the lock so the long
				// OpenAI call doesn't block other WebSocket clients.
				stateCopy := GameState{
					Inventory: append([]string(nil), sess.State.Inventory...),
					Effects:   append([]string(nil), sess.State.Effects...),
				}
				historyCopy := append([]Message(nil), sess.Messages...)
				mu.Unlock()

				aiResp := ProcessPlayerMessage(req.Text, stateCopy, historyCopy)

				mu.Lock()
				// Re-fetch session in case it was removed while we were unlocked.
				sess, exists = sessions[req.SessionID]
				if !exists {
					break
				}

				// Update inventory
				sess.State.Inventory = append(sess.State.Inventory, aiResp.InventoryAdd...)
				sess.State.Inventory = removeFromSlice(sess.State.Inventory, aiResp.InventoryRemove)

				// Update effects
				sess.State.Effects = append(sess.State.Effects, aiResp.EffectsAdd...)
				sess.State.Effects = removeFromSlice(sess.State.Effects, aiResp.EffectsRemove)

				sess.Messages = append(sess.Messages, Message{Sender: "ai", Text: aiResp.Text})
			}
		}
		mu.Unlock()

		sendUpdate(ws)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
