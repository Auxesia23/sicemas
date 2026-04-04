import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import * as v from "valibot";
import { SitusPublikDetailResponseSchema } from "$lib/schemas/sites";

export const load: PageLoad = async ({ params, fetch, parent }) => {
  await parent();
  const PageFetch = createPageLoadFetch(fetch);

  let apiUrl = `/api/public/situs/${params.id}`;

  const res = await PageFetch(apiUrl);

  if (!res.ok) {
    throw error(res.status, await res.text());
  }

  const rawData = await res.json();

  const parsed = v.safeParse(SitusPublikDetailResponseSchema, rawData);

  if (!parsed.success) {
    console.error(parsed.issues);
    throw error(500, "Format data dari server tidak sesuai");
  }

  return {
    site: parsed.output,
  };
};
