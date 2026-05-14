import { fetcher } from '$lib/fetcher';
const API_URL = process.env.API_URL
import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';

export const DELETE: RequestHandler = async ({ params, cookies }) => {
	const { id } = params;

	const res = await fetcher(`${API_URL}/scrobbles/${id}`, 'DELETE', cookies);

	return json({ success: true });
};
