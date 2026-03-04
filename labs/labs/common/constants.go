package common

// FieldControl specifies what UI element should render a field
// Values correspond directly to HTML input types
type FieldControl string

const (
	ControlRange     FieldControl = "range"     // Slider input
	ControlNumber    FieldControl = "number"    // Number input
	ControlCheckbox  FieldControl = "checkbox"  // Checkbox input
	ControlSelect    FieldControl = "select"    // Select dropdown
	ControlText      FieldControl = "text"      // Text input
	ControlNoControl FieldControl = "nocontrol" // Special case for displaying a label only
)

type ChartType string

const (
	ChartTypeLine    ChartType = "line"
	ChartTypeBar     ChartType = "bar"
	ChartTypeScatter ChartType = "scatter"
	ChartTypeBubble  ChartType = "bubble"
	ChartTypePie     ChartType = "pie"
)

type AxisConfig string

const (
	LinearAxis      AxisConfig = "linear"
	LogarithmicAxis AxisConfig = "logarithmic"
	TimeAxis        AxisConfig = "time"
	CategoryAxis    AxisConfig = "category"
)

const (
	Color1  = "#2563eb" // bad color
	Color2  = "#dc2626"
	Color3  = "#fbbf24"
	Color4  = "#22c55e"
	Color5  = "#7c3aed"
	Color6  = "#8b5cf6"
	Color7  = "#e879f9"
	Color8  = "#f97316"
	Color9  = "#a855f7"
	Color10 = "#e11d48"
	Color11 = "#10b981"

	ColorTransparent = "rgba(0, 0, 0, 0.1)"
)
