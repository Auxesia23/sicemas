import type { PageLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { auth } from "$lib/states/auth.svelte";
import { createPageLoadFetch } from "$lib/api";
import { DashboardResponseSchema } from "$lib/schemas/dashboard";
import { goto } from "$app/navigation";

import * as v from "valibot";

export const load: PageLoad = async ({ parent, depends, fetch }) => {
  depends("data:dashboard");
  await parent();

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }
  if (!auth.hasAllPermissions(["dashboard:read"])) {
    goto("/dashboard/situs");
  }

  const PageFetch = createPageLoadFetch(fetch);
  const res = await PageFetch("/api/dashboard");
  if (!res.ok) {
    throw error(res.status, await res.text());
  }

  const rawData = await res.json();

  const parsed = v.safeParse(DashboardResponseSchema, rawData);

  if (!parsed.success) {
    throw error(500, "Gagal memvalidasi struktur data Dashboard dari server");
  }

  return {
    dashboard: parsed.output,
  };
};
