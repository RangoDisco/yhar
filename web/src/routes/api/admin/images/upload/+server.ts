import { fetcher } from '$lib/fetcher';
import { API_URL } from '$app/env/private';
import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';

export const POST: RequestHandler = async ({ request, cookies }) => {
	const formData = await request.formData();
	const image = formData.get('image') as File;
	const contentID = formData.get('contentID');
	const contentType = formData.get('contentType');

	const uploadFormData = new FormData();
	uploadFormData.append('image', image);

	const uploadRes = await fetcher(`${API_URL}/images`, 'POST', cookies, uploadFormData);

	const updateRes = await fetcher(
		`${API_URL}/${contentType}/${contentID}`,
		'PATCH',
		cookies,
		JSON.stringify({
			image_id: uploadRes.id
		})
	);

	return json({ success: true });
};
