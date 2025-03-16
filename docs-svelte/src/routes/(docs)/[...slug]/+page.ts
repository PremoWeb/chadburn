// This file is needed to ensure proper typing for the data loaded by the server
// No actual loading happens here since we're using +page.server.ts

import type { PageLoad } from './$types';

// Disable prerendering for this route
export const prerender = false;

// Export the load function with the correct type
export const load: PageLoad = async ({ data }) => {
	// Pass through the data from the server load function
	return {
		slug: data.slug,
		content: data.content
	};
}; 