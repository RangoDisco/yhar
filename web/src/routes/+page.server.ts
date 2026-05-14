import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
const API_URL = process.env.API_URL
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies }) => {
	if (cookies.get('access_token')) {
		// Get current user and redirect to their profile
		const response = await fetcher(`${API_URL}/users/me`, 'GET', cookies, null);
		const userID = response.user.id;

		redirect(302, `/users/${userID}/top`);
	}

	// Otherwise, redirect to login
	redirect(302, `/auth/login`);
};
