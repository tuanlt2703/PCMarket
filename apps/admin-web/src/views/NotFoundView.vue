<template>
  <div class="error-page text-center py-5">
    <h2 class="display-4 text-warning">404</h2>
    <h3><i class="bi bi-exclamation-triangle text-warning"></i> Page not found.</h3>
    <p>
      The page you were looking for does not exist.
      <RouterLink :to="{ name: 'dashboard' }">Return to dashboard</RouterLink>.
    </p>
  </div>
</template>
