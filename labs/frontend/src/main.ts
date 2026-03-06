import "./style.css";
import "./app.css";

import { ActiveLabChangeEvent } from "./events";
import { createLabTabs, updateAllFieldLabels } from "./lab-init";

import { GetLabs } from "../wailsjs/go/main/App";
import { charting } from "../wailsjs/go/models";
import { Chart, registerables } from "chart.js";
import ChartDataLabels from "chartjs-plugin-datalabels";
import { EventsOn } from "../wailsjs/runtime";
import { renderChart } from "./chart-render";
import { registry } from "./registry";

Chart.register(...registerables, ChartDataLabels);

// State

window.addEventListener("load", () => {
  GetLabs().then((labsResponse) => {
    console.log("Loaded labs:", labsResponse.labs);

    // Register all labs in the registry without initializing UI
    labsResponse.labs.forEach((lab) => {
      registry.addLab(lab.ID, lab);
    });

    // Create lab tabs
    createLabTabs(labsResponse.labs);

    // Initialize first lab by dispatching activeLabChange event
    if (labsResponse.labs.length > 0) {
      window.dispatchEvent(new ActiveLabChangeEvent(labsResponse.labs[0].ID));
    }
  });
});

EventsOn("renderComplete", (data: charting.RenderResponse) => {
  console.log("Render complete:", data);
  if (data.error) {
    console.error("Render error:", data.error);
    const container = document.getElementById("error-container");
    if (container) {
      container.innerHTML = `<div style="color: red; padding: 20px;">Error: ${data.error.Message}</div>`;
    }
  } else {
    console.log("Render completed successfully,", window.activeChartId);

    const activeChartData = data.charts[window.activeChartId!];
    renderChart(activeChartData);

    // Update all field labels to reflect actual rendered values
    updateAllFieldLabels(activeChartData);

    const container = document.getElementById("error-container");
    if (container) {
      container.innerHTML = "";
    }
  }
});

// Listen for render errors
EventsOn("renderError", (data: charting.RenderResponse) => {
  console.error("Render error:", data);
  const container = document.getElementById("error-container");
  if (container) {
    container.innerHTML = `<div style="color: red; padding: 20px;">Error: ${data.error.Message}</div>`;
  }
});
