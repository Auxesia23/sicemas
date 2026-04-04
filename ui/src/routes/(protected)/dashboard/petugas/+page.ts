import type { PageLoad } from "./$types";
import { createPageLoadFetch } from "$lib/api";
import { error } from "@sveltejs/kit";
import { UserResponseSchema, RoleResponseSchema } from "$lib/schemas/user";
import * as v from "valibot";
import { auth } from "$lib/states/auth.svelte";

export const load: PageLoad = async ({ fetch, parent, depends }) => {
  await parent();
  depends("data:users");

  if (!auth.user) {
    throw error(401, "Sesi tidak valid atau belum login");
  }
  if (!auth.hasAllPermissions(["user:read"])) {
    throw error(403, "Anda tidak memiliki izin ke halaman ini");
  }

  const PageFetch = createPageLoadFetch(fetch);

  const [usersRes, rolesRes] = await Promise.all([
    PageFetch("/api/users"),
    PageFetch("/api/roles"),
  ]);

  if (!usersRes.ok) {
    throw error(usersRes.status, await usersRes.text());
  }
  if (!rolesRes.ok) {
    throw error(rolesRes.status, await rolesRes.text());
  }

  const rawUsers = await usersRes.json();
  const rawRoles = await rolesRes.json();

  const parsedUsers = v.safeParse(v.array(UserResponseSchema), rawUsers);
  const parsedRoles = v.safeParse(v.array(RoleResponseSchema), rawRoles);

  if (!parsedUsers.success) {
    throw error(500, "Gagal memvalidasi struktur data Users dari server");
  }
  if (!parsedRoles.success) {
    throw error(500, "Gagal memvalidasi struktur data Roles dari server");
  }

  return {
    users: parsedUsers.output,
    roles: parsedRoles.output,
  };
};
