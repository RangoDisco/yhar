export function load({ cookies }) {
	return {
		isLoggedIn: cookies.get('token')
	};
}
