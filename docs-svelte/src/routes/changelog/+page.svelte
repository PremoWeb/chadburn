<script lang="ts">
    import type { GitHubRelease } from '$lib/types';
    import { onMount } from 'svelte';
    import { marked } from 'marked';
    import { formatDate, formatRelativeTime } from '$lib/utils/date';
    
    export let data: { releases: GitHubRelease[] };
    
    let isLoading = true;
    let error: string | null = null;
    
    $: releases = data.releases;
    $: processedReleases = releases.map(release => {
        console.log('Release assets:', release.tag_name, release.assets);
        return {
            ...release,
            formattedDate: formatDate(new Date(release.published_at || release.created_at)),
            relativeTime: formatRelativeTime(new Date(release.published_at || release.created_at)),
            processedBody: marked.parse(release.body || 'No release notes provided.')
        };
    });
    
    onMount(() => {
        isLoading = false;
        if (releases.length === 0) {
            error = "Unable to load releases. Please try again later.";
        }
    });
    
    // Function to format version display
    function formatVersion(release: GitHubRelease): string {
        const tagName = release.tag_name || '';
        const name = release.name || '';
        
        // If name is just the tag name with a 'v' prefix or without it, use just one
        if (tagName && name) {
            if (name === tagName || 
                name === `v${tagName}` || 
                tagName === `v${name}`) {
                return name.startsWith('v') ? name : `v${name}`;
            }
            
            // If name contains the tag, just use the name
            if (name.includes(tagName)) {
                return name;
            }
            
            // Otherwise show both
            return `${name} (${tagName})`;
        }
        
        // Fallback to whatever is available
        return name || tagName;
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
    
    // Function to copy text to clipboard
    function copyToClipboard(text: string) {
        navigator.clipboard.writeText(text);
    }
    
    // Track which buttons have been clicked for copy feedback
    let copiedStates: Record<string, boolean> = {};
    
    function handleCopy(text: string, id: string) {
        navigator.clipboard.writeText(text);
        copiedStates[id] = true;
        
        // Reset the copied state after 2 seconds
        setTimeout(() => {
            copiedStates[id] = false;
        }, 2000);
    }
</script>

<style>
    /* Custom styles for the release cards */
    .release-card {
        margin-bottom: 2.5rem;
        view-transition-name: release-card;
    }
    
    /* Style for the release date */
    .release-date {
        display: flex;
        align-items: center;
        font-size: 0.875rem;
        color: #4b5563;
        gap: 0.5rem;
    }
    
    /* Vertical line indicator */
    .version-indicator {
        width: 4px;
        height: 100%;
        border-radius: 9999px;
        view-transition-name: version-indicator;
    }
    
    /* Improved prose styling */
    :global(.release-content .prose) {
        font-size: 0.9375rem;
        line-height: 1.6;
    }
    
    :global(.release-content .prose h1, .release-content .prose h2) {
        display: none; /* Hide the first heading which is redundant */
    }
    
    :global(.release-content .prose h3) {
        font-size: 1rem;
        margin-top: 1.25rem;
        margin-bottom: 0.75rem;
        color: #374151;
    }
    
    :global(.release-content .prose ul) {
        margin-top: 0.75rem;
        margin-bottom: 0.75rem;
    }
    
    :global(.release-content .prose li) {
        margin-top: 0.25rem;
        margin-bottom: 0.25rem;
    }
    
    /* Page description area */
    .page-header {
        position: relative;
        padding-bottom: 2rem;
        margin-bottom: 2.5rem;
        width: 100%;
        view-transition-name: page-header;
    }
    
    .page-header::after {
        content: '';
        position: absolute;
        bottom: 0;
        left: 50%;
        transform: translateX(-50%);
        width: 100px;
        height: 2px;
        background: linear-gradient(90deg, transparent, #3b82f6 50%, transparent);
    }
</style>

<svelte:head>
    <title>Changelog | Chadburn</title>
    <meta name="description" content="Changelog for Chadburn - View all releases and updates" />
</svelte:head>

<div class="w-full mx-auto px-8 sm:px-6 py-10" style="view-transition-name: changelog-container;">
    <div class="page-header text-center">
        <h1 class="text-3xl font-bold mb-4 text-gray-900" style="view-transition-name: changelog-title;">Changelog</h1>
        <p class="text-gray-600 mx-auto text-base" style="view-transition-name: changelog-description;">
            Track Chadburn's evolution through our release history. This page lists all updates, improvements, and fixes.
        </p>
    </div>
    
    {#if isLoading}
        <div class="flex justify-center items-center py-12">
            <div class="w-6 h-6 border-2 border-blue-200 border-t-blue-600 rounded-full animate-spin"></div>
        </div>
    {:else if error}
        <div class="bg-red-50 p-6 rounded-lg text-center max-w-3xl mx-auto">
            <svg class="w-5 h-5 mx-auto text-red-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
            </svg>
            <p class="text-red-600 text-sm font-medium">{error}</p>
        </div>
    {:else if processedReleases.length === 0}
        <div class="bg-gray-50 p-6 rounded-lg text-center max-w-3xl mx-auto">
            <svg class="w-5 h-5 mx-auto text-gray-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            <p class="text-gray-600 text-sm">No releases found.</p>
        </div>
    {:else}
        <div class="space-y-20">
            {#each processedReleases as release, index}
                <div class="release-card bg-white rounded-lg overflow-hidden" style="view-transition-name: release-card-{index};">
                    <!-- Header -->
                    <div class="flex items-stretch">
                        <!-- Version indicator -->
                        <div class={`version-indicator ${getVersionColor(index)}`} style="view-transition-name: version-indicator-{index};"></div>
                        
                        <div class="flex-1 px-12 py-8">
                            <!-- Version and date -->
                            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-4">
                                <!-- Version -->
                                <div class="flex items-center">
                                    <h2 class="text-xl font-bold text-gray-900" style="view-transition-name: release-title-{index};">
                                        <a href={release.html_url} target="_blank" rel="noopener noreferrer" 
                                           class="hover:text-blue-600 transition-colors">
                                            {formatVersion(release)}
                                        </a>
                                    </h2>
                                    {#if release.prerelease}
                                        <span class="ml-2 px-2 py-0.5 text-xs bg-yellow-100 text-yellow-800 rounded-full font-medium">
                                            Pre-release
                                        </span>
                                    {/if}
                                </div>
                                
                                <!-- Date and GitHub link -->
                                <div class="flex flex-col items-end gap-2">
                                    <div class="release-date flex items-center gap-2">
                                        <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                                        </svg>
                                        <time datetime={release.published_at || release.created_at} class="font-medium">
                                            {release.formattedDate}
                                        </time>
                                        <span class="text-gray-400">•</span>
                                        <span class="text-gray-500">{release.relativeTime}</span>
                                    </div>
                                    
                                    <a href={release.html_url} target="_blank" rel="noopener noreferrer" 
                                       class="text-blue-600 hover:underline text-sm flex items-center font-medium">
                                        <svg class="w-3.5 h-3.5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path>
                                        </svg>
                                        View on GitHub
                                    </a>
                                </div>
                            </div>
                            
                            {#if release.author}
                            <div class="mb-5">
                                <div class="flex items-center text-sm gap-2"> <!-- Added gap-5 -->
                                    <span class="text-gray-500">Released by</span> <!-- Removed mr-3 -->
                                    <a href={release.author.html_url} target="_blank" rel="noopener noreferrer" 
                                       class="flex items-center hover:text-blue-600 transition-colors gap-2"> <!-- Added gap-3 -->
                                        <img src={release.author.avatar_url} alt={release.author.login} 
                                             class="w-5 h-5 rounded-full" /> <!-- Removed mr-2 -->
                                        <span class="font-medium">{release.author.login}</span>
                                    </a>
                                </div>
                            </div>
                            {/if}
                            
                            <!-- Content -->
                            <div class="release-content mt-4 rounded-md">
                                <div class="prose prose-sm prose-headings:font-semibold prose-a:text-blue-600 prose-ul:list-disc prose-li:ml-4 max-w-none">
                                    {@html release.processedBody}
                                </div>
                            </div>
                            
                            <!-- Assets -->
                            {#if release.assets && release.assets.length > 0}
                                <div class="mt-8 pt-5 border-t border-gray-100">
                                    <h3 class="text-base font-semibold text-gray-900 mb-4 flex items-center">
                                        <svg class="w-4 h-4 mr-2 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
                                        </svg>
                                        Downloads ({release.assets.length})
                                    </h3>
                                    
                                    <!-- Group assets by platform -->
                                    {#if release.assets.some(asset => asset.name.includes('.tar.gz') || asset.name.includes('.zip'))}
                                        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
                                            <!-- Linux Downloads -->
                                            <div class="bg-gray-50 p-4 rounded-lg border border-gray-200">
                                                <div class="flex items-center mb-3">
                                                    <svg class="w-5 h-5 mr-2 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                                                    </svg>
                                                    <h4 class="font-medium text-gray-900">Linux</h4>
                                                </div>
                                                <div class="space-y-2">
                                                    {#each release.assets.filter(asset => asset.name.includes('linux')) as asset}
                                                        <a href={asset.browser_download_url} 
                                                           class="flex items-center justify-between p-3 bg-white rounded border border-gray-200 hover:bg-blue-50 hover:border-blue-200 transition-colors group"
                                                           download>
                                                            <div class="flex items-center">
                                                                <svg class="w-4 h-4 text-blue-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"></path>
                                                                </svg>
                                                                <div>
                                                                    <div class="text-sm font-medium text-gray-900">{asset.name.replace('chadburn-', '').replace('.tar.gz', '')}</div>
                                                                    <div class="text-xs text-gray-500">{(asset.size / (1024 * 1024)).toFixed(2)} MB</div>
                                                                </div>
                                                            </div>
                                                            <span class="text-xs bg-blue-100 text-blue-800 py-1 px-2 rounded-full opacity-0 group-hover:opacity-100 transition-opacity">Download</span>
                                                        </a>
                                                    {/each}
                                                </div>
                                            </div>
                                            
                                            <!-- macOS Downloads -->
                                            <div class="bg-gray-50 p-4 rounded-lg border border-gray-200">
                                                <div class="flex items-center mb-3">
                                                    <svg class="w-5 h-5 mr-2 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                                                    </svg>
                                                    <h4 class="font-medium text-gray-900">macOS</h4>
                                                </div>
                                                <div class="space-y-2">
                                                    {#each release.assets.filter(asset => asset.name.includes('darwin')) as asset}
                                                        <a href={asset.browser_download_url} 
                                                           class="flex items-center justify-between p-3 bg-white rounded border border-gray-200 hover:bg-blue-50 hover:border-blue-200 transition-colors group"
                                                           download>
                                                            <div class="flex items-center">
                                                                <svg class="w-4 h-4 text-blue-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"></path>
                                                                </svg>
                                                                <div>
                                                                    <div class="text-sm font-medium text-gray-900">{asset.name.replace('chadburn-', '').replace('.tar.gz', '')}</div>
                                                                    <div class="text-xs text-gray-500">{(asset.size / (1024 * 1024)).toFixed(2)} MB</div>
                                                                </div>
                                                            </div>
                                                            <span class="text-xs bg-blue-100 text-blue-800 py-1 px-2 rounded-full opacity-0 group-hover:opacity-100 transition-opacity">Download</span>
                                                        </a>
                                                    {/each}
                                                </div>
                                            </div>
                                            
                                            <!-- Windows Downloads -->
                                            <div class="bg-gray-50 p-4 rounded-lg border border-gray-200">
                                                <div class="flex items-center mb-3">
                                                    <svg class="w-5 h-5 mr-2 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                                                    </svg>
                                                    <h4 class="font-medium text-gray-900">Windows</h4>
                                                </div>
                                                <div class="space-y-2">
                                                    {#each release.assets.filter(asset => asset.name.includes('windows')) as asset}
                                                        <a href={asset.browser_download_url} 
                                                           class="flex items-center justify-between p-3 bg-white rounded border border-gray-200 hover:bg-blue-50 hover:border-blue-200 transition-colors group"
                                                           download>
                                                            <div class="flex items-center">
                                                                <svg class="w-4 h-4 text-blue-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"></path>
                                                                </svg>
                                                                <div>
                                                                    <div class="text-sm font-medium text-gray-900">{asset.name.replace('chadburn-', '').replace('.zip', '')}</div>
                                                                    <div class="text-xs text-gray-500">{(asset.size / (1024 * 1024)).toFixed(2)} MB</div>
                                                                </div>
                                                            </div>
                                                            <span class="text-xs bg-blue-100 text-blue-800 py-1 px-2 rounded-full opacity-0 group-hover:opacity-100 transition-opacity">Download</span>
                                                        </a>
                                                    {/each}
                                                </div>
                                            </div>
                                        </div>
                                        
                                        <!-- Checksums -->
                                        {#if release.assets.some(asset => asset.name.includes('checksums.txt'))}
                                            <div class="mt-4">
                                                <div class="flex items-center mb-2">
                                                    <svg class="w-4 h-4 mr-2 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
                                                    </svg>
                                                    <h4 class="text-sm font-medium text-gray-900">Checksums</h4>
                                                </div>
                                                {#each release.assets.filter(asset => asset.name.includes('checksums.txt')) as asset}
                                                    <a href={asset.browser_download_url} 
                                                       class="text-sm text-blue-600 hover:underline flex items-center"
                                                       download>
                                                        <svg class="w-3.5 h-3.5 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"></path>
                                                        </svg>
                                                        Download SHA256 checksums
                                                    </a>
                                                {/each}
                                            </div>
                                        {/if}
                                    {:else}
                                        <!-- Display other assets if no binaries are found -->
                                        <div class="flex flex-wrap gap-2">
                                            {#each release.assets as asset}
                                                <a href={asset.browser_download_url} 
                                                   class="text-xs px-3 py-2 bg-gray-100 rounded hover:bg-gray-200 transition-colors flex items-center gap-1"
                                                   download>
                                                    <svg class="w-3 h-3 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"></path>
                                                    </svg>
                                                    <span class="text-gray-700">{asset.name}</span>
                                                    <span class="text-gray-500 text-xs">({Math.round(asset.size / 1024)} KB)</span>
                                                </a>
                                            {/each}
                                        </div>
                                    {/if}
                                </div>
                            {:else}
                                <div class="mt-8 pt-5 border-t border-gray-100">
                                    <div class="text-sm text-gray-500 italic">No downloadable assets available for this release.</div>
                                </div>
                            {/if}
                            
                            <!-- Docker Images -->
                            {#if release.tag_name}
                                <div class="mt-10 pt-6 border-t border-gray-100">
                                    <h3 class="text-base font-semibold text-gray-900 mb-6 flex items-center">
                                        <svg class="w-5 h-5 mr-3 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"></path>
                                        </svg>
                                        Docker Images
                                    </h3>
                                    <div class="bg-white border border-gray-200 p-6 rounded-md shadow-sm">
                                        <!-- Docker Pull Commands -->
                                        <div class="mb-8">
                                            <div class="text-sm font-medium text-gray-700 mb-3">Pull this specific version:</div>
                                            <div class="bg-blue-50 border border-blue-100 text-gray-800 p-4 rounded-md font-mono text-sm overflow-x-auto whitespace-nowrap flex items-center justify-between group">
                                                <code class="px-2">docker pull ghcr.io/chadburn-container/chadburn:{release.tag_name.replace(/^v/, '')}</code>
                                                <button class="opacity-70 group-hover:opacity-100 transition-opacity bg-blue-100 hover:bg-blue-200 text-blue-800 p-2 rounded-md ml-4 flex items-center justify-center" 
                                                        on:click={() => handleCopy(`docker pull ghcr.io/chadburn-container/chadburn:${release.tag_name.replace(/^v/, '')}`, `pull-${index}`)}>
                                                    {#if copiedStates[`pull-${index}`]}
                                                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                                                        </svg>
                                                    {:else}
                                                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"></path>
                                                        </svg>
                                                    {/if}
                                                </button>
                                            </div>
                                        </div>
                                        
                                        <div class="mb-8">
                                            <div class="text-sm font-medium text-gray-700 mb-3">Or use the latest tag:</div>
                                            <div class="bg-blue-50 border border-blue-100 text-gray-800 p-4 rounded-md font-mono text-sm overflow-x-auto whitespace-nowrap flex items-center justify-between group">
                                                <code class="px-2">docker pull ghcr.io/chadburn-container/chadburn:latest</code>
                                                <button class="opacity-70 group-hover:opacity-100 transition-opacity bg-blue-100 hover:bg-blue-200 text-blue-800 p-2 rounded-md ml-4 flex items-center justify-center" 
                                                        on:click={() => handleCopy('docker pull ghcr.io/chadburn-container/chadburn:latest', `latest-${index}`)}>
                                                    {#if copiedStates[`latest-${index}`]}
                                                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                                                        </svg>
                                                    {:else}
                                                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"></path>
                                                        </svg>
                                                    {/if}
                                                </button>
                                            </div>
                                        </div>
                                        
                                        <!-- Supported Architectures -->
                                        <div class="mt-8 mb-8">
                                            <div class="text-sm font-medium text-gray-700 mb-4">Supported Architectures:</div>
                                            <div class="flex flex-wrap gap-4">
                                                <span class="inline-flex items-center px-4 py-2 rounded-md text-sm font-medium bg-blue-50 text-blue-800 border border-blue-200">
                                                    <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                                                    </svg>
                                                    amd64
                                                </span>
                                                <span class="inline-flex items-center px-4 py-2 rounded-md text-sm font-medium bg-green-50 text-green-800 border border-green-200">
                                                    <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                                                    </svg>
                                                    arm64
                                                </span>
                                                <span class="inline-flex items-center px-4 py-2 rounded-md text-sm font-medium bg-purple-50 text-purple-800 border border-purple-200">
                                                    <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                                                    </svg>
                                                    armv7
                                                </span>
                                            </div>
                                        </div>
                                        
                                        <!-- Usage Example -->
                                        <div class="mt-8 mb-8">
                                            <div class="text-sm font-medium text-gray-700 mb-3">Basic usage example:</div>
                                            <div class="bg-blue-50 border border-blue-100 text-gray-800 p-4 rounded-md font-mono text-sm overflow-x-auto whitespace-nowrap flex items-center justify-between group">
                                                <code class="px-2">docker run -d --name chadburn -v /path/to/config:/config ghcr.io/chadburn-container/chadburn:{release.tag_name.replace(/^v/, '')}</code>
                                                <button class="opacity-70 group-hover:opacity-100 transition-opacity bg-blue-100 hover:bg-blue-200 text-blue-800 p-2 rounded-md ml-4 flex items-center justify-center" 
                                                        on:click={() => handleCopy(`docker run -d --name chadburn -v /path/to/config:/config ghcr.io/chadburn-container/chadburn:${release.tag_name.replace(/^v/, '')}`, `run-${index}`)}>
                                                    {#if copiedStates[`run-${index}`]}
                                                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                                                        </svg>
                                                    {:else}
                                                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"></path>
                                                        </svg>
                                                    {/if}
                                                </button>
                                            </div>
                                        </div>
                                        
                                        <!-- Links -->
                                        <div class="mt-8 text-sm text-gray-700 flex flex-wrap gap-6">
                                            <a href="https://github.com/chadburn-container/chadburn/pkgs/container/chadburn" target="_blank" rel="noopener noreferrer" class="text-blue-700 hover:underline flex items-center gap-2">
                                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                                                </svg>
                                                View all available tags
                                            </a>
                                            <a href="https://github.com/chadburn-container/chadburn/blob/main/README.md" target="_blank" rel="noopener noreferrer" class="text-blue-700 hover:underline flex items-center gap-2">
                                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                                                </svg>
                                                Documentation
                                            </a>
                                        </div>
                                    </div>
                                </div>
                            {/if}
                            
                            <!-- Footer -->
                            <div class="mt-8 pt-4">
                                <!-- Empty footer for spacing -->
                            </div>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div> 