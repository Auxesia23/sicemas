import { redirect } from '@sveltejs/kit';

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals, url }) {
	return {
		user: locals.user
	};
}
