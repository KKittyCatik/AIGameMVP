import { useEffect, useRef, useState } from 'react';

export default function ChatWindow({ session, onSend }) {
  const [input, setInput] = useState('');
  const bottomRef = useRef(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [session?.messages]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!input.trim() || !session) return;
    onSend(input.trim());
    setInput('');
  };

  if (!session) {
    return (
      <div className="panel chat-panel empty">
        <p>Select or create a session to start playing.</p>
      </div>
    );
  }

  return (
    <div className="panel chat-panel">
      <h2>Chat — {session.name}</h2>
      <div className="messages">
        {session.messages.map((msg, i) => (
          <div key={i} className={`message ${msg.sender}`}>
            <span className="sender">{msg.sender === 'ai' ? '🧙 GM' : '🧑 You'}:</span>{' '}
            {msg.text}
          </div>
        ))}
        <div ref={bottomRef} />
      </div>
      <form className="chat-form" onSubmit={handleSubmit}>
        <input
          className="chat-input"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="What do you do?"
        />
        <button className="btn-primary" type="submit" disabled={!input.trim()}>
          Send
        </button>
      </form>
    </div>
  );
}
