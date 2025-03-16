<script lang="ts">
	import { page } from '$app/state';
	import { base } from '$app/paths';
	import { afterNavigate, beforeNavigate } from '$app/navigation';
	import { marked } from 'marked';
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import TableOfContents from '$lib/components/TableOfContents.svelte';
	import TabCodeBlock from '$lib/components/TabCodeBlock.svelte';
	import { createHighlighter, type Highlighter } from 'shiki';
	
	// Track if this is the initial page load
	let isInitialLoad = true;
	
	// Track current navigation state
	let isNavigating = $state(false);
	let currentHtmlContent = $state('');
	let isLoading = $state(true);
	
	// Track navigation direction for slide transitions
	let slideDirection = $state('right'); // 'left' or 'right'
	let previousPath = $state('');
	let useSlideTransition = $state(false); // Whether to use slide or vertical transition
	
	// Main navigation sections in order
	const mainNavSections = ['about', 'docs', 'metrics', 'api'];
	
	// Function to check if a path is a sidebar link (any documentation page)
	function isSidebarLink(path: string): boolean {
		// Simple check - any path that starts with /docs/ is a sidebar link
		return path.startsWith('/docs/') || path === '/docs';
	}
	
	// Function to determine navigation direction
	function determineSlideDirection(fromPath: string, toPath: string) {
		// Check if both paths are sidebar links
		const fromIsSidebarLink = isSidebarLink(fromPath);
		const toIsSidebarLink = isSidebarLink(toPath);
		
		// If both are sidebar links, use fade transition
		if (fromIsSidebarLink && toIsSidebarLink) {
			useSlideTransition = false;
			return 'none';
		}
		
		// Otherwise, we're navigating between main sections, use slide transition
		useSlideTransition = true;
		
		// Extract the first segment of the path for main navigation
		const fromSegment = fromPath.split('/')[1] || '';
		const toSegment = toPath.split('/')[1] || '';
		
		// Get positions in the navigation order
		const fromIndex = mainNavSections.indexOf(fromSegment);
		const toIndex = mainNavSections.indexOf(toSegment);
		
		// If either path is not in the nav order, default to right
		if (fromIndex === -1 || toIndex === -1) {
			return 'right';
		}
		
		// Determine direction based on position
		const direction = fromIndex < toIndex ? 'right' : 'left';
		return direction;
	}
	
	// Shiki highlighter instance
	let highlighter: Highlighter | null = $state(null);
	
	// State to track if TOC has content
	let tocHasContent = $state(false);
	
	// Initialize Shiki highlighter
	async function initHighlighter() {
		try {
			highlighter = await createHighlighter({
				themes: ['github-light'],
				langs: [
					// Web languages
					'javascript', 'typescript', 'svelte', 'html', 'css', 'json', 
					// Server languages
					'go', 'rust', 'python', 'ruby', 'php', 'java', 'c', 'cpp', 'csharp',
					// Shell/CLI
					'bash', 'shell', 'powershell', 
					// Data formats
					'markdown', 'yaml', 'toml', 'xml', 'sql',
					// Config files
					'ini', 'dockerfile', 'nginx',
					// Other
					'diff', 'graphql'
				]
			});
			
			// Re-process content if it's already loaded
			if (htmlContent) {
				processContent();
			}
		} catch (error) {
			// Critical error - highlighter initialization failed
		}
	}
	
	// Initial content processing
	onMount(() => {
		initHighlighter();
		
		// Handle hash navigation on initial load
		handleHashNavigation();
		
		// Initialize tabbed code blocks
		setTimeout(initTabbedCodeBlocks, 300);
		
		// Process content on initial load
		if (isInitialLoad) {
			processContent();
		}
		
		// After content is rendered, add IDs to headings in the DOM
		setTimeout(addIdsToHeadingsInDOM, 100);
		
		// Setup copy buttons after content is rendered
		setTimeout(() => {
			setupCodeBlockCopyButtons();
		}, 200);
		
		// Setup hover effects for code blocks
		setTimeout(() => {
			setupCodeBlockHoverEffects();
		}, 250);
		
		// Initialize tabbed code blocks
		setTimeout(initTabbedCodeBlocks, 300);
	});
	
	// Get data from the page loader
	// In Svelte 5, data is passed directly to the component
	const props = $props();
	
	// Create a reactive reference to the data that updates when props change
	$effect(() => {
		processContent();
	});
	
	// Get the current path and hash from the URL - this will update reactively
	const currentPath = $derived(page.url.pathname);
	const currentHash = $derived(page.url.hash);
	
	// Handle hash navigation after content is processed
	$effect(() => {
		if (htmlContent && currentHash) {
			handleHashNavigation();
		}
	});
	
	// Function to handle hash navigation
	function handleHashNavigation() {
		if (currentHash) {
			// Remove the # from the hash
			const targetId = currentHash.substring(1);
			
			// Wait a bit for the DOM to update
			setTimeout(() => {
				const targetElement = document.getElementById(targetId);
				if (targetElement) {
					// Scroll to the element with a small offset for the header
					const headerOffset = 80; // Adjust based on your header height
					const elementPosition = targetElement.getBoundingClientRect().top;
					const offsetPosition = elementPosition + window.pageYOffset - headerOffset;
					
					window.scrollTo({
						top: offsetPosition,
						behavior: 'smooth'
					});
					
					// Highlight the element briefly
					targetElement.classList.add('highlight-section');
					setTimeout(() => {
						targetElement.classList.remove('highlight-section');
					}, 2000);
				}
			}, 100);
		}
	}
	
	// Handle navigation events
	afterNavigate(() => {
		processContent();
		
		// After content is rendered, add IDs to headings in the DOM
		setTimeout(addIdsToHeadingsInDOM, 100);
		
		// Setup copy buttons after content is rendered
		setTimeout(() => {
			setupCodeBlockCopyButtons();
		}, 200);
		
		// Setup hover effects for code blocks
		setTimeout(() => {
			setupCodeBlockHoverEffects();
		}, 250);
		
		// Initialize tabbed code blocks after content is rendered
		setTimeout(initTabbedCodeBlocks, 300);
		
		// Scroll to top on navigation
		window.scrollTo({ top: 0, behavior: 'smooth' });
	});
	
	// Handle special case for /installation URL
	let actualSlug = $state('');
	
	// Handle title generation
	let title = $state('Documentation');
	
	// Processed HTML content
	let htmlContent = $state('');
	
	// Extract data from props
	function getData() {
		const data = props.data || props;
		return {
			slug: data.slug || '',
			content: data.content || ''
		};
	}
	
	// Process the content whenever the data or path changes
	function processContent() {
		const { slug, content } = getData();
		
		// Determine the title based on the URL path
		if (currentPath === '/installation') {
			actualSlug = 'installation';
			title = 'Installation';
		} else if (currentPath === '/') {
			actualSlug = '';
			title = 'Documentation';
		} else {
			actualSlug = slug || '';
			
			// Update title
			if (!actualSlug) {
				title = 'Documentation';
			} else {
				// Split by dash or slash and capitalize each word
				title = actualSlug
					.split(/[-\/]/)
					.filter(word => word.length > 0)
					.map(word => word.charAt(0).toUpperCase() + word.slice(1))
					.join(' ');
			}
		}
		
		// Process the markdown content
		if (content) {
			try {
				// Remove front matter if present
				const cleanedMarkdown = content.replace(/^---[\s\S]*?---\n/, '');
				
				// Convert markdown to HTML
				const html = marked.parse(cleanedMarkdown) as string;
				
				// Process the HTML to add IDs to headings and enhance code blocks
				const withIds = addIdsToHeadings(html);
				const processed = processCodeBlocks(withIds);
				
				// Update the content with transition
				if (isInitialLoad) {
					// On initial load, prepare the content but don't show it immediately
					// This prevents layout shifts by ensuring all processing is done before display
					currentHtmlContent = processed;
					
					// Use a short timeout to ensure the DOM has time to prepare
					setTimeout(() => {
						htmlContent = processed;
						isInitialLoad = false;
						isLoading = false;
						
						// Handle hash navigation after initial content load
						if (currentHash) {
							setTimeout(() => handleHashNavigation(), 200);
						}
					}, 50); // Reduced from 100ms to 50ms for faster rendering
				} else {
					// For navigation, update the current content for transition
					currentHtmlContent = processed;
					
					// After a short delay, update the actual content
					setTimeout(() => {
						htmlContent = processed;
						isNavigating = false;
						isLoading = false;
						
						// Handle hash navigation after content update
						if (currentHash) {
							setTimeout(() => handleHashNavigation(), 200);
						}
					}, 100);
				}
			} catch (error) {
				// Critical error - markdown processing failed
				htmlContent = `<div class="error">Error processing markdown: ${error instanceof Error ? error.message : String(error)}</div>`;
				currentHtmlContent = htmlContent;
				isNavigating = false;
				isLoading = false;
			}
		} else {
			// Critical error - no content provided
			htmlContent = '<div class="error">No content available</div>';
			currentHtmlContent = htmlContent;
			isNavigating = false;
			isLoading = false;
		}
	}
	
	// Function to generate a slug from text
	function generateSlug(text: string): string {
		return text
			.toLowerCase()
			.replace(/[^\w\s-]/g, '') // Remove special characters
			.replace(/\s+/g, '-')     // Replace spaces with hyphens
			.replace(/-+/g, '-');     // Replace multiple hyphens with single hyphen
	}
	
	// Function to add IDs to headings
	function addIdsToHeadings(html: string): string {
		// Use a regular expression to find all headings
		return html.replace(/<(h[1-6])>(.*?)<\/h[1-6]>/g, (match, tag, content) => {
			const id = generateSlug(content);
			return `<${tag} id="${id}">${content}</${tag}>`;
		});
	}
	
	// Function to map markdown language identifiers to Shiki language identifiers
	function mapLanguage(lang: string): string {
		// Common language aliases
		const langMap: Record<string, string> = {
			// JavaScript and TypeScript
			'js': 'javascript',
			'jsx': 'javascript',
			'ts': 'typescript',
			'tsx': 'typescript',
			'mjs': 'javascript',
			'cjs': 'javascript',
			
			// Shell/CLI
			'sh': 'bash',
			'shell': 'bash',
			'zsh': 'bash',
			'console': 'bash',
			'terminal': 'bash',
			
			// Web
			'htm': 'html',
			'xhtml': 'html',
			'svg': 'xml',
			'scss': 'css',
			'sass': 'css',
			'less': 'css',
			'stylus': 'css',
			
			// Data formats
			'md': 'markdown',
			'yml': 'yaml',
			'conf': 'ini',
			
			// C-family
			'cc': 'cpp',
			'h': 'c',
			'hpp': 'cpp',
			'cs': 'csharp',
			
			// Other
			'rs': 'rust',
			'py': 'python',
			'rb': 'ruby',
			'pl': 'perl',
			'dockerfile': 'docker',
			'makefile': 'make',
			'mk': 'make'
		};
		
		return langMap[lang.toLowerCase()] || lang;
	}
	
	// Function to get a friendly language name
	function getFriendlyLanguageName(lang: string): string {
		const languageNames: Record<string, string> = {
			'js': 'JavaScript',
			'javascript': 'JavaScript',
			'jsx': 'JSX',
			'ts': 'TypeScript',
			'typescript': 'TypeScript',
			'tsx': 'TSX',
			'svelte': 'Svelte',
			'html': 'HTML',
			'css': 'CSS',
			'scss': 'SCSS',
			'sass': 'Sass',
			'less': 'Less',
			'json': 'JSON',
			'yaml': 'YAML',
			'yml': 'YAML',
			'markdown': 'Markdown',
			'md': 'Markdown',
			'bash': 'Bash',
			'shell': 'Shell',
			'sh': 'Shell',
			'zsh': 'Shell',
			'python': 'Python',
			'py': 'Python',
			'ruby': 'Ruby',
			'rb': 'Ruby',
			'go': 'Go',
			'rust': 'Rust',
			'rs': 'Rust',
			'java': 'Java',
			'c': 'C',
			'cpp': 'C++',
			'csharp': 'C#',
			'cs': 'C#',
			'php': 'PHP',
			'sql': 'SQL',
			'graphql': 'GraphQL',
			'xml': 'XML',
			'toml': 'TOML',
			'ini': 'INI',
			'dockerfile': 'Dockerfile',
			'docker': 'Dockerfile',
			'nginx': 'Nginx',
			'diff': 'Diff'
		};
		
		return languageNames[lang.toLowerCase()] || lang.charAt(0).toUpperCase() + lang.slice(1);
	}
	
	// Function to process code blocks
	function processCodeBlocks(html: string): string {
		if (!highlighter) {
			return html.replace(/<pre><code class="language-(\w+)">([\s\S]*?)<\/code><\/pre>/g, 
				(match, language, code) => {
					const friendlyName = getFriendlyLanguageName(language);
					return `<div class="code-block-wrapper">
						<div class="code-block-header">
							<span class="code-language-label">${friendlyName}</span>
							<button class="copy-button" aria-label="Copy code to clipboard">
								<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
									<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
									<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
								</svg>
							</button>
						</div>
						<pre class="language-${language}"><code class="language-${language}">${code}</code></pre>
					</div>`;
				}
			);
		}
		
		const hl = highlighter; // Create a non-null reference
		
		return html.replace(/<pre><code class="language-(\w+)">([\s\S]*?)<\/code><\/pre>/g, 
			(match, language, code) => {
				try {
					// Decode HTML entities
					const decodedCode = code
						.replace(/&lt;/g, '<')
						.replace(/&gt;/g, '>')
						.replace(/&amp;/g, '&')
						.replace(/&quot;/g, '"')
						.replace(/&#39;/g, "'");
					
					// Map language to Shiki language
					const shikiLang = mapLanguage(language);
					const friendlyName = getFriendlyLanguageName(language);
					
					// Check if language is supported
					const loadedLanguages = hl.getLoadedLanguages();
					if (!loadedLanguages.includes(shikiLang as any)) {
						return `<div class="code-block-wrapper">
							<div class="code-block-header">
								<span class="code-language-label">${friendlyName}</span>
								<button class="copy-button" aria-label="Copy code to clipboard">
									<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
										<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
										<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
									</svg>
								</button>
							</div>
							<pre class="language-${language}"><code class="language-${language}">${code}</code></pre>
						</div>`;
					}
					
					// Highlight the code
					const highlighted = hl.codeToHtml(decodedCode, {
						lang: shikiLang as any,
						theme: 'github-light'
					});
					
					// Add original language class for consistency and wrap in a container for the copy button
					const withClass = highlighted.replace('<pre class="shiki', `<pre class="shiki language-${language}`);
					
					// Add a copy button wrapper and language label
					return `<div class="code-block-wrapper">
						<div class="code-block-header">
							<span class="code-language-label">${friendlyName}</span>
							<button class="copy-button" aria-label="Copy code to clipboard">
								<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
									<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
									<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
								</svg>
							</button>
						</div>
						${withClass}
					</div>`;
				} catch (error) {
					// Critical error - code highlighting failed
					const friendlyName = getFriendlyLanguageName(language);
					return `<div class="code-block-wrapper">
						<div class="code-block-header">
							<span class="code-language-label">${friendlyName}</span>
							<button class="copy-button" aria-label="Copy code to clipboard">
								<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
									<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
									<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
								</svg>
							</button>
						</div>
						<pre class="language-${language}"><code class="language-${language}">${code}</code></pre>
					</div>`;
				}
			}
		);
	}
	
	// Function to add IDs to headings in the DOM
	function addIdsToHeadingsInDOM() {
		const contentElement = document.querySelector('.markdown-content');
		if (!contentElement) {
			return;
		}
		
		const headings = contentElement.querySelectorAll('h1, h2, h3, h4, h5, h6');
		
		headings.forEach(heading => {
			if (!heading.id) {
				heading.id = generateSlug(heading.textContent || '');
			}
		});
	}
	
	// Function to add copy functionality to code blocks
	function setupCodeBlockCopyButtons() {
		if (typeof window === 'undefined') return; // Skip during SSR
		
		// Find all copy buttons
		const copyButtons = document.querySelectorAll('.copy-button');
		
		copyButtons.forEach(button => {
			// Remove existing event listeners if any
			const newButton = button.cloneNode(true);
			button.parentNode?.replaceChild(newButton, button);
			
			newButton.addEventListener('click', async (event) => {
				// Prevent default browser behavior
				event.preventDefault();
				event.stopPropagation();
				
				// Find the code element
				const wrapper = (newButton as Element).closest('.code-block-wrapper');
				if (!wrapper) return;
				
				// Get the code content
				const pre = wrapper.querySelector('pre');
				if (!pre) return;
				
				// Extract text content, preserving line breaks but removing extra whitespace
				let code = '';
				const codeLines = pre.querySelectorAll('.line');
				
				if (codeLines.length) {
					// If using Shiki with line divs
					code = Array.from(codeLines)
						.map(line => (line as HTMLElement).textContent || '')
						.join('\n');
				} else {
					// Fallback to getting all text
					code = pre.textContent || '';
				}
				
				try {
					await navigator.clipboard.writeText(code);
					
					// Visual feedback
					const button = newButton as HTMLElement;
					const originalInnerHTML = button.innerHTML;
					button.innerHTML = `
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<polyline points="20 6 9 17 4 12"></polyline>
						</svg>
					`;
					button.classList.add('copied');
					
					// Reset after 2 seconds
					setTimeout(() => {
						button.innerHTML = originalInnerHTML;
						button.classList.remove('copied');
					}, 2000);
				} catch (err) {
					// Silent fail - clipboard operations may not be available
				}
				
				// Return false to prevent any default behavior
				return false;
			});
		});
	}
	
	// Function to setup hover effects for code blocks
	function setupCodeBlockHoverEffects() {
		if (typeof window === 'undefined') return; // Skip during SSR
		
		// Find all code block wrappers
		const codeBlocks = document.querySelectorAll('.code-block-wrapper');
		const tocContainer = document.getElementById('toc-container');
		
		if (!tocContainer) return;
		
		// Store the original z-index of the TOC
		const originalTocZIndex = window.getComputedStyle(tocContainer).zIndex;
		
		codeBlocks.forEach(block => {
			block.addEventListener('mouseenter', () => {
				// Set the code block to a very high z-index
				(block as HTMLElement).style.zIndex = '9999';
				// Lower the TOC z-index
				tocContainer.style.zIndex = '1';
			});
			
			block.addEventListener('mouseleave', () => {
				// Reset the z-index values
				(block as HTMLElement).style.zIndex = '';
				tocContainer.style.zIndex = originalTocZIndex;
			});
		});
	}
	
	// Function to initialize tabbed code blocks
	function initTabbedCodeBlocks() {
		if (typeof window === 'undefined') return; // Skip during SSR
		
		// Find all tabbed code blocks
		const tabbedCodeBlocks = document.querySelectorAll('.tabbed-code');
		
		tabbedCodeBlocks.forEach((block, blockIndex) => {
			// Set an ID if not present
			if (!block.id) {
				block.id = `tabbed-code-${blockIndex}`;
			}
			
			const tabHeader = block.querySelector('.tab-header');
			const tabContent = block.querySelector('.tab-content');
			
			if (!tabHeader || !tabContent) return;
			
			const tabButtons = tabHeader.querySelectorAll('.tab-button');
			const tabPanes = tabContent.querySelectorAll('.tab-pane');
			
			// Add click event listeners to tab buttons
			tabButtons.forEach((button) => {
				// Remove existing event listeners if any
				const newButton = button.cloneNode(true);
				button.parentNode?.replaceChild(newButton, button);
				
				newButton.addEventListener('click', () => {
					// Get the tab ID
					const tabId = (newButton as HTMLElement).getAttribute('data-tab');
					
					// Remove active class from all buttons and panes
					tabButtons.forEach(btn => btn.classList.remove('active'));
					tabPanes.forEach(pane => pane.classList.remove('active'));
					
					// Add active class to current button and pane
					(newButton as HTMLElement).classList.add('active');
					const targetPane = tabContent.querySelector(`#${tabId}`);
					if (targetPane) {
						targetPane.classList.add('active');
					}
					
					// Save the active tab in localStorage
					try {
						localStorage.setItem(`tabbed-code-${block.id}`, tabId || '');
					} catch (e) {
						// Silent fail - localStorage may not be available
					}
				});
			});
			
			// Try to restore active tab from localStorage
			try {
				const savedTabId = localStorage.getItem(`tabbed-code-${block.id}`);
				if (savedTabId) {
					const savedButton = tabHeader.querySelector(`[data-tab="${savedTabId}"]`);
					if (savedButton) {
						(savedButton as HTMLElement).click();
					}
				}
			} catch (e) {
				// Silent fail - localStorage may not be available
			}
		});
	}
	
	// Handle before navigation events
	beforeNavigate(({ from, to }) => {
		if (!from || !to) return;
		
		isNavigating = true;
		isLoading = true;
		
		// Store previous path for direction calculation
		previousPath = from.url.pathname || '';
		
		// Determine slide direction based on navigation order
		slideDirection = determineSlideDirection(from.url.pathname, to.url.pathname);
	});
	
	// Handle navigation events
	afterNavigate(() => {
		processContent();
		
		// After content is rendered, add IDs to headings in the DOM
		setTimeout(addIdsToHeadingsInDOM, 100);
		
		// Setup copy buttons after content is rendered
		setTimeout(() => {
			setupCodeBlockCopyButtons();
		}, 200);
		
		// Setup hover effects for code blocks
		setTimeout(() => {
			setupCodeBlockHoverEffects();
		}, 250);
		
		// Initialize tabbed code blocks after content is rendered
		setTimeout(initTabbedCodeBlocks, 300);
		
		// Scroll to top on navigation
		window.scrollTo({ top: 0, behavior: 'smooth' });
	});
	
	// Function to handle TOC content status
	function handleTocContentStatus(hasContent: boolean) {
		tocHasContent = hasContent;
	}
</script>

<svelte:head>
	<title>{title} - Chadburn</title>
</svelte:head>

<!-- Add a wrapper div with fixed dimensions to prevent layout shifts -->
<div class="page-wrapper">
	<div class="content-container">
		<!-- Always render a placeholder for the TOC to maintain consistent layout -->
		<div class="toc-container" id="toc-container" class:hidden={!tocHasContent} class:placeholder={isLoading}>
			{#if !isLoading}
				<TableOfContents 
					contentSelector=".markdown-content" 
					title="On this page" 
					on:contentStatus={e => handleTocContentStatus(e.detail.hasContent)}
				/>
			{:else}
				<!-- Empty placeholder to maintain layout -->
				<div class="toc-placeholder"></div>
			{/if}
		</div>

		<!-- Content area with consistent width -->
		<div class="content-area">
			{#if isNavigating}
				<!-- Outgoing content with appropriate transition -->
				{#if useSlideTransition}
					<div class="markdown-content" 
						out:fly={{ 
							x: slideDirection === 'right' ? -300 : 300, 
							y: 0,
							duration: 300, 
							easing: cubicOut 
						}}>
						{@html htmlContent}
					</div>
				{:else}
					<div class="markdown-content" 
						out:fly={{ 
							x: 0, 
							y: -20,
							duration: 300, 
							easing: cubicOut 
						}}>
						{@html htmlContent}
					</div>
				{/if}
			{:else if currentHtmlContent && !isLoading}
				<!-- Incoming content with appropriate transition -->
				{#if useSlideTransition}
					<div class="markdown-content" 
						in:fly={{ 
							x: slideDirection === 'right' ? 300 : -300, 
							y: 0, 
							duration: 300, 
							delay: 150, 
							easing: cubicOut 
						}} 
						out:fade={{ duration: 200 }}>
						{@html currentHtmlContent}
					</div>
				{:else}
					<div class="markdown-content" 
						in:fly={{ 
							x: 0, 
							y: 20, 
							duration: 300, 
							delay: 150, 
							easing: cubicOut 
						}} 
						out:fade={{ duration: 200 }}>
						{@html currentHtmlContent}
					</div>
				{/if}
			{:else}
				<!-- Loading placeholder with similar dimensions to maintain layout -->
				<div class="markdown-content loading-placeholder">
					<div class="loading" in:fade={{ duration: 200 }}>
						<div class="loading-spinner"></div>
						<p>Loading content...</p>
					</div>
					<!-- Placeholder elements to maintain layout -->
					<div class="placeholder-heading"></div>
					<div class="placeholder-paragraph"></div>
					<div class="placeholder-paragraph"></div>
					<div class="placeholder-heading-small"></div>
					<div class="placeholder-paragraph"></div>
					<div class="placeholder-paragraph"></div>
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	/* Wrapper to maintain consistent page width */
	.page-wrapper {
		width: 100%;
		max-width: 100%;
		margin: 0 auto;
	}
	
	.content-container {
		position: relative;
		width: 100%;
		max-width: 100%;
		padding-top: 0.5rem;
		display: flex;
		flex-direction: row-reverse;
		align-items: flex-start;
		gap: 2rem;
	}
	
	/* Content area with fixed width to prevent layout shifts */
	.content-area {
		flex: 1;
		min-width: calc(100% - 250px); /* Account for TOC width (220px) + gap (30px) */
		max-width: 100%;
	}
	
	.toc-container {
		position: sticky;
		top: var(--header-height, 60px);
		width: 220px;
		flex-shrink: 0;
		flex-grow: 0;
		z-index: 5;
		pointer-events: auto;
		background-color: transparent;
		border-radius: 0;
		padding: 0;
		margin-top: 0;
		box-shadow: none;
		border: none;
		max-height: none !important;
		overflow: visible !important;
	}
	
	/* TOC placeholder for consistent layout during loading */
	.toc-container.placeholder {
		opacity: 0.4;
		pointer-events: none;
	}
	
	.toc-placeholder {
		height: 200px;
		background-color: #f6f8fa;
		border-radius: 6px;
		opacity: 0.3;
	}
	
	/* Override TableOfContents component styles */
	.toc-container :global(.table-of-contents) {
		padding: 0.75rem 1rem;
		border: 1px solid rgba(0, 0, 0, 0.05);
		background-color: rgba(255, 255, 255, 0.95);
		border-radius: 6px;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
		overflow: visible !important; /* Override the scrollbar */
		max-height: none !important; /* Remove max-height constraint */
		margin-top: 0; /* Ensure no top margin */
	}
	
	/* Ensure the TOC nav doesn't scroll either */
	.toc-container :global(.table-of-contents nav),
	.toc-container :global(.table-of-contents ul),
	.toc-container :global(.table-of-contents .toc-list) {
		overflow: visible !important;
		max-height: none !important;
	}
	
	.markdown-content {
		position: relative;
		line-height: 1.6;
		overflow-wrap: break-word;
		word-wrap: break-word;
		will-change: transform, opacity;
		flex: 1;
		min-width: 0; /* Prevent flex items from overflowing */
	}
	
	/* Remove float-related styles since we're using flexbox now */
	.markdown-content :global(h1:first-child) {
		margin-top: 0;
		padding-top: 0;
	}
	
	/* Ensure all block elements respect the container width */
	.markdown-content :global(h1),
	.markdown-content :global(h2),
	.markdown-content :global(h3),
	.markdown-content :global(h4),
	.markdown-content :global(h5),
	.markdown-content :global(h6),
	.markdown-content :global(p),
	.markdown-content :global(ul),
	.markdown-content :global(ol),
	.markdown-content :global(blockquote),
	.markdown-content :global(pre),
	.markdown-content :global(table),
	.markdown-content :global(hr) {
		width: auto;
		max-width: 100%;
		box-sizing: border-box;
		position: relative;
		z-index: 1;
		transition: z-index 0s 0.3s, background-color 0.3s ease;
	}
	
	/* Hover effect to bring elements to the foreground */
	.markdown-content :global(h1:hover),
	.markdown-content :global(h2:hover),
	.markdown-content :global(h3:hover),
	.markdown-content :global(h4:hover),
	.markdown-content :global(h5:hover),
	.markdown-content :global(h6:hover),
	.markdown-content :global(p:hover),
	.markdown-content :global(ul:hover),
	.markdown-content :global(ol:hover),
	.markdown-content :global(blockquote:hover),
	.markdown-content :global(table:hover),
	.markdown-content :global(hr:hover) {
		z-index: 10;
		position: relative;
		transition: z-index 0s 0.3s, background-color 0.3s ease;
		background-color: rgba(255, 255, 255, 0.8);
	}
	
	.markdown-content :global(h1) {
		font-size: 2.2rem;
		margin-top: 0;
		margin-bottom: 1.5rem;
	}
	
	.markdown-content :global(h2) {
		font-size: 1.8rem;
		margin-top: 2rem;
		margin-bottom: 1rem;
		padding-bottom: 0.3rem;
		border-bottom: 1px solid #eaecef;
	}
	
	.markdown-content :global(h3) {
		font-size: 1.5rem;
		margin-top: 1.5rem;
		margin-bottom: 0.75rem;
	}
	
	.markdown-content :global(h4) {
		font-size: 1.25rem;
		margin-top: 1.25rem;
		margin-bottom: 0.5rem;
	}
	
	.markdown-content :global(p) {
		margin-top: 0;
		margin-bottom: 1rem;
	}
	
	.markdown-content :global(ul), 
	.markdown-content :global(ol) {
		margin-top: 0;
		margin-bottom: 1rem;
		padding-left: 2rem;
	}
	
	.markdown-content :global(li) {
		margin-bottom: 0.25rem;
	}
	
	.markdown-content :global(code) {
		padding: 0.2em 0.4em;
		background-color: rgba(27, 31, 35, 0.05);
		border-radius: 3px;
		font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
		font-size: 85%;
	}
	
	.markdown-content :global(pre) {
		padding: 1rem;
		overflow: auto;
		font-size: 85%;
		line-height: 1.45;
		background-color: #f6f8fa;
		border-radius: 3px;
		margin-top: 0;
		margin-bottom: 1rem;
	}
	
	.markdown-content :global(pre code) {
		padding: 0;
		background-color: transparent;
		border-radius: 0;
		font-size: 100%;
	}
	
	/* Shiki syntax highlighting styles */
	.markdown-content :global(.shiki) {
		background-color: #f6f8fa !important;
		padding: 0 !important;
		margin: 0 !important;
		font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace !important;
		font-size: 100% !important;
		line-height: 1.45 !important;
		overflow: visible !important;
		tab-size: 2 !important;
	}
	
	.markdown-content :global(pre.shiki) {
		overflow: auto !important;
	}
	
	.markdown-content :global(.shiki .line) {
		min-height: 1.45em !important;
	}
	
	.markdown-content :global(.shiki .line::before) {
		content: '';
	}
	
	.markdown-content :global(blockquote) {
		margin: 0 0 1rem;
		padding: 0 1rem;
		color: #6a737d;
		border-left: 0.25rem solid #dfe2e5;
	}
	
	.markdown-content :global(hr) {
		width: auto;
		margin: 1.5rem 0;
		border: 0;
		border-top: 1px solid #eaecef;
		overflow: visible;
		display: block;
	}
	
	.markdown-content :global(table) {
		border-collapse: collapse;
		width: 100%;
		margin-bottom: 1rem;
	}
	
	.markdown-content :global(th), 
	.markdown-content :global(td) {
		padding: 0.5rem;
		border: 1px solid #dfe2e5;
	}
	
	.markdown-content :global(th) {
		background-color: #f6f8fa;
		font-weight: 600;
	}
	
	.error {
		color: #721c24;
		background-color: #f8d7da;
		padding: 1rem;
		border-radius: 4px;
		margin-bottom: 1rem;
	}
	
	.loading {
		color: #6a737d;
		font-style: italic;
		text-align: center;
		padding: 4rem 0;
		min-height: 200px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		grid-area: content;
	}
	
	.loading-spinner {
		display: inline-block;
		width: 40px;
		height: 40px;
		border: 4px solid rgba(0, 0, 0, 0.1);
		border-radius: 50%;
		border-top-color: var(--primary-color, #0366d6);
		animation: spin 1s ease-in-out infinite;
		margin: 1rem auto;
	}
	
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	
	@media (max-width: 1024px) {
		.content-container {
			flex-direction: column;
			gap: 1rem;
		}
		
		.toc-container {
			width: 100%;
			margin-bottom: 1.5rem;
			position: relative;
			top: 0;
			max-height: none;
			overflow-y: visible;
		}
		
		.toc-container :global(.table-of-contents) {
			border-bottom: 1px solid rgba(0, 0, 0, 0.1);
			border-radius: 0;
			padding: 0.5rem 0;
			box-shadow: none;
			background-color: transparent;
			border: none;
		}
		
		.markdown-content :global(.code-block-wrapper) {
			max-width: 100%;
		}
	}
	
	/* Code block wrapper and copy button styles */
	.markdown-content :global(.code-block-wrapper) {
		position: relative;
		margin: 1.5rem 0;
		border-radius: 6px;
		overflow: hidden;
		box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
		border: 1px solid #e1e4e8;
		transition: box-shadow 0.3s ease, transform 0.2s ease;
		z-index: 1;
		width: 100%;
		max-width: 100%; /* Allow code blocks to use full width since we're using flexbox */
	}
	
	/* Hover effect with transform instead of position changes */
	.markdown-content :global(.code-block-wrapper:hover) {
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
		z-index: 9999;
		transform: translateZ(1px); /* Force a new stacking context */
		overflow: visible;
	}
	
	/* Remove the placeholder that was causing layout shifts */
	.markdown-content :global(.code-block-wrapper::before) {
		content: none;
	}
	
	.markdown-content :global(.code-block-wrapper:hover::before) {
		content: none;
	}
	
	.markdown-content :global(.code-block-header) {
		position: relative;
		z-index: 2;
	}
	
	.markdown-content :global(.code-block-wrapper:hover .copy-button) {
		opacity: 1;
	}
	
	.markdown-content :global(.code-language-label) {
		color: #6a737d;
		font-weight: 500;
		user-select: none;
	}
	
	.markdown-content :global(.copy-button) {
		position: relative;
		z-index: 10;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 28px;
		height: 28px;
		padding: 0;
		background-color: transparent;
		border: none;
		border-radius: 4px;
		opacity: 0.7;
		transition: opacity 0.2s ease, background-color 0.2s ease;
		cursor: pointer;
		color: #6a737d;
		/* Ensure the button doesn't trigger form submission */
		type: button;
	}
	
	.markdown-content :global(.copy-button:hover) {
		opacity: 1;
		background-color: rgba(0, 0, 0, 0.05);
		color: #333;
	}
	
	.markdown-content :global(.copy-button:active) {
		background-color: rgba(0, 0, 0, 0.1);
	}
	
	.markdown-content :global(.copy-button.copied) {
		color: #28a745;
		opacity: 1;
	}
	
	.markdown-content :global(.code-block-wrapper pre) {
		margin: 0;
		border-radius: 0;
		border: none;
		padding: 1.25rem;
	}
	
	.markdown-content :global(.code-block-wrapper pre.shiki) {
		margin: 0;
		border-radius: 0;
		border: none;
		padding: 1rem 0;
	}
	
	.markdown-content :global(.code-block-wrapper .line) {
		padding: 0 1.25rem !important;
	}
	
	.markdown-content :global(.code-block-wrapper pre code) {
		padding: 0;
		background-color: transparent;
	}
	
	:global(.highlight-section) {
		animation: highlight-fade 2s ease-out;
	}
	
	@keyframes highlight-fade {
		0% {
			background-color: rgba(255, 230, 0, 0.2);
			box-shadow: 0 0 0 4px rgba(255, 230, 0, 0.2);
		}
		100% {
			background-color: transparent;
			box-shadow: none;
		}
	}
	
	.markdown-content :global(a),
	.markdown-content :global(button),
	.markdown-content :global(.code-block-wrapper),
	.markdown-content :global(h1),
	.markdown-content :global(h2),
	.markdown-content :global(h3),
	.markdown-content :global(h4),
	.markdown-content :global(h5),
	.markdown-content :global(h6) {
		transition-delay: 0.2s;
	}
	
	.markdown-content :global(.code-block-wrapper:hover .copy-button) {
		opacity: 1;
	}
	
	.markdown-content :global(.code-block-header) {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1.25rem;
		background-color: #f6f8fa;
		border-bottom: 1px solid #e1e4e8;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
		font-size: 0.85rem;
	}
	
	/* Add a specific style to ensure code blocks can overlay the TOC */
	@media (min-width: 1025px) {
		.markdown-content {
			/* Create a new stacking context with higher z-index than TOC when needed */
			isolation: isolate;
		}
		
		.markdown-content :global(.code-block-wrapper:hover) {
			/* Only apply these styles on larger screens where TOC is visible */
			position: relative;
			overflow: visible;
		}
	}
	
	.toc-container.hidden {
		display: none;
	}
	
	.loading-placeholder {
		min-height: 70vh;
		position: relative;
	}
	
	.loading {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		color: #6a737d;
		font-style: italic;
		text-align: center;
		padding: 4rem 0;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		z-index: 10;
		background-color: rgba(255, 255, 255, 0.7);
	}
	
	/* Placeholder elements to maintain layout during loading */
	.placeholder-heading {
		height: 2.5rem;
		margin-bottom: 1.5rem;
		background-color: #f6f8fa;
		border-radius: 4px;
		opacity: 0.5;
		width: 70%;
	}
	
	.placeholder-heading-small {
		height: 1.8rem;
		margin-top: 2rem;
		margin-bottom: 1rem;
		background-color: #f6f8fa;
		border-radius: 4px;
		opacity: 0.5;
		width: 50%;
	}
	
	.placeholder-paragraph {
		height: 1rem;
		margin-bottom: 1rem;
		background-color: #f6f8fa;
		border-radius: 4px;
		opacity: 0.5;
		width: 100%;
	}
</style> 