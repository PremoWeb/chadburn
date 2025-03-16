<script lang="ts">
	import { page } from '$app/state';
	import { base } from '$app/paths';
	import { afterNavigate, beforeNavigate } from '$app/navigation';
	import { marked } from 'marked';
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import ApiSidebar from '$lib/components/ApiSidebar.svelte';
	import TableOfContents from '$lib/components/TableOfContents.svelte';
	
	// Get data from the page load function
	let { data } = $props();
	
	// Track if this is the initial page load
	let isInitialLoad = true;
	
	// Track current navigation state
	let isNavigating = $state(false);
	let currentHtmlContent = $state('');
	let isLoading = $state(true);
	
	// Process markdown content
	$effect(() => {
		if (data.content) {
			// Parse markdown to HTML
			try {
				currentHtmlContent = marked.parse(data.content);
			} catch (e) {
				console.error('Error parsing markdown:', e);
				currentHtmlContent = '<p>Error parsing content</p>';
			}
			isLoading = false;
		}
	});
	
	// Handle navigation events
	beforeNavigate(() => {
		isNavigating = true;
	});
	
	afterNavigate(() => {
		isNavigating = false;
		isInitialLoad = false;
	});
	
	// Extract headings for table of contents
	function extractHeadings() {
		if (typeof document === 'undefined') return [];
		
		const headings: { id: string; text: string; level: number }[] = [];
		const contentElement = document.querySelector('.markdown-content');
		
		if (!contentElement) return [];
		
		// Get all headings h2 and h3
		const headingElements = contentElement.querySelectorAll('h2, h3');
		
		headingElements.forEach((el) => {
			const id = el.id;
			const text = el.textContent || '';
			const level = parseInt(el.tagName.substring(1), 10);
			
			headings.push({ id, text, level });
		});
		
		return headings;
	}
	
	// Handle heading extraction after content is rendered
	let headings = $state<{ id: string; text: string; level: number }[]>([]);
	
	onMount(() => {
		// Extract headings after initial render
		setTimeout(() => {
			headings = extractHeadings();
		}, 100);
	});
	
	// Re-extract headings when content changes
	$effect(() => {
		if (currentHtmlContent && !isLoading && typeof document !== 'undefined') {
			setTimeout(() => {
				headings = extractHeadings();
			}, 100);
		}
	});
</script>

<svelte:head>
	<title>Chadburn API - {data.slug.split('/').pop()?.replace(/-/g, ' ') || 'API Documentation'}</title>
</svelte:head>

<div class="api-page">
	<div class="sidebar-container">
		<ApiSidebar />
	</div>
	
	<div class="content-container">
		{#if isLoading}
			<div class="loading-container" in:fade={{ duration: 200 }}>
				<div class="loading-spinner"></div>
				<p>Loading API documentation...</p>
			</div>
		{:else}
			<article 
				class="markdown-content"
				in:fade={{ duration: 200, delay: 100 }}
			>
				{@html currentHtmlContent}
			</article>
		{/if}
	</div>
	
	{#if headings.length > 0}
		<div class="toc-container">
			<TableOfContents headings={headings} />
		</div>
	{/if}
</div>

<style>
	.api-page {
		display: grid;
		grid-template-columns: 200px 1fr 250px;
		gap: 2rem;
		position: relative;
		max-width: 1200px;
		margin: 0 auto;
		padding: 0 1rem;
	}
	
	.sidebar-container {
		grid-column: 1;
	}
	
	.content-container {
		grid-column: 2;
		min-height: 500px;
		position: relative;
	}
	
	.toc-container {
		grid-column: 3;
		position: sticky;
		top: 1rem;
		align-self: start;
		max-height: calc(100vh - 2rem);
		overflow-y: auto;
	}
	
	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 300px;
		color: #666;
	}
	
	.loading-spinner {
		width: 40px;
		height: 40px;
		border: 3px solid rgba(21, 87, 153, 0.1);
		border-radius: 50%;
		border-top-color: #155799;
		animation: spin 1s ease-in-out infinite;
		margin-bottom: 1rem;
	}
	
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	
	.markdown-content {
		line-height: 1.6;
		color: #333;
	}
	
	.markdown-content h1 {
		font-size: 2rem;
		margin-top: 0;
		margin-bottom: 1.5rem;
		color: #155799;
	}
	
	.markdown-content h2 {
		font-size: 1.5rem;
		margin-top: 2rem;
		margin-bottom: 1rem;
		padding-bottom: 0.5rem;
		border-bottom: 1px solid #eee;
		color: #155799;
	}
	
	.markdown-content h3 {
		font-size: 1.25rem;
		margin-top: 1.5rem;
		margin-bottom: 0.75rem;
		color: #333;
	}
	
	.markdown-content p {
		margin-bottom: 1rem;
	}
	
	.markdown-content a {
		color: #155799;
		text-decoration: none;
	}
	
	.markdown-content a:hover {
		text-decoration: underline;
	}
	
	.markdown-content code {
		background-color: #f5f5f5;
		padding: 0.2rem 0.4rem;
		border-radius: 3px;
		font-family: monospace;
		font-size: 0.9em;
	}
	
	.markdown-content pre {
		background-color: #f5f5f5;
		padding: 1rem;
		border-radius: 5px;
		overflow-x: auto;
		margin-bottom: 1.5rem;
	}
	
	.markdown-content pre code {
		background-color: transparent;
		padding: 0;
		border-radius: 0;
	}
	
	.markdown-content ul, .markdown-content ol {
		margin-bottom: 1rem;
		padding-left: 1.5rem;
	}
	
	.markdown-content table {
		width: 100%;
		border-collapse: collapse;
		margin-bottom: 1.5rem;
	}
	
	.markdown-content table th, .markdown-content table td {
		padding: 0.5rem;
		border: 1px solid #ddd;
		text-align: left;
	}
	
	.markdown-content table th {
		background-color: #f5f5f5;
		font-weight: 600;
	}
	
	.markdown-content blockquote {
		border-left: 4px solid #155799;
		padding-left: 1rem;
		margin-left: 0;
		margin-right: 0;
		color: #666;
	}
	
	@media (max-width: 1024px) {
		.api-page {
			grid-template-columns: 1fr;
		}
		
		.sidebar-container {
			display: none;
		}
		
		.content-container {
			grid-column: 1;
		}
		
		.toc-container {
			display: none;
		}
	}
</style> 