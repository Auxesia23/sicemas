// src/routes/error-handler/+page.js
import { error } from '@sveltejs/kit';

export function load({ url }) {
	const code = Number(url.searchParams.get('code')) || 500;
	const message = url.searchParams.get('message') || 'Terjadi Kesalahan Sistem';

	// Melempar error di dalam load() akan otomatis mengaktifkan +error.svelte terdekat
	throw error(code, message);
}
