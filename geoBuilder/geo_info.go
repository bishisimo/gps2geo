//@Time : 2020/9/23 上午9:43
//@Author : bishisimo
package geoBuilder

import (
	"encoding/json"
	"fmt"
	"github.com/bishisimo/errlog"
	"os"
)

type Info struct {
	Type string    `json:"type"`
	Features []PIO `json:"features"`
}

func NewInfo() *Info {
	return &Info{
		Type:     "",
		Features: make([]PIO, 0),
	}
}
type PIO struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Parent struct {
	Adcode int `json:"adcode"`
}

func NewParent() *Parent {
	return &Parent{
		Adcode:0,
	}
}
type Properties struct {
	Adcode          int    `json:"adcode"`
	Name            string `json:"name"`
	Center          Point  `json:"center"`
	Centroid        Point  `json:"centroid"`
	ChildrenNum     int    `json:"childrenNum"`
	Level           string `json:"level"`
	SubFeatureIndex int    `json:"subFeatureIndex"`
	Acroutes        Point  `json:"acroutes"`
	Parent          Parent `json:"type"`
}

func NewProperties() *Properties {
	return &Properties{
		Adcode:          0,
		Name:            "",
		Center:          Point{},
		Centroid:        Point{},
		ChildrenNum:     0,
		Level:           "",
		SubFeatureIndex: 0,
		Acroutes:        Point{},
		Parent:          *NewParent(),
	}
}
type Geometry struct {
	Type string             `json:"type"`
	Coordinates [][][]Point `json:"coordinates"`
}

func NewGeometry() *Geometry {
	return &Geometry{
		Type: "",
		Coordinates: make([][][]Point,0),
	}
}

type Point =[2]float64

func  NewInfoFromJsonFile(filePath string) *Info {
	fp,err:=os.Open(filePath)
	if errlog.Debug(err) {
		fmt.Println("NewInfoFromJsonFile io err")
	}
	if fp != nil {
		defer fp.Close()
	}
	info:= NewInfo()
	decoder:=json.NewDecoder(fp)
	_ = decoder.Decode(&info)
	return info
}