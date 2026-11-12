import type { PageServerLoad } from './$types';
import { fetcher } from '#lib/fetcher.js';
import { API_URL } from '$app/env/private';
import type { Paginated } from '#lib/types/pagination.js';
import type { Scrobble, Track, TimelineResult } from '#lib/types/content.js';
import { timeDay } from 'd3-time';

export const load: PageServerLoad = async ({ params, cookies, locals }) => {
	const { userID, trackID } = params;

	const tracks = await fetcher<Paginated<Track>>(
		`${API_URL}/users/${userID}/stats/tracks/top?period=overall&tracks=${trackID}&limit=1`,
		'GET',
		cookies,
		null
	);

	const history = await fetcher<Paginated<Scrobble>>(
		`${API_URL}/users/${userID}/scrobbles/history?tracks=${trackID}&period=overall&limit=8`,
		'GET',
		cookies,
		null
	);

	const rawTimeline = await fetcher<TimelineResult[]>(
		`${API_URL}/users/${userID}/stats/tracks/${trackID}/timeline?period=year&interval=day`,
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

	return {
		track: tracks.results[0],
		history,
		user: locals.user,
		timeline
	};
};
