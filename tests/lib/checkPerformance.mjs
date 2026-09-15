import { readFileSync } from "node:fs";
import yaml from "js-yaml";

const results = JSON.parse(readFileSync(new URL("../performance-results.json", import.meta.url), "utf8"));
const budgets = yaml.load(readFileSync(new URL("../golden.yaml", import.meta.url), "utf8"));

let fatalReached = false;

for (const metric of Object.keys(budgets)) {
  if (!(metric in results)) {
    console.error(`FATAL ${metric}: no performance result recorded`);
    fatalReached = true;
  }
}

for (const [metric, elapsedMs] of Object.entries(results)) {
  const budget = budgets[metric];
  if (!budget) {
    console.log(`WARN ${metric}: ${elapsedMs}ms — no golden budget recorded yet`);
    continue;
  }
  const relativeWarnMs = budget.relative_warn_percent
    ? budget.baseline_p50_ms * (1 + budget.relative_warn_percent / 100)
    : Infinity;
  const warnAtMs = Math.min(budget.warn_ms ?? Infinity, relativeWarnMs);
  if (budget.fatal_ms && elapsedMs >= budget.fatal_ms) {
    console.error(`FATAL ${metric}: ${elapsedMs}ms >= ${budget.fatal_ms}ms budget`);
    fatalReached = true;
    continue;
  }
  if (elapsedMs >= warnAtMs) {
    console.log(`WARNING ${metric}: ${elapsedMs}ms >= ${Math.round(warnAtMs)}ms (run still passes)`);
    continue;
  }
  if (budget.info_ms && elapsedMs >= budget.info_ms) {
    console.log(`INFO ${metric}: ${elapsedMs}ms >= ${budget.info_ms}ms`);
    continue;
  }
  console.log(`PASS ${metric}: ${elapsedMs}ms (baseline ${budget.baseline_p50_ms}ms)`);
}

process.exit(fatalReached ? 1 : 0);
