import { env } from '$env/dynamic/private';
import { PUBLIC_API_URL } from '$env/static/public';

/** @type {import('./$types').PageServerLoad} */
export async function load({ fetch, url }) {
	const lat = url.searchParams.get('lat');
	const lng = url.searchParams.get('lng');

	let apiUrl = `${PUBLIC_API_URL}/public/situs`;
	if (lat && lng) {
		apiUrl += `?lat=${lat}&lng=${lng}`;
	}
	try {
		const response = await fetch(apiUrl);
		if (!response.ok) {
			if (response.status === 404) {
				throw error(404, 'Situs tidak ditemukan');
			}
			throw error(500, 'Gagal memuat data situs');
		}
		const sites = await response.json();
		return { sites };
	} catch (err) {
		if (err.status) {
			throw err;
		}
		throw error(500, 'Terjadi kesalahan saat memuat data');
	}
}
