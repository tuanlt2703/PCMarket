import { createApp } from 'vue'
import { createPinia } from 'pinia'
import VueApexCharts from 'vue3-apexcharts'

import '@fontsource/source-sans-3/400.css'
import '@fontsource/source-sans-3/600.css'
import '@fontsource/source-sans-3/700.css'
import 'bootstrap-icons/font/bootstrap-icons.min.css'
import '@adminlte/css/adminlte.css'
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import './assets/app.css'

import App from './App.vue'
import router from './router'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(VueApexCharts)
app.mount('#app')
