<template>
  <nav class="app-header navbar navbar-expand bg-body">
    <div class="container-fluid">
      <ul class="navbar-nav">
        <li class="nav-item">
          <a class="nav-link" data-lte-toggle="sidebar" href="#" role="button" aria-label="Toggle sidebar">
            <i class="bi bi-list"></i>
          </a>
        </li>
        <li class="nav-item d-none d-md-block">
          <RouterLink :to="{ name: 'dashboard' }" class="nav-link">Home</RouterLink>
        </li>
        <li class="nav-item d-none d-md-block">
          <RouterLink :to="{ name: 'orders' }" class="nav-link">Orders</RouterLink>
        </li>
      </ul>

      <ul class="navbar-nav ms-auto">
        <li class="nav-item">
          <span class="nav-link">
            <span class="badge text-bg-warning" title="Data is stored in this browser only">
              <i class="bi bi-database-fill-gear me-1"></i>Mock API
            </span>
          </span>
        </li>

        <li class="nav-item">
          <a class="nav-link" href="#" data-lte-toggle="fullscreen" aria-label="Toggle fullscreen">
            <i data-lte-icon="maximize" class="bi bi-arrows-fullscreen"></i>
            <i data-lte-icon="minimize" class="bi bi-fullscreen-exit d-none"></i>
          </a>
        </li>

        <li class="nav-item dropdown">
          <a class="nav-link" href="#" data-bs-toggle="dropdown" aria-expanded="false" aria-label="Toggle color scheme">
            <i class="bi bi-sun-fill" data-lte-theme-icon="light"></i>
            <i class="bi bi-moon-fill d-none" data-lte-theme-icon="dark"></i>
            <i class="bi bi-circle-half d-none" data-lte-theme-icon="auto"></i>
          </a>
          <ul class="dropdown-menu dropdown-menu-end" style="--bs-dropdown-min-width: 8rem">
            <li v-for="t in themes" :key="t.value">
              <button type="button" class="dropdown-item d-flex align-items-center" :data-bs-theme-value="t.value">
                <i :class="['bi', t.icon, 'me-2']"></i>{{ t.label }}
                <i class="bi bi-check-lg ms-auto d-none"></i>
              </button>
            </li>
          </ul>
        </li>

        <li class="nav-item dropdown user-menu">
          <a href="#" class="nav-link dropdown-toggle" data-bs-toggle="dropdown">
            <img :src="avatar" class="user-image rounded-circle shadow" alt="User" />
            <span class="d-none d-md-inline">{{ auth.user?.full_name }}</span>
          </a>
          <ul class="dropdown-menu dropdown-menu-lg dropdown-menu-end">
            <li class="user-header text-bg-primary">
              <img :src="avatar" class="rounded-circle shadow" alt="User" />
              <p>
                {{ auth.user?.full_name }} - {{ auth.user?.role }}
                <small>@{{ auth.user?.username }}</small>
              </p>
            </li>
            <li class="user-footer">
              <button type="button" class="btn btn-outline-warning" @click="confirmReset = true">
                <i class="bi bi-arrow-counterclockwise me-1"></i>Reset mock data
              </button>
              <button type="button" class="btn btn-outline-danger float-end" @click="logout">Sign out</button>
            </li>
          </ul>
        </li>
      </ul>
    </div>
  </nav>

  <ConfirmDialog
    v-if="confirmReset"
    title="Reset mock data"
    message="All changes made in this browser (products, lookups, orders, ...) will be replaced by the initial mock data."
    confirm-text="Reset"
    variant="warning"
    @confirm="resetData"
    @cancel="confirmReset = false"
  />
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { resetDb } from '@/services/mockDb'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import avatar from '@adminlte/assets/img/user2-160x160.jpg'

const themes = [
  { value: 'light', label: 'Light', icon: 'bi-sun-fill' },
  { value: 'dark', label: 'Dark', icon: 'bi-moon-fill' },
  { value: 'auto', label: 'Auto', icon: 'bi-circle-half' },
]

const auth = useAuthStore()
const router = useRouter()
const confirmReset = ref(false)

function logout() {
  auth.logout()
  router.push({ name: 'login' })
}

function resetData() {
  resetDb()
  window.location.reload()
}
</script>
