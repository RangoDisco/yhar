import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
import { API_URL } from '$env/static/private';
import type { Album } from '$lib/types/content';
import type { Paginated } from '$lib/types/pagination';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;
	const page = url.searchParams.get('page') ?? "1";
	const period = url.searchParams.get('period') ?? 'week';
	const artist = url.searchParams.get('artist');

	let queryUrl = `${API_URL}/users/${userID}/scrobbles/top/albums?period=${period}&page=${page}&limit=10`;

	if (artist != null) {
		queryUrl = queryUrl + `&artist=${artist}`;
	}

	const albums: Paginated<Album> = await fetcher(queryUrl, 'GET', cookies, null);

	return {
		period,
		albums,
		page
	};
};
