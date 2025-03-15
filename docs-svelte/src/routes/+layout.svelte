<script lang="ts">
	import { base } from '$app/paths';
	import ParentNav from '$lib/components/ParentNav.svelte';
	import SubNav from '$lib/components/SubNav.svelte';
	import '../app.css';
</script>

<div class="app">
	<header>
		<div class="header-content">
			<div class="logo">
				<a href="{base}/">
					<img src="{base}/images/chadburn-logo.svg" alt="Chadburn Logo" />
					<span class="site-title">Chadburn</span>
				</a>
			</div>
			<ParentNav />
		</div>
	</header>

	<SubNav />

	<main>
		<div class="content">
			<slot />
		</div>
	</main>

	<footer>
		<div class="footer-content">
			<p>
				Chadburn is licensed under the <a
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

	.logo {
		display: flex;
		align-items: center;
	}

	.logo a {
		display: flex;
		align-items: center;
		text-decoration: none;
		color: white;
	}

	.logo img {
		height: 40px;
		margin-right: 0.5rem;
	}

	.site-title {
		font-size: 1.5rem;
		font-weight: bold;
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

	@media (max-width: 768px) {
		.header-content {
			flex-direction: column;
			align-items: center;
		}

		.logo {
			margin-bottom: 1rem;
		}
	}
</style>
