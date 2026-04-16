---
name: wails-chart-adapter
description: Expert guidance for developing the Go-to-Chart.js rendering pipeline. Use when adding new scientific labs, modifying chart datasets, or updating frontend-backend visualization variables.
---

# Wails Chart Adapter

This skill provides expert patterns for the Go-to-Chart.js adapter using the GenericProvider architecture.

## 📡 Core API Contract
The communication between Go and JS happens via `RenderRequest` and `RenderResponse`.
- **RenderRequest**: Contains `LabID`, `ChartID`, and `MutableField` values.
- **RenderResponse**: A map of `Chart` objects keyed by ID.
See [references/api-contract.md](references/api-contract.md) for full JSON schemas.

## 📊 Dataset Implementation
Choose the correct dataset type based on the visualization:
1. **GridDataset**: XY Scatter, Line, or Bubble charts.
2. **CategoricalDataset**: Bar, Pie, or Doughnut charts (using labels).
3. **HeatmapDataset**: 3D data visualizations.
See [references/dataset-guide.md](references/dataset-guide.md) for data structure details.

## 🛠 Adding a New Lab (Simplified)
Follow the standardized 3-step workflow:
1. **Implementation**: Create a package in `labs/` with a `Config` using `charting.NewLabConfig`.
2. **Logic**: Implement a `Render` function and assign it to the chart's `RenderFunc` in an `init()` block.
3. **Registration**: Add `a.registry[labID] = charting.NewProvider(labPackage.Config)` to `app.go`.
See [references/workflow.md](references/workflow.md) for a detailed checklist.

## ⚙️ Mutable Fields (UI Controls)
- `ControlRange`: standard range (slider) input.
- `ControlNumber`: numeric input field.
- `ControlSelect`: requires `Options` slice (`[]string`). Frontend sends the **index** as `float64`.
- `ControlNoControl`: Read-only statistics. Update by modifying `Label` on the chart copy during `Render`.
