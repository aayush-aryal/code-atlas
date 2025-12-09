import { useEffect, useState } from "react";
import type { VisualGraphResponse } from "./types";
import { ReactFlow, Background, Controls } from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import ListCalls from "./components/ListCalls";
import { useSmartEdges } from "./hooks/useSmartEdges";
import { CodeDrawer } from "./components/CodeDrawer";

function App() {
  const [error, setError] = useState<string>("");
  const [visualResponse, setVisualResponse] =
    useState<VisualGraphResponse | null>(null);
  const [hoveredNode, setHoveredNode] = useState<string | null>(null);
  const [selectedNode, setSelectedNode] = useState<string | null>(null);

  useEffect(() => {
    async function loadGraph() {
      try {
        const resp = await fetch("http://localhost:8090/visual");
        if (!resp.ok) {
          setError("Something went wrong while loading data");
          return;
        }
        const data: VisualGraphResponse = await resp.json();
        setVisualResponse(data);
      } catch {
        setError("Failed to fetch graph");
      }
    }
    loadGraph();
  }, []);

  const edges = useSmartEdges({
    edges: visualResponse?.edges ?? [],
    hoveredNode,
  });

  // Calculate Hover Data
  const activeHoverNode = visualResponse?.nodes.find(
    (n) => n.id === hoveredNode
  );
  const incomingCalls = visualResponse?.edges
    .filter((edge) => edge.target === hoveredNode)
    .map((edge) => edge.source);
  const outgoingCalls = visualResponse?.edges
    .filter((edge) => edge.source === hoveredNode)
    .map((edge) => edge.target);

  // Calculate Selection Data
  const selectedNodeData = visualResponse?.nodes.find(
    (n) => n.id === selectedNode
  );

  return (
    <>
      <div>
        <p style={{ color: "red" }}>{error}</p>
      </div>
      <div
        style={{
          width: `100vw`,
          height: `100vh`,
          backgroundColor: "black",
          position: "relative",
        }}
      >
        {visualResponse && (
          <ReactFlow
            nodes={visualResponse.nodes}
            edges={edges}
            onNodeMouseEnter={(_, node) => setHoveredNode(node.id)}
            onNodeMouseLeave={() => setHoveredNode(null)}
            // Handle Selection
            onNodeClick={(_, node) => setSelectedNode(node.id)}
            // Handle Deselection (Clicking empty space)
            onPaneClick={() => setSelectedNode(null)}
            style={{ backgroundColor: "#111" }}
            fitView
            colorMode="dark"
            minZoom={0.1}
          >
            <Background color="#555" gap={20} />
            <Controls />
          </ReactFlow>
        )}

        {/* --- HOVER PANEL (Only show if nothing is selected) --- */}
        {activeHoverNode && !selectedNode && (
          <div style={panelStyle}>
            <h3 style={{ margin: "0 0 10px 0", color: "#63b3ed" }}>
              {activeHoverNode.data.label as unknown as string}
            </h3>
            <hr style={{ borderColor: "#4a5568", marginBottom: 10 }} />
            <div style={{ marginBottom: 15 }}>
              <strong style={{ color: "#a0aec0", fontSize: "0.8rem" }}>
                CALLED BY:
              </strong>
              <ListCalls calls={incomingCalls ?? []} />
            </div>
            <div>
              <strong style={{ color: "#a0aec0", fontSize: "0.8rem" }}>
                CALLS:
              </strong>
              <ListCalls calls={outgoingCalls ?? []} />
            </div>
          </div>
        )}

        {/* --- CODE DRAWER (Shows when a node is selected) --- */}
        {selectedNodeData && (
          <CodeDrawer
            node={selectedNodeData}
            onClose={() => setSelectedNode(null)}
          />
        )}
      </div>
    </>
  );
}

// --- Styles ---

const panelStyle: React.CSSProperties = {
  position: "absolute",
  top: 20,
  right: 20,
  width: 300,
  backgroundColor: "rgba(26, 32, 44, 0.95)", // Slightly more opaque
  border: "1px solid #4a5568",
  borderRadius: 8,
  padding: 20,
  color: "white",
  zIndex: 10,
  pointerEvents: "none", // Let clicks pass through if needed
  boxShadow: "0 4px 6px -1px rgba(0, 0, 0, 0.5)",
};

export default App;
