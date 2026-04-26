# AIGameMVP

A minimal viable prototype of an AI-driven text-based game.

## Tech Stack

- **Backend**: Go + gorilla/websocket + OpenAI-compatible API (OpenRouter by default)
- **Frontend**: React + Vite (served via nginx in Docker)
- **Communication**: WebSocket
- **Infrastructure**: Docker, docker-compose, Makefile

## Project Structure

```
AIGameMVP/
├── backend/
│   ├── main.go        # WebSocket server, session manager, OpenAI AI engine
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── App.jsx
│   │   └── components/
│   │       ├── SessionList.jsx
│   │       ├── ChatWindow.jsx
│   │       └── StatePanel.jsx
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml
├── Makefile
├── .env.example
└── README.md
```

## Quick Start (Docker – recommended)

1. Copy `.env.example` to `.env` and add your API key:

   ```bash
   cp .env.example .env
   # edit .env and set AI_API_KEY=your-api-key-here
   ```

2. Build and start everything:

   ```bash
   make build
   make up
   ```

3. Open http://localhost:3000 in your browser.

4. To stop:

   ```bash
   make down
   ```

## Running locally (without Docker)

### Backend

```bash
cd backend
export AI_API_KEY=your-api-key-here   # required for real AI responses
export AI_BASE_URL=https://openrouter.ai/api/v1  # optional, this is the default
export AI_MODEL=google/gemma-4-26b-a4b-it:free   # optional, this is the default
export SYSTEM_PROMPT="..."            # optional, uses built-in default
go run main.go
# Server starts on http://localhost:8080
```

### Frontend

```bash
cd frontend
echo "VITE_WS_URL=ws://localhost:8080/ws" > .env.local
npm install
npm run dev
# App starts on http://localhost:5173
```

Open http://localhost:5173 in your browser.

## Configuration

| Variable         | Where          | Description                                           |
|------------------|----------------|-------------------------------------------------------|
| `AI_API_KEY`     | Backend / .env | **Required** for real AI. Without it, a fallback message is shown. |
| `AI_BASE_URL`    | Backend / .env | Base URL of the OpenAI-compatible API. Default: `https://openrouter.ai/api/v1`. |
| `AI_MODEL`       | Backend / .env | Model identifier to use. Default: `google/gemma-4-26b-a4b-it:free`. |
| `SYSTEM_PROMPT`  | Backend / .env | Optional custom system prompt for the Game Master AI. |
| `VITE_WS_URL`    | Frontend dev   | WebSocket URL (e.g. `ws://localhost:8080/ws`). Omit in Docker. |

## How to play

1. Click **+ New Session** to start a game (up to 3 sessions).
2. Type anything in the chat – the AI Game Master will respond and update your **Inventory** and **Buffs/Debuffs** automatically.
3. Your Inventory and active effects update in the right panel.

## Game State

- **Inventory**: items carried by the player (e.g. *Torch*, *Sword*).
- **Buffs/Debuffs**: active status effects (e.g. *drunk*, *slowed*, *heavy backpack*).

The AI returns a JSON response that may add or remove items and effects based on your actions.

