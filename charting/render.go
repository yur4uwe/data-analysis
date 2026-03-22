package charting

import (
	"labs/labs/render"
	"strings"
	"time"
)

// RenderRequest is sent from frontend when user requests new data
type RenderRequest struct {
	LabID             string                        `json:"LabID"`
	ChartID           string                        `json:"ChartID"`
	ChartVariables    map[string]float64            `json:"ChartVariables"`    // Global field values
	GraphVariables    map[string]map[string]float64 `json:"GraphVariables"`    // Graph-specific params (graphID -> params)
	DatasetVisibility map[string]bool               `json:"DatasetVisibility"` // what datasets are visible on the graph (to alleviate the load)
}

func (rr *RenderRequest) GetChartVariable(chartId string, variableId string) (float64, bool) {
	if chartId == "" || variableId == "" {
		return 0, false
	}
	if val, ok := rr.ChartVariables[strings.Join([]string{chartId, variableId}, "-")]; ok {
		return val, true
	}
	return 0, false
}

func (rr *RenderRequest) GetGraphVariable(chartId string, graphId string, variableId string) (float64, bool) {
	if chartId == "" || graphId == "" || variableId == "" {
		return 0, false
	}
	if graphVars, ok := rr.GraphVariables[chartId]; ok {
		if val, ok := graphVars[strings.Join([]string{graphId, variableId}, "-")]; ok {
			return val, true
		}
	}
	return 0, false
}

type CachePolicy int

const (
	CachePolicyNone CachePolicy = iota
	CachePolicyDontCache
	CachePolicyCacheOnly
	CachePolicyWithExpiration
)

// RenderResponse contains the updated chart data
type RenderResponse struct {
	Charts map[string]Chart `json:"charts"`
	Error  error            `json:"error,omitempty"`
	// backend only fields for caching should not be sent to frontend
	CachePolicy  CachePolicy `json:"-"`
	CachedAt     int32       `json:"-"`
	ExpirationMS int32       `json:"-"`
}

func NewRenderResponse() *RenderResponse {
	return &RenderResponse{
		Charts: make(map[string]Chart),
	}
}

func (rr *RenderResponse) AddChart(chartId string, c *Chart) {
	if rr == nil {
		rr = NewRenderResponse()
	}
	rr.Charts[chartId] = *c
}

func (rr *RenderResponse) NewErrorf(format string, args ...any) *RenderResponse {
	if rr == nil {
		rr = &RenderResponse{}
	}
	rr.Error = render.NewRenderErrorf(format, args...)
	return rr
}

func (rr *RenderResponse) NewError(message string) *RenderResponse {
	if rr == nil {
		rr = &RenderResponse{}
	}
	rr.CachePolicy = CachePolicyDontCache
	rr.Error = render.NewRenderError(message)
	return rr
}

func (rr *RenderResponse) IsExpired() bool {
	if rr.CachePolicy != CachePolicyWithExpiration {
		return false
	}
	currentTime := time.Now().UnixMilli()
	return currentTime > int64(rr.CachedAt)+int64(rr.ExpirationMS)
}
