import type { PageServerLoad } from './$types';
import { API_URL } from '$app/env/private';
import { fetcher } from '$lib/fetcher';
import type { Paginated } from '$lib/types/pagination';
import type { Album, Artist, Scrobble, Track } from '$lib/types/content';
import { Period } from '$lib/types/period';

export const load: PageServerLoad = async ({ url, params, cookies }) => {
	const { userID } = params;

	const getStreamedPeriodData = async (
		period: Period
	): Promise<{
		artists: Paginated<Artist>;
		albums: Paginated<Album>;
		tracks: Paginated<Track>;
	}> => ({
		artists: await fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/artists?period=${period}&limit=9`,
			'GET',
			cookies,
			null
		),
		albums: await fetcher(
			`${API_URL}/users/${userID}/scrobbles/top/albums?&period=${period}&limit=9`,
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

	const overall = getStreamedPeriodData(Period.overall).catch(() => null);
	const year = getStreamedPeriodData(Period.year).catch(() => null);
	const month = getStreamedPeriodData(Period.month).catch(() => null);
	const week = await getStreamedPeriodData(Period.week);

	const history: Paginated<Scrobble> = await fetcher(
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
