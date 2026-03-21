import { ChartTypeRegistry } from "chart.js";
import { charting } from "../wailsjs/go/models";
import { getDataLabels } from "./chart-render";

export const defaultChartOptions = (title: string) => ({
	responsive: true,
	maintainAspectRatio: false,
	resizeDelay: 100,
	animation: false,
	plugins: {
		title: {
			display: !!title,
			text: title || "",
			color: "#000000",
			font: {
				size: 18,
				weight: "bold",
			},
			padding: {
				top: 10,
				bottom: 20,
			},
		},
		legend: {
			display: true,
			labels: {
				color: "#000000",
				font: {
					size: 13,
				},
				padding: 15,
				usePointStyle: true,
			},
		},
		tooltip: {
			backgroundColor: "rgba(0, 0, 0, 0.9)",
			titleColor: "#ffffff",
			bodyColor: "#ffffff",
			borderColor: "#ffffff",
			borderWidth: 1,
			padding: 12,
			displayColors: true,
		},
		zoom: {
			zoom: {
				wheel: {
					enabled: true,
					speed: 0.02,
					modifierKey: "ctrl",
				},
				pinch: { enabled: true },
				mode: "xy",
			},
			pan: {
				enabled: true,
				mode: "xy",
			},
		},
	},
})

export const processDataset = (hasScales: boolean, chartType: keyof ChartTypeRegistry) => (dataset: charting.ChartDataset) => {
	if (!dataset) return { label: "unknown", data: [] };

	let data: any;

	// For pie/doughnut/polarArea charts, use simple array values
	if (!hasScales && dataset.data) {
		data = dataset.data;
	}
	// For charts with scales, use pointData y-values or data array
	else if (dataset.pointData) {
		data = dataset.pointData;
	} else if (dataset.data) {
		data = dataset.data;
	} else {
		console.warn(`Empty data in dataset ${dataset.label}`);
		data = [];
	}

	const datalabels = getDataLabels(dataset.pointLabels, chartType);

	return {
		label: dataset.label || "Unnamed dataset",
		data: data,
		borderColor: dataset.borderColor || "#000000",
		backgroundColor: dataset.backgroundColor ?? dataset.borderColor ?? "#000000",
		tension: dataset.tension ?? 0,
		fill: dataset.fill ?? false,
		hidden: dataset.hidden ?? false,
		pointRadius: dataset.pointRadius ?? 0,
		borderWidth: dataset.borderWidth ?? 2,
		showLine: dataset.showLine !== false,
		togglable: dataset.togglable !== false,
		pointStyle: dataset.pointStyle ?? undefined,
		datalabels: datalabels,
	};
}

export function newScales(chartConfig: charting.Chart, hasContinuousAxes: boolean) {
	const xAxisType = chartConfig.xAxisConfig || (hasContinuousAxes ? "linear" : "category");
	const yAxisType = chartConfig.yAxisConfig || (hasContinuousAxes ? "linear" : "linear");

	return {
		x: {
			type: xAxisType as any,
			border: {
				display: !hasContinuousAxes,
			},
			title: {
				display: !!chartConfig.xAxisLabel,
				text: chartConfig.xAxisLabel ?? "",
				color: "#000000",
				font: {
					size: 14,
					weight: "bold",
				},
			},
			ticks: {
				color: "#000000",
				font: {
					size: 12,
				},
			},
			grid: {
				color: (ctx: any) =>
					hasContinuousAxes && ctx.tick?.value === 0
						? "#000000"
						: "rgba(0,0,0,0.1)",
				lineWidth: (ctx: any) =>
					hasContinuousAxes && ctx.tick?.value === 0 ? 2 : 1,
			},
		},
		y: {
			type: yAxisType as any,
			border: {
				display: !hasContinuousAxes,
			},
			title: {
				display: !!chartConfig.yAxisLabel,
				text: chartConfig.yAxisLabel ?? "",
				color: "#000000",
				font: {
					size: 14,
					weight: "bold",
				},
			},
			ticks: {
				color: "#000000",
				font: {
					size: 12,
				},
			},
			grid: {
				color: (ctx: any) =>
					hasContinuousAxes && ctx.tick?.value === 0
						? "#000000"
						: "rgba(0,0,0,0.1)",
				lineWidth: (ctx: any) =>
					hasContinuousAxes && ctx.tick?.value === 0 ? 2 : 1,
			},
		},
	};
}
