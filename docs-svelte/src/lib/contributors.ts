import { base } from '$app/paths';
import type { Contributor } from './types';

/**
 * Loads contributors data from the static JSON file
 * @returns A promise that resolves to an array of contributors and an error message if any
 */
export async function loadContributors(): Promise<{ 
  contributors: Contributor[]; 
  error: string | null;
}> {
  try {
    console.log(`Fetching contributors from ${base}/data/contributors.json`);
    const response = await fetch(`${base}/data/contributors.json`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch contributors: ${response.status} ${response.statusText}`);
    }
    
    const contributors = await response.json();
    console.log(`Loaded ${contributors.length} contributors`);
    
    return { 
      contributors, 
      error: null 
    };
  } catch (error) {
    console.error('Error loading contributors:', error);
    return { 
      contributors: [], 
      error: error instanceof Error ? error.message : 'Unknown error loading contributors' 
    };
  }
} 