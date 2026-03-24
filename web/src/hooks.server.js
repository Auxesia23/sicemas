import { env } from '$env/dynamic/public';

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
	if (event.url.pathname.startsWith('/api')) {
		return await resolve(event);
	}

	let accessToken = event.cookies.get('access_token');
	const refreshToken = event.cookies.get('refresh_token');

	const deviceId = event.cookies.get('device_id') || '';
	const clientIP = event.request.headers.get('x-forwarded-for') || event.getClientAddress();

	let contextHeaders = {
		'User-Agent': event.request.headers.get('user-agent') || '',
		'X-Forwarded-For': clientIP,
		'X-True-Client-IP': clientIP,
		Accept: 'application/json',
		Cookie: event.request.headers.get('cookie') || '',
		'X-Device-Id': deviceId
	};

	const fetchProfile = async () => {
		return await fetch(`${env.PUBLIC_API_URL}/users/me`, {
			headers: contextHeaders
		});
	};

	const tryRefreshAndSync = async () => {
		try {
			const refreshRes = await fetch(`${env.PUBLIC_API_URL}/auth/refresh`, {
				method: 'POST',
				headers: contextHeaders
			});

			if (refreshRes.ok) {
				const setCookieHeaders = refreshRes.headers.getSetCookie();

				let newTokens = {};

				setCookieHeaders.forEach((cookieString) => {
					const [nameValue] = cookieString.split(';');

					const splitIndex = nameValue.indexOf('=');
					if (splitIndex === -1) return;

					const cookieName = nameValue.substring(0, splitIndex).trim();
					const cookieVal = nameValue.substring(splitIndex + 1).trim();

					event.cookies.set(cookieName, cookieVal, {
						path: '/',
						httpOnly: true,
						sameSite: 'lax',
						secure: process.env.NODE_ENV === 'production',
						maxAge: cookieName === 'refresh_token' ? 60 * 60 * 24 * 7 : 60 * 15
					});

					newTokens[cookieName] = cookieVal;
				});

				if (newTokens.access_token) {
					contextHeaders['Cookie'] =
						`access_token=${newTokens.access_token}; refresh_token=${newTokens.refresh_token || refreshToken}`;
				}

				return true;
			}
		} catch (err) {
			console.log('Refresh attempt failed in hooks:', err);
		}
		return false;
	};

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
				const body = await response.json();

				event.locals.user = body.data || body;

				return resolve(event);
			}
		} catch (err) {
			console.log('Error fetching profile:', err);
		}
	}

	if (!accessToken && refreshToken) {
		const refreshed = await tryRefreshAndSync();
		if (refreshed) {
			try {
				const response = await fetchProfile();
				if (response.ok) {
					const body = await response.json();
					event.locals.user = body;
				}
			} catch (err) {
				console.log('Error after silent refresh:', err);
			}
		}
	}

	event.locals.user = event.locals.user || null;

	if (!event.locals.user && (accessToken || refreshToken)) {
		event.cookies.delete('access_token', { path: '/' });
		event.cookies.delete('refresh_token', { path: '/' });
	}

	return resolve(event);
}
