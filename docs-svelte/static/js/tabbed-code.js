/**
 * Tabbed Code Blocks
 * 
 * This script handles the functionality for tabbed code blocks in the documentation.
 */

// Initialize tabbed code blocks
function initTabbedCodeBlocks() {
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
      button.addEventListener('click', () => {
        // Get the tab ID
        const tabId = button.getAttribute('data-tab');
        
        // Remove active class from all buttons and panes
        tabButtons.forEach(btn => btn.classList.remove('active'));
        tabPanes.forEach(pane => pane.classList.remove('active'));
        
        // Add active class to current button and pane
        button.classList.add('active');
        const targetPane = tabContent.querySelector(`#${tabId}`);
        if (targetPane) {
          targetPane.classList.add('active');
        }
        
        // Save the active tab in localStorage
        try {
          localStorage.setItem(`tabbed-code-${block.id}`, tabId);
        } catch (e) {
          // Silently fail if localStorage is not available
        }
      });
    });
    
    // Try to restore active tab from localStorage
    try {
      const savedTabId = localStorage.getItem(`tabbed-code-${block.id}`);
      if (savedTabId) {
        const savedButton = tabHeader.querySelector(`[data-tab="${savedTabId}"]`);
        if (savedButton) {
          savedButton.click();
        }
      }
    } catch (e) {
      // Silently fail if localStorage is not available
    }
  });
}

// Initialize when the DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  initTabbedCodeBlocks();
}); 