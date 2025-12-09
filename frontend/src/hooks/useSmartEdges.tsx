import { useMemo } from "react";
import { getDownstreamEdges } from "../utils";
import type { VisualEdge } from "../types";

interface UseSmartEdgesProps {
  edges: VisualEdge[];
  hoveredNode: string | null;
}

export function useSmartEdges({ edges, hoveredNode }: UseSmartEdgesProps) {
  return useMemo(() => {
    if (!edges) return [];

    const traceSet = hoveredNode
      ? getDownstreamEdges(hoveredNode, edges)
      : new Set();
    return edges.map((edge) => {
      if (!hoveredNode) {
        return {
          ...edge,
          animated: false,
          style: { ...edge.style, opacity: 0.1, stroke: "#555" },
        };
      }

      // either directly connected OR part of trace
      const isSourceOrTarget =
        edge.source === hoveredNode || edge.target === hoveredNode;
      const isPartOfTrace = traceSet.has(edge.id);

      // 4. Return Highlighted Style
      if (isSourceOrTarget || isPartOfTrace) {
        return {
          ...edge,
          animated: true,
          zIndex: 999,
          style: {
            ...edge.style,
            opacity: 1,
            stroke: "#63b3ed",
            strokeWidth: 2,
          },
        };
      }
      return {
        ...edge,
        animated: false,
        zIndex: 0,
        style: {
          ...edge.style,
          opacity: 0.05,
          stroke: "#333",
        },
      };
    });
  }, [edges, hoveredNode]);
}
