import adapter from '@sveltejs/adapter-cloudflare';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: [vitePreprocess()],

	kit: {
		// Using adapter-cloudflare for Cloudflare Pages deployment
		adapter: adapter({
			// Cloudflare adapter options
			routes: {
				include: ['/*'],
				exclude: ['<all>']
			}
		}),
		paths: {
			base: '', // Use root path for all environments
			relative: false // Use root-relative URLs for assets
		}
	}
};

export default config;
