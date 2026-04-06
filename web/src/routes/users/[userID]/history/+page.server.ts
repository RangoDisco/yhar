import type { PageServerLoad } from './$types';
import { API_URL } from '$env/static/private';
import { fetcher } from '$lib/fetcher';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;
	const page = url.searchParams.get('page') ?? 1;
	const artist = url.searchParams.get('artist') ?? '';

	const history = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/history?artist=${artist}&page=${page}&period=overall&limit=20`,
		'GET',
		cookies,
		null
	);

	return {
		history
	};
};
