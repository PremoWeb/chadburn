<script lang="ts">
    import type { GitHubRelease } from '$lib/types';
    import { onMount } from 'svelte';
    import { marked } from 'marked';
    import { formatDate, formatRelativeTime } from '$lib/utils/date';
    
    export let data: { releases: GitHubRelease[] };
    
    let isLoading = true;
    let error: string | null = null;
    
    $: releases = data.releases;
    $: processedReleases = releases.map(release => ({
        ...release,
        formattedDate: formatDate(new Date(release.published_at || release.created_at)),
        relativeTime: formatRelativeTime(new Date(release.published_at || release.created_at)),
        processedBody: marked.parse(release.body || 'No release notes provided.')
    }));
    
    onMount(() => {
        isLoading = false;
        if (releases.length === 0) {
            error = "Unable to load releases. Please try again later.";
        }
    });
</script>

<svelte:head>
    <title>Changelog | Chadburn</title>
    <meta name="description" content="Changelog for Chadburn - View all releases and updates" />
</svelte:head>

<div class="container mx-auto px-4 py-8">
    <h1 class="text-4xl font-bold mb-8">Changelog</h1>
    
    {#if isLoading}
        <div class="bg-gray-100 dark:bg-gray-800 p-6 rounded-lg">
            <p class="text-center text-gray-600 dark:text-gray-400">Loading releases...</p>
        </div>
    {:else if error}
        <div class="bg-red-50 dark:bg-red-900/20 p-6 rounded-lg border border-red-200 dark:border-red-800">
            <p class="text-center text-red-600 dark:text-red-400">{error}</p>
        </div>
    {:else if processedReleases.length === 0}
        <div class="bg-gray-100 dark:bg-gray-800 p-6 rounded-lg">
            <p class="text-center text-gray-600 dark:text-gray-400">No releases found.</p>
        </div>
    {:else}
        <div class="space-y-8">
            {#each processedReleases as release}
                <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
                    <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                        <h2 class="text-2xl font-bold text-primary">
                            <a href={release.html_url} target="_blank" rel="noopener noreferrer" 
                               class="hover:underline">
                                {release.name || release.tag_name}
                            </a>
                        </h2>
                        <div class="flex flex-col md:flex-row md:items-center mt-2 md:mt-0">
                            <span class="text-sm text-gray-600 dark:text-gray-400">
                                {release.formattedDate}
                            </span>
                            <span class="md:ml-2 text-xs text-gray-500 dark:text-gray-500">
                                ({release.relativeTime})
                            </span>
                            {#if release.prerelease}
                                <span class="mt-1 md:mt-0 md:ml-2 px-2 py-1 text-xs bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100 rounded-full">
                                    Pre-release
                                </span>
                            {/if}
                        </div>
                    </div>
                    
                    <div class="prose dark:prose-invert max-w-none">
                        {@html release.processedBody}
                    </div>
                    
                    <div class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700 flex justify-between items-center">
                        <a href={release.html_url} target="_blank" rel="noopener noreferrer" 
                           class="text-primary hover:underline text-sm">
                            View on GitHub
                        </a>
                        
                        {#if release.assets && release.assets.length > 0}
                            <div class="flex space-x-2">
                                {#each release.assets as asset}
                                    <a href={asset.browser_download_url} 
                                       class="text-sm px-3 py-1 bg-gray-100 dark:bg-gray-700 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                                       download>
                                        {asset.name}
                                    </a>
                                {/each}
                            </div>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div> 