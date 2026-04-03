import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import { JenisSitusSchema } from "$lib/schemas/sites";
import * as v from "valibot";
import { auth } from "$lib/states/auth.svelte";

export const load: PageLoad = async ({ fetch, parent }) => {
  await parent();

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }
  if (!auth.hasAllPermissions(["situs:create"])) {
    throw error(403, "Anda tidak memiliki izin ke halaman ini");
  }

  const PageFetch = createPageLoadFetch(fetch);
  const res = await PageFetch("/api/jenis-situs");

  if (!res.ok) throw error(res.status, await res.text());

  const rawData = await res.json();
  const parsed = v.safeParse(v.array(JenisSitusSchema), rawData);

  if (!parsed.success) throw error(500, "Gagal memvalidasi struktur data");

  return { jenisSitusList: parsed.output };
};
