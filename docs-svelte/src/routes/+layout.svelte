<script lang="ts">
	import { base } from '$app/paths';
	import ParentNav from '$lib/components/ParentNav.svelte';
	import { version } from '$lib/stores/version';
	import WarningBanner from '$lib/components/WarningBanner.svelte';
	import '../app.css';
</script>

<div class="app">
	<header>
		<div class="header-content">
			<div class="site-header">
				<a href="{base}/">
					<span class="site-title">Chadburn</span>
					{#if !$version.loading && $version.version}
						<span class="version-badge">v{$version.version}</span>
					{/if}
				</a>
			</div>
			<ParentNav />
		</div>
	</header>

	<WarningBanner />

	<main>
		<div class="content">
			<slot />
		</div>
	</main>

	<footer>
		<div class="footer-content">
			<p>
				Chadburn
				{#if !$version.loading && $version.version}
					<span class="version">v{$version.version}</span>
				{/if}
				is licensed under the <a
					href="https://github.com/PremoWeb/Chadburn/blob/main/LICENSE"
					target="_blank"
					rel="noopener noreferrer">MIT License</a
				>.
			</p>
			<p>
				Documentation built with <a
					href="https://kit.svelte.dev/"
					target="_blank"
					rel="noopener noreferrer">SvelteKit</a
				>.
			</p>
		</div>
	</footer>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	header {
		background-color: #155799;
		background-image: linear-gradient(120deg, #155799, #159957);
		color: white;
		padding: 1rem 0;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.header-content {
		max-width: 1200px;
		margin: 0 auto;
		padding: 0 1rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		flex-wrap: wrap;
	}

	.site-header {
		display: flex;
		align-items: center;
	}

	.site-header a {
		display: flex;
		align-items: center;
		text-decoration: none;
		color: white;
	}

	.site-title {
		font-size: 1.5rem;
		font-weight: bold;
	}

	.version-badge {
		font-size: 0.8rem;
		background-color: rgba(255, 255, 255, 0.2);
		padding: 0.1rem 0.4rem;
		border-radius: 0.25rem;
		margin-left: 0.5rem;
		font-weight: normal;
	}

	main {
		flex-grow: 1;
		padding: 2rem 0;
	}

	.content {
		max-width: 1200px;
		margin: 0 auto;
		padding: 0 1rem;
	}

	footer {
		background-color: #f5f5f5;
		padding: 1rem 0;
		border-top: 1px solid #e0e0e0;
	}

	.footer-content {
		max-width: 1200px;
		margin: 0 auto;
		padding: 0 1rem;
		text-align: center;
		color: #666;
	}

	footer a {
		color: #155799;
		text-decoration: none;
	}

	footer a:hover {
		text-decoration: underline;
	}

	.version {
		font-size: 0.8rem;
		color: #888;
	}

	@media (max-width: 768px) {
		.header-content {
			flex-direction: column;
			align-items: center;
		}

		.site-header {
			margin-bottom: 1rem;
		}
	}
</style>
