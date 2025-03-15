<script lang="ts">
	import { onMount } from 'svelte';
	import { renderMarkdown, fetchAndRenderMarkdown } from '$lib/markdown';
	
	export let content: string | undefined = undefined;
	export let path: string | undefined = undefined;
	
	let html = '';
	let loading = true;
	let error: string | null = null;
	
	onMount(async () => {
		try {
			loading = true;
			
			if (content) {
				html = renderMarkdown(content);
			} else if (path) {
				html = await fetchAndRenderMarkdown(path);
			} else {
				error = 'No content or path provided';
			}
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
			console.error('Error rendering markdown:', e);
		} finally {
			loading = false;
		}
	});
</script>

{#if loading}
	<div class="loading">Loading content...</div>
{:else if error}
	<div class="error">
		<p>Error loading content: {error}</p>
	</div>
{:else}
	<div class="markdown-content">
		{@html html}
	</div>
{/if}

<style>
	.loading {
		padding: 2rem;
		text-align: center;
		color: #666;
	}
	
	.error {
		padding: 1rem;
		color: #721c24;
		background-color: #f8d7da;
		border: 1px solid #f5c6cb;
		border-radius: 0.25rem;
	}
	
	.markdown-content :global(pre) {
		padding: 1rem;
		background-color: #f6f8fa;
		border-radius: 0.25rem;
		overflow-x: auto;
	}
	
	.markdown-content :global(code) {
		font-family: SFMono-Regular, Consolas, 'Liberation Mono', Menlo, monospace;
	}
	
	.markdown-content :global(table) {
		border-collapse: collapse;
		width: 100%;
		margin-bottom: 1rem;
	}
	
	.markdown-content :global(th),
	.markdown-content :global(td) {
		padding: 0.5rem;
		border: 1px solid #ddd;
	}
	
	.markdown-content :global(th) {
		background-color: #f6f8fa;
	}
	
	.markdown-content :global(tr:nth-child(even)) {
		background-color: #f9f9f9;
	}
</style> 