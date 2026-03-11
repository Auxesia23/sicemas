// API Service - Centralized HTTP client for all API calls
import config from '$lib/config';
import { browser } from '$app/environment'; // Pastikan ini di-import

class ApiService {
	constructor() {
		// State untuk mengunci dan mengantrekan request
		this.isRefreshing = false;
		this.failedQueue = [];
	}

	// Fungsi untuk memproses semua request yang ngantri
	processQueue(error) {
		this.failedQueue.forEach((prom) => {
			if (error) {
				prom.reject(error);
			} else {
				prom.resolve();
			}
		});
		this.failedQueue = [];
	}

	async request(endpoint, options = {}) {
		const url = `${config.apiUrl}${endpoint}`;
		const requestOptions = {
			...config.defaultFetchOptions,
			...options,
			headers: {
				...config.defaultHeaders,
				...options.headers
			}
		};

		let response = await fetch(url, requestOptions);

		// ==========================================
		// 1. PENANGANAN 401 UNAUTHORIZED (Refresh Token)
		// ==========================================
		if (response.status === 401 && endpoint !== '/auth/refresh') {
			// JIKA SEDANG REFRESH: Masukkan request ini ke antrean (Queue)
			if (this.isRefreshing) {
				return new Promise((resolve, reject) => {
					this.failedQueue.push({ resolve, reject });
				})
					.then(() => {
						// Setelah token berhasil di-refresh oleh request lain, jalankan ulang fetch-nya
						return fetch(url, requestOptions);
					})
					.catch((err) => {
						return Promise.reject(err);
					});
			}

			// KUNCI PROSES REFRESH (Request pertama yang dapat 401 akan masuk ke sini)
			this.isRefreshing = true;

			const success = await this.refreshToken();

			if (success) {
				// BUKA KUNCI dan jalankan semua request yang ada di antrean
				this.isRefreshing = false;
				this.processQueue(null);

				// Ulangi request utama yang gagal tadi
				response = await fetch(url, requestOptions);
			} else {
				// Refresh gagal total, tolak semua request di antrean
				this.isRefreshing = false;
				this.processQueue(new Error('Refresh token failed'));
			}
		}

		// ==========================================
		// 2. PENANGANAN 403 & 500 SECARA GLOBAL
		// ==========================================
		if (response.status === 403 || response.status >= 500) {
			if (browser) {
				let errorMsg =
					response.status === 403 ? 'Akses Ditolak (Forbidden)' : 'Terjadi Kesalahan Server';

				try {
					// Karena backend selalu me-return plaintext saat error,
					// kita ekstrak menggunakan .text() bukan .json()
					const textData = await response.clone().text();

					// Jika textData tidak kosong, gunakan pesan dari backend
					if (textData && textData.trim() !== '') {
						errorMsg = textData.trim();
					}
				} catch (e) {
					// Abaikan jika gagal baca text
				}

				// Dinamis import goto untuk navigasi client-side SvelteKit
				const { goto } = await import('$app/navigation');
				goto(`/error-handler?code=${response.status}&message=${encodeURIComponent(errorMsg)}`);
			}

			// Return promise yang "menggantung" agar sisa kode di komponen Svelte
			// berhenti dieksekusi dan tidak memunculkan alert/error popup lokal.
			return new Promise(() => {});
		}

		return response;
	}

	async refreshToken() {
		try {
			const res = await fetch(`${config.apiUrl}/auth/refresh`, {
				method: 'POST',
				credentials: 'include'
			});

			if (res.ok) {
				return true;
			}
		} catch (e) {
			console.error('Refresh token failed', e);
		}

		window.location.reload();
		return false;
	}

	/** GET request */
	get(endpoint, headers = {}) {
		return this.request(endpoint, { method: 'GET', headers });
	}

	/** POST request */
	post(endpoint, body, headers = {}) {
		return this.request(endpoint, {
			method: 'POST',
			body: typeof body === 'object' ? JSON.stringify(body) : body,
			headers: { 'Content-Type': 'application/json', ...headers }
		});
	}

	/** PUT request */
	put(endpoint, body, headers = {}) {
		return this.request(endpoint, {
			method: 'PUT',
			body: typeof body === 'object' ? JSON.stringify(body) : body,
			headers: { 'Content-Type': 'application/json', ...headers }
		});
	}

	/** DELETE request */
	delete(endpoint, body = null, headers = {}) {
		return this.request(endpoint, {
			method: 'DELETE',
			body: body ? JSON.stringify(body) : null,
			headers: { 'Content-Type': 'application/json', ...headers }
		});
	}
}

// Create a single instance of ApiService
const apiService = new ApiService();
export default apiService;
