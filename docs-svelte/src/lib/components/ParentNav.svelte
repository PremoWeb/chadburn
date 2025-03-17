<script lang="ts">
	import { page } from '$app/stores';
	import { base } from '$app/paths';
	
	// Define the main navigation categories
	const mainNavItems = [
		{ name: 'Home', path: '/' },
		{ name: 'Documentation', path: '/docs' },
		{ name: 'API', path: '/api' },
		{ name: 'Changelog', path: '/changelog' },
		{ name: 'GitHub', url: 'https://github.com/PremoWeb/Chadburn', external: true }
	];
</script>

<nav class="parent-nav">
	<ul>
		{#each mainNavItems as item}
			{#if item.external}
				<li>
					<a href={item.url} target="_blank" rel="noopener noreferrer">
						{item.name}
						{#if item.name === 'GitHub'}
							<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
								<path
									d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
								/>
							</svg>
						{/if}
					</a>
				</li>
			{:else}
				<li class:active={$page.url.pathname.startsWith(`${base}${item.path}`)}>
					<a href="{base}{item.path}">{item.name}</a>
				</li>
			{/if}
		{/each}
	</ul>
</nav>

<style>
	.parent-nav {
		flex-grow: 0;
	}

	.parent-nav ul {
		display: flex;
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.parent-nav li {
		margin: 0 1rem 0 0;
	}

	.parent-nav a {
		color: white;
		text-decoration: none;
		padding: 0.5rem;
		border-radius: 4px;
		transition: background-color 0.2s;
		display: flex;
		align-items: center;
		font-weight: 500;
	}

	.parent-nav a svg {
		margin-left: 0.5rem;
	}

	.parent-nav a:hover {
		background-color: rgba(255, 255, 255, 0.1);
	}

	.parent-nav li.active a {
		background-color: rgba(255, 255, 255, 0.2);
		font-weight: bold;
	}

	@media (max-width: 768px) {
		.parent-nav ul {
			justify-content: center;
			width: 100%;
		}
		
		.parent-nav li {
			margin: 0 0.5rem;
		}
	}
</style> 