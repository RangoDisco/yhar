import {type Cookies, error, redirect} from '@sveltejs/kit';
import {jwtDecode} from 'jwt-decode';

import {API_URL} from '$env/static/private';

export const fetcher = async (
    url: string,
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE',
    cookies: Cookies,
    body?: string | FormData | null
) => {
    // Build query and fetch
    const headers: Headers = new Headers();

    const token = cookies?.get('access_token');
    if (token) {
        headers.set('Authorization', `Bearer ${token}`);
    }

    const response = await fetch(url, {method, headers, body});

    // In case response contains no error, directly returns data
    if (response.ok) {
        const json = await response.json();
        return json.data;
    }

    if (response.status !== 401 || url.endsWith('/auth/login')) {
        error(response.status, response.statusText);
    }

    /* Otherwise if a 401 was returned
     * - If there's a refresh token in the user's local storage, use it to get a new token and retry
     * - If not redirect to the login screen
     */

    cookies.delete('access_token', {
        path: '/'
    });

    const refresh = cookies?.get('refresh_token');
    if (!refresh) {
        redirect(302, `/auth/login`);
    }

    const refreshResponse = await fetch(`${API_URL}/auth/refresh`, {
        method: 'POST',
        headers: new Headers({
            'Content-Type': 'application/json'
        }),
        body: JSON.stringify({refresh_token: refresh})
    });

    if (!refreshResponse.ok) {
        redirect(302, `/auth/login`);
    }

    const json = await refreshResponse.json();
    const decodedToken = jwtDecode(json.data.access_token);
    if (!decodedToken || !decodedToken.exp) {
        redirect(302, `/auth/login`);
    }

    cookies.set('access_token', json.data.access_token, {
        path: '/',
        httpOnly: true,
        sameSite: 'strict'
    });

    // Finally retry query
    headers.set('Authorization', `Bearer ${json.data.access_token}`);
    const retryResponse = await fetch(url, {method, headers, body});

    if (!retryResponse.ok) {
        cookies.delete('refresh_token', {
            path: '/'
        });
        error(retryResponse.status, retryResponse.statusText);
    }

    const retryJson = await retryResponse.json();
    return retryJson.data;
};
