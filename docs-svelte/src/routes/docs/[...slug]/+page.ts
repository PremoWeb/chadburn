import type { PageLoad } from './$types';

export const load: PageLoad = async ({ data }) => {
	return {
		slug: data.slug,
		content: data.content,
		section: data.section
	};
};

// Enable SSR for documentation pages
export const ssr = true; 