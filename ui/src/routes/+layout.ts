import { auth } from "$lib/states/auth.svelte";
import type { LayoutLoad } from "./$types";

export const ssr = false;

export const load: LayoutLoad = async ({ fetch, url }) => {
  if (url.pathname === "/verify-stepup") {
    return {};
  }
  await auth.initDevice();
  await auth.fetchUser(fetch);

  return {};
};
