import { error } from '@sveltejs/kit';
import fs from 'fs';
import path from 'path';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, url }) => {
	try {
		// Ensure the slug is properly defined
		const slug = params.slug || '';
		
		// Handle special cases
		let targetSlug = slug;
		
		// If the URL path is /installation, use the installation slug directly
		if (url.pathname === '/installation') {
			targetSlug = 'installation';
		} else if (!targetSlug) {
			// Default to introduction if no slug is provided
			targetSlug = 'introduction';
		}
		
		// Try multiple file paths to find the correct markdown file
		let filePath = '';
		let content = '';
		let fileFound = false;
		
		// Possible file paths to try in order
		const possiblePaths = [
			// Direct match (e.g., metrics/quick-start.md for /metrics/quick-start)
			path.join(process.cwd(), 'static', 'markdown', `${targetSlug}.md`),
			
			// Nested path with index.md (e.g., metrics/index.md for /metrics)
			path.join(process.cwd(), 'static', 'markdown', targetSlug, 'index.md'),
			
			// Hyphenated path (e.g., metrics-quick-start.md for /metrics/quick-start)
			path.join(process.cwd(), 'static', 'markdown', `${targetSlug.replace(/\//g, '-')}.md`),
			
			// For root path, try index.md
			path.join(process.cwd(), 'static', 'markdown', 'index.md'),
			
			// Alternative extension
			path.join(process.cwd(), 'static', 'markdown', `${targetSlug}.markdown`),
		];
		
		for (const tryPath of possiblePaths) {
			if (fs.existsSync(tryPath)) {
				filePath = tryPath;
				fileFound = true;
				break;
			}
		}
		
		// If no file was found, try one more approach for nested paths
		if (!fileFound && targetSlug.includes('/')) {
			// For paths like 'metrics/quick-start', try both 'metrics/quick-start.md' and 'metrics-quick-start.md'
			const segments = targetSlug.split('/');
			const lastSegment = segments.pop() || '';
			const parentPath = segments.join('/');
			
			const nestedPath = path.join(process.cwd(), 'static', 'markdown', parentPath, `${lastSegment}.md`);
			
			if (fs.existsSync(nestedPath)) {
				filePath = nestedPath;
				fileFound = true;
			}
		}
		
		// If still no file found, throw a 404 error
		if (!fileFound) {
			throw error(404, `Page not found: ${targetSlug}`);
		}
		
		// Read the file content
		content = fs.readFileSync(filePath, 'utf-8');
		
		// Return the slug and content
		const result = { slug: targetSlug, content };
		
		return result;
	} catch (e: any) {
		if (e.status === 404) {
			throw e; // Re-throw 404 errors
		}
		throw error(500, `Error loading page content: ${e.message || 'Unknown error'}`);
	}
}; 