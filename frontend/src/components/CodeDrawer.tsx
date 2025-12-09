import React from "react";
import { Code } from "./Code";
import type { Node } from "@xyflow/react";

interface CodeDrawerProps {
  node: Node | null | undefined;
  onClose: () => void;
}

export function CodeDrawer({ node, onClose }: CodeDrawerProps) {
  // If no node is selected, don't render anything
  if (!node) return null;

  return (
    <div style={drawerStyle}>
      {/* Header Section */}
      <div style={headerStyle}>
        <div>
          <span style={labelStyle}>Function Source</span>
          <h2 style={titleStyle}>{node.data.label as string}</h2>
        </div>
        <button onClick={onClose} style={closeButtonStyle}>
          ✕
        </button>
      </div>

      {/* Content Section */}
      <div style={{ flex: 1, overflow: "hidden" }}>
        {" "}
        <Code functionName={node.data.label as string} />
      </div>
    </div>
  );
}

// --- Clean Styles ---

const drawerStyle: React.CSSProperties = {
  position: "absolute",
  top: 0,
  right: 0,
  width: "40vw",
  minWidth: "450px",
  height: "100%",
  backgroundColor: "#171923",
  borderLeft: "1px solid #4a5568",
  padding: 25,
  color: "white",
  zIndex: 20,
  boxShadow: "-5px 0 20px rgba(0,0,0,0.7)",
  display: "flex",
  flexDirection: "column",
};

const headerStyle: React.CSSProperties = {
  display: "flex",
  justifyContent: "space-between",
  alignItems: "center",
  marginBottom: 20,
};

const labelStyle: React.CSSProperties = {
  color: "#a0aec0",
  fontSize: "0.8rem",
  textTransform: "uppercase",
  letterSpacing: "0.5px",
};

const titleStyle: React.CSSProperties = {
  margin: "5px 0 0 0",
  color: "#fff",
  fontSize: "1.2rem",
};

const closeButtonStyle: React.CSSProperties = {
  background: "transparent",
  border: "1px solid #4a5568",
  color: "#a0aec0",
  cursor: "pointer",
  fontSize: "1rem",
  borderRadius: "4px",
  width: "32px",
  height: "32px",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  transition: "background 0.2s",
};
