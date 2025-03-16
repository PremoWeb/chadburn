<script lang="ts">
	import { page } from '$app/state';
	import { base } from '$app/paths';
	
	// Props
	const { navItems = [], title = 'Documentation' } = $props<{
		navItems: Array<{
			section: string;
			items: Array<{ title: string; path: string }>;
		}>;
		title?: string;
	}>();
	
	// Reactive values using runes
	let baseUrl = $derived(base);
	let currentPath = $derived(page.url.pathname);
	
	// Mobile sidebar state
	let isMobileSidebarOpen = $state(false);
	
	// Toggle mobile sidebar
	function toggleMobileSidebar() {
		isMobileSidebarOpen = !isMobileSidebarOpen;
	}
	
	// Check if a path is active
	function isActive(path: string): boolean {
		// Exact match
		if (currentPath === `${baseUrl}${path}`) {
			return true;
		}
		
		// Check if it's a parent path (but not the root)
		if (path !== '/' && path !== '/api' && path !== '/api/overview' && 
			path !== '/docs' && path !== '/docs/overview' && 
			currentPath.startsWith(`${baseUrl}${path}/`)) {
			return true;
		}
		
		// Special case for index pages
		if ((path === '/api/overview' && (currentPath === `${baseUrl}/api` || currentPath === `${baseUrl}/api/`)) ||
			(path === '/docs/overview' && (currentPath === `${baseUrl}/docs` || currentPath === `${baseUrl}/docs/`))) {
			return true;
		}
		
		return false;
	}
</script>

<div class="sidebar-container" class:mobile-open={isMobileSidebarOpen}>
	<button class="mobile-toggle" on:click={toggleMobileSidebar}>
		{isMobileSidebarOpen ? 'Close Menu' : 'Open Menu'}
	</button>
	
	<nav class="sidebar">
		<div class="sidebar-header">
			<h2>{title}</h2>
		</div>
		
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
	</nav>
</div>

<style>
	.sidebar-container {
		width: 100%;
		position: relative;
	}
	
	.sidebar {
		position: sticky;
		top: 1rem;
		max-height: calc(100vh - 2rem);
		overflow-y: auto;
		padding-right: 1rem;
		scrollbar-width: thin;
	}
	
	.sidebar-header {
		margin-bottom: 1rem;
	}
	
	.sidebar-header h2 {
		font-size: 1.2rem;
		font-weight: 600;
		color: #155799;
		margin: 0;
	}
	
	.nav-section {
		margin-bottom: 1.5rem;
	}
	
	.nav-section h3 {
		font-size: 0.9rem;
		font-weight: 600;
		color: #333;
		margin: 0 0 0.5rem 0;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	
	.nav-section ul {
		list-style: none;
		padding: 0;
		margin: 0;
	}
	
	.nav-section li {
		margin-bottom: 0.25rem;
		border-radius: 4px;
	}
	
	.nav-section li a {
		display: block;
		padding: 0.4rem 0.5rem;
		color: #555;
		text-decoration: none;
		font-size: 0.9rem;
		border-radius: 4px;
		transition: all 0.2s ease;
	}
	
	.nav-section li a:hover {
		background-color: rgba(21, 87, 153, 0.05);
		color: #155799;
	}
	
	.nav-section li.active a {
		background-color: rgba(21, 87, 153, 0.1);
		color: #155799;
		font-weight: 500;
	}
	
	.mobile-toggle {
		display: none;
		width: 100%;
		padding: 0.5rem;
		background-color: #155799;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		margin-bottom: 1rem;
	}
	
	@media (max-width: 768px) {
		.mobile-toggle {
			display: block;
		}
		
		.sidebar {
			display: none;
			position: static;
			max-height: none;
		}
		
		.mobile-open .sidebar {
			display: block;
		}
	}
</style> 