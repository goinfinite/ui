UiToolset.RegisterAlpineState(() => {
  Alpine.data("multiColumnRepeatableFieldset", (columns, newRowFunc) => {
    const rowKeys = new WeakMap();
    let nextRowKey = 1;

    return {
      columns: Array.isArray(columns) ? columns : [],
      newRowFunc: typeof newRowFunc === "string" ? newRowFunc : "",

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
          if (Object.hasOwn(generatedRow, column.id)) {
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
  });
});
