import { delay, ApiError } from './mockDb'

// Mock identity-service. Only a single admin account exists in this first iteration.
const MOCK_ADMIN = { username: 'admin', password: 'admin123', full_name: 'Administrator', role: 'ADMIN' }

export const authApi = {
  async login(username, password) {
    await delay(300)
    if (username !== MOCK_ADMIN.username || password !== MOCK_ADMIN.password) {
      throw new ApiError(401, 'Invalid username or password')
    }
    return {
      token: `mock-${crypto.randomUUID()}`,
      user: { username: MOCK_ADMIN.username, full_name: MOCK_ADMIN.full_name, role: MOCK_ADMIN.role },
    }
  },
}
