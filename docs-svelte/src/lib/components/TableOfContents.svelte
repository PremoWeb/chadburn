<script lang="ts">
	import { onMount, createEventDispatcher } from 'svelte';
	import { fly, fade, crossfade } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';

	// Create event dispatcher
	const dispatch = createEventDispatcher();

	// Props
	let {
		contentSelector = '.doc-content',
		headingSelector = 'h2, h3, h4',
		title = 'On this page',
		maxDepth = 6
	} = $props<{
		contentSelector?: string;
		headingSelector?: string;
		title?: string;
		maxDepth?: number;
	}>();

	// State
	let headings = $state<{ id: string; text: string; level: number }[]>([]);
	let activeId = $state<string | null>(null);
	let content = $state<HTMLElement | null>(null);
	let observer: MutationObserver | null = $state(null);
	let intersectionObserver: IntersectionObserver | null = $state(null);
	let contentChecked = $state(false);
	let isVisible = $state(false);

	// Watch for changes in headings and dispatch event
	$effect(() => {
		dispatch('contentStatus', { hasContent: headings.length > 0 });
	});

	// Generate a unique ID for a heading if it doesn't have one
	function generateId(text: string): string {
		return text
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/(^-|-$)/g, '');
	}

	// Set visibility after a short delay to ensure DOM is updated
	function showContent() {
		setTimeout(() => {
			isVisible = true;
		}, 50);
	}

	// Extract headings from the content
	function extractHeadings() {
		if (!content) {
			content = document.querySelector(contentSelector);
			if (!content) {
				dispatch('contentStatus', { hasContent: false });
				return;
			}
		}

		const headingElements = content.querySelectorAll(headingSelector);
		
		// Reset visibility to trigger animations
		isVisible = false;
		
		headings = Array.from(headingElements).map((heading) => {
			const level = parseInt(heading.tagName.substring(1));
			if (!heading.id) {
				heading.id = generateId(heading.textContent || '');
			}
			return {
				id: heading.id,
				text: heading.textContent || '',
				level
			};
		});
		
		// Dispatch event with content status
		dispatch('contentStatus', { hasContent: headings.length > 0 });
		
		// Show content after a short delay
		showContent();
	}

	// Debounced scroll handler
	let lastScrollTime = 0;
	const scrollDebounceMs = 100;

	function updateActiveHeading() {
		const now = Date.now();
		if (now - lastScrollTime < scrollDebounceMs) return;
		lastScrollTime = now;

		if (!headings.length) return;

		const scrollPosition = window.scrollY + 100;

		for (let i = headings.length - 1; i >= 0; i--) {
			const heading = headings[i];
			const element = document.getElementById(heading.id);

			if (element && heading.level <= maxDepth && element.offsetTop <= scrollPosition) {
				activeId = heading.id;
				return;
			}
		}

		activeId = headings[0]?.id || null;
	}

	// Scroll to a heading when clicked
	function scrollToHeading(id: string) {
		const element = document.getElementById(id);
		if (element) {
			window.scrollTo({
				top: element.offsetTop - 20,
				behavior: 'smooth'
			});
		}
	}

	// Set up content observation
	function setupContentObserver() {
		if (!content) {
			content = document.querySelector(contentSelector);
			if (!content) {
				return;
			}
		}

		extractHeadings();
		updateActiveHeading();

		// Use MutationObserver instead of polling
		if (observer) {
			observer.disconnect();
		}
		
		observer = new MutationObserver(() => {
			extractHeadings();
			updateActiveHeading();
			setupIntersectionObserver();
		});

		observer.observe(content, {
			childList: true,
			subtree: true,
			characterData: true
		});

		return () => {
			if (observer) {
				observer.disconnect();
				observer = null;
			}
		};
	}
	
	// Set up intersection observer for headings
	function setupIntersectionObserver() {
		if (intersectionObserver) {
			intersectionObserver.disconnect();
		}
		
		intersectionObserver = new IntersectionObserver(
			(entries) => {
				entries.forEach((entry) => {
					if (entry.isIntersecting) {
						const id = entry.target.id;
						if (headings.some((h) => h.id === id)) {
							activeId = id;
						}
					}
				});
			},
			{
				rootMargin: '-100px 0px -50% 0px',
				threshold: 0
			}
		);

		// Observe heading elements
		headings.forEach((heading) => {
			const element = document.getElementById(heading.id);
			if (element && intersectionObserver) {
				intersectionObserver.observe(element);
			}
		});
	}

	// Check for content periodically until found
	function pollForContent() {
		const checkInterval = setInterval(() => {
			if (!content) {
				content = document.querySelector(contentSelector);
				if (content) {
					clearInterval(checkInterval);
					extractHeadings();
					updateActiveHeading();
					setupContentObserver();
					setupIntersectionObserver();
				}
			} else {
				clearInterval(checkInterval);
			}
		}, 200);
		
		// Clear interval after 5 seconds to prevent infinite polling
		setTimeout(() => clearInterval(checkInterval), 5000);
	}

	onMount(() => {
		// Initial content setup
		content = document.querySelector(contentSelector);

		if (content) {
			extractHeadings();
			updateActiveHeading();
			setupContentObserver();
			setupIntersectionObserver();
			
			// Show content after a short delay
			showContent();
		} else {
			pollForContent();
		}

		// Fallback scroll listener (debounced)
		window.addEventListener('scroll', updateActiveHeading, { passive: true });

		return () => {
			if (observer) {
				observer.disconnect();
				observer = null;
			}
			if (intersectionObserver) {
				intersectionObserver.disconnect();
				intersectionObserver = null;
			}
			window.removeEventListener('scroll', updateActiveHeading);
		};
	});
	
	// Use $effect to check for content changes after updates
	$effect(() => {
		// Skip if we've already checked once during this render cycle
		if (!contentChecked && !content) {
			contentChecked = true;
			content = document.querySelector(contentSelector);
			if (content) {
				extractHeadings();
				updateActiveHeading();
				setupContentObserver();
				setupIntersectionObserver();
				
				// Show content after a short delay
				showContent();
			}
		}
	});
	
	// Calculate animation delay based on index and level
	function getAnimationDelay(index: number, level: number): number {
		// Base delay of 30ms per item
		const baseDelay = 30;
		// Add extra delay for deeper levels
		const levelDelay = (level - 2) * 10;
		return baseDelay * index + levelDelay;
	}
	
	// Custom transition that combines fade and fly
	function fadeAndFly(node: Element, options: { 
		delay?: number; 
		duration?: number; 
		y?: number;
		easing?: (t: number) => number;
	}) {
		const {
			delay = 0,
			duration = 300,
			y = 20,
			easing = cubicOut
		} = options;
		
		return {
			delay,
			duration,
			easing,
			css: (t: number, u: number) => `
				transform: translateY(${u * y}px);
				opacity: ${t};
			`
		};
	}
</script>

<!-- Template with animations -->
<div class="table-of-contents" in:fade={{ duration: 300, delay: 100 }}>
	{#if headings.length > 0}
		<div class="toc-title" 
			in:fadeAndFly={{ 
				y: 10, 
				duration: 300, 
				delay: 150,
				easing: cubicOut 
			}}
		>
			{title}
		</div>
		<nav>
			<ul class="toc-list">
				{#each headings as heading, i}
					{#if heading.level <= maxDepth}
						{#if isVisible}
							<li 
								class="level-{heading.level}" 
								class:active={heading.id === activeId}
							>
								<div 
									class="toc-item-wrapper"
									in:fadeAndFly={{ 
										y: 20, 
										duration: 400, 
										delay: getAnimationDelay(i, heading.level) + 200,
										easing: cubicOut 
									}}
								>
									<a
										href="#{heading.id}"
										onclick={(e) => {
											e.preventDefault();
											scrollToHeading(heading.id);
										}}
									>
										{heading.text}
									</a>
								</div>
							</li>
						{/if}
					{/if}
				{/each}
			</ul>
		</nav>
	{:else if isVisible}
		<div class="empty-toc" in:fade={{ duration: 200, delay: 150 }}>
			<p>No headings found</p>
		</div>
	{/if}
</div>

<style>
	.table-of-contents {
		position: relative;
		max-height: none !important;
		overflow: visible !important;
		overflow-x: visible !important;
		overflow-y: visible !important;
		width: 100%;
		padding: 0.5rem 0.75rem;
		/* padding-top: 0.25rem; */
		border: 1px solid #eee;
		background-color: white;
		border-radius: 4px;
		box-sizing: border-box;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
		font-size: 0.9rem;
		will-change: transform, opacity;
	}

	nav {
		overflow: visible !important;
		max-height: none !important;
	}
	
	.toc-list {
		overflow: visible !important;
		max-height: none !important;
	}

	.toc-title {
		font-weight: 600;
		margin-bottom: 8px;
		margin-top: 0.25rem;
		color: #495057;
		font-size: 1rem;
		padding-bottom: 6px;
		border-bottom: 1px solid #eee;
	}

	nav ul {
		list-style: none;
		padding: 0;
		margin: 0;
	}

	nav li {
		margin: 6px 0;
		line-height: 1.4;
	}
	
	.toc-item-wrapper {
		will-change: transform, opacity;
	}

	nav a {
		color: #495057;
		text-decoration: none;
		display: block;
		padding: 3px 0;
		border-left: 2px solid transparent;
		padding-left: 10px;
		margin-left: -12px;
		transition: all 0.2s ease;
		font-size: 0.95rem;
	}

	nav a:hover {
		color: var(--primary-color, #0366d6);
	}

	nav li.active > a {
		color: var(--primary-color, #0366d6);
		border-left-color: var(--primary-color, #0366d6);
		font-weight: 500;
	}

	.empty-toc {
		color: #6c757d;
		font-style: italic;
		font-size: 0.9rem;
		padding: 0.5rem 0;
	}

	/* Indentation for different heading levels */
	.level-2 {
		margin-left: 0;
	}

	.level-3 {
		margin-left: 12px;
	}

	.level-4 {
		margin-left: 24px;
	}

	.level-5 {
		margin-left: 36px;
	}

	.level-6 {
		margin-left: 48px;
	}

	@media (max-width: 1024px) {
		.table-of-contents {
			position: relative;
			top: 0;
			max-height: none;
			border-radius: 4px;
			margin-bottom: 1rem;
		}
	}
</style>
