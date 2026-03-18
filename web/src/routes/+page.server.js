import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

/** @type {import('./$types').PageServerLoad} */
export async function load({ fetch }) {
	let apiUrl = `${env.PUBLIC_API_URL}/public/stats`;
	try {
		const response = await fetch(apiUrl);
		if (!response.ok) {
			throw error(500, 'Gagal memuat data situs');
		}
		const stats = await response.json();
		return { stats };
	} catch (err) {
		if (err.status) {
			throw err;
		}
		throw error(500, 'Terjadi kesalahan saat memuat data');
	}
}
