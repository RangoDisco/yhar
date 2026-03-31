import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
import { API_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID, albumID } = params;
	const albums = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/albums?&period=overall&album=${albumID}&limit=1`,
		'GET',
		cookies,
		null
	);

	const tracks = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/tracks?period=overall&album=${albumID}&limit=20`,
		'GET',
		cookies,
		null
	);

	return {
		album: albums.results[0] ?? null,
		tracks
	};
};
