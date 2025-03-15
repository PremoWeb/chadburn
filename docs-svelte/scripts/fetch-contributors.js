#!/usr/bin/env node

/**
 * This script fetches contributors from the GitHub API and saves them to a JSON file
 * It's meant to be run by the GitHub Action when new commits are pushed
 */

const fs = require('fs');
const path = require('path');

async function fetchContributors() {
    try {
        console.log('Fetching contributors from GitHub API...');
        
        // Use GitHub token from environment if available (for CI)
        const headers = {};
        if (process.env.GITHUB_TOKEN) {
            console.log('Using GitHub token for authentication');
            headers.Authorization = `token ${process.env.GITHUB_TOKEN}`;
        }
        
        const response = await fetch('https://api.github.com/repos/PremoWeb/Chadburn/contributors', { headers });
        
        // Check for rate limiting
        const rateLimit = {
            limit: response.headers.get('x-ratelimit-limit'),
            remaining: response.headers.get('x-ratelimit-remaining'),
            reset: response.headers.get('x-ratelimit-reset')
        };
        
        console.log(`GitHub API Rate Limit: ${rateLimit.remaining}/${rateLimit.limit} remaining`);
        
        if (!response.ok) {
            if (response.status === 403 && rateLimit.remaining === '0') {
                const resetDate = new Date(rateLimit.reset * 1000);
                throw new Error(`GitHub API rate limit exceeded. Resets at ${resetDate.toLocaleString()}`);
            }
            throw new Error(`GitHub API error: ${response.status} - ${await response.text()}`);
        }
        
        let contributors = await response.json();
        console.log(`Found ${contributors.length} contributors`);
        
        // Filter out bots or other non-user accounts if needed
        contributors = contributors.filter(contributor => 
            contributor.type === 'User' && !contributor.login.includes('[bot]')
        );
        
        // Add Nick Maietta as the lead developer
        // We'll check if he's already in the list and update his entry
        const nickIndex = contributors.findIndex(c => c.login.toLowerCase() === 'maietta');
        
        if (nickIndex >= 0) {
            // Update existing entry
            contributors[nickIndex] = {
                ...contributors[nickIndex],
                isLeadDeveloper: true,
                name: 'Nick Maietta',
                social: {
                    x: 'https://x.com/maietta',
                    github: 'https://github.com/maietta'
                }
            };
        } else {
            // Add Nick as a new entry at the beginning
            contributors.unshift({
                login: 'maietta',
                id: 0,
                node_id: '',
                avatar_url: 'https://avatars.githubusercontent.com/u/maietta',
                gravatar_id: '',
                url: 'https://api.github.com/users/maietta',
                html_url: 'https://github.com/maietta',
                followers_url: '',
                following_url: '',
                gists_url: '',
                starred_url: '',
                subscriptions_url: '',
                organizations_url: '',
                repos_url: '',
                events_url: '',
                received_events_url: '',
                type: 'User',
                site_admin: false,
                contributions: 1000, // High number to ensure top placement
                isLeadDeveloper: true,
                name: 'Nick Maietta',
                social: {
                    x: 'https://x.com/maietta',
                    github: 'https://github.com/maietta'
                }
            });
        }
        
        // Create the static directory if it doesn't exist
        const staticDir = path.join(__dirname, '..', 'static', 'data');
        if (!fs.existsSync(staticDir)) {
            fs.mkdirSync(staticDir, { recursive: true });
        }
        
        // Save the contributors to a JSON file
        const filePath = path.join(staticDir, 'contributors.json');
        fs.writeFileSync(filePath, JSON.stringify(contributors, null, 2));
        console.log(`Contributors saved to ${filePath}`);
    } catch (error) {
        console.error('Error fetching contributors:', error);
        process.exit(1);
    }
}

fetchContributors(); 