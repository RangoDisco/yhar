import type { PageServerLoad } from './$types';
import { API_URL } from '$env/static/private';
import { fetcher } from '$lib/fetcher';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;

	const getStreamedPeriodData = async (period: 'week' | 'month' | 'year' | 'overall') => ({
		artists: await fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/artists?period=${period}&limit=6`,
			'GET',
			cookies,
			null
		),
		albums: await fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/albums?&period=${period}&limit=6`,
			'GET',
			cookies,
			null
		),
		tracks: await fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/tracks?period=${period}&limit=6`,
			'GET',
			cookies,
			null
		)
	});

	const overall = getStreamedPeriodData('overall').catch(() => null);
	const year = getStreamedPeriodData('year').catch(() => null);
	const month = getStreamedPeriodData('month').catch(() => null);
	const week = await getStreamedPeriodData('week');

	const history = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/history`,
		'GET',
		cookies,
		null
	);

	return {
		overall: overall,
		year: year,
		month: month,
		week,
		history
	};
};
