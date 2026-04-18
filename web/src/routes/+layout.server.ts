export function load({ cookies, locals }) {
	const user = locals.user;
	return {
		user
	};
}
