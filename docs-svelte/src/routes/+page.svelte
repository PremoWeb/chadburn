<script lang="ts">
	import { base } from '$app/paths';
	import { onMount } from 'svelte';
	import Contributors from '$lib/components/Contributors.svelte';
	import { loadContributors } from '$lib/contributors';
	import type { Contributor } from '$lib/types';
	
	let contributors: Contributor[] = [];
	let loading = true;
	let error: string | null = null;
	
	onMount(async () => {
		const result = await loadContributors();
		contributors = result.contributors;
		error = result.error;
		loading = false;
	});
</script>

<svelte:head>
	<title>Chadburn - Modern job scheduler for Docker environments</title>
	<meta name="description" content="Chadburn is a lightweight job scheduler designed for Docker environments. It serves as a modern replacement for traditional cron, with enhanced features for container orchestration." />
</svelte:head>

<div class="home-container">
	<section class="hero">
		<h1>Chadburn</h1>
		<p class="tagline">Modern job scheduler for Docker environments</p>
		<div class="cta-buttons">
			<a href="{base}/docs" class="btn primary">Get Started</a>
			<a href="{base}/api" class="btn secondary">API Reference</a>
			<a href="https://github.com/PremoWeb/Chadburn" class="btn outline" target="_blank" rel="noopener noreferrer">GitHub</a>
		</div>
	</section>

	<section class="features">
		<h2>Key Features</h2>
		<div class="feature-grid">
			<div class="feature-card">
				<h3>Docker Integration</h3>
				<p>Native support for Docker containers with the ability to execute commands inside running containers or create new containers on schedule.</p>
			</div>
			<div class="feature-card">
				<h3>Dynamic Configuration</h3>
				<p>Configure jobs using Docker labels for a flexible, dynamic approach that doesn't require restarts when configuration changes.</p>
			</div>
			<div class="feature-card">
				<h3>Multiple Job Types</h3>
				<p>Support for various job types including local execution, container execution, and container lifecycle events.</p>
			</div>
			<div class="feature-card">
				<h3>Container Lifecycle Events</h3>
				<p>Execute commands when containers start or stop, enabling powerful automation workflows.</p>
			</div>
			<div class="feature-card">
				<h3>Notifications</h3>
				<p>Built-in support for Slack, Email, and Gotify notifications to keep you informed about job executions.</p>
			</div>
			<div class="feature-card">
				<h3>Metrics</h3>
				<p>Prometheus-compatible metrics endpoint for monitoring job executions and performance.</p>
			</div>
		</div>
	</section>

	<section class="api-highlight">
		<div class="api-content">
			<div class="api-text">
				<h2>New: REST API Documentation</h2>
				<p>Control and monitor your Chadburn instance programmatically with our comprehensive REST API.</p>
				<ul>
					<li>Create and manage jobs via API calls</li>
					<li>Monitor job status and execution history</li>
					<li>Trigger jobs on demand</li>
					<li>Integrate with your existing systems</li>
				</ul>
				<a href="{base}/api" class="btn primary">Explore the API</a>
			</div>
			<div class="api-image">
				<div class="code-snippet">
					{@html `<pre><code>// Example: Create a new job
fetch('/api/jobs', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    name: "backup-database",
    type: "job-local",
    schedule: "@daily",
    command: "/usr/local/bin/backup.sh"
  })
})</code></pre>`}
				</div>
			</div>
		</div>
	</section>

	<Contributors {contributors} {loading} {error} />

	<section class="quick-start">
		<h2>Quick Start</h2>
		<div class="code-example">
			<pre><code>docker run -d --name chadburn \
  -v /var/run/docker.sock:/var/run/docker.sock:ro,z \
  -v /path/to/config.ini:/etc/chadburn.conf \
  premoweb/chadburn:latest daemon</code></pre>
		</div>
		<p>Check out the <a href="{base}/docs/quick-start">Getting Started Guide</a> for more detailed instructions.</p>
	</section>

	<section class="documentation">
		<h2>Documentation</h2>
		<div class="doc-links">
			<a href="{base}/docs/quick-start" class="doc-link">
				<h3>Getting Started</h3>
				<p>Installation, quick start, and basic concepts</p>
			</a>
			<a href="{base}/docs/guides/configuration" class="doc-link">
				<h3>Configuration</h3>
				<p>Configuration file format, environment variables, and command line options</p>
			</a>
			<a href="{base}/docs/concepts/jobs" class="doc-link">
				<h3>Job Types</h3>
				<p>Different job types and their configuration options</p>
			</a>
			<a href="{base}/api" class="doc-link highlight">
				<h3>API Reference</h3>
				<p>Comprehensive REST API documentation for programmatic control</p>
			</a>
			<a href="{base}/docs/docker-integration" class="doc-link">
				<h3>Docker Integration</h3>
				<p>Docker labels, container lifecycle events, and dynamic configuration</p>
			</a>
			<a href="{base}/docs/examples" class="doc-link">
				<h3>Examples</h3>
				<p>Common use cases and Docker Compose examples</p>
			</a>
		</div>
	</section>
</div>

<style>
	.home-container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 0 1rem;
	}

	.hero {
		text-align: center;
		padding: 3rem 0;
	}

	.hero h1 {
		font-size: 3rem;
		margin: 0;
		color: var(--primary-color);
	}

	.tagline {
		font-size: 1.5rem;
		color: var(--light-text-color);
		margin-bottom: 2rem;
	}

	.cta-buttons {
		display: flex;
		justify-content: center;
		gap: 1rem;
		margin-top: 2rem;
		flex-wrap: wrap;
	}

	.btn {
		display: inline-block;
		padding: 0.75rem 1.5rem;
		border-radius: 4px;
		font-weight: bold;
		text-decoration: none;
		transition: all 0.2s;
	}

	.btn.primary {
		background-color: var(--primary-color);
		color: white;
	}

	.btn.primary:hover {
		background-color: #0e4377;
		text-decoration: none;
	}

	.btn.secondary {
		background-color: #159957;
		color: white;
	}

	.btn.secondary:hover {
		background-color: #0e7040;
		text-decoration: none;
	}

	.btn.outline {
		background-color: transparent;
		color: var(--primary-color);
		border: 1px solid var(--primary-color);
	}

	.btn.outline:hover {
		background-color: #f5f5f5;
		text-decoration: none;
	}

	section {
		margin-bottom: 4rem;
	}

	h2 {
		text-align: center;
		margin-bottom: 2rem;
	}

	.feature-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 2rem;
	}

	.feature-card {
		background-color: #f9f9f9;
		border-radius: 8px;
		padding: 1.5rem;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		transition: transform 0.2s, box-shadow 0.2s;
	}

	.feature-card:hover {
		transform: translateY(-5px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	}

	.feature-card h3 {
		margin-top: 0;
		color: var(--primary-color);
	}

	.api-highlight {
		background-color: #f0f7ff;
		border-radius: 12px;
		padding: 2rem;
		margin: 3rem 0;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
	}

	.api-content {
		display: flex;
		gap: 2rem;
		align-items: center;
	}

	.api-text {
		flex: 1;
	}

	.api-text h2 {
		text-align: left;
		color: var(--primary-color);
		margin-top: 0;
	}

	.api-text ul {
		margin-bottom: 1.5rem;
		padding-left: 1.5rem;
	}

	.api-text li {
		margin-bottom: 0.5rem;
	}

	.api-image {
		flex: 1;
	}

	.code-snippet {
		background-color: #2d2d2d;
		color: #f8f8f2;
		border-radius: 8px;
		padding: 1.5rem;
		overflow-x: auto;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
	}

	.code-snippet pre {
		margin: 0;
	}

	.code-example {
		background-color: #f6f8fa;
		border-radius: 8px;
		padding: 1.5rem;
		margin-bottom: 1.5rem;
		overflow-x: auto;
	}

	.code-example pre {
		margin: 0;
		background-color: transparent;
		border: none;
	}

	.doc-links {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.doc-link {
		display: block;
		padding: 1.5rem;
		background-color: #f9f9f9;
		border-radius: 8px;
		text-decoration: none;
		color: inherit;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		transition: transform 0.2s, box-shadow 0.2s;
	}

	.doc-link:hover {
		transform: translateY(-5px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
		text-decoration: none;
	}

	.doc-link h3 {
		margin-top: 0;
		color: var(--primary-color);
	}

	.doc-link p {
		color: var(--light-text-color);
		margin-bottom: 0;
	}

	.doc-link.highlight {
		background-color: #f0f7ff;
		border: 1px solid #d0e8ff;
	}

	@media (max-width: 768px) {
		.hero h1 {
			font-size: 2.5rem;
		}

		.tagline {
			font-size: 1.25rem;
		}

		.feature-grid,
		.doc-links {
			grid-template-columns: 1fr;
		}

		.api-content {
			flex-direction: column;
		}

		.api-text, .api-image {
			width: 100%;
		}
	}
</style>
