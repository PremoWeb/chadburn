<script lang="ts">
	import { onMount } from 'svelte';
	
	// Props
	const { tabs = [] } = $props<{ tabs: string[] }>();
	
	// State
	let activeTab = $state(0);
	let tabsElement: HTMLElement;
	let mounted = $state(false);
	
	// Set the active tab
	function setActiveTab(index: number) {
		activeTab = index;
		// Save the active tab in localStorage if available
		if (typeof window !== 'undefined' && window.localStorage) {
			try {
				const tabGroupId = tabsElement.dataset.id || '';
				if (tabGroupId) {
					localStorage.setItem(`tabgroup-${tabGroupId}`, index.toString());
				}
			} catch (e) {
				// Silently fail if localStorage is not available
			}
		}
	}
	
	// Initialize from localStorage if available
	onMount(() => {
		if (typeof window !== 'undefined' && window.localStorage) {
			try {
				const tabGroupId = tabsElement.dataset.id || '';
				if (tabGroupId) {
					const savedTab = localStorage.getItem(`tabgroup-${tabGroupId}`);
					if (savedTab !== null) {
						activeTab = parseInt(savedTab, 10);
					}
				}
			} catch (e) {
				// Silently fail if localStorage is not available
			}
		}
		mounted = true;
	});
	
	// Generate a unique ID for the tab group
	const tabGroupId = $derived(tabs.join('-').toLowerCase().replace(/\s+/g, '-'));
</script>

<div class="tab-group" bind:this={tabsElement} data-id={tabGroupId}>
	<div class="tab-buttons">
		{#each tabs as tab, i}
			<button 
				class="tab-button" 
				class:active={activeTab === i}
				on:click={() => setActiveTab(i)}
				aria-selected={activeTab === i}
				role="tab"
				id={`tab-${i}`}
				aria-controls={`panel-${i}`}
			>
				{tab}
			</button>
		{/each}
	</div>
	
	<div class="tab-content">
		{#if mounted}
			{#if activeTab === 0}
				<slot name="tab-0"></slot>
			{:else if activeTab === 1}
				<slot name="tab-1"></slot>
			{:else if activeTab === 2}
				<slot name="tab-2"></slot>
			{:else if activeTab === 3}
				<slot name="tab-3"></slot>
			{:else if activeTab === 4}
				<slot name="tab-4"></slot>
			{/if}
		{/if}
	</div>
</div>

<style>
	.tab-group {
		margin: 1.5rem 0;
		border-radius: 6px;
		overflow: hidden;
		border: 1px solid #e1e4e8;
		background-color: #fff;
	}
	
	.tab-buttons {
		display: flex;
		border-bottom: 1px solid #e1e4e8;
		background-color: #f6f8fa;
		overflow-x: auto;
		scrollbar-width: thin;
	}
	
	.tab-button {
		padding: 0.75rem 1.25rem;
		border: none;
		background: transparent;
		font-size: 0.9rem;
		font-weight: 500;
		color: #6a737d;
		cursor: pointer;
		white-space: nowrap;
		border-bottom: 2px solid transparent;
		transition: all 0.2s ease;
	}
	
	.tab-button:hover {
		color: var(--primary-color, #0366d6);
		background-color: rgba(0, 0, 0, 0.05);
	}
	
	.tab-button.active {
		color: var(--primary-color, #0366d6);
		border-bottom-color: var(--primary-color, #0366d6);
		background-color: #fff;
	}
	
	.tab-content {
		padding: 0;
	}
	
	.tab-content :global(.code-block-wrapper) {
		margin: 0;
		border: none;
		border-radius: 0;
		box-shadow: none;
	}
	
	.tab-content :global(pre) {
		margin: 0;
		border-radius: 0;
	}
	
	@media (max-width: 768px) {
		.tab-button {
			padding: 0.5rem 1rem;
			font-size: 0.85rem;
		}
	}
</style> 