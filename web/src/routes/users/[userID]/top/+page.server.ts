import type { PageServerLoad } from './$types';
import { API_URL } from '$env/static/private';
import { fetcher } from '$lib/fetcher';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;

	const getStreamedPeriodData = async (period: 'month' | 'year' | 'overall') => ({
		artists: await fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/artists?period=${period}&limit=6`,
			'GET',
			null,
			cookies
		),
		albums: await fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/albums?&period=${period}&limit=6`,
			'GET',
			null,
			cookies
		),
		tracks: await fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/tracks?period=${period}&limit=6`,
			'GET',
			null,
			cookies
		)
	});

	const [wArtists, wAlbums, wTracks, history] = await Promise.all([
		fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/artists?period=week&limit=6`,
			'GET',
			null,
			cookies
		),
		fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/albums?&period=week&limit=6`,
			'GET',
			null,
			cookies
		),
		fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/tracks?period=week&limit=6`,
			'GET',
			null,
			cookies
		),
		fetcher(`${API_URL}/users/${userID}/scrobbles/history`, 'GET', null, cookies)
	]);

	return {
		overall: getStreamedPeriodData('overall'),
		year: getStreamedPeriodData('year'),
		month: getStreamedPeriodData('month'),
		week: {
			artists: wArtists,
			albums: wAlbums,
			tracks: wTracks
		},
		history
	};
};
