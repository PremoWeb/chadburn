<script lang="ts">
  import { version } from '$lib/stores/version';
</script>

<footer>
  <div class="footer-content">
    <p>
      Chadburn
      {#if !$version.loading && $version.version}
        <span class="version">v{$version.version}</span>
      {/if}
      is licensed under the <a
        href="https://github.com/PremoWeb/Chadburn/blob/main/LICENSE"
        target="_blank"
        rel="noopener noreferrer">MIT License</a
      >.
    </p>
    <p>
      Documentation built with <a
        href="https://kit.svelte.dev/"
        target="_blank"
        rel="noopener noreferrer">SvelteKit</a
      >.
    </p>
  </div>
</footer>

<style>
  footer {
    background-color: #f5f5f5;
    padding: 1rem 0;
    border-top: 1px solid #e0e0e0;
  }

  .footer-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1rem;
    text-align: center;
    color: #666;
  }

  footer a {
    color: #155799;
    text-decoration: none;
  }

  footer a:hover {
    text-decoration: underline;
  }

  .version {
    font-size: 0.8rem;
    color: #888;
  }
</style> 