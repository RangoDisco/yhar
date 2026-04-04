import { redirect } from '@sveltejs/kit';

export const actions = {
	default: ({ cookies }) => {
		cookies.delete('token', { path: '/' });
		redirect(303, '/');
	}
};
