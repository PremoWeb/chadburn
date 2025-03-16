import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

// Use the recommended approach with query parameter instead of deprecated 'as' option
const markdownModules = import.meta.glob('/static/markdown/api/**/*.md', { query: '?raw', import: 'default', eager: true });

export const load: PageServerLoad = async ({ params, url }) => {
	try {
		// Ensure the slug is properly defined
		const slug = params.slug || '';
		
		// Handle special cases
		let targetSlug = slug;
		
		// If the URL path is /api, use the overview slug directly
		if (url.pathname === '/api') {
			targetSlug = 'overview';
		} else if (!targetSlug) {
			// Default to overview if no slug is provided
			targetSlug = 'overview';
		}
		
		// Debug: Log all available markdown files
		console.log('Available API markdown modules:', Object.keys(markdownModules));
		
		// Function to find a markdown file by slug
		function findMarkdownContent(slug: string): string | null {
			// Try different path patterns
			const patterns = [
				// Direct match
				`/static/markdown/api/${slug}.md`,
				// Nested index
				`/static/markdown/api/${slug}/index.md`,
				// Hyphenated
				`/static/markdown/api/${slug.replace(/\//g, '-')}.md`
			];
			
			// For the root path
			if (!slug || slug === '') {
				patterns.push('/static/markdown/api/overview.md');
			}
			
			// Try each pattern
			for (const pattern of patterns) {
				if (markdownModules[pattern]) {
					console.log('Found file at path:', pattern);
					return markdownModules[pattern] as string;
				}
			}
			
			// If slug contains slashes, try a nested approach
			if (slug.includes('/')) {
				const segments = slug.split('/');
				const lastSegment = segments.pop() || '';
				const parentPath = segments.join('/');
				
				const nestedPattern = `/static/markdown/api/${parentPath}/${lastSegment}.md`;
				if (markdownModules[nestedPattern]) {
					console.log('Found file at nested path:', nestedPattern);
					return markdownModules[nestedPattern] as string;
				}
			}
			
			return null;
		}
		
		// Find the markdown file
		const content = findMarkdownContent(targetSlug);
		
		// If no file found, throw a 404 error
		if (!content) {
			console.error(`API page not found: ${targetSlug}`);
			throw error(404, `API page not found: ${targetSlug}`);
		}
		
		// Return the slug and content
		return { slug: targetSlug, content, section: 'api' };
	} catch (e: any) {
		console.error('Error in API load function:', e);
		if (e.status === 404) {
			throw e; // Re-throw 404 errors
		}
		throw error(500, `Error loading API page content: ${e.message || 'Unknown error'}`);
	}
}; 