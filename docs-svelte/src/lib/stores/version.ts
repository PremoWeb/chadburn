import { readable } from 'svelte/store';
import { browser } from '$app/environment';
import { base } from '$app/paths';

type VersionData = {
  version: string;
  loading: boolean;
  error: string | null;
};

export const version = readable<VersionData>(
  { version: '', loading: true, error: null },
  (set) => {
    if (browser) {
      fetch(`${base}/data/version.json`)
        .then((response) => {
          if (!response.ok) {
            throw new Error(`Failed to load version data: ${response.status}`);
          }
          return response.json();
        })
        .then((data) => {
          set({ version: data.version, loading: false, error: null });
        })
        .catch((error) => {
          console.error('Error loading version:', error);
          set({ version: '', loading: false, error: error.message });
        });
    }

    return () => {
      // Cleanup function (if needed)
    };
  }
); 