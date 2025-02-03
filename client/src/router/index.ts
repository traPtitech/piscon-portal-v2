import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/instances',
      name: 'instances',
      component: () => import('@/views/InstancesView.vue'),
    },
    {
      path: '/benches',
      name: 'benches',
      component: () => import('@/views/BenchesView.vue'),
    },
    {
      path: '/benches/:id',
      name: 'bench',
      component: () => import('@/views/BenchView.vue'),
    },
    {
      path: '/docs',
      name: 'docs',
      component: () => import('@/views/DocsView.vue'),
    },
    {
      path: '/team',
      name: 'team',
      component: () => import('@/views/TeamView.vue'),
    },
    {
      path: '/admin/instances',
      name: 'admin-instances',
      component: () => import('@/views/AdminInstancesView.vue'),
    },
    {
      path: '/admin/benches',
      name: 'admin-benches',
      component: () => import('@/views/AdminBenchesView.vue'),
    },
    {
      path: '/admin/benches/:id',
      name: 'admin-bench',
      component: () => import('@/views/AdminBenchView.vue'),
    },
    {
      path: '/admin/teams',
      name: 'admin-teams',
      component: () => import('@/views/AdminTeamsView.vue'),
    },
    {
      path: '/admin/docs',
      name: 'admin-docs',
      component: () => import('@/views/AdminDocsView.vue'),
    },
    {
      path: '/admin/permissions',
      name: 'admin-permissions',
      component: () => import('@/views/AdminPermissionsView.vue'),
    },
  ],
})

export default router
