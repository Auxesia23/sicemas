import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import * as v from "valibot";
import { ProfileResponseSchema } from "$lib/schemas/user";
import { auth } from "$lib/states/auth.svelte";

export const load: PageLoad = async ({ fetch, parent }) => {
  await parent();

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }

  const PageFetch = createPageLoadFetch(fetch);
  const res = await PageFetch("/api/users/profile");

  if (!res.ok) {
    throw error(res.status, await res.text());
  }

  const rawProfile = await res.json();
  const parsed = v.safeParse(ProfileResponseSchema, rawProfile);

  if (!parsed.success) {
    console.error("Valibot Error (Profile):", parsed.issues);
    throw error(500, "Gagal memvalidasi struktur data Profil dari server");
  }

  return {
    profile: parsed.output,
  };
};
