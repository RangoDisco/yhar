import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
import { API_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;
	const page = url.searchParams.get('page') ?? 1;
	const period = url.searchParams.get('period') ?? 'week';

	const artists = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/artists?period=${period}&page=${page}&limit=10`,
		'GET',
		null,
		cookies
	);

	return {
		period,
		artists
	};
};
