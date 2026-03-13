import { error } from '@sveltejs/kit';
import { hasAllPermissions } from '$lib/permissions';

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals }) {
	if (!hasAllPermissions(locals.user.permissions, ['situs:create'])) {
		throw error(403);
	}
}
