package labs

import (
	"labs/charting"
	"labs/labs/forecasting"
)

const ()

type Lab7Provider struct{}

var _ charting.LabProvider = Lab7Provider{}

func NewLab7() *Lab7Provider {
	return &Lab7Provider{}
}

func (lp Lab7Provider) GetMetadata() charting.LabMetadata {
	return forecasting.Metadata
}

func (lp Lab7Provider) GetConfig() charting.LabConfig {
	return forecasting.Config
}

func (lp Lab7Provider) Render(req *charting.RenderRequest) *charting.RenderResponse {
	res := &charting.RenderResponse{}
	if req == nil {
		return res.NewError("request is nil")
	}

	chart, ok := forecasting.Config.Charts[req.ChartID]
	if !ok {
		return res.NewErrorf("chart with id %q not found", req.ChartID)
	}

	return chart.RenderFunc(req)
}
