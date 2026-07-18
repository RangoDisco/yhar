import { defineEnvVars } from '@sveltejs/kit/hooks';
import * as v from 'valibot';

export const variables = defineEnvVars({
	API_URL: {
		schema: v.optional(v.string())
	}
});
