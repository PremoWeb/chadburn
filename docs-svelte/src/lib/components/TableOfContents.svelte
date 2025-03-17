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
	
	// Offset for scrolling (to account for sticky headers)
	const offset = 80; // Adjust this value based on your layout
	
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
		// Reset the active ID when headings change (page navigation)
		activeId = '';
		
		// Process the new headings
		nestedHeadings = processHeadings();
		
		// If we're in the browser, set up the intersection observer for the new headings
		if (typeof window !== 'undefined') {
			// Wait a bit for the DOM to be fully rendered with the new content
			setTimeout(() => {
				setupIntersectionObserver();
			}, 200);
		}
	});
	
	// Set up intersection observer and page navigation listeners
	onMount(() => {
		// Wait a bit for the DOM to be fully rendered
		setTimeout(() => {
			setupIntersectionObserver();
		}, 200);
		
		// Track the current path to detect navigation
		let previousPath = page.url.pathname;
		
		// Create an interval to check for page navigation
		const navigationCheckInterval = setInterval(() => {
			const currentPath = page.url.pathname;
			if (currentPath !== previousPath) {
				console.log(`Page navigation detected: ${previousPath} -> ${currentPath}`);
				previousPath = currentPath;
				
				// Reset active ID on page navigation
				activeId = '';
				
				// Wait for the new page content to render
				setTimeout(() => {
					setupIntersectionObserver();
				}, 200);
			}
		}, 100);
		
		return () => {
			// Clean up
			if (observer) {
				observer.disconnect();
			}
			clearInterval(navigationCheckInterval);
		};
	});
	
	// Setup the intersection observer
	function setupIntersectionObserver() {
		console.log('Setting up intersection observer for TOC...');
		
		// First, disconnect any existing observer
		if (observer) {
			console.log('Disconnecting existing observer');
			observer.disconnect();
			observer = null as unknown as IntersectionObserver;
		}
		
		// If there are no headings, don't bother setting up an observer
		if (headings.length === 0) {
			console.warn('No headings to observe');
			return;
		}
		
		// Get all heading elements from the DOM
		const headingElements: Element[] = [];
		
		// Find all heading elements in the document
		headings.forEach((heading: Heading) => {
			const element = document.getElementById(heading.id);
			if (element) {
				headingElements.push(element);
			} else {
				console.warn(`Heading element with ID "${heading.id}" not found in the DOM`);
			}
		});
		
		if (headingElements.length === 0) {
			console.warn('No heading elements found in the DOM');
			return;
		}
		
		console.log(`Found ${headingElements.length} heading elements in the DOM`);
		
		// Check if there's a hash in the URL and scroll to that section
		if (window.location.hash) {
			const id = window.location.hash.substring(1);
			const element = document.getElementById(id);
			if (element) {
				console.log(`Found element for hash: #${id}`);
				setTimeout(() => {
					const rect = element.getBoundingClientRect();
					const offsetPosition = rect.top + window.scrollY - offset;
					
					window.scrollTo({
						top: offsetPosition,
						behavior: 'smooth'
					});
					
					activeId = id;
				}, 100);
			} else {
				console.warn(`Element with ID "${id}" from URL hash not found`);
			}
		}
		
		// Configure the intersection observer
		const options = {
			// This rootMargin creates a zone in the viewport where headings become active
			// -80px from top (accounts for fixed headers)
			// 0px from right and left
			// -60% from bottom (so headings become active before they reach the middle of the screen)
			rootMargin: '-80px 0px -60% 0px',
			// Multiple thresholds for better accuracy
			threshold: [0, 0.1, 0.25, 0.5, 0.75, 1]
		};
		
		observer = new IntersectionObserver((entries) => {
			// Get all currently visible headings
			const visibleHeadings = entries
				.filter(entry => entry.isIntersecting)
				.map(entry => ({
					id: entry.target.id,
					ratio: entry.intersectionRatio,
					y: entry.boundingClientRect.y
				}))
				.sort((a, b) => {
					// First sort by intersection ratio (higher is better)
					if (b.ratio !== a.ratio) {
						return b.ratio - a.ratio;
					}
					// If ratios are the same, sort by position (higher in the page is better)
					return a.y - b.y;
				});
			
			// If we have visible headings, update the active ID
			if (visibleHeadings.length > 0) {
				// Use the heading with the highest intersection ratio or the topmost one if tied
				const newActiveId = visibleHeadings[0].id;
				
				// Only update if the active ID has changed
				if (newActiveId !== activeId) {
					console.log(`Setting active heading to: ${newActiveId} (ratio: ${visibleHeadings[0].ratio.toFixed(2)})`);
					activeId = newActiveId;
				}
			}
		}, options);
		
		// Observe all heading elements
		headingElements.forEach(element => {
			observer.observe(element);
		});
		
		console.log('Intersection observer setup complete');
	}
	
	// Handle click on TOC item
	function handleTocClick(e: Event, id: string) {
		e.preventDefault();
		console.log(`TOC click handler called for ID: "${id}"`);
		
		// Find the element with the given ID
		const element = document.getElementById(id);
		
		if (element) {
			console.log(`Found element with ID "${id}"`);
			
			// Update active ID
			activeId = id;
			
			// Calculate position with offset
			const rect = element.getBoundingClientRect();
			const elementPosition = rect.top;
			const offsetPosition = elementPosition + window.scrollY - offset;
			
			// Update URL hash without triggering scroll
			const currentPath = window.location.pathname + window.location.search;
			history.pushState(null, '', `${currentPath}#${id}`);
			
			// Scroll to the element
			console.log(`Scrolling to position: ${offsetPosition}px`);
			window.scrollTo({
				top: offsetPosition,
				behavior: 'smooth'
			});
		} else {
			console.error(`Element with ID "${id}" not found!`);
			
			// Debug: List all elements with IDs
			console.log('All elements with IDs:');
			const allElementsWithIds = document.querySelectorAll('[id]');
			allElementsWithIds.forEach(el => {
				console.log(`- Element: ${el.tagName}, ID: "${el.id}", Text: "${el.textContent?.trim()}"`);
			});
			
			// Debug: List all headings
			console.log('All headings:');
			const allHeadings = document.querySelectorAll('h1, h2, h3, h4');
			allHeadings.forEach(el => {
				console.log(`- Heading: ${el.tagName}, ID: "${el.id}", Text: "${el.textContent?.trim()}"`);
			});
		}
	}
	
	// Check if a heading or any of its children is active
	function isHeadingActive(heading: NestedHeading): boolean {
		// Check if this exact heading is active
		if (activeId === heading.id) return true;
		
		// Check if any of its children are active
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
						data-heading-id={h2.id}
						data-toc-link="true"
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
										data-heading-id={h3.id}
										data-toc-link="true"
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
