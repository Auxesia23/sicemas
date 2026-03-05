// API Service - Centralized HTTP client for all API calls
import config from '$lib/config';

class ApiService {
	/**
	 * Makes an HTTP request to the API
	 * @param {string} endpoint - The API endpoint (e.g., '/auth/login')
	 * @param {Object} options - Request options (method, body, headers, etc.)
	 * @returns {Promise} The raw response from the API
	 */
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

		// Jika Unauthorized (401), coba refresh token
		if (response.status === 401 && endpoint !== '/auth/refresh') {
			const success = await this.refreshToken();
			if (success) {
				// Ulangi request yang tadi gagal dengan token baru
				response = await fetch(url, requestOptions);
			}
		}

		return response;
	}

	async refreshToken() {
		try {
			// Panggil endpoint refresh (backend harus kirim access_token baru)
			// Cookie refresh_token akan otomatis terkirim karena credentials: include
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

		// Jika gagal refresh, paksa logout
		window.location.href = '/login';
		return false;
	}
	/**
	 * GET request
	 */
	get(endpoint, headers = {}) {
		return this.request(endpoint, {
			method: 'GET',
			headers
		});
	}

	/**
	 * POST request
	 */
	post(endpoint, body, headers = {}) {
		return this.request(endpoint, {
			method: 'POST',
			body: typeof body === 'object' ? JSON.stringify(body) : body,
			headers: {
				'Content-Type': 'application/json',
				...headers
			}
		});
	}

	/**
	 * PUT request
	 */
	put(endpoint, body, headers = {}) {
		return this.request(endpoint, {
			method: 'PUT',
			body: typeof body === 'object' ? JSON.stringify(body) : body,
			headers: {
				'Content-Type': 'application/json',
				...headers
			}
		});
	}

	/**
	 * DELETE request
	 */
	delete(endpoint, body = null, headers = {}) {
		return this.request(endpoint, {
			method: 'DELETE',
			body: body ? JSON.stringify(body) : null,
			headers: {
				'Content-Type': 'application/json',
				...headers
			}
		});
	}
}

// Create a single instance of ApiService
const apiService = new ApiService();

export default apiService;
