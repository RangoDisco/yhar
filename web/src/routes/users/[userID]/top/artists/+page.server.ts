import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
import { API_URL } from '$app/env/private';
import type { Paginated } from '$lib/types/pagination';
import type { Artist } from '$lib/types/content';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;
	const page = url.searchParams.get('page') ?? '1';
	const period = url.searchParams.get('period') ?? 'week';

	const artists: Paginated<Artist> = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/artists?period=${period}&page=${page}&limit=10`,
		'GET',
		cookies,
		null
	);

	return {
		period,
		artists,
		page
	};
};
