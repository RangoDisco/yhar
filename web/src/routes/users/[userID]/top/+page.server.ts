import type { PageServerLoad } from './$types';
import { API_URL } from '$env/static/private';
import { fetcher } from '$lib/fetcher';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;

	const getStreamedPeriodData = async (period: 'week' | 'month' | 'year' | 'overall') => ({
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

	const overall = getStreamedPeriodData('overall');
	const year = getStreamedPeriodData('year');
	const month = getStreamedPeriodData('month');
	const week = await getStreamedPeriodData('week');

	const history = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/history`,
		'GET',
		null,
		cookies
	);

	return {
		overall: overall,
		year: year,
		month: month,
		week,
		history
	};
};
