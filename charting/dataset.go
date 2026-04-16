package charting

import (
	"fmt"
)

type HeatmapPoint struct {
	DataPoint
	Value *float64 `json:"v"`
}

type Dataset interface {
	UpdateData([]any) // panics if data is not of the correct type
	UpdateLabel(string)
	UpdateVariableLabel(id string, label string)
	GetData() []any
	GetFields() []MutableField
	Copy() Dataset
	GetBase() *BaseDataset
	GetType() GraphType
}

type BaseDataset struct {
	Label          string         `json:"label"`
	Type           GraphType      `json:"type,omitempty"`
	BorderColor    Color          `json:"borderColor"`
	BorderWidth    int            `json:"borderWidth"`
	Hidden         bool           `json:"hidden"`
	Togglable      bool           `json:"togglable"`
	DataLabels     []string       `json:"dataLabels,omitempty"`
	GraphVariables []MutableField `json:"fields,omitempty"`
}

func (bd *BaseDataset) GetType() GraphType {
	return bd.Type
}

func (bd *BaseDataset) GetBase() *BaseDataset {
	return bd
}

func (bd *BaseDataset) GetFields() []MutableField {
	return bd.GraphVariables
}

func (bd *BaseDataset) UpdateVariableLabel(id string, label string) {
	for i := range bd.GraphVariables {
		if bd.GraphVariables[i].ID == id {
			bd.GraphVariables[i].Label = label
			return
		}
	}
	panic(fmt.Sprintf("variable ID %s not found in dataset", id))
}

func (bd *BaseDataset) UpdateLabel(new_label string) {
	bd.Label = new_label
}

func (bd *BaseDataset) CopyBase() BaseDataset {
	newBD := *bd
	if bd.DataLabels != nil {
		newBD.DataLabels = make([]string, len(bd.DataLabels))
		copy(newBD.DataLabels, bd.DataLabels)
	}
	if bd.GraphVariables != nil {
		newBD.GraphVariables = make([]MutableField, len(bd.GraphVariables))
		copy(newBD.GraphVariables, bd.GraphVariables)
	}
	return newBD
}

type GridDataset struct {
	BaseDataset
	Data            []DataPoint `json:"data,omitempty"`
	BackgroundColor Color       `json:"backgroundColor,omitempty"`
	PointRadius     int         `json:"pointRadius"`
	PointStyle      string      `json:"pointStyle,omitempty"`
	HideLine        bool        `json:"hideLine"`
}

var _ Dataset = &GridDataset{}

func (gd *GridDataset) UpdateData(data []any) {
	gd.Data = make([]DataPoint, len(data))
	for i, v := range data {
		if v == nil {
			// If Y is nullable, we can't just have a nil DataPoint in a value slice.
			// We should probably skip or provide a DataPoint with nil Y.
			// Given the user's intent, we assume 'v' should be a DataPoint or convertible.
			gd.Data[i] = DataPoint{X: float64(i), Y: nil}
			continue
		}

		if p, ok := v.(*DataPoint); ok {
			gd.Data[i] = *p
		} else if p, ok := v.(DataPoint); ok {
			gd.Data[i] = p
		} else {
			panic(fmt.Errorf("invalid data type for GridDataset: expected *DataPoint or DataPoint, got %T", v))
		}
	}
}

func (gd *GridDataset) GetData() []any {
	res := make([]any, len(gd.Data))
	for i, v := range gd.Data {
		res[i] = v
	}
	return res
}

func (gd *GridDataset) Copy() Dataset {
	newGD := *gd
	newGD.BaseDataset = gd.CopyBase()
	if gd.Data != nil {
		newGD.Data = make([]DataPoint, len(gd.Data))
		copy(newGD.Data, gd.Data)
	}
	return &newGD
}

type CategoricalDataset struct {
	BaseDataset
	Data            []*float64 `json:"data,omitempty"`
	BackgroundColor []Color    `json:"backgroundColor,omitempty"`
}

var _ Dataset = &CategoricalDataset{}

func (cd *CategoricalDataset) UpdateData(data []any) {
	cd.Data = make([]*float64, len(data))
	for i, v := range data {
		if v == nil {
			cd.Data[i] = nil
			continue
		}
		if f, ok := v.(*float64); ok {
			cd.Data[i] = f
		} else if f, ok := v.(float64); ok {
			val := f
			cd.Data[i] = &val
		} else {
			panic("invalid data type for CategoricalDataset: expected *float64 or float64")
		}
	}
}

func (cd *CategoricalDataset) GetData() []any {
	res := make([]any, len(cd.Data))
	for i, v := range cd.Data {
		res[i] = v
	}
	return res
}

func (cd *CategoricalDataset) Copy() Dataset {
	newCD := *cd
	newCD.BaseDataset = cd.CopyBase()
	if cd.Data != nil {
		newCD.Data = make([]*float64, len(cd.Data))
		copy(newCD.Data, cd.Data)
	}
	if cd.BackgroundColor != nil {
		newCD.BackgroundColor = make([]Color, len(cd.BackgroundColor))
		copy(newCD.BackgroundColor, cd.BackgroundColor)
	}
	return &newCD
}

type HeatmapDataset struct {
	BaseDataset
	Data            []HeatmapPoint `json:"pointData,omitempty"`
	BackgroundColor []Color        `json:"backgroundColor,omitempty"`
}

var _ Dataset = &HeatmapDataset{}

func (hd *HeatmapDataset) UpdateData(data []any) {
	hd.Data = make([]HeatmapPoint, len(data))
	for i, v := range data {
		if v == nil {
			hd.Data[i] = HeatmapPoint{DataPoint: DataPoint{X: float64(i), Y: nil}, Value: nil}
			continue
		}
		if p, ok := v.(*HeatmapPoint); ok {
			hd.Data[i] = *p
		} else if p, ok := v.(HeatmapPoint); ok {
			hd.Data[i] = p
		} else {
			panic("invalid data type for HeatmapDataset: expected *HeatmapPoint or HeatmapPoint")
		}
	}
}

func (hd *HeatmapDataset) GetData() []any {
	res := make([]any, len(hd.Data))
	for i, v := range hd.Data {
		res[i] = v
	}
	return res
}

func (hd *HeatmapDataset) Copy() Dataset {
	newHD := *hd
	newHD.BaseDataset = hd.CopyBase()
	if hd.Data != nil {
		newHD.Data = make([]HeatmapPoint, len(hd.Data))
		copy(newHD.Data, hd.Data)
	}
	if hd.BackgroundColor != nil {
		newHD.BackgroundColor = make([]Color, len(hd.BackgroundColor))
		copy(newHD.BackgroundColor, hd.BackgroundColor)
	}
	return &newHD
}
