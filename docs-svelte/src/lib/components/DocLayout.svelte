<script lang="ts">
	import { fade } from 'svelte/transition';
	import DocSidebar from './DocSidebar.svelte';
	import TableOfContents from './TableOfContents.svelte';
	import { processMarkdown, extractHeadings, addIdsToHeadings, type Heading } from '$lib/utils/markdown';
	import { onMount } from 'svelte';
	
	// Props
	const { 
		content = '', 
		title = '',
		navItems = [],
		sidebarTitle = 'Documentation',
		isLoading = false
	} = $props<{
		content: string;
		title?: string;
		navItems: Array<{
			section: string;
			items: Array<{ title: string; path: string }>;
		}>;
		sidebarTitle?: string;
		isLoading?: boolean;
	}>();
	
	// State
	let currentHtmlContent = $state('');
	let headings = $state<Heading[]>([]);
	let contentProcessed = $state(false);
	
	// Process markdown content
	$effect(() => {
		if (content && !isLoading) {
			// Process markdown asynchronously
			processMarkdown(content).then(html => {
				// Add IDs to headings for anchor links
				currentHtmlContent = addIdsToHeadings(html);
				contentProcessed = true;
			});
		}
	});
	
	// Extract headings after content is rendered
	onMount(() => {
		// Extract headings after initial render
		setTimeout(() => {
			headings = extractHeadings();
		}, 100);
	});
	
	// Re-extract headings when content changes
	$effect(() => {
		if (currentHtmlContent && contentProcessed && typeof document !== 'undefined') {
			setTimeout(() => {
				headings = extractHeadings();
			}, 100);
		}
	});
</script>

<svelte:head>
	<title>Chadburn - {title}</title>
</svelte:head>

<div class="doc-layout">
	<div class="sidebar-container">
		<DocSidebar navItems={navItems} title={sidebarTitle} />
	</div>
	
	<div class="content-container">
		{#if isLoading}
			<div class="loading-container" in:fade={{ duration: 200 }}>
				<div class="loading-spinner"></div>
				<p>Loading content...</p>
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
	.doc-layout {
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
		.doc-layout {
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