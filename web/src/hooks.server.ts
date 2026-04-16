import { error, type Handle } from '@sveltejs/kit';

export const handle = (async ({ event, resolve }) => {
	let user = null;
	const rawUser = event.cookies?.get('user');
	if (rawUser !== undefined) {
		user = JSON.parse(rawUser);
	}

	if (event?.route?.id?.includes('/(protected)/') && !user) {
		throw error(401, 'requires authentication');
	}
	if (event?.route?.id?.includes('/(proctectedAdmin)/') && user?.role !== 'ADMIN') {
		throw error(401, 'requires admin');
	}

	return resolve(event);
}) satisfies Handle;
