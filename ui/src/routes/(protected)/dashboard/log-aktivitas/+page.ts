import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import { ActivityResponseSchema } from "$lib/schemas/activity";
import * as v from "valibot";
import { auth } from "$lib/states/auth.svelte";

export const load: PageLoad = async ({ fetch, parent, depends, url }) => {
  await parent();
  depends("data:activities");

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }

  if (!auth.hasAllPermissions(["activity:read"])) {
    throw error(403, "Anda tidak memiliki izin ke halaman ini");
  }

  let targetDate = url.searchParams.get("date");
  if (!targetDate) {
    const today = new Date();
    const year = today.getFullYear();
    const month = String(today.getMonth() + 1).padStart(2, "0");
    const day = String(today.getDate()).padStart(2, "0");
    targetDate = `${year}-${month}-${day}`;
  }

  const PageFetch = createPageLoadFetch(fetch);

  const res = await PageFetch(`/api/activities?date=${targetDate}`);

  if (!res.ok) {
    throw error(res.status, await res.text());
  }

  const rawData = await res.json();
  const parsed = v.safeParse(v.array(ActivityResponseSchema), rawData);

  if (!parsed.success) {
    console.error(parsed.issues);
    throw error(
      500,
      "Gagal memvalidasi struktur data Log Aktivitas dari server",
    );
  }

  return {
    activities: parsed.output,
    selectedDate: targetDate,
  };
};
