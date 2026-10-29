UiToolset.RegisterAlpineState(() => {
  const bulkEditCheckedTexts = ["true", "1", "yes", "on"];

  const multiColumnRepeatableFieldsetStateFactory = (
    columns,
    newRowFunc,
    bulkEditSettings,
  ) => {
    const rowKeys = new WeakMap();
    let nextRowKey = 1;
    let resolvedBulkEditSettings = {};
    if (bulkEditSettings !== null && typeof bulkEditSettings === "object") {
      resolvedBulkEditSettings = bulkEditSettings;
    }

    return {
      columns: Array.isArray(columns) ? columns : [],
      newRowFunc: typeof newRowFunc === "string" ? newRowFunc : "",
      bulkEditSettings: resolvedBulkEditSettings,
      isBulkEditing: false,
      bulkEditText: "",

      get bulkEditToggleLabelResolver() {
        if (this.isBulkEditing) {
          return this.bulkEditSettings.tableButtonLabel;
        }
        return this.bulkEditSettings.buttonLabel;
      },

      get bulkEditToggleIconResolver() {
        if (this.isBulkEditing) {
          return this.bulkEditSettings.tableButtonIcon || "ph-table";
        }
        return this.bulkEditSettings.buttonIcon || "ph-text-align-left";
      },

      rowKeyResolver(row) {
        const rawRow = Alpine.raw(row);
        if (!rowKeys.has(rawRow)) {
          rowKeys.set(rawRow, nextRowKey);
          nextRowKey += 1;
        }
        return rowKeys.get(rawRow);
      },

      buildDefaultRow() {
        const defaultRow = {};
        for (const column of this.columns) {
          defaultRow[column.id] = column.inputType === "checkbox" ? false : "";
        }
        if (this.newRowFunc === "") {
          return defaultRow;
        }
        let generatedRow = null;
        try {
          generatedRow = Alpine.evaluate(this.$el, this.newRowFunc);
        } catch (evaluateError) {
          console.error(
            `MultiColumnRepeatableFieldsetNewRowFuncFailed: ${evaluateError.message}`,
          );
          return defaultRow;
        }
        if (generatedRow === null || typeof generatedRow !== "object") {
          return defaultRow;
        }
        for (const column of this.columns) {
          if (
            Object.hasOwn(generatedRow, column.id) &&
            generatedRow[column.id] !== undefined
          ) {
            defaultRow[column.id] = generatedRow[column.id];
          }
        }
        return defaultRow;
      },

      appendRow(rows) {
        rows.push(this.buildDefaultRow());
      },

      duplicateRow(rows, rowIndex) {
        rows.splice(rowIndex + 1, 0, { ...rows[rowIndex] });
      },

      moveRow(rows, rowIndex, offset) {
        const targetIndex = rowIndex + offset;
        if (targetIndex < 0 || targetIndex >= rows.length) {
          return;
        }
        const [movedRow] = rows.splice(rowIndex, 1);
        rows.splice(targetIndex, 0, movedRow);
      },

      removeRow(rows, rowIndex) {
        rows.splice(rowIndex, 1);
      },

      parseBulkEditValue(column, rawPart) {
        const cellText = String(rawPart ?? "").trim();
        if (column.inputType === "checkbox") {
          return bulkEditCheckedTexts.includes(cellText.toLowerCase());
        }
        return cellText;
      },

      resolveBulkEditLineParts(line, delimiter) {
        const lineParts = line.split(delimiter);
        if (lineParts.length <= this.columns.length) {
          return lineParts;
        }
        const leadingParts = lineParts.slice(0, this.columns.length - 1);
        const lastColumnText = lineParts
          .slice(this.columns.length - 1)
          .join(delimiter);
        return leadingParts.concat(lastColumnText);
      },

      parseBulkEditRows(text) {
        const delimiter = this.bulkEditSettings.delimiter || "=";
        const parsedRows = [];
        const textLines = String(text ?? "").split("\n");
        for (const rawLine of textLines) {
          const line = rawLine.replace(/\r$/, "");
          if (line.trim() === "") {
            continue;
          }
          const lineParts = this.resolveBulkEditLineParts(line, delimiter);
          const row = {};
          for (
            let columnIndex = 0;
            columnIndex < this.columns.length;
            columnIndex += 1
          ) {
            const column = this.columns[columnIndex];
            row[column.id] = this.parseBulkEditValue(
              column,
              lineParts[columnIndex],
            );
          }
          parsedRows.push(row);
        }
        return parsedRows;
      },

      formatBulkEditValue(column, row) {
        const value = row[column.id];
        if (column.inputType === "checkbox") {
          return value ? "true" : "false";
        }
        if (value === undefined || value === null) {
          return "";
        }
        return String(value);
      },

      serializeBulkEditRows(rows) {
        const delimiter = this.bulkEditSettings.delimiter || "=";
        const textLines = [];
        for (const row of rows) {
          const cellTexts = [];
          let hasValue = false;
          for (const column of this.columns) {
            const cellText = this.formatBulkEditValue(column, row);
            cellTexts.push(cellText);
            if (cellText !== "") {
              hasValue = true;
            }
          }
          if (!hasValue) {
            continue;
          }
          textLines.push(cellTexts.join(delimiter));
        }
        return textLines.join("\n");
      },

      replaceRowsFromBulkText(rows, text) {
        const parsedRows = this.parseBulkEditRows(text);
        rows.splice(0, rows.length);
        for (const row of parsedRows) {
          rows.push(row);
        }
      },

      reconcileRowsFromBulkText(rows) {
        if (!this.isBulkEditing) {
          return;
        }
        if (this.serializeBulkEditRows(rows) === this.bulkEditText) {
          return;
        }
        this.replaceRowsFromBulkText(rows, this.bulkEditText);
      },

      toggleBulkEdit(rows) {
        if (this.isBulkEditing) {
          this.reconcileRowsFromBulkText(rows);
          this.isBulkEditing = false;
          return;
        }
        this.bulkEditText = this.serializeBulkEditRows(rows);
        this.isBulkEditing = true;
      },

      columnValueIsMissing(column, row) {
        const value = row[column.id];
        if (column.inputType === "checkbox") {
          return column.isRequired && !value;
        }
        const valueIsEmpty =
          value === undefined || value === null || String(value).trim() === "";
        return column.isRequired && valueIsEmpty;
      },

      columnIsInvalid(row, columnId) {
        const column = this.columns.find((entry) => entry.id === columnId);
        if (column === undefined) {
          return false;
        }
        return this.columnValueIsMissing(column, row);
      },
    };
  };

  Alpine.data(
    "multiColumnRepeatableFieldset",
    multiColumnRepeatableFieldsetStateFactory,
  );
});
