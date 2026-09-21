import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import yaml from "js-yaml";

const levelRank = { fast: 1, standard: 2, exhaustive: 3 };
const registryPath = fileURLToPath(
  new URL("../registry.yaml", import.meta.url),
);

function fail(message) {
  console.error(`registry: ${message} (try: tests.sh list)`);
  process.exit(2);
}

function parseFilters(argv) {
  const filters = { level: null, feature: null, scope: null };
  for (const arg of argv) {
    const separator = arg.indexOf("=");
    const key = arg.slice(0, separator).replace(/^-{1,2}/, "");
    if (!(key in filters)) fail(`unknown argument: ${arg}`);
    filters[key] = arg
      .slice(separator + 1)
      .split(",")
      .map((value) => value.trim());
  }
  return filters;
}

function validateFilters(registry, filters) {
  const featureEntries = Object.entries(registry.features);
  const known = {
    level: new Set(registry.levels),
    scope: new Set(
      featureEntries.flatMap(([, entries]) =>
        entries.map((entry) => entry.scope),
      ),
    ),
    feature: new Set(featureEntries.map(([feature]) => feature)),
  };
  for (const [key, values] of Object.entries(filters)) {
    for (const value of values ?? []) {
      if (!known[key].has(value)) fail(`unknown ${key}: ${value}`);
    }
  }
}

function selectMatchingNodes(registry, filters, defaultRank) {
  const maxRank = filters.level ? levelRank[filters.level[0]] : defaultRank;
  const nodes = [];
  for (const [feature, entries] of Object.entries(registry.features)) {
    if (filters.feature && !filters.feature.includes(feature)) continue;
    for (const entry of entries) {
      if (filters.scope && !filters.scope.includes(entry.scope)) continue;
      if (levelRank[entry.level] > maxRank) continue;
      nodes.push({ feature, ...entry });
    }
  }
  return nodes.sort((left, right) =>
    `${left.feature}:${left.scope}/${left.level}`.localeCompare(
      `${right.feature}:${right.scope}/${right.level}`,
    ),
  );
}

function printList(nodes) {
  for (const node of nodes) {
    const marker = node.provisional ? "*" : "";
    const label = `${node.feature}:${node.scope}/${node.level}${marker}`;
    console.log(`${label.padEnd(36)} ${node.run}`);
  }
  if (nodes.some((node) => node.provisional)) {
    console.error(
      "* provisional: implementation-coupled suite awaiting refactoring",
    );
  }
}

function printSelection(nodes) {
  for (const node of nodes) {
    const label = `${node.feature}:${node.scope}/${node.level}`;
    console.log(`${label}\t${node.run}\t${node.requires_demo === true}`);
  }
}

const [command, ...argv] = process.argv.slice(2);
if (!["list", "select"].includes(command ?? ""))
  fail(`unknown command: ${command}`);

const registry = yaml.load(readFileSync(registryPath, "utf8"));
const filters = parseFilters(argv);
validateFilters(registry, filters);

if (command === "list") {
  printList(selectMatchingNodes(registry, filters, levelRank.exhaustive));
  process.exit(0);
}
printSelection(selectMatchingNodes(registry, filters, levelRank.standard));
