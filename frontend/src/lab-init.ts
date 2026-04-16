import { charting } from "../wailsjs/go/models";
import { toggleDatasetVisibility } from "./chart-render";
import {
  ActiveChartChangeEvent,
  ActiveLabChangeEvent,
  RerenderEvent,
} from "./events";
import { registry } from "./registry";
import { SafeChart } from "./types";

export function buildInputFieldName(
  chartId: string,
  graphId: string | null,
  field: string,
) {
  if (!graphId) {
    return `${chartId}-${field}`;
  }
  return `${chartId}-${graphId}-${field}`;
}

function CreateInput(
  field: charting.MutableField,
  parentElement: HTMLElement,
  chartId: string,
  graphId: string | null = null,
) {
  const label = document.createElement("label");
  label.textContent = field.label;
  label.setAttribute("data-field-label", field.label);
  label.setAttribute(
    "data-input-id",
    buildInputFieldName(chartId, graphId, field.id),
  );
  parentElement.appendChild(label);

  let fieldInput: HTMLSelectElement | HTMLInputElement;
  if (field.control === "select") {
    fieldInput = document.createElement("select") as HTMLSelectElement;
  } else {
    fieldInput = document.createElement("input") as HTMLInputElement;

    fieldInput.type = field.control;
    fieldInput.min = field.min.toString();
    fieldInput.max = field.max.toString();
    fieldInput.step = field.step.toString();
  }

  fieldInput.setAttribute("data-input-type", field.control);
  fieldInput.id = buildInputFieldName(chartId, graphId, field.id);
  fieldInput.value = field.default.toString();
  const updateLabel = () => {
    label.textContent = `${field.label} (${fieldInput.value})`;
  };

  switch (field.control) {
    case "range":
      updateLabel();
      fieldInput.oninput = updateLabel;
      fieldInput.onchange = updateLabel;
      break;
    case "nocontrol":
      fieldInput.hidden = true;
      break;
    case "select":
      if (!field.options || field.options.length === 0) {
        console.warn("no options for select field");
        break;
      }
      for (const opt of field.options) {
        const option = document.createElement("option");
        option.value = field.options.indexOf(opt).toString();
        option.innerText = opt;

        fieldInput.appendChild(option);
      }
      break;
  }
  parentElement.appendChild(fieldInput);
}

export function updateAllFieldLabels(chart: SafeChart) {
  if (!chart) {
    return;
  }

  const chartId = chart.id;

  const changeLabel =
    (chartId: string, graphId: string | null) =>
    (field: charting.MutableField) => {
      const inputId = buildInputFieldName(chartId, graphId, field.id);
      const input = document.getElementById(inputId) as HTMLInputElement;
      const label = document.querySelector(
        `label[data-input-id="${inputId}"]`,
      ) as HTMLLabelElement;

      if (label && input && input.getAttribute("data-input-type") === "nocontrol") {
        label.innerText = field.label;
      }
    };

  // Update chart variables
  if (chart.chartVariables) {
    chart.chartVariables.forEach(changeLabel(chartId, null));
  }

  // Update dataset-specific variables if they exist
  if (chart.datasets) {
    Object.entries(chart.datasets).forEach(([datasetKey, dataset]) => {
      if (dataset.fields) {
        dataset.fields.forEach(changeLabel(chartId, datasetKey));
      }
    });
  }
}

export function InitializeChart(chartId: string) {
  const labId = window.activeLabId;
  if (!labId) {
    return;
  }

  const lab = registry.getLab(labId);
  if (!lab) {
    return;
  }

  const chart = lab.Charts[chartId];
  if (!chart) {
    return;
  }

  console.log("Initializing chart:", chartId);
  console.log("Chart to render:", chart);

  // Update active chart tab styling
  const allChartTabs = document.querySelectorAll("#chart-list li");
  allChartTabs.forEach((tab) => {
    tab.classList.remove("active-chart");
    tab.classList.add("inactive-chart");
  });
  const activeTab = document.getElementById(`chart-tab-${chartId}`);
  if (activeTab) {
    activeTab.classList.add("active-chart");
    activeTab.classList.remove("inactive-chart");
  }

  const appContainer = document.getElementById("app");
  if (!appContainer) {
    throw new Error("App container not found");
  }

  let chartVariablesContainer = document.getElementById("chart-variables");
  if (!chartVariablesContainer) {
    chartVariablesContainer = document.createElement("div");
    chartVariablesContainer.id = "chart-variables";
    appContainer.appendChild(chartVariablesContainer);
  }

  chartVariablesContainer.innerHTML = "";

  const chartRow = document.createElement("div");
  chartRow.className = "chart-variables-row";

  const chartParamsColumn = document.createElement("div");
  chartParamsColumn.className = "chart-variables-column";

  const datasetsColumn = document.createElement("div");
  datasetsColumn.className = "chart-variables-column";

  const divider = document.createElement("div");
  divider.className = "chart-variables-divider";

  chartRow.appendChild(chartParamsColumn);
  chartRow.appendChild(divider);
  chartRow.appendChild(datasetsColumn);
  chartVariablesContainer.appendChild(chartRow);

  if (chart.chartVariables && chart.chartVariables.length > 0) {
    const chartVarsTitle = document.createElement("div");
    chartVarsTitle.className = "chart-section-title";
    chartVarsTitle.textContent = "Chart Parameters";
    chartParamsColumn.appendChild(chartVarsTitle);
  }

  chart.chartVariables?.forEach((field) =>
    CreateInput(field, chartParamsColumn, chartId),
  );

  const datasetToggles = Object.entries(chart.graphVariables);
  if (datasetToggles.length > 0) {
    const datasetsTitle = document.createElement("div");
    datasetsTitle.className = "chart-section-title";
    datasetsTitle.textContent = "Datasets";
    datasetsColumn.appendChild(datasetsTitle);

    // Create a map of dataset keys to their indices
    // graphVariables used as a trick, i can add empty arrays of
    // MutableFields and the dataset will appear here
    const datasetKeys = Object.keys(chart.graphVariables);

    for (const [graphId, fields] of datasetToggles) {
      // Find the actual index of this dataset in the datasets array
      const datasetIndex = datasetKeys.indexOf(graphId);

      if (datasetIndex === -1) {
        console.warn(`Dataset ${graphId} not found in chart.datasets`);
        continue;
      }

      const toggleContainer = document.createElement("div");
      toggleContainer.className = "dataset-toggle-container";

      const checkbox = document.createElement("input");
      checkbox.type = "checkbox";
      checkbox.id = `toggle-${chartId}-${graphId}`;
      checkbox.checked = true;

      const label = document.createElement("label");
      label.htmlFor = checkbox.id;
      label.className = "dataset-toggle-label";
      label.textContent = graphId;

      checkbox.addEventListener("change", () => {
        console.log(
          `Checkbox changed - dataset: ${graphId}, index: ${datasetIndex}, checked: ${checkbox.checked}`,
        );
        console.log(
          `Calling toggleDatasetVisibility with chartId: ${chartId}, datasetIndex: ${datasetIndex}`,
        );
        toggleDatasetVisibility(chartId, datasetIndex);
        label.classList.toggle("hidden", !checkbox.checked);
      });

      toggleContainer.appendChild(checkbox);
      toggleContainer.appendChild(label);
      datasetsColumn.appendChild(toggleContainer);

      if (fields && fields.length > 0) {
        const fieldGroup = document.createElement("div");
        fieldGroup.style.marginLeft = "30px";
        fieldGroup.style.marginBottom = "10px";

        fields.forEach((field) =>
          CreateInput(field, fieldGroup, chartId, graphId),
        );
        datasetsColumn.appendChild(fieldGroup);
      }
    }
  }

  let rerenderBtn = document.getElementById(
    "rerender-btn",
  ) as HTMLButtonElement | null;
  if (!rerenderBtn) {
    rerenderBtn = document.createElement("button");
    rerenderBtn.id = "rerender-btn";
    rerenderBtn.textContent = "Rerender";
    appContainer.appendChild(rerenderBtn);
  }

  rerenderBtn.onclick = () => {
    rerenderBtn?.classList.add("loading");
    window.dispatchEvent(new RerenderEvent(window.activeChartId!));
  };
}

// Create lab tabs (call once with all labs)
export function createLabTabs(labsMetadata: charting.LabMetadata[]) {
  let labList = document.getElementById("lab-list");
  if (!labList) {
    labList = document.createElement("ul");
    labList.id = "lab-list";
    labList.classList.add("tabs");
    document.body.insertBefore(labList, document.body.firstChild);
  }

  labsMetadata.forEach((lab) => {
    const li = document.createElement("li");
    li.id = `lab-tab-${lab.ID}`;
    li.textContent = lab.ID;
    li.style.cursor = "pointer";
    li.addEventListener("click", () => {
      window.dispatchEvent(new ActiveLabChangeEvent(lab.ID));
    });
    labList!.appendChild(li);
  });
}

// Initialize a specific lab (call on lab change)
export function InitializeLab(labId: string) {
  const lab = registry.getLab(labId);
  if (!lab) {
    console.error(`Lab ${labId} not found in registry`);
    return;
  }

  console.log("Initializing lab:", lab.ID);

  // Update active tab styling
  const allLabTabs = document.querySelectorAll("#lab-list li");
  allLabTabs.forEach((tab) => {
    tab.classList.remove("active-lab");
    tab.classList.add("inactive-lab");
  });
  const activeTab = document.getElementById(`lab-tab-${labId}`);
  if (activeTab) {
    activeTab.classList.add("active-lab");
    activeTab.classList.remove("inactive-lab");
  }

  const appContainer = document.getElementById("app");
  if (!appContainer) {
    throw new Error("App container not found");
  }

  // Clear and rebuild chart list
  let chartList = document.getElementById("chart-list");
  if (!chartList) {
    chartList = document.createElement("ul");
    chartList.id = "chart-list";
    chartList.classList.add("tabs");
    appContainer.appendChild(chartList);
  } else {
    chartList.innerHTML = "";
  }

  // Create chart tabs for this lab
  for (const [chartId, chart] of Object.entries(lab.Charts)) {
    const chartElement = document.createElement("li");
    chartElement.id = `chart-tab-${chartId}`;
    chartElement.textContent = chart.title;
    chartElement.style.cursor = "pointer";
    chartElement.addEventListener("click", () => {
      window.dispatchEvent(new ActiveChartChangeEvent(chartId));
    });
    chartList.appendChild(chartElement);
  }
}
