<template>
  <aside class="app-sidebar bg-body-secondary shadow" data-bs-theme="dark">
    <div class="sidebar-brand">
      <RouterLink :to="{ name: 'dashboard' }" class="brand-link">
        <img :src="logo" alt="PCMarket Logo" class="brand-image opacity-75 shadow" />
        <span class="brand-text fw-light">PCMarket Admin</span>
      </RouterLink>
    </div>
    <div class="sidebar-wrapper">
      <nav class="mt-2" aria-label="Main navigation">
        <ul class="nav sidebar-menu flex-column" role="navigation">
          <template v-for="item in menu" :key="item.key">
            <li v-if="item.header" class="nav-header">{{ item.header }}</li>
            <li v-else class="nav-item" :class="{ 'menu-open': isOpen(item) }">
              <RouterLink
                v-if="!item.children"
                :to="item.to"
                class="nav-link"
                :class="{ active: isActive(item) }"
              >
                <i :class="['nav-icon', 'bi', item.icon]"></i>
                <p>
                  {{ item.label }}
                  <span v-if="item.badge" class="nav-badge badge text-bg-warning me-3">{{ item.badge }}</span>
                </p>
              </RouterLink>
              <template v-else>
                <a href="#" class="nav-link" :class="{ active: isActive(item) }" @click.prevent="toggle(item)">
                  <i :class="['nav-icon', 'bi', item.icon]"></i>
                  <p>
                    {{ item.label }}
                    <i class="nav-arrow bi bi-chevron-right"></i>
                  </p>
                </a>
                <ul class="nav nav-treeview">
                  <li v-for="child in item.children" :key="child.key" class="nav-item">
                    <RouterLink :to="child.to" class="nav-link" :class="{ active: isActive(child) }">
                      <i :class="['nav-icon', 'bi', child.icon ?? 'bi-circle']"></i>
                      <p>{{ child.label }}</p>
                    </RouterLink>
                  </li>
                </ul>
              </template>
            </li>
          </template>
        </ul>
      </nav>
    </div>
  </aside>
</template>

<script setup>
import { computed, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { useLookupStore } from '@/stores/lookups'
import { LOOKUP_GROUPS } from '@/config/lookupTables'
import { getSchema } from '@/config/componentSchemas'
import logo from '@adminlte/assets/img/AdminLTELogo.png'

const route = useRoute()
const lookups = useLookupStore()
const manuallyToggled = reactive({})

const menu = computed(() => [
  { key: 'dashboard', label: 'Dashboard', icon: 'bi-speedometer', to: { name: 'dashboard' }, match: ['dashboard'] },
  { key: 'h-catalog', header: 'CATALOG' },
  {
    key: 'products',
    label: 'Products',
    icon: 'bi-box-seam-fill',
    match: ['products', 'product-create', 'product-detail', 'product-edit'],
    children: [
      { key: 'p-all', label: 'All Products', icon: 'bi-grid', to: { name: 'products' }, match: ['products'], query: { type: undefined } },
      ...lookups.options('component_types').map((t) => ({
        key: `p-${t.code}`,
        label: t.name,
        icon: getSchema(t.code)?.icon ?? 'bi-box',
        to: { name: 'products', query: { type: t.code } },
        match: ['products'],
        query: { type: t.code },
      })),
    ],
  },
  {
    key: 'lookups',
    label: 'Lookup Data',
    icon: 'bi-list-ul',
    match: ['lookups', 'lookup-table'],
    children: [
      { key: 'l-all', label: 'All Tables', icon: 'bi-table', to: { name: 'lookups' }, match: ['lookups'], query: { group: undefined } },
      ...LOOKUP_GROUPS.map((g) => ({
        key: `l-${g.key}`,
        label: g.label,
        icon: g.icon,
        to: { name: 'lookups', query: { group: g.key } },
        match: ['lookups'],
        query: { group: g.key },
      })),
    ],
  },
  { key: 'h-store', header: 'STORE' },
  {
    key: 'inventory',
    label: 'Stock',
    icon: 'bi-boxes',
    match: ['inventory', 'stock-movements'],
    children: [
      { key: 'i-levels', label: 'Stock Levels', icon: 'bi-bar-chart-steps', to: { name: 'inventory' }, match: ['inventory'] },
      { key: 'i-moves', label: 'Stock Movements', icon: 'bi-arrow-left-right', to: { name: 'stock-movements' }, match: ['stock-movements'] },
    ],
  },
  { key: 'orders', label: 'Orders', icon: 'bi-cart-check-fill', to: { name: 'orders' }, match: ['orders', 'order-detail'] },
  { key: 'customers', label: 'Customers', icon: 'bi-people-fill', to: { name: 'customers' }, match: ['customers', 'customer-detail'] },
  { key: 'warranty', label: 'Warranty', icon: 'bi-shield-check', to: { name: 'warranty' }, match: ['warranty'] },
])

// `query` on an item narrows the active match (undefined = param must be absent)
function isActive(item) {
  if (!item.match?.includes(route.name)) return false
  if (!item.query) return true
  return Object.entries(item.query).every(([k, v]) => (route.query[k] ?? undefined) === v)
}

function isOpen(item) {
  if (!item.children) return false
  return manuallyToggled[item.key] ?? item.match.includes(route.name)
}

function toggle(item) {
  manuallyToggled[item.key] = !isOpen(item)
}
</script>
