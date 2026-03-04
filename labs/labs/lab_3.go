package labs

import (
	"labs/labs/common"
	"labs/labs/polyapprox"
	"labs/labs/render"
)

type Lab3Provider struct{}

var _ common.LabProvider = Lab3Provider{}

func NewLab3() *Lab3Provider {
	return &Lab3Provider{}
}

func (lp Lab3Provider) GetMetadata() common.LabMetadata {
	return polyapprox.Metadata
}

func (lp Lab3Provider) Render(req *common.RenderRequest) *common.RenderResponse {
	if req == nil {
		return &common.RenderResponse{Error: render.NewRenderError("empty render request")}
	}

	switch req.ChartID {
	case polyapprox.RandomFitsID:
		return polyapprox.RenderRandomFits(req)
	case polyapprox.SampleDataID:
		return polyapprox.RenderSampleData(req)
	case polyapprox.RandomMSEID:
		return polyapprox.RenderRandomPolynomialMSE(req)
	case polyapprox.SampleMSEID:
		return polyapprox.RenderSamplePolynomialMSE(req)
	default:
		return &common.RenderResponse{Error: render.NewRenderErrorf("unrecognised Chart: %s", req.ChartID)}
	}
}

func (lp Lab3Provider) GetConfig() common.LabConfig {
	return common.LabConfig{
		Lab: polyapprox.Metadata,
		Charts: map[string]*common.Chart{
			polyapprox.RandomFitsID: &polyapprox.RandomFitsChart,
			polyapprox.SampleDataID: &polyapprox.SampleDataChart,
			polyapprox.RandomMSEID:  &polyapprox.RandomMSEChart,
			polyapprox.SampleMSEID:  &polyapprox.SampleMSEChart,
		},
	}
}
