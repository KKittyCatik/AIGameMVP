import { useState, useEffect } from 'react';

export default function PromptPanel({ globalPrompt, sessionPrompt, onUpdateGlobal, onUpdateSession, hasActiveSession }) {
  const [globalInput, setGlobalInput] = useState(globalPrompt);
  const [sessionInput, setSessionInput] = useState(sessionPrompt);

  // Sync textarea values when backend state changes
  useEffect(() => { setGlobalInput(globalPrompt); }, [globalPrompt]);
  useEffect(() => { setSessionInput(sessionPrompt); }, [sessionPrompt]);

  return (
    <div className="prompt-panel">
      <div className="prompt-section">
        <label className="prompt-label">🌐 Global Prompt</label>
        <textarea
          className="prompt-textarea"
          value={globalInput}
          onChange={(e) => setGlobalInput(e.target.value)}
          rows={3}
          placeholder="System prompt applied to all sessions…"
        />
        <button className="btn-primary prompt-btn" onClick={() => onUpdateGlobal(globalInput)}>
          Apply
        </button>
      </div>
      <div className="prompt-section">
        <label className="prompt-label">🎭 Session Prompt</label>
        <textarea
          className="prompt-textarea"
          value={sessionInput}
          onChange={(e) => setSessionInput(e.target.value)}
          rows={2}
          disabled={!hasActiveSession}
          placeholder={hasActiveSession ? 'Additional prompt for this session…' : 'Select a session first'}
        />
        <button
          className="btn-primary prompt-btn"
          onClick={() => onUpdateSession(sessionInput)}
          disabled={!hasActiveSession}
        >
          Apply
        </button>
      </div>
    </div>
  );
}
