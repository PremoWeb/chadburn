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
    
    // Function to extract version number from tag name or release name
    function getVersionDisplay(release: GitHubRelease & { formattedDate: string }): string {
        // If name is different from tag_name and doesn't start with the tag, use both
        if (release.name && release.tag_name && 
            release.name !== release.tag_name && 
            !release.name.startsWith(release.tag_name)) {
            return `${release.name} (${release.tag_name})`;
        }
        // Otherwise just use the name or tag_name
        return release.name || release.tag_name;
    }
</script>

<svelte:head>
    <title>Changelog | Chadburn</title>
    <meta name="description" content="Changelog for Chadburn - View all releases and updates" />
</svelte:head>

<div class="container mx-auto px-4 py-8">
    <h1 class="text-4xl font-bold mb-2">Changelog</h1>
    <p class="text-gray-600 dark:text-gray-400 mb-8">A history of changes and updates to Chadburn</p>
    
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
                    <div class="border-b border-gray-200 dark:border-gray-700 pb-4 mb-4">
                        <div class="flex flex-col md:flex-row md:items-center justify-between">
                            <h2 class="text-2xl font-bold text-primary">
                                <a href={release.html_url} target="_blank" rel="noopener noreferrer" 
                                   class="hover:underline">
                                    {getVersionDisplay(release)}
                                </a>
                                {#if release.prerelease}
                                    <span class="ml-2 px-2 py-1 text-xs align-middle bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100 rounded-full">
                                        Pre-release
                                    </span>
                                {/if}
                            </h2>
                            <div class="flex items-center mt-2 md:mt-0">
                                <time datetime={release.published_at || release.created_at} class="text-sm text-gray-600 dark:text-gray-400">
                                    {release.formattedDate}
                                </time>
                                <span class="ml-2 text-xs text-gray-500 dark:text-gray-500">
                                    ({release.relativeTime})
                                </span>
                            </div>
                        </div>
                        
                        {#if release.author}
                            <div class="mt-2 flex items-center">
                                <span class="text-sm text-gray-600 dark:text-gray-400">Released by</span>
                                <a href={release.author.html_url} target="_blank" rel="noopener noreferrer" 
                                   class="flex items-center ml-2 text-sm text-primary hover:underline">
                                    <img src={release.author.avatar_url} alt={release.author.login} 
                                         class="w-5 h-5 rounded-full mr-1" />
                                    {release.author.login}
                                </a>
                            </div>
                        {/if}
                    </div>
                    
                    <div class="prose dark:prose-invert max-w-none">
                        {@html release.processedBody}
                    </div>
                    
                    <div class="mt-6 pt-4 border-t border-gray-200 dark:border-gray-700 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                        <a href={release.html_url} target="_blank" rel="noopener noreferrer" 
                           class="text-primary hover:underline text-sm flex items-center">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                            </svg>
                            View on GitHub
                        </a>
                        
                        {#if release.assets && release.assets.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each release.assets as asset}
                                    <a href={asset.browser_download_url} 
                                       class="text-sm px-3 py-1 bg-gray-100 dark:bg-gray-700 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors flex items-center"
                                       download>
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                                        </svg>
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