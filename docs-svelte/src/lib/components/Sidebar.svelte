<script lang="ts">
	import { page } from '$app/state';
	import { base } from '$app/paths';
	
	// Navigation items
	const navItems = [
		{
			section: 'Getting Started',
			items: [
				{ title: 'Introduction', path: '/introduction' },
				{ title: 'Installation', path: '/installation' },
				{ title: 'Configuration', path: '/configuration' }
			]
		},
		{
			section: 'Metrics',
			items: [
				{ title: 'Overview', path: '/metrics' },
				{ title: 'Quick Start', path: '/metrics/quick-start' },
				{ title: 'Test Data Generator', path: '/metrics/test-data-generator' }
			]
		},
		{
			section: 'Deployment',
			items: [
				{ title: 'Overview', path: '/deployment' },
				{ title: 'Docker', path: '/deployment/docker' },
				{ title: 'Kubernetes', path: '/deployment/kubernetes' }
			]
		},
		{
			section: 'Integrations',
			items: [
				{ title: 'Overview', path: '/integrations' },
				{ title: 'Traefik', path: '/integrations/traefik' },
				{ title: 'Caddy', path: '/integrations/caddy' }
			]
		},
		{
			section: 'Help',
			items: [
				{ title: 'FAQ', path: '/faq' },
				{ title: 'Troubleshooting', path: '/troubleshooting' }
			]
		}
	];
	
	// Reactive values using runes
	let baseUrl = $derived(base);
	let currentPath = $derived(page.url.pathname);
	
	// Function to find the most specific active path
	function findMostSpecificActivePath(): string {
		// Flatten all navigation items into a single array
		const allPaths = navItems.flatMap(section => 
			section.items.map(item => item.path)
		);
		
		// Sort paths by length (longest/most specific first)
		const sortedPaths = [...allPaths].sort((a, b) => b.length - a.length);
		
		// Find the most specific path that matches the current path
		for (const path of sortedPaths) {
			// Exact match
			if (currentPath === path) {
				return path;
			}
			
			// Path with trailing slash
			if (currentPath.startsWith(path + '/')) {
				return path;
			}
		}
		
		// If no match found, return empty string
		return '';
	}
	
	// Get the most specific active path
	let mostSpecificPath = $derived(findMostSpecificActivePath());
	
	// Function to check if a navigation item is active
	function isActive(path: string): boolean {
		// Handle hash links specially
		if (path.includes('#')) {
			const [basePath, hash] = path.split('#');
			return currentPath === basePath || currentPath.startsWith(basePath + '/');
		}
		
		// For regular paths, only the most specific path is active
		return path === mostSpecificPath;
	}
	
	// Toggle mobile sidebar
	let isMobileSidebarOpen = $state(false);
	function toggleMobileSidebar() {
		isMobileSidebarOpen = !isMobileSidebarOpen;
	}
</script>

<div class="sidebar-container">
	<button class="mobile-toggle" onclick={toggleMobileSidebar}>
		{isMobileSidebarOpen ? 'Close Menu' : 'Open Menu'}
	</button>
	
	<aside class="sidebar" class:open={isMobileSidebarOpen}>
		{#each navItems as section}
			<div class="nav-section">
				<h3>{section.section}</h3>
				<ul>
					{#each section.items as item}
						<li class:active={isActive(item.path)}>
							<a href="{baseUrl}{item.path}">{item.title}</a>
						</li>
					{/each}
				</ul>
			</div>
		{/each}
	</aside>
</div>

<style>
	.sidebar-container {
		position: relative;
		padding-top: 0.5rem;
	}
	
	.mobile-toggle {
		display: none;
		width: 100%;
		padding: 0.5rem;
		background-color: var(--primary-color);
		color: white;
		border: none;
		border-radius: 4px;
		font-weight: 600;
		cursor: pointer;
		margin-bottom: 0.75rem;
		font-size: 0.9rem;
	}
	
	.sidebar {
		width: 200px;
		padding-right: 1rem;
		position: sticky;
		top: calc(1rem + var(--header-height, 60px));
		max-height: calc(100vh - var(--header-height, 60px) - 2rem);
		overflow-y: auto;
	}
	
	.nav-section {
		margin-bottom: 1rem;
	}
	
	.nav-section h3 {
		font-size: 0.8rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: var(--light-text-color);
		margin-bottom: 0.25rem;
		margin-top: 0.5rem;
	}
	
	.nav-section ul {
		list-style: none;
		padding: 0;
		margin: 0;
	}
	
	.nav-section li {
		margin-bottom: 0.1rem;
	}
	
	.nav-section a {
		display: block;
		padding: 0.3rem 0;
		color: var(--text-color);
		text-decoration: none;
		font-size: 0.85rem;
		border-left: 2px solid transparent;
		padding-left: 0.5rem;
		transition: all 0.2s ease;
	}
	
	.nav-section a:hover {
		color: var(--primary-color);
		border-left-color: var(--light-border-color);
	}
	
	.nav-section li.active a {
		color: var(--primary-color);
		border-left-color: var(--primary-color);
		font-weight: 600;
	}
	
	@media (max-width: 768px) {
		.mobile-toggle {
			display: block;
		}
		
		.sidebar {
			display: none;
			position: fixed;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background-color: white;
			z-index: 1000;
			padding: 1rem;
			overflow-y: auto;
			box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
		}
		
		.sidebar.open {
			display: block;
		}
	}
</style> 