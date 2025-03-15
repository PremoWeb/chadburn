<script lang="ts">
	import { onMount } from 'svelte';
	import type { Contributor } from '$lib/types';
	
	export let contributors: Contributor[] = [];
	export let loading: boolean = true;
	export let error: string | null = null;
	
	// Sort contributors: lead developer first, then by contributions
	$: sortedContributors = [...contributors].sort((a, b) => {
		if (a.isLeadDeveloper && !b.isLeadDeveloper) return -1;
		if (!a.isLeadDeveloper && b.isLeadDeveloper) return 1;
		return b.contributions - a.contributions;
	});
</script>

<section class="contributors">
	<h2>Thank You to Our Contributors</h2>
	
	{#if loading}
		<div class="loading">Loading contributors...</div>
	{:else if error}
		<div class="error">
			<p>Error loading contributors: {error}</p>
		</div>
	{:else if contributors.length === 0}
		<div class="no-contributors">
			<p>No contributors found. Be the first to contribute!</p>
		</div>
	{:else}
		<div class="contributors-grid">
			{#each sortedContributors as contributor}
				<div class="contributor-card {contributor.isLeadDeveloper ? 'lead-developer' : ''}">
					<img src={contributor.avatar_url} alt={contributor.name || contributor.login} class="avatar" />
					<div class="contributor-info">
						<h3>
							{contributor.name || contributor.login}
							{#if contributor.isLeadDeveloper}
								<span class="lead-badge">Lead Developer</span>
							{/if}
						</h3>
						<p class="contributions">Contributions: {contributor.contributions}</p>
						
						{#if contributor.social}
							<div class="social-links">
								{#if contributor.social.github}
									<a 
										href={contributor.social.github} 
										target="_blank" 
										rel="noopener noreferrer" 
										title="GitHub" 
										aria-label="GitHub profile for {contributor.name || contributor.login}"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" aria-hidden="true" focusable="false">
											<path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
										</svg>
									</a>
								{:else}
									<a 
										href={contributor.html_url} 
										target="_blank" 
										rel="noopener noreferrer" 
										title="GitHub" 
										aria-label="GitHub profile for {contributor.name || contributor.login}"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" aria-hidden="true" focusable="false">
											<path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
										</svg>
									</a>
								{/if}
								
								{#if contributor.social.x}
									<a 
										href={contributor.social.x} 
										target="_blank" 
										rel="noopener noreferrer" 
										title="X (Twitter)" 
										aria-label="X (Twitter) profile for {contributor.name || contributor.login}"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" aria-hidden="true" focusable="false">
											<path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
										</svg>
									</a>
								{/if}
								
								{#if contributor.social.linkedin}
									<a 
										href={contributor.social.linkedin} 
										target="_blank" 
										rel="noopener noreferrer" 
										title="LinkedIn" 
										aria-label="LinkedIn profile for {contributor.name || contributor.login}"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" aria-hidden="true" focusable="false">
											<path d="M19 0h-14c-2.761 0-5 2.239-5 5v14c0 2.761 2.239 5 5 5h14c2.762 0 5-2.239 5-5v-14c0-2.761-2.238-5-5-5zm-11 19h-3v-11h3v11zm-1.5-12.268c-.966 0-1.75-.79-1.75-1.764s.784-1.764 1.75-1.764 1.75.79 1.75 1.764-.783 1.764-1.75 1.764zm13.5 12.268h-3v-5.604c0-3.368-4-3.113-4 0v5.604h-3v-11h3v1.765c1.396-2.586 7-2.777 7 2.476v6.759z"/>
										</svg>
									</a>
								{/if}
								
								{#if contributor.social.website}
									<a 
										href={contributor.social.website} 
										target="_blank" 
										rel="noopener noreferrer" 
										title="Website" 
										aria-label="Personal website for {contributor.name || contributor.login}"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" aria-hidden="true" focusable="false">
											<path d="M12 0c-6.627 0-12 5.373-12 12s5.373 12 12 12 12-5.373 12-12-5.373-12-12-12zm1 16.057v-3.057h2.994c-.059 1.143-.212 2.24-.456 3.279-.823-.12-1.674-.188-2.538-.222zm1.957 2.162c-.499 1.33-1.159 2.497-1.957 3.456v-3.62c.666.028 1.319.081 1.957.164zm-1.957-7.219v-3.015c.868-.034 1.721-.103 2.548-.224.238 1.027.389 2.111.446 3.239h-2.994zm0-5.014v-3.661c.806.969 1.471 2.15 1.971 3.496-.642.084-1.3.137-1.971.165zm2.703-3.267c1.237.496 2.354 1.228 3.29 2.146-.642.234-1.311.442-2.019.607-.344-.992-.775-1.91-1.271-2.753zm-7.241 13.56c-.244-1.039-.398-2.136-.456-3.279h2.994v3.057c-.865.034-1.714.102-2.538.222zm2.538 1.776v3.62c-.798-.959-1.458-2.126-1.957-3.456.638-.083 1.291-.136 1.957-.164zm-2.994-7.055c.057-1.128.207-2.212.446-3.239.827.121 1.68.19 2.548.224v3.015h-2.994zm1.024-5.179c.5-1.346 1.165-2.527 1.97-3.496v3.661c-.671-.028-1.329-.081-1.97-.165zm-2.005-.35c-.708-.165-1.377-.373-2.018-.607.937-.918 2.053-1.65 3.29-2.146-.496.844-.927 1.762-1.272 2.753zm-.549 1.918c-.264 1.151-.434 2.36-.492 3.611h-3.933c.165-1.658.739-3.197 1.617-4.518.88.361 1.816.67 2.808.907zm.009 9.262c-.988.236-1.92.542-2.797.9-.89-1.328-1.471-2.879-1.637-4.551h3.934c.058 1.265.231 2.488.5 3.651zm.553 1.917c.342.976.768 1.881 1.257 2.712-1.223-.49-2.326-1.211-3.256-2.115.636-.229 1.299-.435 1.999-.597zm9.924 0c.7.163 1.362.367 1.999.597-.931.903-2.034 1.625-3.257 2.116.489-.832.915-1.737 1.258-2.713zm.553-1.917c.27-1.163.442-2.386.501-3.651h3.934c-.167 1.672-.748 3.223-1.638 4.551-.877-.358-1.81-.664-2.797-.9zm.501-5.651c-.058-1.251-.229-2.46-.492-3.611.992-.237 1.929-.546 2.809-.907.877 1.321 1.451 2.86 1.616 4.518h-3.933z"/>
										</svg>
									</a>
								{/if}
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</section>

<style>
	.contributors {
		margin-bottom: 4rem;
	}
	
	h2 {
		text-align: center;
		margin-bottom: 2rem;
	}
	
	.loading, .error, .no-contributors {
		text-align: center;
		padding: 2rem;
		background-color: #f9f9f9;
		border-radius: 8px;
		margin-bottom: 1rem;
	}
	
	.error {
		color: #721c24;
		background-color: #f8d7da;
		border: 1px solid #f5c6cb;
	}
	
	.contributors-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}
	
	.contributor-card {
		display: flex;
		align-items: center;
		padding: 1.5rem;
		background-color: #f9f9f9;
		border-radius: 8px;
		color: inherit;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		transition: transform 0.2s, box-shadow 0.2s;
	}
	
	.contributor-card:hover {
		transform: translateY(-5px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	}
	
	.contributor-card.lead-developer {
		background-color: #f0f7ff;
		border: 1px solid #d0e3ff;
		grid-column: 1 / -1;
	}
	
	.avatar {
		width: 80px;
		height: 80px;
		border-radius: 50%;
		margin-right: 1.5rem;
		object-fit: cover;
	}
	
	.lead-developer .avatar {
		width: 100px;
		height: 100px;
	}
	
	.contributor-info {
		flex-grow: 1;
	}
	
	.contributor-info h3 {
		margin: 0 0 0.5rem 0;
		color: var(--primary-color);
		font-size: 1.2rem;
		display: flex;
		align-items: center;
		flex-wrap: wrap;
		gap: 0.5rem;
	}
	
	.lead-developer .contributor-info h3 {
		font-size: 1.5rem;
	}
	
	.lead-badge {
		font-size: 0.7rem;
		background-color: #155799;
		color: white;
		padding: 0.2rem 0.5rem;
		border-radius: 4px;
		font-weight: normal;
		display: inline-block;
		vertical-align: middle;
	}
	
	.contributions {
		margin: 0 0 0.5rem 0;
		font-size: 0.9rem;
		color: var(--light-text-color);
	}
	
	.social-links {
		display: flex;
		gap: 1rem;
		margin-top: 0.5rem;
	}
	
	.social-links a {
		color: #666;
		transition: color 0.2s;
	}
	
	.social-links a:hover {
		color: var(--primary-color);
	}
	
	@media (max-width: 768px) {
		.contributors-grid {
			grid-template-columns: 1fr;
		}
		
		.contributor-card {
			padding: 1rem;
		}
		
		.avatar {
			width: 60px;
			height: 60px;
			margin-right: 1rem;
		}
		
		.lead-developer .avatar {
			width: 80px;
			height: 80px;
		}
	}
</style> 