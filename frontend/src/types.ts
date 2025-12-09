/*
type VisualNode struct{
	ID string `json:"id"`
	Type string `json:"type,omitempty"`
	Position Position `json:"position"`
	Data map[string]interface{} `json:"data"`
	Style map[string]interface{} `json:"style,omitempty"`
}


type Position struct{
	X int `json:"x"`
	Y int `json:"y"`
}

type VisualEdge struct{
	ID string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Animated bool `json:"animated"`
	Style map[string]interface{} `json:"style,omitempty"`


}

type VisualGraphResponse struct{
	Nodes []VisualNode `json:"nodes"`
	Edges []VisualEdge `json:"edges"`

}
*/

export type VisualNode = {
  id: string;
  type: string;
  position: Position;
  data: Record<string, object>;
  Style: Record<string, object>;
};

export type Position = {
  x: number;
  y: number;
};

export type VisualEdge = {
  id: string;
  source: string;
  target: string;
  animated: boolean;
  style: Record<string, object>;
};

export type VisualGraphResponse = {
  nodes: VisualNode[];
  edges: VisualEdge[];
};
