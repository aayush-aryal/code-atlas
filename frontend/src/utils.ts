import type { VisualEdge } from "./types";

//Traces the full path of outgoing calls
export const getDownstreamEdges = (
  startNodeId: string,
  allEdges: VisualEdge[]
) => {
  const activeEdgeIds = new Set<string>();
  const visitedNodes = new Set<string>();
  const queue = [startNodeId];

  while (queue.length > 0) {
    const currentNode = queue.shift()!;
    if (visitedNodes.has(currentNode)) continue;
    visitedNodes.add(currentNode);
    const outgoing = allEdges.filter((e) => e.source === currentNode);
    outgoing.forEach((edge) => {
      activeEdgeIds.add(edge.id);
      queue.push(edge.target);
    });
  }

  return activeEdgeIds;
};
