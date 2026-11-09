import { fetcher } from '#lib/fetcher.js';
import { API_URL } from '$app/env/private';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request, cookies }) => {
	const formData = await request.formData();
	const image = formData.get('image') as File;
	const contentID = formData.get('contentID');
	const contentType = formData.get('contentType');

	const uploadFormData = new FormData();
	uploadFormData.append('image', image);

	const uploadRes = await fetcher<{ id: string }>(
		`${API_URL}/images`,
		'POST',
		cookies,
		uploadFormData
	);

	await fetcher(
		`${API_URL}/${contentType}/${contentID}`,
		'PATCH',
		cookies,
		JSON.stringify({
			image_id: uploadRes.id
		})
	);

	return Response.json({ success: true });
};
