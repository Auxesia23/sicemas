import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import { JenisSitusSchema } from "$lib/schemas/sites";
import { auth } from "$lib/states/auth.svelte";
import * as v from "valibot";

export const load: PageLoad = async ({ fetch, parent, depends }) => {
  await parent();
  depends("data:jenis-situs");

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }
  if (!auth.hasAllPermissions(["jenis-situs:read"])) {
    throw error(403, "Anda tidak memiliki izin ke halaman ini");
  }

  const PageFetch = createPageLoadFetch(fetch);
  const res = await PageFetch("/api/jenis-situs");

  if (!res.ok) {
    throw error(res.status, await res.text());
  }

  const rawData = await res.json();
  const parsed = v.safeParse(v.array(JenisSitusSchema), rawData);

  if (!parsed.success) {
    throw error(500, "Gagal memvalidasi struktur data Jenis Situs dari server");
  }

  return {
    jenisSitus: parsed.output,
  };
};
