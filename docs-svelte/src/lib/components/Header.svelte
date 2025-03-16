<script lang="ts">
  import { page } from '$app/state';
  import { base } from '$app/paths';
  
  // Reactive values using runes
  let baseUrl = $derived(base);
  let currentPath = $derived(page.url.pathname);
  
  // Function to check if a path is active
  function isActive(path: string): boolean {
    return currentPath.startsWith(path);
  }
</script>

<header>
  <div class="header-content">
    <div class="logo">
      <a href="{baseUrl}/">Chadburn</a>
    </div>
    
    <nav>
      <ul>
        <li class:active={isActive('/introduction')}>
          <a href="{baseUrl}/introduction">About</a>
        </li>
        <li class:active={isActive('/installation')}>
          <a href="{baseUrl}/installation">Documentation</a>
        </li>
        <li class:active={isActive('/metrics')}>
          <a href="{baseUrl}/metrics">Metrics</a>
        </li>
        <li>
          <a href="https://github.com/PremoWeb/Chadburn" target="_blank" rel="noopener noreferrer">GitHub</a>
        </li>
      </ul>
    </nav>
  </div>
</header>

<style>
  header {
    background-color: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    position: sticky;
    top: 0;
    z-index: 1000;
    width: 100%;
  }
  
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem;
  }
  
  .logo a {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--primary-color);
    text-decoration: none;
  }
  
  nav ul {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: 1.5rem;
  }
  
  nav a {
    color: var(--text-color);
    text-decoration: none;
    font-weight: 500;
    padding: 0.5rem 0;
    position: relative;
  }
  
  nav a:hover {
    color: var(--primary-color);
  }
  
  nav li.active a {
    color: var(--primary-color);
  }
  
  nav li.active a::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 2px;
    background-color: var(--primary-color);
  }
  
  @media (max-width: 768px) {
    .header-content {
      flex-direction: column;
      gap: 1rem;
    }
    
    nav ul {
      gap: 1rem;
    }
  }
</style> 