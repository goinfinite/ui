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

      bulkEditDelimiterResolver() {
        const configuredDelimiter = this.bulkEditSettings.delimiter || "";
        const delimiterIsValid =
          Array.from(configuredDelimiter).length === 1 &&
          !configuredDelimiter.startsWith("n") &&
          !configuredDelimiter.startsWith("\\") &&
          !configuredDelimiter.includes("\n") &&
          !configuredDelimiter.includes("\r");
        if (delimiterIsValid) {
          return configuredDelimiter;
        }
        if (this.columns.length > 2) {
          return "|";
        }
        return "=";
      },

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

      splitBulkEditLine(line, delimiter) {
        if (delimiter === "") {
          return [line];
        }
        const lineParts = [];
        let currentPart = "";
        let characterIndex = 0;
        while (characterIndex < line.length) {
          const character = line[characterIndex];
          if (character === "\\") {
            if (line.startsWith("\\\\", characterIndex)) {
              currentPart += "\\";
              characterIndex += 2;
              continue;
            }
            if (line.startsWith("\\n", characterIndex)) {
              currentPart += "\n";
              characterIndex += 2;
              continue;
            }
            if (line.startsWith(`\\${delimiter}`, characterIndex)) {
              currentPart += delimiter;
              characterIndex += delimiter.length + 1;
              continue;
            }
            currentPart += character;
            characterIndex += 1;
            continue;
          }
          if (line.startsWith(delimiter, characterIndex)) {
            lineParts.push(currentPart);
            currentPart = "";
            characterIndex += delimiter.length;
            continue;
          }
          currentPart += character;
          characterIndex += 1;
        }
        lineParts.push(currentPart);
        return lineParts;
      },

      resolveBulkEditLineParts(line, delimiter) {
        const lineParts = this.splitBulkEditLine(line, delimiter);
        if (lineParts.length <= this.columns.length) {
          return lineParts;
        }
        const lastColumnIndex = Math.max(this.columns.length - 1, 0);
        const leadingParts = lineParts.slice(0, lastColumnIndex);
        const lastColumnText = lineParts.slice(lastColumnIndex).join(delimiter);
        return leadingParts.concat(lastColumnText);
      },

      parseBulkEditRows(text) {
        const delimiter = this.bulkEditDelimiterResolver();
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
              lineParts[columnIndex] ?? "",
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

      escapeBulkEditValue(rawText, delimiter, shouldEscapeDelimiter) {
        let escapedText = String(rawText ?? "").replaceAll("\\", "\\\\");
        if (shouldEscapeDelimiter && delimiter !== "") {
          escapedText = escapedText.replaceAll(delimiter, `\\${delimiter}`);
        }
        escapedText = escapedText.replaceAll("\r\n", "\n");
        escapedText = escapedText.replaceAll("\r", "\n");
        escapedText = escapedText.replaceAll("\n", "\\n");
        return escapedText;
      },

      serializeBulkEditRows(rows) {
        const delimiter = this.bulkEditDelimiterResolver();
        const textLines = [];
        for (const row of rows) {
          const cellTexts = [];
          let hasValue = false;
          for (
            let columnIndex = 0;
            columnIndex < this.columns.length;
            columnIndex += 1
          ) {
            const column = this.columns[columnIndex];
            const cellText = this.formatBulkEditValue(column, row);
            const isLastColumn = columnIndex === this.columns.length - 1;
            cellTexts.push(
              this.escapeBulkEditValue(cellText, delimiter, !isLastColumn),
            );
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

      mergeBulkEditRow(existingRow, parsedRow) {
        for (const column of this.columns) {
          existingRow[column.id] = parsedRow[column.id];
        }
        return existingRow;
      },

      replaceRowsFromBulkText(rows, text) {
        const parsedRows = this.parseBulkEditRows(text);
        const rowIndexesByText = new Map();
        for (let rowIndex = 0; rowIndex < rows.length; rowIndex += 1) {
          const rowText = this.serializeBulkEditRows([rows[rowIndex]]);
          if (!rowIndexesByText.has(rowText)) {
            rowIndexesByText.set(rowText, []);
          }
          rowIndexesByText.get(rowText).push(rowIndex);
        }
        const matchedRowIndexesByParsedIndex = new Array(parsedRows.length);
        const usedRowIndexes = new Set();
        for (
          let parsedRowIndex = 0;
          parsedRowIndex < parsedRows.length;
          parsedRowIndex += 1
        ) {
          const parsedRowText = this.serializeBulkEditRows([
            parsedRows[parsedRowIndex],
          ]);
          const matchingRowIndexes = rowIndexesByText.get(parsedRowText) ?? [];
          const matchedRowIndex = matchingRowIndexes.find(
            (rowIndex) => !usedRowIndexes.has(rowIndex),
          );
          if (matchedRowIndex === undefined) {
            continue;
          }
          matchedRowIndexesByParsedIndex[parsedRowIndex] = matchedRowIndex;
          usedRowIndexes.add(matchedRowIndex);
        }
        if (parsedRows.length === rows.length) {
          const unmatchedParsedRowIndexes = [];
          for (
            let parsedRowIndex = 0;
            parsedRowIndex < parsedRows.length;
            parsedRowIndex += 1
          ) {
            if (matchedRowIndexesByParsedIndex[parsedRowIndex] === undefined) {
              unmatchedParsedRowIndexes.push(parsedRowIndex);
            }
          }
          const unusedRowIndexes = [];
          for (let rowIndex = 0; rowIndex < rows.length; rowIndex += 1) {
            if (!usedRowIndexes.has(rowIndex)) {
              unusedRowIndexes.push(rowIndex);
            }
          }
          for (
            let unmatchedIndex = 0;
            unmatchedIndex < unmatchedParsedRowIndexes.length;
            unmatchedIndex += 1
          ) {
            const parsedRowIndex = unmatchedParsedRowIndexes[unmatchedIndex];
            matchedRowIndexesByParsedIndex[parsedRowIndex] =
              unusedRowIndexes[unmatchedIndex];
          }
        }
        const reconciledRows = [];
        for (
          let parsedRowIndex = 0;
          parsedRowIndex < parsedRows.length;
          parsedRowIndex += 1
        ) {
          const parsedRow = parsedRows[parsedRowIndex];
          const matchedRowIndex =
            matchedRowIndexesByParsedIndex[parsedRowIndex];
          if (matchedRowIndex === undefined) {
            reconciledRows.push(parsedRow);
            continue;
          }
          reconciledRows.push(
            this.mergeBulkEditRow(rows[matchedRowIndex], parsedRow),
          );
        }
        rows.splice(0, rows.length, ...reconciledRows);
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
