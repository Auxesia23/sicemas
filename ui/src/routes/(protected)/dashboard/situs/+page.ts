import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import { SitusKeagamaanSchema, JenisSitusSchema } from "$lib/schemas/sites";
import * as v from "valibot";
import { auth } from "$lib/states/auth.svelte";

export const load: PageLoad = async ({ fetch, parent, depends }) => {
  await parent();
  depends("data:situs");

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }
  if (!auth.hasAnyPermission(["situs:read_all", "situs:read_own"])) {
    throw error(403, "Anda tidak memiliki izin ke halaman ini");
  }

  const PageFetch = createPageLoadFetch(fetch);

  const [situsRes, jenisSitusRes] = await Promise.all([
    PageFetch("/api/situs"),
    PageFetch("/api/jenis-situs"),
  ]);

  if (!situsRes.ok) {
    throw error(situsRes.status, await situsRes.text());
  }
  if (!jenisSitusRes.ok) {
    throw error(jenisSitusRes.status, await jenisSitusRes.text());
  }

  const rawSitus = await situsRes.json();
  const rawJenisSitus = await jenisSitusRes.json();

  const parsedSitus = v.safeParse(v.array(SitusKeagamaanSchema), rawSitus);
  const parsedJenisSitus = v.safeParse(
    v.array(JenisSitusSchema),
    rawJenisSitus,
  );

  if (!parsedSitus.success) {
    throw error(
      500,
      "Gagal memvalidasi struktur data Situs Keagamaan dari server",
    );
  }

  if (!parsedJenisSitus.success) {
    throw error(500, "Gagal memvalidasi struktur data Jenis Situs dari server");
  }

  return {
    situs: parsedSitus.output,
    jenisSitus: parsedJenisSitus.output,
  };
};
