<script lang="ts">
	import { page } from '$app/state';
	import DocLayout from '$lib/components/DocLayout.svelte';
	
	// Define the data type
	type PageData = {
		content: string;
	};
	
	// Get data from the page load function
	let { data } = $props<{ data: PageData }>();
	
	// Define documentation navigation items
	const docsNavItems = [
		{
			section: 'Getting Started',
			items: [
				{ title: 'Introduction', path: '/docs/introduction' },
				{ title: 'Installation', path: '/docs/installation' },
				{ title: 'Quick Start', path: '/docs/quick-start' }
			]
		},
		{
			section: 'Core Concepts',
			items: [
				{ title: 'Jobs', path: '/docs/concepts/jobs' },
				{ title: 'Schedules', path: '/docs/concepts/schedules' },
				{ title: 'Execution', path: '/docs/concepts/execution' },
				{ title: 'Metrics', path: '/docs/concepts/metrics' }
			]
		},
		{
			section: 'Guides',
			items: [
				{ title: 'Configuration', path: '/docs/guides/configuration' },
				{ title: 'Security', path: '/docs/guides/security' },
				{ title: 'Monitoring', path: '/docs/guides/monitoring' },
				{ title: 'Scaling', path: '/docs/guides/scaling' }
			]
		},
		{
			section: 'Advanced',
			items: [
				{ title: 'Clustering', path: '/docs/advanced/clustering' },
				{ title: 'High Availability', path: '/docs/advanced/high-availability' },
				{ title: 'Custom Executors', path: '/docs/advanced/custom-executors' }
			]
		}
	];
</script>

<DocLayout 
	content={data.content}
	title="Documentation"
	navItems={docsNavItems}
	sidebarTitle="Documentation"
/> 