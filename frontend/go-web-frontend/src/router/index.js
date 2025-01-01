import { createRouter, createWebHistory } from 'vue-router'
import UserList from '../components/UserList.vue'
import TableManager from '../components/TableManager.vue'
import RedisManager from '../components/RedisManager.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'users',
      component: UserList
    },
    {
      path: '/tables',
      name: 'tables',
      component: TableManager
    },
    {
      path: '/redis',
      name: 'redis',
      component: RedisManager
    }
  ]
})

export default router
