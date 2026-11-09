import { fetcher } from '#lib/fetcher.js';
import { API_URL } from '$app/env/private';
import type { RequestHandler } from './$types';

export const DELETE: RequestHandler = async ({ params, cookies }) => {
	const { id } = params;

	const res = await fetcher(`${API_URL}/scrobbles/${id}`, 'DELETE', cookies);

	return Response.json({ success: true });
};
