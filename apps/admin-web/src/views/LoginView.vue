<template>
  <div class="login-box">
    <div class="card card-outline card-primary">
      <div class="card-header text-center">
        <span class="link-dark text-center link-offset-2 link-opacity-100 link-opacity-50-hover">
          <h1 class="mb-0"><b>PC</b>Market</h1>
        </span>
      </div>
      <div class="card-body login-card-body">
        <p class="login-box-msg">Sign in to the store administration</p>
        <div v-if="error" class="alert alert-danger py-2">{{ error }}</div>
        <form @submit.prevent="submit">
          <div class="input-group mb-3">
            <div class="form-floating">
              <input id="loginUser" v-model.trim="username" type="text" class="form-control" placeholder="" autocomplete="username" required />
              <label for="loginUser">Username</label>
            </div>
            <div class="input-group-text"><span class="bi bi-person"></span></div>
          </div>
          <div class="input-group mb-3">
            <div class="form-floating">
              <input id="loginPassword" v-model="password" type="password" class="form-control" placeholder="" autocomplete="current-password" required />
              <label for="loginPassword">Password</label>
            </div>
            <div class="input-group-text"><span class="bi bi-lock-fill"></span></div>
          </div>
          <div class="d-grid gap-2">
            <button type="submit" class="btn btn-primary" :disabled="busy">
              <span v-if="busy" class="spinner-border spinner-border-sm me-1"></span>Sign In
            </button>
          </div>
        </form>
        <p class="text-secondary small mt-3 mb-0">
          <i class="bi bi-info-circle me-1"></i>Mock account: <code>admin</code> / <code>admin123</code>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const BODY_CLASSES = ['login-page', 'bg-body-secondary']

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const username = ref('')
const password = ref('')
const error = ref(null)
const busy = ref(false)

async function submit() {
  busy.value = true
  error.value = null
  try {
    await auth.login(username.value, password.value)
    const redirect = route.query.redirect
    // Only allow in-app paths to avoid open redirects
    const target = typeof redirect === 'string' && redirect.startsWith('/') && !redirect.startsWith('//') ? redirect : '/'
    router.replace(target)
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

onMounted(() => document.body.classList.add(...BODY_CLASSES))
onBeforeUnmount(() => document.body.classList.remove(...BODY_CLASSES))
</script>
