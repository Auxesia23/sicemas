import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import * as v from "valibot";
import { auth } from "$lib/states/auth.svelte";
import { RoleResponseSchema } from "$lib/schemas/user";
import { PolicySchema } from "$lib/schemas/policy";

export const load: PageLoad = async ({ fetch, parent, depends }) => {
  await parent();
  depends("data:permissions");

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }
  if (!auth.hasAllPermissions(["role:read", "policy:read"])) {
    throw error(403, "Anda tidak memiliki izin ke halaman ini");
  }

  const PageFetch = createPageLoadFetch(fetch);

  const [rolesRes, policiesRes] = await Promise.all([
    PageFetch("/api/roles"),
    PageFetch("/api/policies"),
  ]);

  if (!rolesRes.ok) throw error(rolesRes.status, await rolesRes.text());
  if (!policiesRes.ok)
    throw error(policiesRes.status, await policiesRes.text());

  const rawRoles = await rolesRes.json();
  const rawPolicies = await policiesRes.json();

  const parsedRoles = v.safeParse(v.array(RoleResponseSchema), rawRoles);
  const parsedPolicies = v.safeParse(v.array(PolicySchema), rawPolicies);

  if (!parsedRoles.success) {
    throw error(500, "Gagal memvalidasi struktur data Roles dari server");
  }
  if (!parsedPolicies.success) {
    throw error(500, "Gagal memvalidasi struktur data Policies dari server");
  }

  return {
    roles: parsedRoles.output,
    policies: parsedPolicies.output,
  };
};
