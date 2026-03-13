import { Chart, ChartTypeRegistry, ScaleChartOptions } from "chart.js";
import ChartDataLabels from "chartjs-plugin-datalabels";
import ZoomPlugin from "chartjs-plugin-zoom";
import { charting } from "../wailsjs/go/models";

Chart.register(ChartDataLabels, ZoomPlugin);

// Store chart instances globally for access
declare global {
  interface Window {
    chartInstances: Map<string, Chart>;
  }
}

if (!window.chartInstances) {
  window.chartInstances = new Map();
}

function getDataLabels(
  pointLabels: string[] | undefined,
  chartType: keyof ChartTypeRegistry,
) {
  if (pointLabels?.length) {
    return {
      display: true,
      formatter: (_: any, ctx: any) => pointLabels[ctx.dataIndex],
      anchor: "end" as const,
      align: "top" as const,
      color: "#000000",
      font: { weight: "bold" as const, size: 11 },
      backgroundColor: "rgba(255,255,255,0.85)",
      borderRadius: 3,
      padding: 4,
    };
  }
  switch (chartType) {
    case "pie":
      console.log(`returning labels for pie`);
    case "doughnut":
      console.log(`returning labels for doughnut`);
      return {
        color: "#ffffff",
        font: { weight: "bold" as const, size: 14 },
        formatter: (value: number, ctx: any) => {
          const total = (ctx.dataset.data as number[]).reduce(
            (a, b) => a + b,
            0,
          );
          const pct = ((value / total) * 100).toFixed(1);
          return `${value}\n(${pct}%)`;
        },
      };

    case "bar":
      console.log(`returning labels for bar`);
      return {
        color: "#ffffff",
        font: { weight: "bold" as const, size: 14 },
      };
    default:
      console.log(`returning labels for default`);
      return {
        display: false,
      };
  }
}

function newScales(chartConfig: charting.Chart, hasContinuousAxes: boolean) {
  return {
    x: {
      ...(hasContinuousAxes && { type: "linear" as const }),
      border: {
        display: !hasContinuousAxes,
      },
      title: {
        display: !!chartConfig.xAxisLabel,
        text: chartConfig.xAxisLabel,
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
      border: {
        display: !hasContinuousAxes,
      },
      title: {
        display: !!chartConfig.yAxisLabel,
        text: chartConfig.yAxisLabel,
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

export function renderChart(chartConfig: charting.Chart) {
  let container = document.getElementById("chart-container");
  if (!container) {
    container = document.createElement("div");
    container.id = "chart-container";
    document.body.appendChild(container);
  }

  // Clear previous content
  container.innerHTML = "";

  if (!chartConfig) {
    console.log("Charts data is null or undefined!");
    return;
  }

  const canvas = document.createElement("canvas");
  canvas.id = `chart-${chartConfig.id}`;

  container.appendChild(canvas);

  const ctx = canvas.getContext("2d");
  if (!ctx) {
    console.log("Canvas context is null or undefined!");
    return;
  }

  const chartType = chartConfig.type as keyof ChartTypeRegistry;

  console.log("Chart Type:", chartType);
  const hasScales = !["pie", "doughnut", "polarArea"].includes(chartType);
  const hasContinuousAxes = ["scatter", "line", "bubble"].includes(chartType);

  const labels: string[] = chartConfig.labels ?? [];

  // Process datasets based on chart type
  const processedDatasets = Object.values(chartConfig.datasets).map(
    (dataset) => {
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
        console.warn(`Empty data in ${dataset.label}`);
        data = [];
      }

      return {
        label: dataset.label,
        data: data,
        borderColor: dataset.borderColor,
        backgroundColor: dataset.backgroundColor ?? dataset.borderColor,
        tension: dataset.tension ?? 0,
        fill: dataset.fill ?? false,
        hidden: dataset.hidden ?? false,
        pointRadius: dataset.pointRadius ?? 0,
        borderWidth: dataset.borderWidth ?? 2,
        showLine: dataset.showLine === true,
        pointStyle: dataset.pointStyle ?? undefined,
        datalabels: getDataLabels(dataset.pointLabels, chartType),
      };
    },
  );

  const chartOptions: any = {
    responsive: true,
    resizeDelay: 100,
    animation: false,
    plugins: {
      title: {
        display: true,
        text: chartConfig.title,
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
        labels: {
          color: "#000000",
          font: {
            size: 13,
          },
          padding: 15,
          usePointStyle: true,
        },
      },
      datalabels: getDataLabels(undefined, chartType),
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
  };

  // Only add scales for charts that use them
  if (hasScales) {
    chartOptions.scales = newScales(
      chartConfig,
      hasContinuousAxes,
    ) as unknown as ScaleChartOptions;
  }

  const chart = new Chart(ctx, {
    type: chartType,
    data: {
      labels: labels,
      datasets: processedDatasets,
    },
    options: chartOptions,
  });

  window.chartInstances.set(chartConfig.id, chart);

  const resetZoom = document.getElementById("reset-zoom-btn");
  if (!resetZoom) {
    return;
  }
  const resetZoomCallback = () => {
    (window.chartInstances.get(chartConfig.id) as any).resetZoom();
  };
  resetZoom.removeEventListener("click", resetZoomCallback);
  resetZoom.addEventListener("click", resetZoomCallback);
}

// Helper functions for dataset control
export function toggleDatasetVisibility(chartId: string, datasetIndex: number) {
  const chart = window.chartInstances.get(chartId);
  if (chart) {
    chart.data.datasets[datasetIndex].hidden =
      !chart.data.datasets[datasetIndex].hidden;

    chart.update();
  }
}

export function updateDatasets(chartId: string, newDatasets: any[]) {
  const chart = window.chartInstances.get(chartId);
  if (chart) {
    chart.data.datasets = newDatasets;
    chart.update();
  }
}

export function addDataset(chartId: string, newDataset: any) {
  const chart = window.chartInstances.get(chartId);
  if (chart) {
    chart.data.datasets.push(newDataset);
    chart.update();
  }
}

export function removeDataset(chartId: string, datasetIndex: number) {
  const chart = window.chartInstances.get(chartId);
  if (chart) {
    chart.data.datasets.splice(datasetIndex, 1);
    chart.update();
  }
}
