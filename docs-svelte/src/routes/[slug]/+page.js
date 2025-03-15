export function load({ params }) {
	return {
		slug: params.slug
	};
}

// Prerender all documentation pages
export const prerender = true;

// Enable client-side routing
export const csr = true;

// Define the list of routes to prerender
export const entries = () => [
  { slug: 'getting-started' },
  { slug: 'configuration' },
  { slug: 'jobs' },
  { slug: 'docker-integration' },
  { slug: 'troubleshooting' },
  { slug: 'advanced-topics' },
  { slug: 'deployment' },
  { slug: 'examples' },
  { slug: 'faq' },
  { slug: 'contributing' }
]; 