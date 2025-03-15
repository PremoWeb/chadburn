import type { Contributor } from './types';

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