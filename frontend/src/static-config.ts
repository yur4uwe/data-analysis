import { charting } from "../wailsjs/go/models";
import { Dataset, isGridDataset, isHeatmapDataset, isCategoricalDataset } from "./types";
import { Layout as PlotlyLayout } from "plotly.js-dist-min";

export const formatScientific = (value: number | string): string => {
	const num = typeof value === "string" ? parseFloat(value) : value;
	if (isNaN(num)) return value.toString();
	if (num === 0) return "0";
	const absNum = Math.abs(num);
	if (absNum >= 1e6 || absNum <= 1e-3) {
		return num.toExponential(2);
	}
	return num.toLocaleString(undefined, { maximumFractionDigits: 4 });
};

export const defaultPlotlyLayout = (title: string, chartType: string, chartConfig?: charting.Chart): Partial<PlotlyLayout> => {
	const is3D = chartType === "surface" || chartType.includes("3d");
	
	const layout: any = {
		title: {
			text: title,
			font: { size: 18, color: "#000000", family: "Nunito, sans-serif" },
			pad: { t: 10, b: 20 }
		},
		autosize: true,
		margin: { l: 80, r: 50, t: 80, b: 80, pad: 10 },
		paper_bgcolor: "rgba(0,0,0,0)",
		plot_bgcolor: "rgba(0,0,0,0)",
		font: { family: "Nunito, sans-serif", size: 12, color: "#000000" },
		showlegend: !chartType.includes("heatmap") && !chartType.includes("surface"),
		legend: {
			orientation: "h",
			yanchor: "bottom",
			y: -0.3,
			xanchor: "center",
			x: 0.5
		}
	};

	if (is3D && chartConfig) {
		layout.scene = {
			xaxis: { title: { text: chartConfig.xAxisLabel || "X" } },
			yaxis: { title: { text: chartConfig.yAxisLabel || "Y" } },
			zaxis: { title: { text: "Value" } }
		};
	}

	return layout;
};

export function newPlotlyAxes(chartConfig: charting.Chart) {
	const mapAxisType = (config: string | undefined) => {
		switch (config) {
			case "logarithmic": return "log";
			case "linear": return "linear";
			case "category": return "category";
			case "time": return "date";
			default: return "linear";
		}
	};

	return {
		xaxis: {
			type: mapAxisType(chartConfig.xAxisConfig),
			title: { 
				text: chartConfig.xAxisLabel || "X",
				font: { size: 14, family: "Nunito, sans-serif" }
			},
			showgrid: true,
			zeroline: true,
			gridcolor: "rgba(0,0,0,0.1)",
			zerolinecolor: "#000000",
			zerolinewidth: 2,
			tickfont: { size: 12 }
		},
		yaxis: {
			type: mapAxisType(chartConfig.yAxisConfig),
			title: { 
				text: chartConfig.yAxisLabel || "Y",
				font: { size: 14, family: "Nunito, sans-serif" }
			},
			showgrid: true,
			zeroline: true,
			gridcolor: "rgba(0,0,0,0.1)",
			zerolinecolor: "#000000",
			zerolinewidth: 2,
			tickfont: { size: 12 }
			// Removed tickformat ".2e" to avoid forcing scientific notation if not needed
		}
	};
}

export const processDatasetToPlotly = (chartType: string, labels: string[]) => (dataset: Dataset): any => {
	if (!dataset) return {};

	const base: any = {
		name: dataset.label || "Dataset",
		visible: dataset.hidden ? "legendonly" : true,
	};

	if (isHeatmapDataset(dataset)) {
		const validData = dataset.pointData.filter(p => p !== null && p.y !== null);
		const xValues = Array.from(new Set(validData.map(p => p.x))).sort((a, b) => a - b);
		const yValues = Array.from(new Set(validData.map(p => p.y as number))).sort((a, b) => a - b);
		
		const zMatrix: (number | null)[][] = yValues.map(() => xValues.map(() => null));
		
		validData.forEach(p => {
			const xi = xValues.indexOf(p.x);
			const yi = yValues.indexOf(p.y as number);
			if (xi !== -1 && yi !== -1) {
				zMatrix[yi][xi] = p.v;
			}
		});

		return {
			...base,
			type: chartType === "surface" ? "surface" : "heatmap",
			x: xValues,
			y: yValues,
			z: zMatrix,
			colorscale: "Viridis",
			showscale: true
		};
	}

	if (isGridDataset(dataset)) {
		const x = dataset.data.map(p => p.x);
		const y = dataset.data.map(p => p.y);
		
		const radius = (dataset as any).pointRadius ?? 6;
		const hideLine = (dataset as any).hideLine ?? false;
		
		let mode = "lines+markers";
		if (chartType === "scatter" || chartType === "bubble") {
			mode = "markers";
		} else if (radius === 0 && hideLine) {
			mode = "none";
		} else if (radius === 0) {
			mode = "lines";
		} else if (hideLine) {
			mode = "markers";
		}

		return {
			...base,
			type: "scatter",
			mode: mode,
			x,
			y: y,
			line: { color: dataset.borderColor, width: dataset.borderWidth },
			marker: { 
				color: dataset.borderColor, 
				size: chartType === "bubble" ? 10 : radius,
				opacity: radius === 0 ? 0 : 1
			}
		};
	}

	if (isCategoricalDataset(dataset)) {
		if (chartType === "pie" || chartType === "doughnut") {
			return {
				...base,
				type: "pie",
				labels: labels.length > 0 ? labels : dataset.data.map((_, i) => `Item ${i+1}`),
				values: dataset.data.filter(v => v !== null),
				marker: { colors: (dataset as any).backgroundColor }
			};
		}
		
		return {
			...base,
			type: "bar",
			x: labels.length > 0 ? labels : dataset.data.map((_, i) => `Item ${i+1}`),
			y: dataset.data,
			marker: { color: dataset.borderColor }
		};
	}

	return base;
};
