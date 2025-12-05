package codebase

import (
	"fmt"
	"math"
	"sort"
)

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

const(
	MinRingRadius=600.0
	ArcPerFile=500.0

	//
	FileNodeDiameter=150
	FuncNodeDiameter=70
	OrbitRadius=200
)

func (p *Project)ComputeVisualGraph()VisualGraphResponse{
	var nodes []VisualNode 
	var edges []VisualEdge 

	var filePaths []string 
	for path:= range p.Graph{
		filePaths=append(filePaths,path)
	}

	// sort files to maintain same order
	sort.Strings(filePaths)
	totalFiles:=len(filePaths)

	mainRingRadius:=float64(totalFiles*ArcPerFile)/(2*math.Pi)
	if MinRingRadius<mainRingRadius{
		// the radius is smaller 
		mainRingRadius=MinRingRadius
	}

	centerX,centerY:=0,0

	// files loop (main ring of files)
	for i,filePath:=range filePaths{
		fileNode:=p.Graph[filePath]

		// angle for this file in radians
		fileAngle:=(2*math.Pi/float64(totalFiles)*float64(i))

		// convert to x and y coordinates
		fileX := centerX + int(mainRingRadius*math.Cos(fileAngle)) 		
		fileY := centerY + int(mainRingRadius*math.Sin(fileAngle))
		nodes = append(nodes, VisualNode{
			ID:       filePath,
			Position: Position{X: fileX, Y: fileY},
			Data:     map[string]interface{}{"label": filePath},
			Style: map[string]interface{}{
				"width":           FileNodeDiameter,
				"height":          FileNodeDiameter,
				"borderRadius":    "50%",
				"backgroundColor": "rgba(66, 153, 225, 0.2)",
				"color":           "#90cdf4",
				"border":          "2px dashed #4299e1",
				"display":         "flex",
				"alignItems":      "center",
				"justifyContent":  "center",
				"textAlign":       "center",
				"fontWeight":      "bold",
				"fontSize":        "14px",
				"zIndex":          10,
			},
		})
		numFun:=len(fileNode.Functions)
		for j,fn:=range fileNode.Functions{
			funcID := fmt.Sprintf("%s::%s", filePath, fn.Name)
			// angle of this func related to file center
			funcAngle:=(2*math.Pi/float64(numFun))*float64(j)

			funcX:=fileX+int(FuncNodeDiameter*math.Cos(funcAngle))
			funcY:=fileY+int(FuncNodeDiameter*math.Sin(funcAngle))
			nodes = append(nodes, VisualNode{
				ID:       funcID,
				Position: Position{X: funcX, Y: funcY},
				Data:     map[string]interface{}{"label": fn.Name},
				Style: map[string]interface{}{
					"width":           FuncNodeDiameter,
					"height":          FuncNodeDiameter,
					"borderRadius":    "50%",
					"backgroundColor": "#1a202c",
					"color":           "#e2e8f0",
					"border":          "2px solid #718096",
					"display":         "flex",
					"alignItems":      "center",
					"justifyContent":  "center",
					"fontSize":        "11px",
					"zIndex":          20,
					"boxShadow":       "0 4px 6px -1px rgba(0, 0, 0, 0.5)",
				},
			})
			// create edges
			for _,callName:=range fn.Calls{
				targets,ok:=p.FunctionTable[callName]
				if ok && len(targets)>0{
					targetFile:=targets[0].Path
					targetID:=fmt.Sprintf("%s::%s",targetFile,callName)

					if targetID!=funcID{
						edges=append(edges, VisualEdge{
							ID:fmt.Sprintf("%s->%s", funcID, targetID),
							Source: funcID,
							Target: targetID,
							Animated: true,
							Style:map[string]interface{}{"stroke": "#a0aec0", "strokeWidth": 2, "opacity": 0.6},
						})
					}
				}
			}
		}
	}
		return VisualGraphResponse{Nodes:nodes,Edges:edges}
	}

