import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue')
  },
  {
    path: '/',
    redirect: '/subscription'
  },
  {
    path: '/subscription',
    name: 'Subscription',
    component: () => import('@/views/subcription/index.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
