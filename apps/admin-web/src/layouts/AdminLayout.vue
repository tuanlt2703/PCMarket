<template>
  <div class="app-wrapper">
    <AppHeader />
    <AppSidebar />
    <main class="app-main">
      <div class="app-content-header">
        <div class="container-fluid">
          <div class="row">
            <div class="col-sm-6">
              <h1 class="mb-0 fs-3">{{ pageTitle }}</h1>
            </div>
            <div class="col-sm-6">
              <nav aria-label="breadcrumb">
                <ol class="breadcrumb float-sm-end">
                  <li class="breadcrumb-item"><RouterLink :to="{ name: 'dashboard' }">Home</RouterLink></li>
                  <li v-if="parentRoute" class="breadcrumb-item">
                    <RouterLink :to="{ name: parentRoute.name }">{{ parentRoute.meta.title }}</RouterLink>
                  </li>
                  <li v-if="route.name !== 'dashboard'" class="breadcrumb-item active" aria-current="page">
                    {{ pageTitle }}
                  </li>
                </ol>
              </nav>
            </div>
          </div>
        </div>
      </div>
      <div class="app-content">
        <div class="container-fluid">
          <div v-if="!lookups.loaded" class="text-center py-5">
            <div class="spinner-border text-primary" role="status"></div>
          </div>
          <RouterView v-else :key="route.path" />
        </div>
      </div>
    </main>
    <footer class="app-footer">
      <div class="float-end d-none d-sm-inline">Mock data mode</div>
      <strong>PCMarket Admin</strong> &middot; built on
      <a href="https://adminlte.io" class="text-decoration-none" target="_blank" rel="noopener">AdminLTE 4</a>
    </footer>
  </div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, provide, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import { useLookupStore } from '@/stores/lookups'
import { initAdminLte } from '@/utils/adminlte'

const BODY_CLASSES = ['layout-fixed', 'sidebar-expand-lg', 'bg-body-tertiary']

const route = useRoute()
const router = useRouter()
const lookups = useLookupStore()

// Views can override the header title (e.g. with the product name)
const titleOverride = ref(null)
provide('setPageTitle', (t) => (titleOverride.value = t))
watch(
  () => route.fullPath,
  () => (titleOverride.value = null),
)

const pageTitle = computed(() => titleOverride.value ?? route.meta.title)
const parentRoute = computed(() =>
  route.meta.parent ? router.getRoutes().find((r) => r.name === route.meta.parent) : null,
)

onMounted(async () => {
  document.body.classList.add(...BODY_CLASSES)
  await initAdminLte()
  await lookups.load()
})

onBeforeUnmount(() => {
  document.body.classList.remove(...BODY_CLASSES, 'sidebar-open', 'sidebar-collapse')
})
</script>
