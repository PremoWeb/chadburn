import { marked } from 'marked';

// Configure marked options
marked.setOptions({
  gfm: true,
  breaks: true
});

// Function to render markdown to HTML
export function renderMarkdown(markdown: string): string {
  // Use parse as a synchronous function
  return marked.parse(markdown) as string;
}

// Function to fetch and render a markdown file
export async function fetchAndRenderMarkdown(path: string): Promise<string> {
  try {
    const response = await fetch(path);
    if (!response.ok) {
      throw new Error(`Failed to fetch markdown: ${response.status} ${response.statusText}`);
    }
    const markdown = await response.text();
    return renderMarkdown(markdown);
  } catch (error) {
    console.error('Error fetching markdown:', error);
    return `<p>Error loading content: ${error instanceof Error ? error.message : String(error)}</p>`;
  }
} 