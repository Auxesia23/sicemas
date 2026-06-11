import type { PageLoad } from "./$types";
import { auth } from "$lib/states/auth.svelte";
import { error } from "@sveltejs/kit";
export const load: PageLoad = async ({ parent }) => {
  await parent();
  if (!auth.user || auth.isStepUpRequired == false) {
    throw error(404);
  }
  return {};
};
