export default function ListCalls({ calls }: { calls: string[] }) {
  return calls && calls.length > 0 ? (
    <ul
      style={{
        paddingLeft: 20,
        margin: "5px 0",
        fontSize: "0.9rem",
      }}
    >
      {calls.map((callerId: string) => (
        <li key={callerId} style={{ color: "#cbd5e0" }}>
          {callerId.split("::").pop()}
        </li>
      ))}
    </ul>
  ) : (
    <div
      style={{
        fontStyle: "italic",
        color: "#718096",
        fontSize: "0.9rem",
      }}
    >
      No callers (Entry Point?)
    </div>
  );
}
