import type { PageServerLoad } from './$types';
import { fetcher } from '#lib/fetcher.js';
import { API_URL } from '$app/env/private';
import type { Paginated } from '#lib/types/pagination.js';
import type { Album, Artist, Scrobble, Track, TimelineResult } from '#lib/types/content.js';
import { timeDay } from 'd3-time';

export const load: PageServerLoad = async ({ url, params, cookies, locals }) => {
	const { userID, artistID } = params;

	// TODO: change
	const artists: Paginated<Artist> = await fetcher(
		`${API_URL}/users/${userID}/stats/artists/top?period=overall&artist=${artistID}&limit=1`,
		'GET',
		cookies,
		null
	);

	const albums: Paginated<Album> = await fetcher(
		`${API_URL}/users/${userID}/stats/albums/top?&period=overall&artist=${artistID}&limit=9`,
		'GET',
		cookies,
		null
	);
	const tracks: Paginated<Track> = await fetcher(
		`${API_URL}/users/${userID}/stats/tracks/top?period=overall&artist=${artistID}&limit=6`,
		'GET',
		cookies,
		null
	);

	const history: Paginated<Scrobble> = await fetcher(
		`${API_URL}/users/${userID}/scrobbles/history?artist=${artistID}&period=overall&limit=8`,
		'GET',
		cookies,
		null
	);

	const rawTimeline: TimelineResult[] = await fetcher(
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
