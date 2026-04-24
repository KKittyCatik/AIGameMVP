export default function StatePanel({ session }) {
  if (!session) {
    return (
      <div className="panel state-panel empty">
        <h2>Player State</h2>
        <p>No active session.</p>
      </div>
    );
  }

  const { hp, inventory } = session.state;
  const hpPercent = Math.max(0, Math.min(100, hp));

  return (
    <div className="panel state-panel">
      <h2>Player State</h2>
      <div className="stat">
        <span>HP</span>
        <div className="hp-bar-bg">
          <div
            className="hp-bar-fill"
            style={{ width: `${hpPercent}%`, backgroundColor: hp > 50 ? '#a6e3a1' : hp > 25 ? '#f9e2af' : '#f38ba8' }}
          />
        </div>
        <span>{hp} / 100</span>
      </div>
      <div className="inventory">
        <h3>Inventory</h3>
        {inventory.length === 0 ? (
          <p className="hint">Empty</p>
        ) : (
          <ul>
            {inventory.map((item, i) => (
              <li key={i}>🎒 {item}</li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
