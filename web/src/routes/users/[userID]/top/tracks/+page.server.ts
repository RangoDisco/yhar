import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
import { API_URL } from '$app/env/private';
import type { Paginated } from '$lib/types/pagination';
import type { Track } from '$lib/types/content';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;
	const page = url.searchParams.get('page') ?? '1';
	const period = url.searchParams.get('period') ?? 'week';
	const artist = url.searchParams.get('artist') ?? '';

	const tracks: Paginated<Track> = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/tracks?period=${period}&page=${page}&artist=${artist}&limit=10`,
		'GET',
		cookies,
		null
	);

	return {
		period,
		tracks,
		page
	};
};
