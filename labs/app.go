package main

import (
	"context"
	"fmt"

	"labs/labs"
	"labs/labs/common"
	"labs/labs/polyapprox"
	"labs/labs/render"
	"labs/labs/stats"
	"labs/labs/visualization"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// LabFactory is a function that creates a lab provider on demand
type LabFactory func() common.LabProvider

// App struct
type App struct {
	ctx      context.Context
	registry map[string]common.LabProvider
	cache    *ResponseCache
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		registry: make(map[string]common.LabProvider),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.cache = NewResponseCache()

	// Register lab factories (no instantiation yet)
	a.registry[labs.Lab2ID] = labs.NewLab2()
	a.registry[polyapprox.LabID] = labs.NewLab3()
	a.registry[visualization.LabID] = labs.NewLab4()
	a.registry[stats.LabID] = labs.NewLab5()

	fmt.Printf("Registered %d labs (lazy initialization)\n", len(a.registry))
}

func (a *App) GetLabs() common.GetLabsResponse {
	resp := common.GetLabsResponse{}
	for labID := range a.registry {
		// Get metadata without initializing the full lab
		resp.Labs = append(resp.Labs, a.registry[labID].GetMetadata())
	}
	return resp
}

// GetLabConfig returns the configuration for a specific lab
func (a *App) GetLabConfig(labID string) (*common.LabConfig, error) {
	provider, ok := a.registry[labID]
	if !ok {
		return nil, fmt.Errorf("lab %q not found", labID)
	}
	config := provider.GetConfig()
	return &config, nil
}

// Render processes the render request asynchronously to avoid blocking the GTK main loop
func (a *App) Render(req *common.RenderRequest) {
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

		a.cache.StoreResponse(req, &res)

		if render.IsRenderError(res.Error) {
			runtime.EventsEmit(a.ctx, "renderError", res)
		}

		runtime.EventsEmit(a.ctx, "renderComplete", res)
	}()
}

// RenderSync renders synchronously and returns the response directly
// This ensures RenderResponse is exported to TypeScript
func (a *App) RenderSync(req *common.RenderRequest) common.RenderResponse {
	provider, ok := a.registry[req.LabID]
	if !ok {
		return common.RenderResponse{
			Error: render.NewRenderError(fmt.Sprintf("lab %q not found", req.LabID)),
		}
	}

	result := provider.Render(req)
	if result.Error != nil {
		return common.RenderResponse{
			Error: render.NewRenderError(fmt.Sprintf("failed to render lab %q: %v", req.LabID, result.Error)),
		}
	}

	return common.RenderResponse{
		Charts: result.Charts,
	}
}
