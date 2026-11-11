import type { PageServerLoad } from './$types';
import { API_URL } from '$app/env/private';
import { fetcher } from '#lib/fetcher.js';
import type { Paginated } from '#lib/types/pagination.js';
import type {
	Album,
	Artist,
	AvgDailyScrobbles,
	Scrobble,
	TimelineResult,
	TotalDistinctCount,
	TotalScrobbleDurationAndCount,
	Track
} from '#lib/types/content.js';
import { Period } from '#lib/types/period.js';
import { timeDay } from 'd3-time';

export const load: PageServerLoad = async ({ params, cookies }) => {
	const { userID } = params;

	// PERIOD DATA //
	const getStreamedPeriodData = async (period: Period) => ({
		artists: await fetcher<Paginated<Artist>>(
			`${API_URL}/users/${userID}/stats/artists/top?period=${period}&limit=8`,
			'GET',
			cookies,
			null
		),
		albums: await fetcher<Paginated<Album>>(
			`${API_URL}/users/${userID}/stats/albums/top?&period=${period}&limit=8`,
			'GET',
			cookies,
			null
		),
		tracks: await fetcher<Paginated<Track>>(
			`${API_URL}/users/${userID}/stats/tracks/top?period=${period}&limit=8`,
			'GET',
			cookies,
			null
		)
	});

	const overall = getStreamedPeriodData(Period.overall).catch(() => null);
	const year = getStreamedPeriodData(Period.year).catch(() => null);
	const month = getStreamedPeriodData(Period.month).catch(() => null);
	const week = await getStreamedPeriodData(Period.week);

	const totalScrobblesStats = await fetcher<TotalScrobbleDurationAndCount>(
		`${API_URL}/users/${userID}/stats/scrobbles/total`,
		'GET',
		cookies,
		null
	);

	const history = await fetcher<Paginated<Scrobble>>(
		`${API_URL}/users/${userID}/scrobbles/history?period=overall&page=1`,
		'GET',
		cookies,
		null
	);

	// HEATMAP //
	const rawTimeline = await fetcher<TimelineResult[]>(
		`${API_URL}/users/${userID}/stats/scrobbles/timeline?period=year&interval=day`,
		'GET',
		cookies,
		null
	);

	const timeline = rawTimeline.map((item) => {
		return {
			date: timeDay.floor(new Date(item.periodInterval)),
			value: item.listenedCount > 0 ? item.listenedCount : null
		};
	});

	const yearlyScrobbleCount = timeline.reduce((n, { value }) => n + (value || 0), 0);

	const avgDailyScrobbles = await fetcher<AvgDailyScrobbles>(
		`${API_URL}/users/${userID}/stats/scrobbles/daily-avg?period=year`,
		'GET',
		cookies,
		null
	);

	const totalDistinctTracksCount = await fetcher<TotalDistinctCount>(
		`${API_URL}/users/${userID}/stats/tracks/total`,
		'GET',
		cookies,
		null
	);

	return {
		overall: overall,
		year: year,
		month: month,
		week,
		history,
		timeline,
		yearlyScrobbleCount,
		totalScrobblesStats,
		avgDailyScrobbles,
		totalDistinctTracksCount
	};
};
