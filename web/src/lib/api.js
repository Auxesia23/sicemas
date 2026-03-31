// API Service - Centralized HTTP client for all API calls
import config from '$lib/config';
import { browser } from '$app/environment';

// Import FingerprintJS for device identification
let fpPromise = null;
if (browser) {
	fpPromise = import('@fingerprintjs/fingerprintjs').then((fp) => fp.default.load());
}

class ApiService {
	constructor() {
		// State untuk mengunci dan mengantrekan request
		this.isRefreshing = false;
		this.failedQueue = [];
		// Cache for Device ID
		this.deviceId = null;
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

	// Get or generate Device ID using FingerprintJS
	async getDeviceId() {
		if (this.deviceId) return this.deviceId;

		if (browser && fpPromise) {
			try {
				const fp = await fpPromise;
				const result = await fp.get();
				this.deviceId = result.visitorId;

				document.cookie = `device_id=${this.deviceId}; path=/; max-age=31536000; SameSite=Lax`;
				// -------------------------------------------------------------

				return this.deviceId;
			} catch (error) {
				console.error('Failed to get fingerprint:', error);
				this.deviceId = `fallback-${Date.now()}-${Math.random().toString(36).substring(7)}`;
				document.cookie = `device_id=${this.deviceId}; path=/; max-age=31536000; SameSite=Lax`;
				return this.deviceId;
			}
		}

		this.deviceId = 'server-side';
		return this.deviceId;
	}

	async request(endpoint, options = {}) {
		const url = `${config.apiUrl}${endpoint}`;

		// Check if body is FormData to handle Content-Type properly
		const isFormData = options.body instanceof FormData;

		// Get Device ID and inject into headers
		const deviceId = await this.getDeviceId();

		// Build headers
		const requestHeaders = {
			// Don't include default Content-Type for FormData
			...(isFormData ? {} : config.defaultHeaders),
			...options.headers,
			// Always include Device ID
			'X-Device-Id': deviceId
		};

		// Determine timeout - allow override for specific requests (e.g., file uploads)
		const timeout = options.timeout || config.timeout;

		// Create AbortController for timeout handling
		const controller = new AbortController();
		const timeoutId = setTimeout(() => {
			controller.abort();
		}, timeout);

		const requestOptions = {
			...config.defaultFetchOptions,
			...options,
			headers: requestHeaders,
			signal: controller.signal
		};

		try {
			let response = await fetch(url, requestOptions);

			// Clear timeout after fetch completes
			clearTimeout(timeoutId);

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
					window.location.href = '/login';
				}
			}

			if (response.status === 429) {
				// Coba ambil pesan dari Go lu, atau pake pesan default
				let errorMsg = 'Terlalu banyak permintaan. Silakan tunggu beberapa saat.';
				try {
					const textData = await response.clone().text();
					if (textData && textData.trim() !== '') {
						errorMsg = textData.trim();
					}
				} catch (e) {
					// Abaikan
				}

				if (browser) {
					const { goto } = await import('$app/navigation');
					goto(`/error-handler?code=429&message=${encodeURIComponent(errorMsg)}`);
					return new Promise(() => {}); // Gantung promise agar eksekusi komponen berhenti
				} else {
					throw new Error(`API Error: 429 - ${errorMsg}`);
				}
			}

			if (response.status === 403 || response.status >= 500) {
				let errorMsg =
					response.status === 403 ? 'Akses Ditolak (Forbidden)' : 'Terjadi Kesalahan Server';

				try {
					const textData = await response.clone().text();
					if (textData && textData.trim() !== '') {
						errorMsg = textData.trim();
					}
				} catch (e) {
					// Abaikan jika gagal baca text
				}

				if (browser) {
					const { goto } = await import('$app/navigation');
					goto(`/error-handler?code=${response.status}&message=${encodeURIComponent(errorMsg)}`);
					return new Promise(() => {});
				} else {
					throw new Error(`API Error: ${response.status} - ${errorMsg}`);
				}
			}

			return response;
		} catch (error) {
			// Clear timeout on error
			clearTimeout(timeoutId);

			// Handle timeout error specifically
			if (error.name === 'AbortError') {
				throw new Error(`Request timeout after ${timeout}ms`);
			}

			// Re-throw other errors
			throw error;
		}
	}

	async refreshToken() {
		try {
			const res = await fetch(`${config.apiUrl}/auth/refresh`, {
				method: 'POST',
				credentials: 'include',
				headers: {
					'X-Device-Id': this.deviceId,
					Accept: 'application/json'
				}
			});

			if (res.ok) {
				return true;
			}
		} catch (e) {
			console.error('Refresh token failed', e);
		}
		return false;
	}

	/** GET request */
	get(endpoint, headers = {}) {
		return this.request(endpoint, { method: 'GET', headers });
	}

	/** POST request */
	post(endpoint, body, headers = {}) {
		const isFormData = body instanceof FormData;

		// Buat salinan headers
		const requestHeaders = { ...headers };

		// Jika BUKAN FormData, baru kita pasang default Content-Type JSON
		if (!isFormData && !requestHeaders['Content-Type']) {
			requestHeaders['Content-Type'] = 'application/json';
		}

		return this.request(endpoint, {
			method: 'POST',
			// Jika FormData, biarkan apa adanya. Jika object biasa, ubah ke JSON string.
			body: isFormData ? body : typeof body === 'object' ? JSON.stringify(body) : body,
			headers: requestHeaders
		});
	}

	/** PUT request */
	put(endpoint, body, headers = {}) {
		const isFormData = body instanceof FormData;
		const requestHeaders = { ...headers };

		if (!isFormData && !requestHeaders['Content-Type']) {
			requestHeaders['Content-Type'] = 'application/json';
		}

		return this.request(endpoint, {
			method: 'PUT',
			body: isFormData ? body : typeof body === 'object' ? JSON.stringify(body) : body,
			headers: requestHeaders
		});
	}

	/** PATCH request */
	patch(endpoint, body, headers = {}) {
		const isFormData = body instanceof FormData;
		const requestHeaders = { ...headers };

		if (!isFormData && !requestHeaders['Content-Type']) {
			requestHeaders['Content-Type'] = 'application/json';
		}

		return this.request(endpoint, {
			method: 'PATCH',
			body: isFormData ? body : typeof body === 'object' ? JSON.stringify(body) : body,
			headers: requestHeaders
		});
	}

	/** DELETE request */
	delete(endpoint, body = null, headers = {}) {
		const isFormData = body instanceof FormData;
		const requestHeaders = { ...headers };

		if (!isFormData && !requestHeaders['Content-Type']) {
			requestHeaders['Content-Type'] = 'application/json';
		}

		return this.request(endpoint, {
			method: 'DELETE',
			body: isFormData ? body : typeof body === 'object' ? JSON.stringify(body) : body,
			headers: requestHeaders
		});
	}
}

// Create a single instance of ApiService
const apiService = new ApiService();
export default apiService;
