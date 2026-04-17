package main

import (
	"context"
	"fmt"

	"labs/charting"
	"labs/labs"
	neuron "labs/labs/1-neuron"
	"labs/labs/cluster"
	"labs/labs/forecasting"
	forecastinglinparab "labs/labs/forecasting-lin-parab"
	"labs/labs/holt"
	"labs/labs/optimizations"
	"labs/labs/polyapprox"
	"labs/labs/render"
	statslab "labs/labs/stats"
	"labs/labs/visualization"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	registry map[string]charting.LabProvider
	cache    *ResponseCache
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		registry: make(map[string]charting.LabProvider),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.cache = NewResponseCache()

	// Register lab factories using GenericProvider
	a.registry[labs.Lab2ID] = charting.NewProvider(labs.Config)
	a.registry[polyapprox.LabID] = charting.NewProvider(polyapprox.Config)
	a.registry[visualization.LabID] = charting.NewProvider(visualization.Config)
	a.registry[statslab.LabID] = charting.NewProvider(statslab.Config)
	a.registry[cluster.LabID] = charting.NewProvider(cluster.Config)
	a.registry[forecasting.LabID] = charting.NewProvider(forecasting.Config)
	a.registry[holt.LabID] = charting.NewProvider(holt.Config)
	a.registry[optimizations.LabID] = charting.NewProvider(optimizations.Config)
	a.registry[forecastinglinparab.LabID] = charting.NewProvider(forecastinglinparab.LinParabConfig)
	a.registry[neuron.LabID] = charting.NewProvider(neuron.Config)

	fmt.Printf("Registered %d labs\n", len(a.registry))
}

func (a *App) GetLabs() charting.GetLabsResponse {
	resp := charting.GetLabsResponse{}
	for labID := range a.registry {
		// Get metadata without initializing the full lab
		resp.Labs = append(resp.Labs, a.registry[labID].GetMetadata())
	}
	return resp
}

// GetLabConfig returns the configuration for a specific lab
func (a *App) GetLabConfig(labID string) (*charting.LabConfig, error) {
	provider, ok := a.registry[labID]
	if !ok {
		return nil, fmt.Errorf("lab %q not found", labID)
	}
	config := provider.GetConfig()
	return &config, nil
}

// Render processes the render request asynchronously to avoid blocking the GTK main loop
func (a *App) Render(req *charting.RenderRequest) {
	if req == nil {
		runtime.EventsEmit(a.ctx, "renderError", render.NewRenderError("request is nil"))
		return
	}
	if cachedRes, found := a.cache.GetResponse(req); found {
		fmt.Printf("Cache hit for lab %q, chart %q\n", req.LabID, req.ChartID)
		if render.IsRenderError(cachedRes.Error) {
			runtime.EventsEmit(a.ctx, "renderError", cachedRes)
		} else {
			runtime.EventsEmit(a.ctx, "renderComplete", cachedRes)
		}
		return
	}

	// Run rendering in a goroutine to prevent blocking the UI thread
	go func() {
		res := a.RenderSync(req)

		a.cache.StoreResponse(req, res)

		if render.IsRenderError(res.Error) {
			runtime.EventsEmit(a.ctx, "renderError", res)
		}

		runtime.EventsEmit(a.ctx, "renderComplete", res)
	}()
}

func (a *App) RenderSync(req *charting.RenderRequest) (res *charting.RenderResponse) {
	provider, ok := a.registry[req.LabID]
	if !ok {
		return res.NewErrorf("lab %q not found", req.LabID)
	}

	if provider == nil {
		return res.NewErrorf("lab %q not found", req.LabID)
	}

	res = provider.Render(req)
	if res.Error != nil {
		return res.NewErrorf("failed to render lab %q: %v", req.LabID, res.Error)
	}

	return res
}
