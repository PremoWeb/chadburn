<script lang="ts">
  // No props or state needed for this simple banner
</script>

<div class="warning-banner">
  <div class="warning-content">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="warning-icon">
      <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
      <line x1="12" y1="9" x2="12" y2="13"></line>
      <line x1="12" y1="17" x2="12.01" y2="17"></line>
    </svg>
    <div class="message-container">
      <p>
        <strong>Development Notice:</strong> Chadburn is undergoing a major rewrite. The code is not fully tested and documentation may not reflect current functionality. Please use prior stable versions for production environments.
      </p>
      <p class="donation-text">
        <strong>Support Needed:</strong> Development costs are mounting. Please consider donating via PayPal or Venmo to <strong>nick@premoweb.com</strong> to help fund AI tools and development time.
      </p>
    </div>
  </div>
</div>

<style>
  .warning-banner {
    background-color: #fff3cd;
    border: 1px solid #ffeeba;
    color: #856404;
    width: 100%;
    box-sizing: border-box;
    margin-bottom: 1rem;
  }

  .warning-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0.75rem 1rem;
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .warning-icon {
    flex-shrink: 0;
    width: 24px;
    height: 24px;
    margin-top: 3px;
  }

  .message-container {
    flex: 1;
  }

  p {
    margin: 0;
    font-size: 0.95rem;
    line-height: 1.4;
  }

  .donation-text {
    margin-top: 0.5rem;
    font-size: 0.95rem;
    color: #721c24;
  }

  @media (max-width: 768px) {
    .warning-content {
      flex-direction: column;
      text-align: center;
      padding: 1rem;
    }
    
    .donation-text {
      margin-top: 0.75rem;
    }
  }
</style> 