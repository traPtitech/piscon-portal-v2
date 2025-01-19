<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'
import { Icon } from '@iconify/vue'
import { useMe } from '@/lib/useServerData'
const route = useRoute()

const { data: me } = useMe()

const links = [
  { icon: 'mdi:crown', name: '順位表', path: '/' },
  { icon: 'mdi:thunder', name: 'ベンチマーク', path: '/benches' },
  { icon: 'mdi:server-network', name: 'インスタンス', path: '/instances' },
  { icon: 'mdi:account-group', name: 'チーム管理', path: '/team' },
  { icon: 'mdi:drive-document', name: 'ドキュメント', path: '/docs' },
]

const adminLinks = [
  { icon: 'mdi:thunder-circle', name: 'ベンチマーク', path: '/admin/benches' },
  { icon: 'mdi:database-cog', name: 'インスタンス', path: '/admin/instances' },
  { icon: 'mdi:account-cog', name: 'チーム管理', path: '/admin/teams' },
  { icon: 'mdi:text-box-edit', name: 'ドキュメント', path: '/admin/docs' },
  { icon: 'mdi:account-lock', name: '権限管理', path: '/admin/permissions' },
]

const isActive = (link: string) =>
  link !== '/' ? route.path.startsWith(link) : route.path === link
</script>

<template>
  <div class="main-layout-container">
    <nav class="nav-container">
      <div class="logo">PISCON</div>
      <div class="nav-links">
        <RouterLink
          v-for="link in links"
          class="nav-link"
          :class="{ active: isActive(link.path) }"
          :key="link.path"
          :to="link.path"
        >
          <Icon :icon="link.icon" width="24" height="24" />
          <span>
            {{ link.name }}
          </span>
        </RouterLink>
        <div class="for-admins-label">管理者向け</div>
        <template v-if="me?.isAdmin">
          <RouterLink
            v-for="link in adminLinks"
            class="nav-link"
            :class="{ active: isActive(link.path) }"
            :key="link.path"
            :to="link.path"
          >
            <Icon :icon="link.icon" width="24" height="24" />
            <span>
              {{ link.name }}
            </span>
          </RouterLink>
        </template>
      </div>
    </nav>
    <div class="main-content">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.main-layout-container {
  display: grid;
  grid-template-columns: 200px 1fr;
  height: 100vh;
}

.nav-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;
  background-color: var(--color-primary-background);
  width: 200px;
  height: 100%;
  overflow-y: auto;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
  text-align: center;
}

.nav-links {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-weight: 500;
  color: var(--ct-slate-600);
  text-decoration: none;
  padding: 0.5rem;
  border-radius: 0.25rem;
  transition: background-color 0.2s;
}

.nav-link.active {
  background-color: rgba(var(--ct-slate-900-rgb), 0.15);
  font-weight: 700;
  color: var(--ct-slate-800);
  opacity: 1;
}

.main-content {
  overflow-y: auto;
  min-width: 0;
}

.for-admins-label {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--ct-slate-500);
  margin-top: 2rem;
  display: flex;
  align-items: center;
}
.for-admins-label::before,
.for-admins-label::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid var(--ct-slate-500);
}
.for-admins-label::before {
  margin-right: 0.5rem;
}
.for-admins-label::after {
  margin-left: 0.5rem;
}
</style>
