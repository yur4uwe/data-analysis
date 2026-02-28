package labs

import (
	"labs/labs/common"
	"labs/labs/stats"
)

type Lab5Provider struct{}

var _ common.LabProvider = Lab5Provider{}

func NewLab5() *Lab5Provider {
	return &Lab5Provider{}
}

func (lp Lab5Provider) GetMetadata() common.LabMetadata {
	return stats.Metadata
}

func (lp Lab5Provider) GetConfig() common.LabConfig {
	return stats.Config
}

func (lp Lab5Provider) Render(req *common.RenderRequest) *common.RenderResponse {
	res := &common.RenderResponse{}
	if req == nil {
		return res.NewError("request is nil")
	}

	switch req.ChartID {
	case stats.ErrorAnalysisChartID:
		return stats.RenderErrorAnalysis(req)
	default:
		return res.NewErrorf("unrecognized chart ID: %s", req.ChartID)
	}
}
