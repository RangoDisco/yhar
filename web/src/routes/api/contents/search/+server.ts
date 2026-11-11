import { fetcher } from '#lib/fetcher.js';
import { API_URL } from '$app/env/private';
import type { RequestHandler } from './$types';
import type { SearchResult } from '#lib/types/content.js';
import type { Paginated } from '#lib/types/pagination.ts';

export const GET: RequestHandler = async ({ url, cookies }) => {
	const value = url.searchParams.get('value');
	const page = url.searchParams.get('page');
	const res = await fetcher<Paginated<SearchResult>>(
		`${API_URL}/contents/search?value=${value}&page=${page}&limit=10`,
		'GET',
		cookies
	);

	return Response.json({ ...res });
};
