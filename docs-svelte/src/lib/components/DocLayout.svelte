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
		isLoading = false,
		debug = false
	} = $props<{
		content: string;
		title?: string;
		navItems: Array<{
			section: string;
			items: Array<{ title: string; path: string }>;
		}>;
		sidebarTitle?: string;
		isLoading?: boolean;
		debug?: boolean;
	}>();
	
	// State
	let currentHtmlContent = $state('');
	let headings = $state<Heading[]>([]);
	let contentProcessed = $state(false);
	let rawContent = $state('');
	let strippedContent = $state('');
	let showDebug = $state(debug);
	let contentAreaTop = $state(0);
	let contentWrapper: HTMLElement;
	
	// Strip YAML frontmatter from markdown content
	function stripFrontmatter(content: string): string {
		// Check if content starts with --- which indicates frontmatter
		if (!content.trim().startsWith('---')) {
			return content;
		}
		
		// Find the second --- which ends the frontmatter
		const parts = content.split('---');
		if (parts.length < 3) {
			return content; // No proper frontmatter found
		}
		
		// Return everything after the second ---
		return parts.slice(2).join('---').trim();
	}
	
	// Process markdown content
	$effect(() => {
		if (content && !isLoading) {
			// Store raw content for debugging
			rawContent = content;
			
			// Strip frontmatter first
			strippedContent = stripFrontmatter(content);
			
			// Process markdown asynchronously
			processMarkdown(strippedContent).then(html => {
				// Add IDs to headings for anchor links
				currentHtmlContent = addIdsToHeadings(html);
				contentProcessed = true;
				
				// Extract headings from the HTML content directly
				// This is more reliable than waiting for DOM to render
				const tempDiv = document.createElement('div');
				tempDiv.innerHTML = currentHtmlContent;
				const headingElements = tempDiv.querySelectorAll('h1, h2, h3, h4');
				
				console.log(`Found ${headingElements.length} headings in the processed HTML content`);
				
				const extractedHeadings: Heading[] = [];
				headingElements.forEach((el, index) => {
					// Ensure the heading has an ID
					const id = el.id || `heading-${index}`;
					if (!el.id) {
						console.warn(`Heading "${el.textContent}" had no ID, assigned: ${id}`);
						el.id = id;
					}
					
					const level = parseInt(el.tagName.substring(1));
					const text = el.textContent || '';
					
					extractedHeadings.push({
						id,
						text,
						level
					});
					
					console.log(`Extracted heading: ${el.tagName} "${text}" with ID "${id}"`);
				});
				
				// Update the currentHtmlContent with any ID changes made during extraction
				currentHtmlContent = tempDiv.innerHTML;
				
				headings = extractedHeadings;
				console.log('Extracted headings:', headings);
			});
		}
	});
	
	// Backup extraction method using DOM
	onMount(() => {
		// Only run if no headings were found in the direct HTML parsing
		setTimeout(() => {
			if (headings.length === 0 && contentProcessed) {
				console.log('No headings found in direct parsing, trying DOM extraction');
				const domHeadings = extractHeadings();
				if (domHeadings.length > 0) {
					console.log('Found headings via DOM extraction:', domHeadings);
					headings = domHeadings;
				} else {
					console.warn('No headings found in DOM extraction either');
				}
			}
		}, 300);
		
		// Get the content area's top position for sticky TOC
		if (contentWrapper) {
			const rect = contentWrapper.getBoundingClientRect();
			contentAreaTop = rect.top + window.scrollY;
			console.log('Content area top position:', contentAreaTop);
		}
	});
	
	// Toggle debug mode
	function toggleDebug() {
		showDebug = !showDebug;
	}
	
	// Debug function to check headings and IDs
	function checkHeadingsAndIds() {
		console.log('Checking headings and IDs in the rendered content...');
		
		// Find all headings in the rendered content
		const contentArea = document.querySelector('.main-content article');
		if (!contentArea) {
			console.error('Content area not found');
			return;
		}
		
		const renderedHeadings = contentArea.querySelectorAll('h1, h2, h3, h4');
		console.log(`Found ${renderedHeadings.length} headings in the rendered content`);
		
		renderedHeadings.forEach((heading, index) => {
			console.log(`Rendered heading ${index + 1}: ${heading.tagName} "${heading.textContent}" (ID: "${heading.id}")`);
			
			// Check if this heading is in our extracted headings
			const matchingHeading = headings.find(h => h.id === heading.id);
			if (matchingHeading) {
				console.log(`  ✓ Found in extracted headings with ID "${matchingHeading.id}"`);
			} else {
				console.error(`  ✗ Not found in extracted headings!`);
				
				// Try to find a heading with similar text
				const similarHeading = headings.find(h => h.text === heading.textContent);
				if (similarHeading) {
					console.log(`  ⚠ Found heading with matching text but different ID: "${similarHeading.id}"`);
				}
			}
			
			// Check if the ID is actually working
			const elementById = document.getElementById(heading.id);
			if (elementById) {
				console.log(`  ✓ Element can be found by ID "${heading.id}"`);
			} else {
				console.error(`  ✗ Element CANNOT be found by ID "${heading.id}"!`);
			}
		});
		
		// Check if all extracted headings are in the rendered content
		console.log('\nChecking if all extracted headings are in the rendered content...');
		headings.forEach((heading, index) => {
			console.log(`Extracted heading ${index + 1}: Level ${heading.level} "${heading.text}" (ID: "${heading.id}")`);
			
			// Try to find this heading in the rendered content
			const elementById = document.getElementById(heading.id);
			if (elementById) {
				console.log(`  ✓ Found in rendered content with ID "${heading.id}"`);
			} else {
				console.error(`  ✗ Not found in rendered content!`);
			}
		});
	}
</script>

<style>
	/* Layout structure */
	.doc-grid {
		display: grid;
		grid-template-columns: 220px minmax(0, 1fr);
		gap: 1.5rem;
	}
	
	.sidebar-area {
		grid-column: 1;
		grid-row: 1;
	}
	
	.content-area {
		grid-column: 2;
		grid-row: 1;
		position: relative;
		padding: 1.25rem;
		background-color: white;
		border-radius: 0.5rem;
		border: 1px solid #e5e7eb;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
	}
	
	/* Content wrapper - this is the positioning context for the TOC */
	.content-wrapper {
		position: relative;
		overflow: visible; /* Allow content to flow around floated elements */
		min-height: 300px; /* Ensure there's enough space for the TOC to be sticky */
		display: flex;
		flex-direction: row-reverse; /* TOC first, then content */
		gap: 1.5rem;
	}
	
	/* Main content styles */
	.main-content {
		flex: 1;
		min-width: 0; /* Allow content to shrink if needed */
	}
	
	/* TOC styles */
	.toc-container {
		width: 240px;
		flex-shrink: 0;
		display: flex;
		flex-direction: column;
	}
	
	.toc-title {
		font-size: 0.8125rem;
		font-weight: 600;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		margin-bottom: 0.75rem;
		background-color: white;
		padding-top: 0.5rem;
		flex-shrink: 0; /* Prevent title from shrinking */
	}
	
	.toc-sticky {
		position: sticky;
		top: 1.5rem;
		max-height: calc(100vh - 6rem);
		overflow-y: auto;
		scrollbar-width: thin;
		scrollbar-color: #cbd5e1 transparent;
		display: flex;
		flex-direction: column;
	}
	
	.toc-content {
		overflow-y: auto;
		max-height: calc(100vh - 10rem);
		scrollbar-width: thin;
		scrollbar-color: #cbd5e1 transparent;
	}
	
	.toc-content::-webkit-scrollbar,
	.toc-sticky::-webkit-scrollbar {
		width: 4px;
	}
	
	.toc-content::-webkit-scrollbar-track,
	.toc-sticky::-webkit-scrollbar-track {
		background: transparent;
	}
	
	.toc-content::-webkit-scrollbar-thumb,
	.toc-sticky::-webkit-scrollbar-thumb {
		background-color: #cbd5e1;
		border-radius: 4px;
	}
	
	/* Responsive adjustments */
	@media (max-width: 1279px) {
		.toc-container {
			width: 200px; /* Slightly smaller on smaller screens */
		}
	}
	
	@media (max-width: 1023px) {
		.doc-grid {
			display: block; /* Switch to block layout on mobile */
		}
		
		.sidebar-area {
			margin-bottom: 1rem;
			width: 100%;
		}
		
		.content-area {
			width: 100%;
		}
		
		.content-wrapper {
			display: block; /* Switch to block layout on mobile */
		}
		
		.toc-container {
			width: 100%;
			margin: 1rem 0;
			padding: 0.75rem;
			background-color: #f9fafb;
			border-radius: 0.375rem;
		}
		
		.toc-sticky {
			position: relative;
			top: 0;
			max-height: none;
		}
	}
	
	/* Prose content styling */
	:global(.prose) {
		max-width: 65ch !important;
		margin-left: 0;
		margin-right: 0;
	}
	
	/* Scrollbar styling */
	:global(.scrollbar-thin::-webkit-scrollbar) {
		width: 4px;
	}
	
	:global(.scrollbar-thin::-webkit-scrollbar-track) {
		background: #f1f5f9;
		border-radius: 2px;
	}
	
	:global(.scrollbar-thin::-webkit-scrollbar-thumb) {
		background: #cbd5e1;
		border-radius: 2px;
	}
	
	:global(.scrollbar-thin::-webkit-scrollbar-thumb:hover) {
		background: #94a3b8;
	}
	
	/* Debug styles */
	.debug-panel {
		position: fixed;
		top: 0;
		right: 0;
		width: 400px;
		max-width: 90vw;
		height: 100vh;
		background: #fff;
		border-left: 1px solid #ddd;
		box-shadow: -2px 0 5px rgba(0, 0, 0, 0.1);
		z-index: 1000;
		overflow-y: auto;
		padding: 1rem;
	}
	
	.debug-panel h3 {
		margin-top: 0;
		border-bottom: 1px solid #ddd;
		padding-bottom: 0.5rem;
	}
	
	.debug-panel pre {
		background: #f5f5f5;
		padding: 0.5rem;
		border-radius: 4px;
		overflow-x: auto;
		font-size: 12px;
		max-height: 200px;
		overflow-y: auto;
	}
	
	.debug-actions {
		margin-bottom: 1rem;
	}
	
	.debug-actions button {
		background: #0066cc;
		color: white;
		border: none;
		padding: 0.5rem 1rem;
		border-radius: 4px;
		cursor: pointer;
	}
	
	.debug-toggle {
		position: fixed;
		bottom: 1rem;
		right: 1rem;
		background: #0066cc;
		color: white;
		border: none;
		padding: 0.5rem 1rem;
		border-radius: 4px;
		cursor: pointer;
		z-index: 999;
	}
	
	.debug-content ul {
		padding-left: 1rem;
	}
	
	.debug-content li {
		margin-bottom: 0.25rem;
		font-size: 12px;
	}
</style>

<svelte:head>
	<title>Chadburn - {title}</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
	<div class="doc-grid">
		<!-- Sidebar -->
		<div class="sidebar-area">
			<div class="sticky top-4">
				<DocSidebar navItems={navItems} title={sidebarTitle} />
			</div>
		</div>
		
		<!-- Main Content Area -->
		<div class="content-area">
			{#if showDebug}
				<div class="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md">
					<div class="flex justify-between items-center mb-2">
						<h3 class="text-sm font-bold text-yellow-800">Debug Information</h3>
						<button 
							class="text-xs bg-yellow-200 hover:bg-yellow-300 text-yellow-800 px-2 py-1 rounded"
							onclick={toggleDebug}
						>
							Hide Debug
						</button>
					</div>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
						<div>
							<h4 class="text-xs font-bold text-yellow-800 mb-1">Raw Content</h4>
							<pre class="text-xs overflow-auto max-h-40 p-2 bg-yellow-100 rounded">{rawContent}</pre>
						</div>
						<div>
							<h4 class="text-xs font-bold text-yellow-800 mb-1">Stripped Content</h4>
							<pre class="text-xs overflow-auto max-h-40 p-2 bg-yellow-100 rounded">{strippedContent}</pre>
						</div>
					</div>
					<div class="mt-3">
						<h4 class="text-xs font-bold text-yellow-800 mb-1">Processed HTML</h4>
						<pre class="text-xs overflow-auto max-h-40 p-2 bg-yellow-100 rounded">{currentHtmlContent}</pre>
					</div>
					<div class="mt-3">
						<h4 class="text-xs font-bold text-yellow-800 mb-1">Headings Found ({headings.length})</h4>
						<pre class="text-xs overflow-auto max-h-40 p-2 bg-yellow-100 rounded">{JSON.stringify(headings, null, 2)}</pre>
					</div>
					<div class="mt-3">
						<h4 class="text-xs font-bold text-yellow-800 mb-1">Content Area Top Position</h4>
						<pre class="text-xs overflow-auto max-h-40 p-2 bg-yellow-100 rounded">{contentAreaTop}px</pre>
					</div>
				</div>
			{/if}
			
			<div class="content-wrapper" bind:this={contentWrapper}>
				<!-- Table of Contents (Desktop) -->
				{#if headings.length > 0}
					<div class="toc-container hidden lg:block">
						<div class="toc-sticky">
							<h3 class="toc-title">On this page</h3>
							<div class="toc-content">
								<TableOfContents headings={headings} />
							</div>
						</div>
					</div>
				{:else if contentProcessed && !isLoading}
					<div class="toc-container hidden lg:block">
						<div class="toc-sticky">
							<h3 class="toc-title">On this page</h3>
							<div class="toc-content">
								<div class="p-3 bg-blue-50 border border-blue-100 rounded text-sm text-blue-800">
									<p class="font-medium mb-1">No headings found</p>
									<p class="text-xs">This document doesn't contain any headings that can be used for navigation.</p>
								</div>
							</div>
						</div>
					</div>
				{/if}
				
				<!-- Main Content -->
				<div class="main-content">
					{#if isLoading}
						<div class="flex flex-col items-center justify-center min-h-[300px] text-gray-500" in:fade={{ duration: 200 }}>
							<div class="w-10 h-10 border-3 border-blue-100 rounded-full border-t-blue-600 animate-spin mb-4"></div>
							<p>Loading content...</p>
						</div>
					{:else}
						<article 
							class="prose prose-blue"
							in:fade={{ duration: 200, delay: 100 }}
						>
							{@html currentHtmlContent}
						</article>
						
						<!-- Mobile TOC (below content) -->
						{#if headings.length > 0}
							<div class="lg:hidden mt-4 pt-4 border-t border-gray-200">
								<h3 class="text-sm font-semibold text-gray-500 uppercase mb-2">On this page</h3>
								<div class="mobile-toc-content">
									<TableOfContents headings={headings} />
								</div>
							</div>
						{/if}
					{/if}
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Debug panel -->
{#if showDebug}
	<div class="debug-panel">
		<h3>Debug Info</h3>
		<div class="debug-actions">
			<button onclick={checkHeadingsAndIds}>Check Headings & IDs</button>
		</div>
		<div class="debug-content">
			<h4>Raw Content</h4>
			<pre>{rawContent}</pre>
			
			<h4>Stripped Content</h4>
			<pre>{strippedContent}</pre>
			
			<h4>Headings ({headings.length})</h4>
			<ul>
				{#each headings as heading}
					<li>
						{heading.level}: {heading.text} (ID: {heading.id})
					</li>
				{/each}
			</ul>
		</div>
	</div>
{/if}

<!-- Debug toggle button -->
{#if debug}
	<button class="debug-toggle" onclick={toggleDebug}>
		{showDebug ? 'Hide Debug' : 'Show Debug'}
	</button>
{/if} 