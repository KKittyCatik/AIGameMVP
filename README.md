# AIGameMVP

A minimal viable prototype of an AI-driven text-based game.

## Tech Stack

- **Backend**: Go + gorilla/websocket
- **Frontend**: React + Vite
- **Communication**: WebSocket

## Project Structure

```
AIGameMVP/
├── backend/
│   ├── main.go        # WebSocket server, session manager, mock AI engine
│   ├── go.mod
│   └── go.sum
└── frontend/
    └── src/
        ├── App.jsx
        └── components/
            ├── SessionList.jsx
            ├── ChatWindow.jsx
            └── StatePanel.jsx
```

## Running the project

### Backend

```bash
cd backend
go run main.go
# Server starts on http://localhost:8080
```

### Frontend

```bash
cd frontend
npm install
npm run dev
# App starts on http://localhost:5173
```

Open http://localhost:5173 in your browser.

## How to play

1. Click **+ New Session** to start a game (up to 3 sessions).
2. Type commands in the chat box:
   - `look` — examine your surroundings
   - `take` / `grab` — pick up an item
   - `attack` / `fight` — engage in combat (costs HP)
   - `heal` / `drink` — restore HP
   - Any other text gets a random atmospheric response.
3. Your HP and Inventory update in the right panel.
