import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import * as v from "valibot";
import { SitusDetailResponseSchema } from "$lib/schemas/sites";
import { auth } from "$lib/states/auth.svelte";

export const load: PageLoad = async ({ params, fetch, parent }) => {
  await parent();

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }
  if (!auth.hasAllPermissions(["situs:update"])) {
    throw error(403, "Anda tidak memiliki izin ke halaman ini");
  }

  const PageFetch = createPageLoadFetch(fetch);
  const res = await PageFetch(`/api/situs/${params.id}`);

  if (!res.ok) throw error(res.status, "Gagal memuat data situs");

  const rawData = await res.json();
  const parsed = v.safeParse(SitusDetailResponseSchema, rawData);

  if (!parsed.success) throw error(500, "Struktur data tidak valid");

  return { situs: parsed.output };
};
