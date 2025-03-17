import { marked } from 'marked';

// Define heading type
export type Heading = {
	id: string;
	text: string;
	level: number;
};

/**
 * Strip YAML frontmatter from markdown content
 * Frontmatter is delimited by --- at the start and end
 */
export function stripFrontmatter(content: string): string {
	// Check if content starts with --- which indicates frontmatter
	if (!content.trim().startsWith('---')) {
		return content;
	}
	
	// Find the second --- which ends the frontmatter
	const parts = content.split('---');
	if (parts.length < 3) {
		return content; // No proper frontmatter found
	}
	
	// Return everything after the second ---
	return parts.slice(2).join('---').trim();
}

/**
 * Process markdown content and convert to HTML
 */
export async function processMarkdown(content: string): Promise<string> {
	try {
		const parsedContent = marked.parse(content);
		if (typeof parsedContent === 'string') {
			return parsedContent;
		} else {
			// Handle the case where marked.parse returns a Promise
			return await parsedContent;
		}
	} catch (e) {
		console.error('Error parsing markdown:', e);
		return '<p>Error parsing content</p>';
	}
}

/**
 * Extract headings from HTML content
 */
export function extractHeadings(selector: string = '.markdown-content'): Heading[] {
	if (typeof document === 'undefined') return [];
	
	const headings: Heading[] = [];
	const contentElement = document.querySelector(selector);
	
	if (!contentElement) return [];
	
	// Get all headings h2 and h3
	const headingElements = contentElement.querySelectorAll('h2, h3');
	
	headingElements.forEach((el) => {
		const id = el.id;
		const text = el.textContent || '';
		const level = parseInt(el.tagName.substring(1), 10);
		
		headings.push({ id, text, level });
	});
	
	return headings;
}

/**
 * Add IDs to headings in HTML for anchor links
 * @param html HTML content
 * @returns HTML with IDs added to headings
 */
export function addIdsToHeadings(html: string): string {
	// Create a DOM parser
	const parser = new DOMParser();
	const doc = parser.parseFromString(html, 'text/html');
	
	// Find all headings
	const headings = doc.querySelectorAll('h1, h2, h3, h4, h5, h6');
	console.log(`Found ${headings.length} headings to process for IDs`);
	
	// Track used IDs to ensure uniqueness
	const usedIds = new Set<string>();
	
	// Process each heading
	headings.forEach((heading, index) => {
		// Get the text content of the heading
		const text = heading.textContent?.trim() || '';
		
		// Check if the heading already has an ID
		if (!heading.id) {
			// Generate an ID from the text content
			let baseId = text
				.toLowerCase()
				.replace(/[^\w\s-]/g, '') // Remove special characters
				.replace(/\s+/g, '-') // Replace spaces with hyphens
				.replace(/-+/g, '-') // Replace multiple hyphens with a single one
				.replace(/^-+|-+$/g, ''); // Remove leading and trailing hyphens
			
			// If the ID is empty (e.g., heading only had special characters), use a fallback
			if (!baseId) {
				baseId = `heading-${index}`;
				console.warn(`Generated fallback ID "${baseId}" for empty heading text`);
			}
			
			// Ensure the ID is unique
			let uniqueId = baseId;
			let counter = 1;
			
			// If the ID is already used, append a number to make it unique
			while (usedIds.has(uniqueId)) {
				uniqueId = `${baseId}-${counter}`;
				counter++;
			}
			
			// Set the ID on the heading
			heading.id = uniqueId;
			usedIds.add(uniqueId);
			
			if (uniqueId !== baseId) {
				console.log(`Generated unique ID "${uniqueId}" for heading "${text}" (base ID "${baseId}" was already used)`);
			} else {
				console.log(`Generated ID "${uniqueId}" for heading "${text}"`);
			}
		} else {
			// If the heading already has an ID, make sure it's unique
			let existingId = heading.id;
			
			if (usedIds.has(existingId)) {
				// If the ID is already used, append a number to make it unique
				let counter = 1;
				let uniqueId = existingId;
				
				while (usedIds.has(uniqueId)) {
					uniqueId = `${existingId}-${counter}`;
					counter++;
				}
				
				console.log(`Changed duplicate ID "${existingId}" to "${uniqueId}" for heading "${text}"`);
				heading.id = uniqueId;
				usedIds.add(uniqueId);
			} else {
				console.log(`Heading "${text}" already has ID "${existingId}"`);
				usedIds.add(existingId);
			}
		}
		
		// Add a data attribute for easier debugging
		heading.setAttribute('data-toc-heading', 'true');
	});
	
	// Return the modified HTML
	return doc.body.innerHTML;
} 