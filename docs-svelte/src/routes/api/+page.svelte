<script lang="ts">
	import { page } from '$app/state';
	import DocLayout from '$lib/components/DocLayout.svelte';
	
	// Define the data type
	type PageData = {
		content: string;
	};
	
	// Get data from the page load function
	let { data } = $props<{ data: PageData }>();
	
	// Define API navigation items
	const apiNavItems = [
		{
			section: 'API Reference',
			items: [
				{ title: 'Overview', path: '/api/overview' },
				{ title: 'Authentication', path: '/api/authentication' },
				{ title: 'Rate Limiting', path: '/api/rate-limiting' }
			]
		},
		{
			section: 'Endpoints',
			items: [
				{ title: 'Jobs', path: '/api/endpoints/jobs' },
				{ title: 'Schedules', path: '/api/endpoints/schedules' },
				{ title: 'Metrics', path: '/api/endpoints/metrics' },
				{ title: 'Webhooks', path: '/api/endpoints/webhooks' },
				{ title: 'Tags', path: '/api/endpoints/tags' },
				{ title: 'Users', path: '/api/endpoints/users' },
				{ title: 'Settings', path: '/api/endpoints/settings' }
			]
		},
		{
			section: 'Examples',
			items: [
				{ title: 'Basic Usage', path: '/api/examples/basic' },
				{ title: 'Advanced Usage', path: '/api/examples/advanced' }
			]
		},
		{
			section: 'SDKs & Libraries',
			items: [
				{ title: 'Overview', path: '/api/sdks' },
				{ title: 'Go Client', path: '/api/sdks/go' },
				{ title: 'Python Client', path: '/api/sdks/python' }
			]
		}
	];
</script>

<DocLayout 
	content={data.content}
	title="API Documentation"
	navItems={apiNavItems}
	sidebarTitle="API Documentation"
/> 