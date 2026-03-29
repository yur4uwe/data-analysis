package charting

import (
	"fmt"
	"regexp"
	"strings"
)

// DataPoint represents a point with x and y coordinates for scatter/bubble charts
type DataPoint struct {
	X float64  `json:"x" csv:"x"`
	Y *float64 `json:"y" csv:"y"`
}

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

type GraphType string

const (
	ChartTypeLine    GraphType = "line"
	ChartTypeBar     GraphType = "bar"
	ChartTypeScatter GraphType = "scatter"
	ChartTypeBubble  GraphType = "bubble"
	ChartTypePie     GraphType = "pie"
	ChartTypeHeatmap GraphType = "heatmap"

	ChartTypeMultiLine    GraphType = "multi-line"
	ChartTypeMultiBar     GraphType = "multi-bar"
	ChartTypeMultiScatter GraphType = "multi-scatter"
	ChartTypeMultiBubble  GraphType = "multi-bubble"
	ChartTypeMultiPie     GraphType = "multi-pie"
	ChartTypeMultiHeatmap GraphType = "multi-heatmap"
)

type AxisConfig string

const (
	LinearAxis      AxisConfig = "linear"
	LogarithmicAxis AxisConfig = "logarithmic"
	TimeAxis        AxisConfig = "time"
	CategoryAxis    AxisConfig = "category"
)

type Color string

func isHex(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	l := c | 32
	return l >= 'a' && l <= 'f'
}

func ToColor(colorlike string) Color {
	if strings.HasPrefix(colorlike, "#") {
		if len(colorlike) != 7 {
			panic(fmt.Sprintf("cannot turn %s into hex color with length %d", colorlike, len(colorlike)))
		}
		for i := range 6 {
			if !isHex(colorlike[i+1]) {
				panic(fmt.Sprintf("cannot turn %q into a hex color, character at index %d isn't hex valid", colorlike, i+1))
			}
		}
		return Color(colorlike)
	} else if strings.HasPrefix(colorlike, "rgb") {
		colorRegex := regexp.MustCompile(`(?i)rgba?\(\s*\d{1,3}\s*(?:,\s*\d{1,3}\s*){2}(?:,\s*(?:\d*\.)?\d+\s*)?\)`)
		if !colorRegex.Match([]byte(colorlike)) {
			panic(fmt.Sprintf("failed to parse %q with an rgba regex", colorlike))
		}

		return Color(colorlike)

	} else {
		panic(fmt.Sprintf("impossible to turn %q into color", colorlike))
	}
}

const (
	ColorBlue        Color = "#1d4ed8"
	ColorRed         Color = "#b91c1c"
	ColorAmber       Color = "#d97706"
	ColorGreen       Color = "#16a34a"
	ColorViolet      Color = "#6d28d9"
	ColorPurple      Color = "#7c3aed"
	ColorFuchsia     Color = "#c026d3"
	ColorOrange      Color = "#ea580c"
	ColorLightPurple Color = "#9333ea"
	ColorCrimson     Color = "#be123c"
	ColorEmerald     Color = "#059669"
	ColorCyan        Color = "#0891b2"
	ColorPink        Color = "#db2777"
	ColorLime        Color = "#65a30d"
	ColorTeal        Color = "#0d9488"
	ColorIndigo      Color = "#4f46e5"
	ColorRose        Color = "#e11d48"
	ColorSky         Color = "#0284c7"
	ColorYellow      Color = "#ca8a04"
	ColorSlate       Color = "#475569"

	ColorTransparent Color = "rgba(0, 0, 0, 0.1)"
)
