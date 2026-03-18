import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

/** @type {import('./$types').PageLoad} */
export async function load({ params, fetch }) {
	try {
		const response = await fetch(`${env.PUBLIC_API_URL}/public/situs/${params.id}`);

		if (!response.ok) {
			if (response.status === 404) {
				throw error(404, 'Situs tidak ditemukan');
			}
			throw error(500, 'Gagal memuat data situs');
		}

		const site = await response.json();

		return {
			site
		};
	} catch (err) {
		if (err.status) {
			throw err;
		}
		throw error(500, 'Terjadi kesalahan saat memuat data');
	}
}
