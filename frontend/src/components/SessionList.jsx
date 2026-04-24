export default function SessionList({ sessions, activeSessionId, onSelect, onCreate }) {
  const sessionList = Object.values(sessions);
  const canCreate = sessionList.length < 3;

  return (
    <div className="panel session-panel">
      <h2>Sessions</h2>
      <button className="btn-primary" onClick={onCreate} disabled={!canCreate}>
        + New Session
      </button>
      <div className="session-list">
        {sessionList.map((s) => (
          <button
            key={s.id}
            className={`session-item${activeSessionId === s.id ? ' active' : ''}`}
            onClick={() => onSelect(s.id)}
          >
            {s.name}
          </button>
        ))}
      </div>
      {!canCreate && <p className="hint">Max 3 sessions reached.</p>}
    </div>
  );
}
