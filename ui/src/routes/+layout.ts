import { auth } from "$lib/states/auth.svelte";
import type { LayoutLoad } from "./$types";

export const ssr = false;

export const load: LayoutLoad = async ({ fetch, url }) => {
  await auth.initDevice();
  await auth.fetchUser(fetch);

  return {};
};
