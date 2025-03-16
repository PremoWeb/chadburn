<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import type { Heading } from '$lib/utils/markdown';
	
	// Props
	const { headings = [], title = 'On this page' } = $props<{
		headings?: Heading[];
		title?: string;
	}>();
	
	// State
	let activeId = $state('');
	let observer: IntersectionObserver;
	let toc: HTMLElement;
	
	// Reactive values
	let currentPath = $derived(page.url.pathname);
	
	// Set up intersection observer to track which heading is currently in view
	onMount(() => {
		const headingElements = headings.map((h: { id: string }) => document.getElementById(h.id)).filter(Boolean);
		
		if (headingElements.length === 0) return;
		
		const options = {
			rootMargin: '0px 0px -80% 0px',
			threshold: 1.0
		};
		
		observer = new IntersectionObserver(entries => {
			entries.forEach(entry => {
				if (entry.isIntersecting) {
					activeId = entry.target.id;
					
					// Scroll the TOC to keep the active item visible
					const activeItem = toc?.querySelector(`a[href="#${activeId}"]`)?.parentElement;
					if (activeItem && toc) {
						const tocRect = toc.getBoundingClientRect();
						const activeRect = activeItem.getBoundingClientRect();
						
						if (activeRect.bottom > tocRect.bottom || activeRect.top < tocRect.top) {
							activeItem.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
						}
					}
				}
			});
		}, options);
		
		headingElements.forEach((el: Element | null) => {
			if (el) observer.observe(el);
		});
		
		return () => {
			if (observer) {
				observer.disconnect();
			}
		};
	});
	
	// Handle click on TOC item
	function handleTocClick(e: MouseEvent, id: string) {
		e.preventDefault();
		const el = document.getElementById(id);
		if (el) {
			activeId = id;
			el.scrollIntoView({ behavior: 'smooth' });
			history.pushState(null, '', `#${id}`);
		}
	}
	
	// Get indentation based on heading level
	function getIndent(level: number): string {
		return `pl-${(level - 2) * 4}`;
	}
</script>

<div class="toc-container" bind:this={toc}>
	<div class="toc-header">
		<h3>{title}</h3>
	</div>
	
	{#if headings.length > 0}
		<ul class="toc-list">
			{#each headings as { id, text, level }}
				{#if level > 1 && level < 4}
					<li class="toc-item {getIndent(level)} {activeId === id ? 'active' : ''}">
						<a href="#{id}" on:click={(e) => handleTocClick(e, id)}>
							{text}
						</a>
					</li>
				{/if}
			{/each}
		</ul>
	{:else}
		<p class="toc-empty">No headings found</p>
	{/if}
</div>

<style>
	.toc-container {
		position: sticky;
		top: 2rem;
		max-height: calc(100vh - 4rem);
		overflow-y: auto;
		padding-left: 1rem;
		scrollbar-width: thin;
		font-size: 0.9rem;
	}
	
	.toc-header {
		margin-bottom: 0.75rem;
	}
	
	.toc-header h3 {
		font-size: 0.9rem;
		font-weight: 600;
		color: #333;
		margin: 0;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}
	
	.toc-list {
		list-style: none;
		padding: 0;
		margin: 0;
	}
	
	.toc-item {
		margin-bottom: 0.5rem;
		line-height: 1.4;
	}
	
	.toc-item a {
		color: #555;
		text-decoration: none;
		display: inline-block;
		transition: all 0.2s ease;
		border-left: 2px solid transparent;
		padding-left: 0.5rem;
	}
	
	.toc-item a:hover {
		color: #155799;
	}
	
	.toc-item.active a {
		color: #155799;
		border-left-color: #155799;
		font-weight: 500;
	}
	
	.toc-empty {
		color: #777;
		font-style: italic;
		margin: 0;
	}
	
	.pl-4 {
		padding-left: 1rem;
	}
	
	.pl-8 {
		padding-left: 2rem;
	}
	
	.pl-12 {
		padding-left: 3rem;
	}
	
	@media (max-width: 768px) {
		.toc-container {
			display: none;
		}
	}
</style>
