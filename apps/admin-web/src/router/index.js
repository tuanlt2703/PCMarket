import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AdminLayout from '@/layouts/AdminLayout.vue'

const routes = [
  { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue'), meta: { public: true, title: 'Sign in' } },
  {
    path: '/',
    component: AdminLayout,
    children: [
      { path: '', name: 'dashboard', component: () => import('@/views/DashboardView.vue'), meta: { title: 'Dashboard' } },

      { path: 'products', name: 'products', component: () => import('@/views/products/ProductListView.vue'), meta: { title: 'Products' } },
      { path: 'products/new', name: 'product-create', component: () => import('@/views/products/ProductFormView.vue'), meta: { title: 'New Product', parent: 'products' } },
      { path: 'products/:id(\\d+)', name: 'product-detail', component: () => import('@/views/products/ProductDetailView.vue'), meta: { title: 'Product Detail', parent: 'products' } },
      { path: 'products/:id(\\d+)/edit', name: 'product-edit', component: () => import('@/views/products/ProductFormView.vue'), meta: { title: 'Edit Product', parent: 'products' } },

      { path: 'lookups', name: 'lookups', component: () => import('@/views/lookups/LookupIndexView.vue'), meta: { title: 'Lookup Data' } },
      { path: 'lookups/:table', name: 'lookup-table', component: () => import('@/views/lookups/LookupTableView.vue'), meta: { title: 'Lookup Table', parent: 'lookups' } },

      { path: 'inventory', name: 'inventory', component: () => import('@/views/inventory/InventoryView.vue'), meta: { title: 'Stock' } },
      { path: 'inventory/movements', name: 'stock-movements', component: () => import('@/views/inventory/StockMovementsView.vue'), meta: { title: 'Stock Movements', parent: 'inventory' } },

      { path: 'orders', name: 'orders', component: () => import('@/views/orders/OrderListView.vue'), meta: { title: 'Orders' } },
      { path: 'orders/:id(\\d+)', name: 'order-detail', component: () => import('@/views/orders/OrderDetailView.vue'), meta: { title: 'Order Detail', parent: 'orders' } },

      { path: 'customers', name: 'customers', component: () => import('@/views/customers/CustomerListView.vue'), meta: { title: 'Customers' } },
      { path: 'customers/:id(\\d+)', name: 'customer-detail', component: () => import('@/views/customers/CustomerDetailView.vue'), meta: { title: 'Customer Detail', parent: 'customers' } },

      { path: 'warranty', name: 'warranty', component: () => import('@/views/warranty/WarrantyListView.vue'), meta: { title: 'Warranty Claims' } },

      { path: ':pathMatch(.*)*', name: 'not-found', component: () => import('@/views/NotFoundView.vue'), meta: { title: 'Not Found' } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isAuthenticated) return { name: 'login', query: { redirect: to.fullPath } }
  if (to.name === 'login' && auth.isAuthenticated) return { name: 'dashboard' }
})

router.afterEach((to) => {
  document.title = `${to.meta.title ?? 'Admin'} | PCMarket Admin`
})

export default router
