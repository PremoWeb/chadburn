<script lang="ts">
	import '../../app.css';
	import { page } from '$app/state';
	import { beforeNavigate, afterNavigate } from '$app/navigation';
	import { onMount } from 'svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import ScrollToTop from '$lib/components/ScrollToTop.svelte';
	import WarningBanner from '$lib/components/WarningBanner.svelte';
	
	let { children } = $props();
	
	// Track current route for transition effects
	let currentPath = $state(page.url.pathname);
	let isNavigating = $state(false);
	
	// Handle navigation events
	beforeNavigate(() => {
		isNavigating = true;
	});
	
	afterNavigate(() => {
		isNavigating = false;
		currentPath = page.url.pathname;
	});
	
	onMount(() => {
		currentPath = page.url.pathname;
	});
</script>

<div class="app">
	<div class="content-wrapper">
		<Sidebar />
		<main class="main-content">
			<WarningBanner />
			<div class="page-transition-container">
				{@render children()}
			</div>
		</main>
	</div>
	<ScrollToTop threshold={400} minHeight={800} />
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	.content-wrapper {
		display: flex;
		flex: 1;
		max-width: 1200px;
		margin: 0 auto;
		width: 100%;
		padding: 0 1rem;
		gap: 1.5rem;
		margin-top: 0.75rem;
		padding-bottom: 2rem;
	}

	.main-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		width: 100%;
		box-sizing: border-box;
		padding: 0;
		max-width: calc(100% - 200px); /* Account for sidebar width */
		position: relative;
	}
	
	.page-transition-container {
		width: 100%;
		position: relative;
		min-height: 300px;
	}

	@media (max-width: 768px) {
		.content-wrapper {
			flex-direction: column;
		}
		
		.main-content {
			max-width: 100%;
		}
	}
</style> 