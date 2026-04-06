import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
import { API_URL } from '$env/static/private';
import type { Paginated } from '$lib/types/pagination';
import type { Album, Track } from '$lib/types/content';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID, artistID } = params;

	// TODO: change
	const artists = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/artists?period=overall&artist=${artistID}&limit=1`,
		'GET',
		cookies,
		null
	);

	const albums: Paginated<Album> = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/albums?&period=overall&artist=${artistID}&limit=6`,
		'GET',
		cookies,
		null
	);
	const tracks: Paginated<Track> = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/tracks?period=overall&artist=${artistID}&limit=6`,
		'GET',
		cookies,
		null
	);

	const history: Paginated<Track> = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/history?artist=${artistID}&period=overall`,
		'GET',
		cookies,
		null
	);

	return {
		artist: artists.results[0],
		albums,
		tracks,
		history
	};
};
