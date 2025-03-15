<script lang="ts">
	import { page } from '$app/stores';
	import { base } from '$app/paths';
	
	// Define the documentation navigation items
	const docNavItems = [
		{ name: 'Getting Started', path: '/getting-started' },
		{ name: 'Configuration', path: '/configuration' },
		{ name: 'Job Types', path: '/jobs' },
		{ name: 'Docker Integration', path: '/docker-integration' },
		{ name: 'Metrics Setup', path: '/metrics-setup' },
		{ name: 'Troubleshooting', path: '/troubleshooting' },
		{ name: 'Advanced Topics', path: '/advanced-topics' },
		{ name: 'Deployment', path: '/deployment' },
		{ name: 'Examples', path: '/examples' },
		{ name: 'FAQ', path: '/faq' },
		{ name: 'Contributing', path: '/contributing' }
	];
	
	// Check if we're on a documentation page
	$: isDocPage = $page.url.pathname !== `${base}/` && 
	               !$page.url.pathname.endsWith('/404');
</script>

{#if isDocPage}
	<nav class="sub-nav">
		<ul>
			{#each docNavItems as item}
				<li class:active={$page.url.pathname === `${base}${item.path}`}>
					<a href="{base}{item.path}">{item.name}</a>
				</li>
			{/each}
		</ul>
	</nav>
{/if}

<style>
	.sub-nav {
		background-color: #f5f5f5;
		border-bottom: 1px solid #e0e0e0;
		padding: 0.5rem 0;
	}

	.sub-nav ul {
		display: flex;
		list-style: none;
		margin: 0;
		padding: 0;
		max-width: 1200px;
		margin: 0 auto;
		padding: 0 1rem;
		flex-wrap: wrap;
	}

	.sub-nav li {
		margin: 0.25rem 0.75rem 0.25rem 0;
	}

	.sub-nav a {
		color: #333;
		text-decoration: none;
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		transition: background-color 0.2s;
		font-size: 0.9rem;
	}

	.sub-nav a:hover {
		background-color: rgba(0, 0, 0, 0.05);
	}

	.sub-nav li.active a {
		background-color: #155799;
		color: white;
		font-weight: bold;
	}

	@media (max-width: 768px) {
		.sub-nav ul {
			flex-direction: row;
			flex-wrap: wrap;
			justify-content: center;
		}
		
		.sub-nav li {
			margin: 0.25rem 0.5rem;
		}
	}
</style> 