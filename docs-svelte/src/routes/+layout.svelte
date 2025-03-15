<script lang="ts">
	import { page } from '$app/stores';
	import { base } from '$app/paths';
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
			<nav>
				<ul>
					<li class:active={$page.url.pathname === `${base}/`}>
						<a href="{base}/">Home</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/getting-started`}>
						<a href="{base}/getting-started">Getting Started</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/configuration`}>
						<a href="{base}/configuration">Configuration</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/jobs`}>
						<a href="{base}/jobs">Job Types</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/docker-integration`}>
						<a href="{base}/docker-integration">Docker Integration</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/troubleshooting`}>
						<a href="{base}/troubleshooting">Troubleshooting</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/advanced-topics`}>
						<a href="{base}/advanced-topics">Advanced Topics</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/deployment`}>
						<a href="{base}/deployment">Deployment</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/examples`}>
						<a href="{base}/examples">Examples</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/faq`}>
						<a href="{base}/faq">FAQ</a>
					</li>
					<li class:active={$page.url.pathname === `${base}/contributing`}>
						<a href="{base}/contributing">Contributing</a>
					</li>
				</ul>
			</nav>
			<div class="github-link">
				<a href="https://github.com/PremoWeb/Chadburn" target="_blank" rel="noopener noreferrer" aria-label="GitHub Repository">
					<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
						<path
							d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
						/>
					</svg>
				</a>
			</div>
		</div>
	</header>

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

	nav {
		flex-grow: 1;
		margin: 0 2rem;
	}

	nav ul {
		display: flex;
		list-style: none;
		margin: 0;
		padding: 0;
		flex-wrap: wrap;
		justify-content: center;
	}

	nav li {
		margin: 0 0.5rem;
	}

	nav a {
		color: white;
		text-decoration: none;
		padding: 0.5rem;
		border-radius: 4px;
		transition: background-color 0.2s;
	}

	nav a:hover {
		background-color: rgba(255, 255, 255, 0.1);
	}

	nav li.active a {
		background-color: rgba(255, 255, 255, 0.2);
		font-weight: bold;
	}

	.github-link a {
		color: white;
		display: flex;
		align-items: center;
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

		nav {
			margin: 1rem 0;
			width: 100%;
		}

		nav ul {
			flex-direction: column;
			align-items: center;
		}

		nav li {
			margin: 0.25rem 0;
		}

		.github-link {
			margin-top: 1rem;
		}
	}
</style>
