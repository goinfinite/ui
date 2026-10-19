UiToolset.RegisterAlpineState(() => {
  const collapsibleTextEllipsis = "…";
  const collapsibleTextMeasurementTolerancePx = 1;

  Alpine.data(
    "collapsibleText",
    (fullText, maxVisibleLines, expandLabel, collapseLabel) => ({
      fullText: String(fullText ?? ""),
      truncatedText: "",
      hasToggle: false,
      isExpanded: false,
      maxVisibleLines,
      expandLabel,
      collapseLabel,
      textResizeObserver: null,
      observedWidth: 0,
      isDestroyed: false,
      truncationRunId: 0,

      get displayText() {
        if (!this.hasToggle || this.isExpanded) {
          return this.fullText;
        }
        return this.truncatedText + collapsibleTextEllipsis;
      },

      get toggleLabel() {
        return this.isExpanded ? this.collapseLabel : this.expandLabel;
      },

      init() {
        this.truncateText();
        this.observeTextResize();
      },

      destroy() {
        this.isDestroyed = true;
        this.textResizeObserver?.disconnect();
      },

      toggleTextDisplay() {
        this.isExpanded = !this.isExpanded;
        this.$nextTick(() => this.truncateText());
      },

      async truncateText() {
        if (this.isDestroyed || this.isExpanded) {
          return;
        }
        const runId = ++this.truncationRunId;
        const isStale = () =>
          this.isDestroyed || runId !== this.truncationRunId;
        const clipElement = this.$refs.textContent;
        if (!clipElement) {
          return;
        }
        const lineHeight = parseFloat(getComputedStyle(clipElement).lineHeight);
        const maxHeightPx = lineHeight * this.maxVisibleLines;
        this.hasToggle = false;
        this.truncatedText = "";
        await this.$nextTick();
        if (isStale()) {
          return;
        }
        if (this.clipElementFits(clipElement, maxHeightPx)) {
          return;
        }
        this.hasToggle = true;
        let lowerBound = 0;
        let upperBound = this.fullText.length;
        let bestFit = 0;
        while (lowerBound <= upperBound) {
          const candidateLength = Math.floor((lowerBound + upperBound) / 2);
          this.truncatedText = this.fullText.slice(0, candidateLength);
          await this.$nextTick();
          if (isStale()) {
            return;
          }
          if (this.clipElementFits(clipElement, maxHeightPx)) {
            bestFit = candidateLength;
            lowerBound = candidateLength + 1;
            continue;
          }
          upperBound = candidateLength - 1;
        }
        while (bestFit > 0) {
          this.truncatedText = this.fullText.slice(0, bestFit);
          await this.$nextTick();
          if (isStale()) {
            return;
          }
          if (this.candidateFits(clipElement, maxHeightPx)) {
            break;
          }
          bestFit -= 1;
        }
        const rawCut = this.fullText.slice(0, bestFit);
        const trimmedCut = this.wordTrimmedCut(bestFit);
        this.truncatedText = trimmedCut;
        if (trimmedCut === rawCut) {
          return;
        }
        await this.$nextTick();
        if (isStale()) {
          return;
        }
        if (!this.candidateFits(clipElement, maxHeightPx)) {
          this.truncatedText = rawCut;
        }
      },

      candidateFits(clipElement, maxHeightPx) {
        if (!this.clipElementFits(clipElement, maxHeightPx)) {
          return false;
        }
        return this.suffixSharesLastTextLine();
      },

      suffixSharesLastTextLine() {
        const textSpan = this.$refs.textSpan;
        const toggleButton = this.$refs.toggleButton;
        if (!textSpan || !toggleButton) {
          return true;
        }
        const textRects = textSpan.getClientRects();
        if (textRects.length === 0) {
          return true;
        }
        const lastTextRect = textRects[textRects.length - 1];
        const toggleRect = toggleButton.getBoundingClientRect();
        return toggleRect.top < lastTextRect.bottom;
      },

      clipElementFits(clipElement, maxHeightPx) {
        return (
          clipElement.offsetHeight <=
          maxHeightPx + collapsibleTextMeasurementTolerancePx
        );
      },

      wordTrimmedCut(cutIndex) {
        const rawCut = this.fullText.slice(0, cutIndex);
        if (cutIndex >= this.fullText.length) {
          return rawCut;
        }
        const trailingWordStart = rawCut.search(/\s+\S*$/);
        if (trailingWordStart < 0) {
          return rawCut;
        }
        return rawCut.slice(0, trailingWordStart);
      },

      observeTextResize() {
        this.textResizeObserver?.disconnect();
        if (typeof ResizeObserver === "undefined") {
          return;
        }
        this.observedWidth = this.$el.getBoundingClientRect().width;
        this.textResizeObserver = new ResizeObserver((resizeEntries) => {
          const currentWidth = resizeEntries[0].contentRect.width;
          if (currentWidth === this.observedWidth) {
            return;
          }
          this.observedWidth = currentWidth;
          this.truncateText();
        });
        this.textResizeObserver.observe(this.$el);
      },
    }),
  );
});
