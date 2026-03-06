import { Chart, ChartTypeRegistry } from "chart.js";
import { charting } from "../wailsjs/go/models";

// Store chart instances globally for access
declare global {
  interface Window {
    chartInstances: Map<string, Chart>;
  }
}

if (!window.chartInstances) {
  window.chartInstances = new Map();
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

  const canvasWrapper = document.createElement("div");
  canvasWrapper.className = "chart-wrapper";

  const canvas = document.createElement("canvas");
  canvas.id = `chart-${chartConfig.id}`;

  canvasWrapper.appendChild(canvas);
  container!.appendChild(canvasWrapper);

  // Initialize Chart.js
  const ctx = canvas.getContext("2d");
  if (!ctx) {
    console.log("Canvas context is null or undefined!");
    return;
  }

  const chartType = chartConfig.type as keyof ChartTypeRegistry;

  console.log("Chart Type:", chartType);
  const hasScales = !["pie", "doughnut", "polarArea"].includes(chartType);

  let labels: string[] = [];
  if (chartConfig.labels) {
    labels = chartConfig.labels;
  } else if (chartConfig.datasets && hasScales) {
    console.log("no appropriate labels found, trying to generate our own");
    let lowBound = Number.MAX_VALUE;
    let highBound = Number.MIN_VALUE;
    let length = 0;
    for (const [_, graph] of Object.entries(chartConfig.datasets)) {
      if (graph.data && graph.data.length > labels.length) {
        labels = graph.data.map((_, idx) => idx.toFixed(2));
        continue;
      }

      if (graph.pointData) {
        for (const val of graph.pointData) {
          lowBound = Math.min(lowBound, val.x);
          highBound = Math.max(highBound, val.x);
        }
        length = Math.max(length, graph.pointData.length);
        continue;
      }

      const error = document.getElementById("error-container");
      if (error) {
        error.innerHTML = `<div style="color: red; padding: 20px;">Error: missing data attributes</div>`;
      }
    }
    if (
      lowBound != Number.MAX_VALUE &&
      highBound != Number.MIN_VALUE &&
      lowBound != highBound
    ) {
      const step = (highBound - lowBound) / length;
      for (let i = 0; i < length; i++) {
        labels[i] = (lowBound + step * i).toFixed(2);
      }
    }
  }

  // Determine if chart type requires scales (axes)

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
        console.warn("Empty data");
        data = [];
      }

      // For pie/doughnut charts, generate color palette if not specified
      let backgroundColor = dataset.backgroundColor;
      if (!hasScales && !backgroundColor && Array.isArray(data)) {
        const colors = [
          "#FF6384",
          "#36A2EB",
          "#FFCE56",
          "#4BC0C0",
          "#9966FF",
          "#FF9F40",
          "#FF6384",
          "#C9CBCF",
          "#4BC0C0",
          "#FF6384",
        ];
        backgroundColor = data.map(
          (_: any, idx: number) => colors[idx % colors.length],
        );
      }

      return {
        label: dataset.label,
        data: data,
        borderColor: dataset.borderColor ?? "#666",
        backgroundColor: backgroundColor ?? "transparent",
        tension: dataset.tension ?? 0,
        fill: dataset.fill ?? false,
        hidden: dataset.hidden ?? false,
        pointRadius: dataset.pointRadius ?? 0,
        borderWidth: dataset.borderWidth ?? 2,
        showLine: dataset.showLine === true,
      };
    },
  );

  const chartOptions: any = {
    responsive: true,
    plugins: {
      title: {
        display: true,
        text: chartConfig.title,
        color: "#ffffff",
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
          color: "#ffffff",
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
    },
  };

  // Only add scales for charts that use them
  if (hasScales) {
    chartOptions.scales = {
      x: {
        title: {
          display: true,
          text: chartConfig.xAxisLabel,
          color: "#ffffff",
          font: {
            size: 14,
            weight: "bold",
          },
        },
        ticks: {
          color: "#ffffff",
          font: {
            size: 12,
          },
        },
        grid: {
          color: "rgba(255, 255, 255, 0.2)",
        },
      },
      y: {
        title: {
          display: true,
          text: chartConfig.yAxisLabel,
          color: "#ffffff",
          font: {
            size: 14,
            weight: "bold",
          },
        },
        ticks: {
          color: "#ffffff",
          font: {
            size: 12,
          },
        },
        grid: {
          color: "rgba(255, 255, 255, 0.2)",
        },
      },
    };
  }

  const chart = new Chart(ctx, {
    type: chartType,
    data: {
      labels: labels,
      datasets: processedDatasets,
    },
    options: chartOptions,
  });

  // Store chart instance for later access
  window.chartInstances.set(chartConfig.id, chart);
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
