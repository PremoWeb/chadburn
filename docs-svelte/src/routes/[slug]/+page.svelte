<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import MarkdownContent from '$lib/components/MarkdownContent.svelte';
	
	// Get the current slug from the page store
	$: slug = $page.params.slug;
	$: title = formatTitle(slug);
	$: markdownPath = `${base}/markdown/${slug}.md`;
	
	// Format the title from the slug
	function formatTitle(slug: string): string {
		return slug
			.split('-')
			.map(word => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	}
</script>

<svelte:head>
	<title>{title} - Chadburn Documentation</title>
</svelte:head>

<div class="doc-page">
	<MarkdownContent path={markdownPath} />
</div>

<style>
	.doc-page {
		max-width: 900px;
		margin: 0 auto;
	}
</style> 