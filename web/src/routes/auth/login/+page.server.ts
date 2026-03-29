import { type Actions, fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { API_URL } from '$env/static/private';
import { fetcher } from '$lib/fetcher';

export const load: PageServerLoad = async ({ cookies }) => {
	const token = cookies?.get('token');

	if (token) {
		redirect(302, '/');
	}
};

export const actions = {
	default: async ({ cookies, request }) => {
		const data = await request.formData();
		const username = data.get('username');
		const password = data.get('password');

		if (!username || username.toString().length < 2) {
			return fail(400, { field: 'username', error: 'Username must be at least 2 characters long' });
		}

		if (!password || password.toString().length < 6) {
			return fail(400, { field: 'password', error: 'Password is required' });
		}
		try {
			const response = await fetcher(
				`${API_URL}/auth/login`,
				'POST',
				cookies,
				JSON.stringify({ username, password })
			);

			cookies.set('token', response.token, { path: '/', httpOnly: true, sameSite: 'strict' });
			redirect(302, '/');
		} catch (error) {
			return fail(401, { field: null, error: 'Invalid username or password' });
		}
	}
} satisfies Actions;
