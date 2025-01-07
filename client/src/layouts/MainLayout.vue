<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'
import { Icon } from '@iconify/vue'
const route = useRoute()

const links = [
  { icon: 'mdi:crown', name: '順位表', path: '/' },
  { icon: 'mdi:thunder', name: 'ベンチマーク', path: '/benches' },
  { icon: 'mdi:server-network', name: 'インスタンス', path: '/instances' },
  { icon: 'mdi:account-group', name: 'チーム管理', path: '/team' },
  { icon: 'mdi:drive-document', name: 'ドキュメント', path: '/docs' },
]
</script>

<template>
  <div class="main-layout-container">
    <nav class="nav-container">
      <div class="logo">PISCON</div>
      <div class="nav-links">
        <RouterLink
          v-for="link in links"
          class="nav-link"
          :class="{ active: link.path === route.path }"
          :key="link.path"
          :to="link.path"
        >
          <Icon :icon="link.icon" width="24" height="24" />
          <span>
            {{ link.name }}
          </span>
        </RouterLink>
      </div>
    </nav>
    <div>
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
</style>
