<script lang="ts">
	import { onMount } from 'svelte';
	
	// Props with defaults
	let {
		threshold = 300, // Show button after scrolling this many pixels
		minHeight = 1000, // Minimum page height to show button at all
		scrollBehavior = 'smooth',
		size = '40px',
		color = 'var(--primary-color, #155799)',
		backgroundColor = 'white',
		borderColor = '#eaecef',
		borderRadius = '50%',
		bottom = '30px',
		right = '30px',
		zIndex = 99
	} = $props<{
		threshold?: number;
		minHeight?: number;
		scrollBehavior?: 'auto' | 'smooth';
		size?: string;
		color?: string;
		backgroundColor?: string;
		borderColor?: string;
		borderRadius?: string;
		bottom?: string;
		right?: string;
		zIndex?: number;
	}>();
	
	// State
	let visible = $state(false);
	let pageIsLongEnough = $state(false);
	
	// Function to scroll to top
	function scrollToTop() {
		window.scrollTo({
			top: 0,
			behavior: scrollBehavior
		});
	}
	
	// Function to check if page is long enough to warrant scrolling
	function checkPageHeight() {
		const documentHeight = Math.max(
			document.body.scrollHeight,
			document.body.offsetHeight,
			document.documentElement.clientHeight,
			document.documentElement.scrollHeight,
			document.documentElement.offsetHeight
		);
		
		const viewportHeight = window.innerHeight;
		pageIsLongEnough = documentHeight > viewportHeight + minHeight;
	}
	
	// Function to check scroll position and update visibility
	function checkScrollPosition() {
		// Only show if we've scrolled down AND the page is long enough
		visible = window.scrollY > threshold && pageIsLongEnough;
	}
	
	// Set up scroll listener on mount
	onMount(() => {
		// Check initial page height
		checkPageHeight();
		
		// Check initial position
		checkScrollPosition();
		
		// Add scroll event listener
		window.addEventListener('scroll', checkScrollPosition, { passive: true });
		
		// Add resize event listener to recheck page height
		window.addEventListener('resize', () => {
			checkPageHeight();
			checkScrollPosition();
		}, { passive: true });
		
		// Clean up on component destruction
		return () => {
			window.removeEventListener('scroll', checkScrollPosition);
			window.removeEventListener('resize', checkPageHeight);
		};
	});
</script>

<!-- Button that appears when scrolled down -->
{#if visible}
	<button 
		class="scroll-to-top" 
		aria-label="Scroll to top"
		onclick={scrollToTop}
		style:--size={size}
		style:--color={color}
		style:--background-color={backgroundColor}
		style:--border-color={borderColor}
		style:--border-radius={borderRadius}
		style:--bottom={bottom}
		style:--right={right}
		style:--z-index={zIndex}
	>
		<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<path d="M18 15l-6-6-6 6"/>
		</svg>
	</button>
{/if}

<style>
	.scroll-to-top {
		position: fixed;
		bottom: var(--bottom);
		right: var(--right);
		width: var(--size);
		height: var(--size);
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: var(--background-color);
		color: var(--color);
		border: 1px solid var(--border-color);
		border-radius: var(--border-radius);
		cursor: pointer;
		box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
		z-index: var(--z-index);
		transition: opacity 0.3s, transform 0.3s;
		opacity: 0.8;
		padding: 0;
	}
	
	.scroll-to-top:hover {
		opacity: 1;
		transform: translateY(-3px);
	}
	
	.scroll-to-top svg {
		width: 60%;
		height: 60%;
	}
	
	@media (max-width: 768px) {
		.scroll-to-top {
			bottom: 20px;
			right: 20px;
		}
	}
</style> 