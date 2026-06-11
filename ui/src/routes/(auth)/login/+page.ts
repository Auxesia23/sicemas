import type { PageLoad } from "./$types";
import { auth } from "$lib/states/auth.svelte";
import { goto } from "$app/navigation";
export const load: PageLoad = async ({ parent }) => {
  await parent();
  if (auth.user) {
    goto("/dashboard");
  }
  return {};
};
