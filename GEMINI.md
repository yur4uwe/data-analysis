# Gemini Project Documentation: Go-Chart.js Adapter & Render Pipeline

This project is a Wails-based application that bridges a Go backend (performing heavy mathematical analysis) with a Chart.js-powered frontend.

## 🏗 Project Structure

- **`/charting`**: Core Go models for the visualization adapter.
  - `chart.go`: Defines `Chart`, `DataPoint`, and `MutableField`.
  - `dataset.go`: Defines the `Dataset` interface and concrete types:
    - `GridDataset`: For Scatter, Line, and Bubble charts (uses `DataPoint`).
    - `CategoricalDataset`: For Bar and Pie charts (uses `[]any` for values).
    - `HeatmapDataset`: For Heatmaps (uses `HeatmapPoint`).
  - `render.go`: Defines `RenderRequest`, `RenderResponse`, and `CachePolicy`.
  - `lab.go`: Defines `LabMetadata`, `LabConfig`, and `GenericProvider`.
- **`/labs`**: Scientific implementations.
  - Sub-packages include: `1-neuron`, `cluster`, `forecasting`, `forecasting-lin-parab`, `holt`, `neural-network`, `optimizations`, `polyapprox`, `stats`, `visualization`.
  - Each contains a `RenderFunc` and a `Config` (type `charting.LabConfig`).
- **`/analysis`**: Low-level mathematical functions (e.g., polynomial fitting).
- **`/uncsv`**: High-performance CSV decoder for `DataPoint` structures.
- **`app.go`**: The Wails Application entry point, handles lab registration and the event-driven bridge.
- **`cache.go`**: Implements `ResponseCache` to store and reuse `RenderResponse` objects based on MD5 hashes of `RenderRequest`.

## 🔄 The Render Flow

The application uses an **Asynchronous Event-Driven Pipeline** with caching to keep the UI responsive.

1.  **Interaction**: User modifies a `MutableField` (slider/select) or clicks **Rerender**.
2.  **JS Request**: `fetchChartData()` gathers all current UI variables into a `RenderRequest`.
3.  **Go Invocation**: `App.Render(req)` is called.
    - It immediately returns to the JS side.
    - **Cache Check**: If an identical request exists in `ResponseCache`, the cached response is emitted immediately.
    - **Execution**: If not cached, a goroutine is spawned to execute the scientific logic.
4.  **Backend Execution**: The `LabProvider` (usually `GenericProvider`) locates the specific `RenderFunc`.
    - Data is fetched (usually from `/data/*.csv` or generated).
    - Math is performed.
    - A `charting.Chart` object is populated with datasets.
5.  **Event Emit**: Once finished, the backend stores the result in cache and emits a `renderComplete` event (or `renderError`) via Wails Runtime.
6.  **JS Receipt**: `EventsOn("renderComplete")` in `main.ts` receives the `RenderResponse`.
7.  **Draw**: Based on the `chart.type`, the frontend calls either `renderChartInto` or `renderMultiChart`.

## 📊 Data Display Modes

### 1. Standard Mode (`type: "line" | "bar" | "scatter" ...`)
Renders a single Chart.js instance into the main container.
- **Labels**: Uses global `chart.labels` for the X-axis.
- **Datasets**: All datasets in `chart.datasets` are rendered on the same axes.

### 2. Multi-Chart Mode (`type: "multi-bar" | "multi-scatter" ...`)
Designed for comparisons (e.g., Silhouette plots where each cluster needs its own scale).
- **Logic**: The frontend strips the `multi-` prefix and creates a **Grid Layout**.
- **Synthetic Charts**: For every dataset in the original chart, a "Synthetic Chart" is created.
- **Isolation**: Each cluster/dataset gets its own individual canvas and coordinate system.

## 🛠 Adapter Safety & Best Practices

- **Typed Datasets**: Always instantiate the correct dataset type for your chart:
  - Use `GridDataset` when you have `{x, y}` coordinates.
  - Use `CategoricalDataset` when you have indexed values corresponding to `chart.Labels`.
- **Labels & length**: Always use `Array.isArray(labels) && labels.length > 0` before processing. The frontend is defensive against `null` labels from the backend.
- **Surgical Updates**: Use helper methods on the `Chart` object to populate data safely:
  - `UpdatePointsForDataset(id, x, y)`: Converts raw float slices to `DataPoint` objects.
  - `UpdateDataForDataset(id, data)`: For `CategoricalDataset` using `[]any`.
  - `GenerateLabels(precision)`: Automatically creates X-axis labels from the dataset with the most points.
- **Datalabels Plugin**: 
  - For **Scatter**: Uses `DataLabels` for per-point tooltips/labels.
  - For **Pie/Doughnut**: Automatically calculates and displays percentages.
- **Colors**: Use `charting.ToColor("#hex")` to ensure color strings are wrapped in the `Color` type.
- **Responsive Height**: All charts use `maintainAspectRatio: false`. The parent container must have a `min-height` (defined in `style.css`).

## 🧬 Architecture & Data Integrity

### 1. Handling Missing Data (The `null` pattern)
- **`[]any` over `[]float64`**: The `CategoricalDataset.Data` field uses `[]any` to allow `nil` values.
- **Visual Integrity**: Always use `nil` for indices where data is unavailable. This prevents Chart.js from jumping to `0`.

### 2. Variable Registration & Metadata
- **Static Declaration**: All `MutableField`s must be declared in the `Chart` template during package initialization.
- **Discovery**: `charting.NewLabConfig` scans templates to generate `LabMetadata`. Fields added dynamically during `Render` will not be visible in the UI sidebars.
- **Graph-Level Stats**: Use the `GraphVariables` map on `ChartMetadata` to display read-only statistics specific to a dataset.

### 3. State Management
- **Immutable Templates**: Never modify the global `Chart` template. Use `charting.CopyChart(Template)` at the start of every `RenderFunc`. `CopyChart` performs a deep copy of configuration while resetting data fields.

### 4. Caching Policy
- **MD5 Keys**: `ResponseCache` uses an MD5 hash of the JSON-marshaled `RenderRequest` as a lookup key.
- **Policy Control**: Use `CachePolicy` in `RenderResponse` to control behavior (`CachePolicyDontCache`, `CachePolicyWithExpiration`).
