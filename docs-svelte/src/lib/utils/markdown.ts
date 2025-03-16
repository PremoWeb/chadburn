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
 * Add IDs to headings in HTML content for anchor links
 */
export function addIdsToHeadings(html: string): string {
	const parser = new DOMParser();
	const doc = parser.parseFromString(html, 'text/html');
	
	// Process h1, h2, h3, h4 headings
	const headings = doc.querySelectorAll('h1, h2, h3, h4');
	
	headings.forEach((heading) => {
		if (!heading.id) {
			const text = heading.textContent || '';
			const id = text
				.toLowerCase()
				.replace(/[^\w\s-]/g, '') // Remove special characters
				.replace(/\s+/g, '-') // Replace spaces with hyphens
				.replace(/-+/g, '-'); // Replace multiple hyphens with single hyphen
			
			heading.id = id;
		}
	});
	
	return new XMLSerializer().serializeToString(doc.body).replace(/<\/?body>/g, '');
} 