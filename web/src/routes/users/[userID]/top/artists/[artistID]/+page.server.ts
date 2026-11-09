import type { PageServerLoad } from './$types';
import { fetcher } from '#lib/fetcher.js';
import { API_URL } from '$app/env/private';
import type { Paginated } from '#lib/types/pagination.js';
import type { Album, Artist, Scrobble, Track, TimelineResult } from '#lib/types/content.js';
import { timeDay } from 'd3-time';

export const load: PageServerLoad = async ({ params, cookies, locals }) => {
	const { userID, artistID } = params;

	// TODO: change
	const artists = await fetcher<Paginated<Artist>>(
		`${API_URL}/users/${userID}/stats/artists/top?period=overall&artists=${artistID}&limit=1`,
		'GET',
		cookies,
		null
	);

	const albums = await fetcher<Paginated<Album>>(
		`${API_URL}/users/${userID}/stats/albums/top?&period=overall&artists=${artistID}&limit=9`,
		'GET',
		cookies,
		null
	);
	const tracks = await fetcher<Paginated<Track>>(
		`${API_URL}/users/${userID}/stats/tracks/top?period=overall&artists=${artistID}&limit=6`,
		'GET',
		cookies,
		null
	);

	const history = await fetcher<Paginated<Scrobble>>(
		`${API_URL}/users/${userID}/scrobbles/history?artists=${artistID}&period=overall&limit=8`,
		'GET',
		cookies,
		null
	);

	const rawTimeline = await fetcher<TimelineResult[]>(
		`${API_URL}/users/${userID}/stats/artists/${artistID}/timeline?period=year&interval=day`,
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
		artist: artists.results[0],
		albums,
		tracks,
		history,
		user: locals.user,
		timeline
	};
};
