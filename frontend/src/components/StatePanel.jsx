export default function StatePanel({ session }) {
  if (!session) {
    return (
      <div className="panel state-panel empty">
        <h2>Player State</h2>
        <p>No active session.</p>
      </div>
    );
  }

  const { inventory = [], effects = [] } = session.state;

  return (
    <div className="panel state-panel">
      <h2>Player State</h2>
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
      <div className="effects">
        <h3>Buffs &amp; Debuffs</h3>
        {effects.length === 0 ? (
          <p className="hint">None</p>
        ) : (
          <ul>
            {effects.map((effect, i) => (
              <li key={i}>✨ {effect}</li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
