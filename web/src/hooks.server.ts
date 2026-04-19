import { error, type Handle } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';
import type { User } from '@lucide/svelte';

export const handle = (async ({ event, resolve }) => {
	let user = null;
	const token = event.cookies?.get('access_token');

	if (token !== undefined) {
		user = jwtDecode(token);
		if (!user || !user.exp) {
			throw error(401, 'Invalid user');
		}
		event.locals.user = user as {
			id: string;
			username: string;
			role: 'USER' | 'ADMIN';
			expiresAt: number;
		};
	}

	if (event?.route?.id?.includes('/(protected)/') && !user) {
		throw error(401, 'requires authentication');
	}
	if (event?.route?.id?.includes('/(proctectedAdmin)/') && user?.role !== 'ADMIN') {
		throw error(401, 'requires admin');
	}

	return resolve(event);
}) satisfies Handle;
