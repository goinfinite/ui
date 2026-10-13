UiToolset.RegisterAlpineState(() => {
  const resourceUsageWindowSize = 20;
  const resourceUsageIntervalMs = 2000;

  window.chartDemoTooltipContent = (row) => {
    const received = Number(row.received ?? 0);
    const transmitted = Number(row.transmitted ?? 0);
    const total = received + transmitted;
    const receivedShare = total === 0 ? 0 : (received / total) * 100;
    return (
      `<p class="font-bold">${row.month}</p>` +
      `<p class="text-neutral-300">Received ${received} GiB</p>` +
      `<p class="text-neutral-300">Transmitted ${transmitted} GiB</p>` +
      `<p class="text-secondary-200 font-bold">${receivedShare.toFixed(1)}% received</p>`
    );
  };

  function nextUsageFraction(previousFraction) {
    const driftFraction = (Math.random() - 0.5) * 0.16;
    return Math.min(Math.max(previousFraction + driftFraction, 0.02), 0.98);
  }

  Alpine.data("chartDemoResourceUsage", () => ({
    resourceUsageRows: [],
    resourceUsageIntervalId: null,

    init() {
      const nowMs = Date.now();
      let cpuFraction = 0.35;
      let memoryFraction = 0.6;
      for (
        let rowIndex = resourceUsageWindowSize;
        rowIndex > 0;
        rowIndex -= 1
      ) {
        cpuFraction = nextUsageFraction(cpuFraction);
        memoryFraction = nextUsageFraction(memoryFraction);
        this.resourceUsageRows.push({
          time: nowMs - rowIndex * resourceUsageIntervalMs,
          cpu: cpuFraction,
          memory: memoryFraction,
        });
      }
      this.resourceUsageIntervalId = setInterval(
        () => this.appendResourceUsageRow(),
        resourceUsageIntervalMs,
      );
    },

    destroy() {
      clearInterval(this.resourceUsageIntervalId);
    },

    appendResourceUsageRow() {
      const lastRow = this.resourceUsageRows[this.resourceUsageRows.length - 1];
      this.resourceUsageRows.push({
        time: Date.now(),
        cpu: nextUsageFraction(lastRow.cpu),
        memory: nextUsageFraction(lastRow.memory),
      });
      if (this.resourceUsageRows.length > resourceUsageWindowSize) {
        this.resourceUsageRows.shift();
      }
    },
  }));
});
