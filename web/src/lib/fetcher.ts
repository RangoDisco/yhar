import { type Cookies, error, redirect } from '@sveltejs/kit';

export const fetcher = async (
	url: string,
	method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE',
	cookies: Cookies,
	body?: string | FormData | null
) => {
	// Build query and fetch
	const headers: Headers = new Headers();

	const token = cookies?.get('token');
	if (token) {
		headers.set('Authorization', `Bearer ${token}`);
	}

	const response = await fetch(url, { method, headers, body });

	// In case response contains no error, directly returns data
	if (response.ok) {
		const json = await response.json();
		return json.data;
	}

	/* Otherwise if a 401 was returned
	 * - If there's a refresh token in the user's local storage, use it to get a new token and retry
	 * - If not redirect to the login screen
	 */
	if (response.status === 401 && !url.endsWith('/auth/login')) {
		cookies.delete('token', {
			path: '/'
		});
		const refresh = cookies?.get('refresh_token');
		if (refresh) {
			// TODO: Get token and try again
		}

		redirect(302, `/auth/login`);
	}

	error(response.status, response.statusText);
};
