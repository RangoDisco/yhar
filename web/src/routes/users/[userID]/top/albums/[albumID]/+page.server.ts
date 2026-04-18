import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
import { API_URL } from '$env/static/private';
import type { Paginated } from '$lib/types/pagination';
import type { Album, Track } from '$lib/types/content';

export const load: PageServerLoad = async ({ url, params, cookies, locals }) => {
	const { userID, albumID } = params;
	const albums: Paginated<Album> = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/albums?&period=overall&album=${albumID}&limit=1`,
		'GET',
		cookies,
		null
	);

	const tracks: Paginated<Track> = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/tracks?period=overall&album=${albumID}&limit=20`,
		'GET',
		cookies,
		null
	);

	return {
		album: albums.results[0] ?? null,
		tracks,
		user: locals.user,
	};
};
