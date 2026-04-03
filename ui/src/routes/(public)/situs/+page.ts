import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import * as v from "valibot";
import { SitusPublikResponseSchema } from "$lib/schemas/sites";

export const load: PageLoad = async ({ url, fetch, parent }) => {
  await parent();
  const PageFetch = createPageLoadFetch(fetch);

  const lat = url.searchParams.get("lat");
  const lng = url.searchParams.get("lng");

  let apiUrl = `/api/public/situs`;
  if (lat && lng) {
    apiUrl += `?lat=${lat}&lng=${lng}`;
  }

  const res = await PageFetch(apiUrl);

  if (!res.ok) {
    throw error(res.status, await res.text());
  }

  const rawData = await res.json();

  const parsed = v.safeParse(SitusPublikResponseSchema, rawData);

  if (!parsed.success) {
    console.error("Valibot Error (Public Sites):", parsed.issues);
    return { sites: [] };
  }

  return {
    sites: parsed.output,
  };
};
