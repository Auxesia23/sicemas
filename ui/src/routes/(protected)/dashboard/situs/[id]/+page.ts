import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import * as v from "valibot";
import { SitusDetailResponseSchema } from "$lib/schemas/sites";
import { auth } from "$lib/states/auth.svelte";

export const load: PageLoad = async ({ params, fetch, parent, depends }) => {
  await parent();
  depends("situs:detail");

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }
  if (!auth.hasAnyPermission(["situs:read_all", "situs:read_own"])) {
    throw error(403, "Anda tidak memiliki izin ke halaman ini");
  }

  const PageFetch = createPageLoadFetch(fetch);
  const res = await PageFetch(`/api/situs/${params.id}`);

  if (!res.ok) {
    throw error(res.status, await res.text());
  }

  const rawData = await res.json();
  const parsed = v.safeParse(SitusDetailResponseSchema, rawData);

  if (!parsed.success) {
    console.error("Valibot Error (Situs Detail):", parsed.issues);
    throw error(500, "Gagal memvalidasi struktur data Situs dari server");
  }

  return {
    situs: parsed.output,
  };
};
