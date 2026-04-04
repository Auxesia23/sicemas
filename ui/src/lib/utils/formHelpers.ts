import * as v from "valibot";

export function addToArray(array: string[], value: string): string {
  if (value.trim() && !array.includes(value.trim())) {
    array.push(value.trim());
    return "";
  }
  return value;
}

export function removeFromArray(array: string[], index: number): void {
  array.splice(index, 1);
}

export function handleEnterKey(
  event: KeyboardEvent,
  actionFn: () => void,
): void {
  if (event.key === "Enter") {
    event.preventDefault();
    actionFn();
  }
}

export function clearNestedError(
  errorsObj: Record<string, any>,
  ...keys: string[]
): void {
  let current = errorsObj;
  for (let i = 0; i < keys.length - 1; i++) {
    if (!current[keys[i]]) return;
    current = current[keys[i]];
  }
  const lastKey = keys[keys.length - 1];
  if (current && current[lastKey]) {
    current[lastKey] = null;
  }
}

export function mapFlatErrors(
  flatIssues: v.FlatErrors<any>,
  targetErrorObj: Record<string, any>,
) {
  if (flatIssues.nested) {
    Object.entries(flatIssues.nested).forEach(([key, messages]) => {
      const message = messages?.[0];
      if (key.includes(".")) {
        const [parent, child] = key.split(".");
        if (!targetErrorObj[parent]) targetErrorObj[parent] = {};
        targetErrorObj[parent][child] = message;
      } else {
        targetErrorObj[key] = message;
      }
    });
  }
}
