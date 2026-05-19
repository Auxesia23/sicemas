import { auth } from "$lib/states/auth.svelte";
import type { LayoutLoad } from "./$types";

export const ssr = false;

export const load: LayoutLoad = async ({ fetch, url }) => {
  if (url.pathname === "/verify-stepup") {
    return {};
  }
  await auth.initDevice();
  console.log("Fetch User  dari load");
  await auth.fetchUser(fetch);
  console.log("fetch user dari load selesai");

  return {};
};
