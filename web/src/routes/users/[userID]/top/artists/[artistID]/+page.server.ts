import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
import { API_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID, artistID } = params;

	// TODO: change
	const artists = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/artists?period=overall&artist=${artistID}&limit=1`,
		'GET',
		cookies,
		null
	);

	const albums = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/albums?&period=overall&artist=${artistID}&limit=6`,
		'GET',
		cookies,
		null
	);
	const tracks = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/tracks?period=overall&artist=${artistID}&limit=6`,
		'GET',
		cookies,
		null
	);

	const history = await fetcher(
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
