import type { PageServerLoad } from './$types';
import { API_URL } from '$app/env/private';
import { fetcher } from '#lib/fetcher.js';
import type { Paginated } from '#lib/types/pagination.js';
import type { Scrobble } from '#lib/types/content.js';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;
	const page = url.searchParams.get('page') ?? '1';
	const artist = url.searchParams.get('artists') ?? '';

	const history = await fetcher<Paginated<Scrobble>>(
		`${API_URL}/users/${userID}/scrobbles/history?artists=${artist}&page=${page}&period=overall&limit=20`,
		'GET',
		cookies,
		null
	);

	return {
		page,
		history
	};
};
