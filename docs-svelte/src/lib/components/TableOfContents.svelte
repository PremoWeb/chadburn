<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import type { Heading } from '$lib/utils/markdown';
	
	// Props
	const { headings = [] } = $props<{
		headings?: Heading[];
	}>();
	
	// State
	let activeId = $state('');
	let observer: IntersectionObserver;
	let toc: HTMLElement;
	
	// Reactive values
	let currentPath = $derived(page.url.pathname);
	
	// Define nested heading type
	type NestedHeading = Heading & { children: Heading[] };
	
	// Organize headings into a nested structure
	const processHeadings = (): NestedHeading[] => {
		const result: NestedHeading[] = [];
		const h2Stack: NestedHeading[] = [];
		
		// Process headings in order
		headings.forEach((heading: Heading) => {
			if (heading.level === 2) {
				// Create a new h2 entry with empty children array
				const newH2: NestedHeading = { ...heading, children: [] };
				h2Stack.push(newH2);
				result.push(newH2);
			} else if (heading.level === 3 && h2Stack.length > 0) {
				// Add h3 to the children of the most recent h2
				h2Stack[h2Stack.length - 1].children.push(heading);
			}
			// Ignore h1 and h4+ headings
		});
		
		return result;
	};
	
	let nestedHeadings = $state<NestedHeading[]>([]);
	
	// Update nested headings when raw headings change
	$effect(() => {
		nestedHeadings = processHeadings();
	});
	
	// Set up intersection observer to track which heading is currently in view
	onMount(() => {
		const headingElements = headings.map((h: { id: string }) => document.getElementById(h.id)).filter(Boolean);
		
		if (headingElements.length === 0) return;
		
		// Check if there's a hash in the URL and scroll to that section
		if (window.location.hash) {
			const id = window.location.hash.substring(1);
			const element = document.getElementById(id);
			if (element) {
				setTimeout(() => {
					element.scrollIntoView({ behavior: 'smooth' });
					activeId = id;
				}, 100);
			}
		}
		
		const options = {
			rootMargin: '0px 0px -65% 0px', // Adjusted for content area
			threshold: 0.1
		};
		
		observer = new IntersectionObserver(entries => {
			// Sort entries by their position in the document
			const visibleEntries = entries
				.filter(entry => entry.isIntersecting)
				.sort((a, b) => {
					const aRect = a.boundingClientRect;
					const bRect = b.boundingClientRect;
					return aRect.top - bRect.top;
				});
			
			// Use the topmost visible heading
			if (visibleEntries.length > 0) {
				const topEntry = visibleEntries[0];
				activeId = topEntry.target.id;
			}
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
			// Set active ID immediately for better UX
			activeId = id;
			
			// Calculate offset to account for sticky header if needed
			const offset = 24; // Adjusted for content area
			const elementPosition = el.getBoundingClientRect().top;
			const offsetPosition = elementPosition + window.pageYOffset - offset;
			
			// Scroll to the element with offset
			window.scrollTo({
				top: offsetPosition,
				behavior: 'smooth'
			});
			
			// Update URL hash without triggering a scroll
			history.pushState(null, '', `#${id}`);
		}
	}
	
	// Check if a heading or any of its children is active
	function isHeadingActive(heading: NestedHeading): boolean {
		if (activeId === heading.id) return true;
		return heading.children.some(child => activeId === child.id);
	}
</script>

<style>
	.toc-container {
		position: relative;
		padding-left: 0.75rem;
		width: 100%; /* Ensure it takes full width of its container */
	}
	
	/* Thin left border with gradient fade */
	.toc-container::before {
		content: '';
		position: absolute;
		left: 0;
		top: 0;
		bottom: 0;
		width: 1px;
		background: linear-gradient(to bottom, 
			transparent 0%, 
			#e5e7eb 10%, 
			#e5e7eb 90%, 
			transparent 100%);
	}
	
	.toc-list {
		list-style: none;
		padding: 0;
		margin: 0;
		width: 100%;
	}
	
	.toc-item {
		margin-bottom: 0.375rem;
		font-size: 0.9375rem;
	}
	
	.toc-link {
		display: block;
		padding: 0.125rem 0.25rem 0.125rem 0.5rem;
		color: #4b5563;
		text-decoration: none;
		border-left: 2px solid transparent;
		transition: all 150ms ease-in-out;
		line-height: 1.4;
		word-break: break-word; /* Allow long words to break */
		hyphens: auto; /* Enable hyphenation */
	}
	
	.toc-link:hover {
		color: #2563eb;
		background-color: rgba(243, 244, 246, 0.5);
	}
	
	.toc-link.active {
		color: #2563eb;
		border-left-color: #2563eb;
		background-color: rgba(239, 246, 255, 0.6);
		font-weight: 500;
	}
	
	.toc-sublist {
		list-style: none;
		padding: 0;
		margin: 0.25rem 0 0.375rem 0.5rem;
	}
	
	.toc-subitem {
		margin-bottom: 0.25rem;
		font-size: 0.875rem;
	}
	
	/* Parent highlight when child is active */
	.parent-active > a {
		color: #4338ca;
		font-weight: 500;
	}
</style>

<div class="toc-container" bind:this={toc}>
	{#if nestedHeadings.length > 0}
		<ul class="toc-list">
			{#each nestedHeadings as h2 (h2.id)}
				<li class="toc-item {h2.children.some(child => child.id === activeId) ? 'parent-active' : ''}">
					<!-- H2 Heading -->
					<a 
						href="#{h2.id}" 
						onclick={(e) => handleTocClick(e, h2.id)}
						class="toc-link {activeId === h2.id ? 'active' : ''}"
					>
						{h2.text}
					</a>
					
					<!-- H3 Children -->
					{#if h2.children.length > 0}
						<ul class="toc-sublist">
							{#each h2.children as h3 (h3.id)}
								<li class="toc-subitem">
									<a 
										href="#{h3.id}" 
										onclick={(e) => handleTocClick(e, h3.id)}
										class="toc-link {activeId === h3.id ? 'active' : ''}"
									>
										{h3.text}
									</a>
								</li>
							{/each}
						</ul>
					{/if}
				</li>
			{/each}
		</ul>
	{:else}
		<p class="text-xs text-gray-500 italic">No headings found</p>
	{/if}
</div>
