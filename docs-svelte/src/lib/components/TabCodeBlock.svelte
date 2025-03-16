<script lang="ts">
	import TabGroup from './TabGroup.svelte';
	
	// Props
	const { tabs = [], languages = [] } = $props<{ 
		tabs: string[],
		languages: string[] 
	}>();
</script>

<TabGroup tabs={tabs}>
	<div slot="tab-0">
		<slot name="tab-0"></slot>
	</div>
	
	<div slot="tab-1">
		<slot name="tab-1"></slot>
	</div>
	
	<div slot="tab-2">
		<slot name="tab-2"></slot>
	</div>
	
	<div slot="tab-3">
		<slot name="tab-3"></slot>
	</div>
	
	<div slot="tab-4">
		<slot name="tab-4"></slot>
	</div>
</TabGroup>

<style>
	/* Additional styles specific to code tabs can be added here */
</style> 