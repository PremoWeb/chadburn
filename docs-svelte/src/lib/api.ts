import type { Contributor, GitHubRelease } from './types';

/**
 * Fetches contributors from the GitHub API
 * This function is meant to be called during build time by the GitHub Action
 * The results are saved to a static JSON file
 */
export async function fetchContributors(): Promise<Contributor[]> {
    try {
        const response = await fetch('https://api.github.com/repos/PremoWeb/Chadburn/contributors');
        
        if (!response.ok) {
            throw new Error(`GitHub API error: ${response.status}`);
        }
        
        const contributors: Contributor[] = await response.json();
        return contributors;
    } catch (error) {
        console.error('Error fetching contributors:', error);
        return [];
    }
}

/**
 * Fetches releases from the GitHub API
 * @param perPage Number of releases to fetch per page
 * @param page Page number to fetch
 * @returns Array of GitHub releases
 */
export async function fetchReleases(perPage: number = 10, page: number = 1): Promise<GitHubRelease[]> {
    try {
        const response = await fetch(
            `https://api.github.com/repos/PremoWeb/Chadburn/releases?per_page=${perPage}&page=${page}`
        );
        
        if (!response.ok) {
            throw new Error(`GitHub API error: ${response.status}`);
        }
        
        const releases: GitHubRelease[] = await response.json();
        return releases;
    } catch (error) {
        console.error('Error fetching releases:', error);
        return [];
    }
} 