package charting

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

	ChartTypeMultiLine    ChartType = "multi-line"
	ChartTypeMultiBar     ChartType = "multi-bar"
	ChartTypeMultiScatter ChartType = "multi-scatter"
	ChartTypeMultiBubble  ChartType = "multi-bubble"
	ChartTypeMultiPie     ChartType = "multi-pie"
)

func Multi(ct ChartType) ChartType {
	return ChartType("multi-" + string(ct))
}

type AxisConfig string

const (
	LinearAxis      AxisConfig = "linear"
	LogarithmicAxis AxisConfig = "logarithmic"
	TimeAxis        AxisConfig = "time"
	CategoryAxis    AxisConfig = "category"
)

const (
	ColorBlue        = "#1d4ed8"
	ColorRed         = "#b91c1c"
	ColorAmber       = "#d97706"
	ColorGreen       = "#16a34a"
	ColorViolet      = "#6d28d9"
	ColorPurple      = "#7c3aed"
	ColorFuchsia     = "#c026d3"
	ColorOrange      = "#ea580c"
	ColorLightPurple = "#9333ea"
	ColorCrimson     = "#be123c"
	ColorEmerald     = "#059669"
	ColorCyan        = "#0891b2"
	ColorPink        = "#db2777"
	ColorLime        = "#65a30d"
	ColorTeal        = "#0d9488"
	ColorIndigo      = "#4f46e5"
	ColorRose        = "#e11d48"
	ColorSky         = "#0284c7"
	ColorYellow      = "#ca8a04"
	ColorSlate       = "#475569"

	ColorTransparent = "rgba(0, 0, 0, 0.1)"
)
