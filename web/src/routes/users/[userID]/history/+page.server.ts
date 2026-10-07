import type { PageServerLoad } from './$types';
import { API_URL } from '$app/env/private';
import { fetcher } from '#lib/fetcher.js';
import type { Paginated } from '#lib/types/pagination.js';
import type { Scrobble } from '#lib/types/content.js';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;
	const page = url.searchParams.get('page') ?? '1';
	const artist = url.searchParams.get('artist') ?? '';

	const history: Paginated<Scrobble> = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/history?artist=${artist}&page=${page}&period=overall&limit=20`,
		'GET',
		cookies,
		null
	);

	return {
		page,
		history
	};
};
