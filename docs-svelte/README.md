# sv

Everything you need to build a Svelte project, powered by [`sv`](https://github.com/sveltejs/cli).

## Creating a project

If you're seeing this, you've probably already done this step. Congrats!

```bash
# create a new project in the current directory
npx sv create

# create a new project in my-app
npx sv create my-app
```

## Developing

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```bash
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```bash
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.

# Chadburn Documentation Site

This is the documentation site for Chadburn, a modern job scheduler for Docker environments. The site is built with SvelteKit and uses Svelte 5 with Runes mode.

## Development

Once you've cloned the project and installed dependencies with `bun install`, start a development server:

```bash
bun run dev

# or start the server and open the app in a new browser tab
bun run dev -- --open
```

## Building

To create a production version of the site:

```bash
bun run build
```

You can preview the production build with `bun run preview`.

## Contributors Section

The site includes a contributors section that displays all code contributors to the project. This section is automatically updated by a GitHub Action whenever there's a new commit to the main branch.

### How it works

1. The GitHub Action runs the `scripts/fetch-contributors.js` script to fetch the latest contributors from the GitHub API
2. The script saves the contributors data to `static/data/contributors.json`
3. The site loads this data at runtime and displays it on the home page

### Manually updating contributors

You can manually update the contributors list by running:

```bash
bun run update-contributors
```

## Deployment

The site is automatically deployed to GitHub Pages when changes are pushed to the main branch.

### Custom Domain Setup

The site is configured to use the custom domain `chadburn.dev`. To set this up:

1. In your DNS provider, add the following records:
   - A record: `@` pointing to `185.199.108.153`
   - A record: `@` pointing to `185.199.109.153`
   - A record: `@` pointing to `185.199.110.153`
   - A record: `@` pointing to `185.199.111.153`
   - CNAME record: `www` pointing to `premoweb.github.io`

2. In your GitHub repository settings:
   - Go to Settings > Pages
   - Under "Custom domain", enter `chadburn.dev`
   - Check "Enforce HTTPS" once the DNS changes have propagated

The CNAME file is already included in the static directory and will be automatically deployed with the site.

### Base Path Configuration

The site is configured to use the root path (`/`) instead of `/Chadburn`. This is set in the `svelte.config.js` file:

```js
paths: {
  base: '' // Use root path for all environments
}
```

# Chadburn Documentation

This directory contains the Svelte-based documentation for the Chadburn project.

## Development

### Prerequisites

- [Bun](https://bun.sh/) (latest version)
- Node.js 18 or later

### Setup

1. Install dependencies:

```bash
cd docs-svelte
bun install
```

2. Start the development server:

```bash
bun run dev
```

3. Open your browser and navigate to `http://localhost:5173`

## Building

To build the documentation site:

```bash
bun run build
```

The built files will be in the `build` directory.

## Deployment

The documentation is automatically deployed to Cloudflare Pages when changes are pushed to the `main` branch.

### Setting up Cloudflare Pages Deployment

To set up the GitHub Actions workflow for deploying to Cloudflare Pages, you need to add the following secrets to your GitHub repository:

1. `CLOUDFLARE_API_TOKEN`: An API token with the "Cloudflare Pages" permission
2. `CLOUDFLARE_ACCOUNT_ID`: Your Cloudflare account ID

#### Creating a Cloudflare API Token

1. Log in to your Cloudflare dashboard
2. Go to "My Profile" > "API Tokens"
3. Click "Create Token"
4. Select "Create Custom Token"
5. Name your token (e.g., "GitHub Actions - Chadburn Docs")
6. Under "Permissions", add:
   - Account > Cloudflare Pages > Edit
   - Account > Account Settings > Read
7. Under "Account Resources", select your account
8. Click "Continue to summary" and then "Create Token"
9. Copy the token value (you won't be able to see it again)

#### Finding Your Cloudflare Account ID

1. Log in to your Cloudflare dashboard
2. Your account ID is in the URL: `https://dash.cloudflare.com/<account-id>`
3. It's also visible in the right sidebar when viewing your account details

#### Adding Secrets to GitHub

1. Go to your GitHub repository
2. Click on "Settings" > "Secrets and variables" > "Actions"
3. Click "New repository secret"
4. Add the `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID` secrets

#### Setting up the Cloudflare Pages Project

Before the first deployment, you need to create a Cloudflare Pages project:

1. Go to your Cloudflare dashboard
2. Navigate to "Pages"
3. Click "Create a project"
4. Choose "Connect to Git"
5. Select your GitHub repository
6. Configure your project:
   - Project name: `chadburn-docs` (must match the `projectName` in the workflow)
   - Production branch: `main`
   - Build command: Leave empty (we build in GitHub Actions)
   - Build directory: Leave empty (we specify in the deployment)
7. Click "Save and Deploy"

After the initial setup, the GitHub Actions workflow will handle all future deployments.
