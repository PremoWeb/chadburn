import { fetchReleases } from '$lib/api';
import type { GitHubRelease } from '$lib/types';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    // Fetch releases from GitHub API
    const releases = await fetchReleases(30); // Get up to 30 releases
    
    return {
        releases
    };
}; 