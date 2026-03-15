import { error } from '@sveltejs/kit';
import { hasAllPermissions, hasAnyPermission } from '$lib/permissions';

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals }) {
	if (!hasAnyPermission(locals.user.permissions, ['situs:read_all', 'situs:read_own'])) {
		throw error(403);
	}
}
