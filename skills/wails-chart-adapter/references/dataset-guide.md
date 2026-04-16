# 📊 Dataset Implementation Guide

This document describes how to implement datasets for each chart type.

## 1. GridDataset (Scatter/Line)
Uses `{X, Y}` coordinates.
- **Go Struct**: `GridDataset`
- **Data Field**: `data: []DataPoint`
- **Helper**: `chart.UpdatePointsForDataset(id, x, y)`
- **Behavior**: X-values determine the labels if `chart.GenerateLabels(precision)` is called.

## 2. CategoricalDataset (Bar/Pie)
Uses indexed values matching `chart.Labels`.
- **Go Struct**: `CategoricalDataset`
- **Data Field**: `data: []*float64` (supports `nil` for missing data).
- **Helper**: `chart.UpdateDataForDataset(id, data)`
- **Constraint**: Must ensure `len(data) == len(chart.Labels)`.

## 3. HeatmapDataset (3D)
Uses `{X, Y, V}` coordinates for intensity mapping.
- **Go Struct**: `HeatmapDataset`
- **Data Field**: `pointData: []HeatmapPoint`
- **Behavior**: `V` (value) determines the color in the heatmap.
- **Visualization Note**: `ChartTypeHeatmap` is preferred over `ChartTypeSurface` when you need to overlay 2D optimization paths or points on top of a 3D function, as it provides a clear top-down view.
