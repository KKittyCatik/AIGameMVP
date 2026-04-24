import { useState, useEffect } from 'react';
import SessionList from './components/SessionList';
import ChatWindow from './components/ChatWindow';
import StatePanel from './components/StatePanel';
import './App.css';

// In dev: connect directly to the Go backend.
// In Docker (behind nginx proxy): use the same host with /ws path.
const WS_URL = import.meta.env.VITE_WS_URL || `ws://${window.location.host}/ws`;

export default function App() {
  const [ws, setWs] = useState(null);
  const [sessions, setSessions] = useState({});
  const [activeSessionId, setActiveSessionId] = useState(null);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const socket = new WebSocket(WS_URL);

    socket.onopen = () => setConnected(true);
    socket.onclose = () => setConnected(false);

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      const incoming = data.sessions || {};
      setSessions(incoming);

      // Auto-select first session if current selection is gone
      setActiveSessionId((prev) => {
        if (prev && incoming[prev]) return prev;
        const keys = Object.keys(incoming);
        return keys.length > 0 ? keys[0] : null;
      });
    };

    setWs(socket);
    return () => socket.close();
  }, []);

  const handleCreateSession = () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ action: 'create_session' }));
    }
  };

  const handleSend = (text) => {
    if (ws && ws.readyState === WebSocket.OPEN && activeSessionId) {
      ws.send(JSON.stringify({ action: 'chat', session_id: activeSessionId, text }));
    }
  };

  const activeSession = sessions[activeSessionId] || null;

  return (
    <div className="app">
      <header className="app-header">
        <span>⚔️ AI Dungeon MVP</span>
        <span className={`status ${connected ? 'connected' : 'disconnected'}`}>
          {connected ? '● Connected' : '○ Disconnected'}
        </span>
      </header>
      <main className="app-body">
        <SessionList
          sessions={sessions}
          activeSessionId={activeSessionId}
          onSelect={setActiveSessionId}
          onCreate={handleCreateSession}
        />
        <ChatWindow session={activeSession} onSend={handleSend} />
        <StatePanel session={activeSession} />
      </main>
    </div>
  );
}
