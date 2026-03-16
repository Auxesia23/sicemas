import { redirect } from '@sveltejs/kit';
import { hasAnyPermission } from '$lib/permissions';

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals, url }) {
	if (!locals.user) {
		throw redirect(303, '/login');
	}

	if (url.pathname === '/dashboard') {
		if (!hasAnyPermission(locals.user.permissions, ['situs:verify', 'user:create'])) {
			throw redirect(303, '/dashboard/situs');
		}
	}
}
