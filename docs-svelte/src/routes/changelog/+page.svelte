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
    
    // Function to get a color for the version badge based on index
    function getVersionColor(index: number): string {
        const colors = [
            'bg-blue-600',
            'bg-teal-600',
            'bg-indigo-600',
            'bg-purple-600',
            'bg-cyan-600'
        ];
        return colors[index % colors.length];
    }
</script>

<svelte:head>
    <title>Changelog | Chadburn</title>
    <meta name="description" content="Changelog for Chadburn - View all releases and updates" />
</svelte:head>

<div class="max-w-4xl mx-auto px-4 py-10">
    <div class="mb-10">
        <h1 class="text-3xl font-bold mb-3 text-center text-gray-900 dark:text-white">Changelog</h1>
        <div class="w-14 h-0.5 bg-blue-600 mx-auto mb-5"></div>
        <p class="text-gray-600 dark:text-gray-400 max-w-2xl mx-auto text-center">
            Track Chadburn's evolution through our release history
        </p>
    </div>
    
    {#if isLoading}
        <div class="flex justify-center items-center py-12">
            <div class="w-6 h-6 border-2 border-blue-200 border-t-blue-600 rounded-full animate-spin"></div>
        </div>
    {:else if error}
        <div class="bg-red-50 dark:bg-red-900/20 p-5 rounded-md border border-red-200 dark:border-red-800 text-center">
            <svg class="w-5 h-5 mx-auto text-red-500 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
            </svg>
            <p class="text-red-600 dark:text-red-400 text-sm font-medium">{error}</p>
        </div>
    {:else if processedReleases.length === 0}
        <div class="bg-gray-50 dark:bg-gray-800 p-6 rounded-md text-center">
            <svg class="w-5 h-5 mx-auto text-gray-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            <p class="text-gray-600 dark:text-gray-400 text-sm">No releases found.</p>
        </div>
    {:else}
        <div class="space-y-8">
            {#each processedReleases as release, index}
                <div class="bg-white dark:bg-gray-800 rounded-md shadow-sm hover:shadow transition-all duration-200 overflow-hidden border border-gray-100 dark:border-gray-700">
                    <!-- Header -->
                    <div class="p-4 pb-3">
                        <div class="flex flex-wrap items-start justify-between gap-3 mb-2">
                            <div class="flex items-center">
                                <div class="w-0.5 h-6 bg-blue-600 rounded-full mr-3"></div>
                                <h2 class="text-lg font-bold text-gray-900 dark:text-white">
                                    <a href={release.html_url} target="_blank" rel="noopener noreferrer" 
                                       class="hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
                                        {getVersionDisplay(release)}
                                    </a>
                                </h2>
                                {#if release.prerelease}
                                    <span class="ml-2 px-1.5 py-0.5 text-xs bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100 rounded-full font-medium">
                                        Pre-release
                                    </span>
                                {/if}
                            </div>
                            
                            <div class="flex items-center text-xs text-gray-500 dark:text-gray-400">
                                <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                                </svg>
                                <time datetime={release.published_at || release.created_at} class="font-medium">
                                    {release.formattedDate}
                                </time>
                                <span class="mx-1.5 text-gray-300 dark:text-gray-600">•</span>
                                <span>{release.relativeTime}</span>
                            </div>
                        </div>
                        
                        {#if release.author}
                            <div class="flex items-center text-xs text-gray-500 dark:text-gray-400">
                                <span class="mr-1.5">Released by</span>
                                <a href={release.author.html_url} target="_blank" rel="noopener noreferrer" 
                                   class="flex items-center hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
                                    <img src={release.author.avatar_url} alt={release.author.login} 
                                         class="w-3.5 h-3.5 rounded-full mr-1 border border-gray-200 dark:border-gray-700" />
                                    <span class="font-medium">{release.author.login}</span>
                                </a>
                            </div>
                        {/if}
                    </div>
                    
                    <!-- Divider -->
                    <div class="h-px bg-gray-100 dark:bg-gray-700"></div>
                    
                    <!-- Content -->
                    <div class="p-4">
                        <div class="prose prose-sm dark:prose-invert prose-headings:font-semibold prose-a:text-blue-600 dark:prose-a:text-blue-400 max-w-none">
                            {@html release.processedBody}
                        </div>
                        
                        <!-- Assets -->
                        {#if release.assets && release.assets.length > 0}
                            <div class="mt-4 pt-4 border-t border-gray-100 dark:border-gray-800">
                                <h3 class="text-xs font-semibold text-gray-900 dark:text-gray-100 mb-2 flex items-center">
                                    <svg class="w-3 h-3 mr-1 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
                                    </svg>
                                    Downloads
                                </h3>
                                <div class="flex flex-wrap gap-2">
                                    {#each release.assets as asset}
                                        <a href={asset.browser_download_url} 
                                           class="text-xs px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors flex items-center"
                                           download>
                                            <span class="text-gray-700 dark:text-gray-300">{asset.name}</span>
                                        </a>
                                    {/each}
                                </div>
                            </div>
                        {/if}
                    </div>
                    
                    <!-- Footer -->
                    <div class="px-4 py-2 bg-gray-50 dark:bg-gray-900/30 border-t border-gray-100 dark:border-gray-800 flex justify-end">
                        <a href={release.html_url} target="_blank" rel="noopener noreferrer" 
                           class="text-blue-600 dark:text-blue-400 hover:underline text-xs flex items-center font-medium">
                            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path>
                            </svg>
                            View on GitHub
                        </a>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div> 