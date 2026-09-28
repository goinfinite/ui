UiToolset.RegisterAlpineState(() => {
  Alpine.data("filterBar", (valuesPath) => ({
    valuesPath: valuesPath,

    filterValuesObject() {
      return this[this.valuesPath] ?? {};
    },

    isFilterActive(filterKey) {
      const filterValue = this.filterValuesObject()[filterKey];
      if (Array.isArray(filterValue)) {
        return filterValue.length > 0;
      }
      if (isRangeFilterValue(filterValue)) {
        return filterValue.min !== "" || filterValue.max !== "";
      }
      return (
        filterValue !== "" && filterValue !== null && filterValue !== undefined
      );
    },

    resolveFilterChipLabel(filterKey) {
      const filterValue = this.filterValuesObject()[filterKey];
      if (Array.isArray(filterValue)) {
        return filterValue.join(", ");
      }
      if (!isRangeFilterValue(filterValue)) {
        return filterValue ?? "";
      }
      return rangeFilterChipLabel(filterValue);
    },

    resetFilter(filterKey) {
      const filterValue = this.filterValuesObject()[filterKey];
      if (Array.isArray(filterValue)) {
        this[this.valuesPath][filterKey] = [];
        return;
      }
      if (isRangeFilterValue(filterValue)) {
        this[this.valuesPath][filterKey] = { min: "", max: "" };
        return;
      }
      this[this.valuesPath][filterKey] = "";
    },

    resetAllFilters() {
      for (const filterKey of Object.keys(this.filterValuesObject())) {
        this.resetFilter(filterKey);
      }
    },

    hasActiveFilters() {
      return Object.keys(this.filterValuesObject()).some((filterKey) =>
        this.isFilterActive(filterKey),
      );
    },
  }));
});

function isRangeFilterValue(filterValue) {
  return (
    !Array.isArray(filterValue) &&
    filterValue !== null &&
    typeof filterValue === "object"
  );
}

function rangeFilterChipLabel(rangeFilterValue) {
  const minimum = rangeFilterValue.min ?? "";
  const maximum = rangeFilterValue.max ?? "";
  if (minimum !== "" && maximum !== "") {
    return `${minimum}–${maximum}`;
  }
  if (minimum !== "") {
    return `≥ ${minimum}`;
  }
  if (maximum !== "") {
    return `≤ ${maximum}`;
  }
  return "";
}
