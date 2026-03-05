/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
	let accessToken = event.cookies.get('access_token');
	const refreshToken = event.cookies.get('refresh_token');

	// Ambil identitas asli untuk Zero Trust di Go
	const contextHeaders = {
		'User-Agent': event.request.headers.get('user-agent') || '',
		'X-Forwarded-For': event.getClientAddress(),
		Accept: 'application/json'
	};

	// Fungsi helper untuk hit API profile
	const fetchProfile = async () => {
		return await event.fetch(`${import.meta.env.VITE_API_URL}/users/me`, {
			headers: contextHeaders
		});
	};

	// Fungsi untuk refresh token dan SINKRONISASI COOKIE
	const tryRefreshAndSync = async () => {
		try {
			const refreshRes = await event.fetch(`${import.meta.env.VITE_API_URL}/auth/refresh`, {
				method: 'POST',
				headers: contextHeaders
			});

			if (refreshRes.ok) {
				// Ambil semua Set-Cookie dari respons Go (Access & Refresh Token baru)
				const setCookieHeaders = refreshRes.headers.getSetCookie();

				setCookieHeaders.forEach((cookieString) => {
					// Split untuk ambil "nama=value"
					const [nameValue] = cookieString.split(';');
					const [name, value] = nameValue.split('=');

					// Set manual ke event.cookies SvelteKit
					// Atribut lain (HttpOnly, Path, dll) sebaiknya disesuaikan
					event.cookies.set(name.trim(), value.trim(), {
						path: '/',
						httpOnly: true,
						sameSite: 'lax',
						secure: process.env.NODE_ENV === 'production'
					});
				});
				return true;
			}
		} catch (err) {
			console.log('Refresh attempt failed in hooks:', err);
		}
		return false;
	};

	// LOGIKA UTAMA
	// 1. Kasus Access Token expired (401)
	if (accessToken) {
		try {
			let response = await fetchProfile();

			if (response.status === 401 && refreshToken) {
				const refreshed = await tryRefreshAndSync();
				if (refreshed) {
					response = await fetchProfile();
				}
			}

			if (response.ok) {
				event.locals.user = await response.json();
				return resolve(event);
			}
		} catch (err) {
			console.log('Error fetching profile:', err);
		}
	}

	// 2. Kasus Access Token hilang/kosong tapi Refresh Token masih ada
	if (!accessToken && refreshToken) {
		const refreshed = await tryRefreshAndSync();
		if (refreshed) {
			try {
				const response = await fetchProfile();
				if (response.ok) {
					event.locals.user = await response.json();
				}
			} catch (err) {
				console.log('Error after silent refresh:', err);
			}
		}
	}

	// Pastikan locals.user tidak undefined
	event.locals.user = event.locals.user || null;

	// Jika gagal semua, hapus cookies agar tidak loop
	if (!event.locals.user && (accessToken || refreshToken)) {
		// Optional: Hapus jika ingin benar-benar clean saat auth gagal
		event.cookies.delete('access_token', { path: '/' });
		event.cookies.delete('refresh_token', { path: '/' });
	}

	return resolve(event);
}
