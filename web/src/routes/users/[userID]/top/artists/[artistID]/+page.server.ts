import type { PageServerLoad } from './$types';
import { fetcher } from '$lib/fetcher';
import { API_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID, artistID } = params;

	// TODO: change
	const artists = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/artists?period=overall&artist=${artistID}&limit=1`,
		'GET',
		null,
		cookies
	);

	const albums = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/albums?&period=overall&artist=${artistID}&limit=6`,
		'GET',
		null,
		cookies
	);
	const tracks = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/top/tracks?period=overall&artist=${artistID}&limit=6`,
		'GET',
		null,
		cookies
	);

	const history = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/history?artist=${artistID}`,
		'GET',
		null,
		cookies
	);

	return {
		artist: artists.result[0],
		albums,
		tracks,
		history
	};
};
