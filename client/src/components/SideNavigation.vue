<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'
import { Icon } from '@iconify/vue'
import { useMe } from '@/lib/useServerData'
import { ref, watch } from 'vue'

const isMobile = window.matchMedia('(max-width: 560px)').matches

const route = useRoute()
const { data: me } = useMe()
const closed = ref(isMobile)

const links = [
  { icon: 'mdi:home', name: 'ホーム', path: '/' },
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

watch(closed, () => {
  setTimeout(() => {
    window.dispatchEvent(new Event('resize'))
  }, 200)
})
</script>

<template>
  <nav :class="['nav-container', { closed }]">
    <div class="nav-links">
      <RouterLink
        v-for="link in links"
        class="nav-link"
        :class="{ active: isActive(link.path) }"
        :key="link.path"
        :to="link.path"
      >
        <Icon :icon="link.icon" width="24" height="24" />
        <span v-if="!closed" class="link-label">
          {{ link.name }}
        </span>
      </RouterLink>
      <template v-if="me?.isAdmin">
        <div class="for-admins-label" v-if="!closed">管理者向け</div>
        <hr class="for-admins-border" v-else />
        <RouterLink
          v-for="link in adminLinks"
          class="nav-link"
          :class="{ active: isActive(link.path) }"
          :key="link.path"
          :to="link.path"
        >
          <Icon :icon="link.icon" width="24" height="24" />
          <span v-if="!closed" class="link-label">
            {{ link.name }}
          </span>
        </RouterLink>
      </template>
    </div>
    <hr class="nav-footer-border" />
    <div class="nav-footer">
      <button class="nav-footer-open-close-button" @click="closed = !closed">
        <Icon icon="mdi:chevron-left" width="24" height="24" class="nav-footer-button-icon" />
      </button>
    </div>
  </nav>
</template>

<style scoped>
.nav-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding-block: 1rem;
  background-color: var(--color-primary-background);
  width: 200px;
  height: 100%;
  min-height: 0;
  transition: width 0.2s;
}

.nav-container.closed {
  width: 56px;
}

.nav-links {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding-inline: 1rem;
  overflow-y: auto;
  overflow-x: hidden;
  flex: 1;
  transition: padding 0.2s;
}

.nav-container.closed .nav-links {
  padding-inline: 0.5rem;
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

.link-label {
  white-space: nowrap;
}

.for-admins-label {
  font-size: 0.9rem;
  height: 1.5rem;
  font-weight: 600;
  color: var(--ct-slate-500);
  margin-top: 2rem;
  display: flex;
  align-items: center;
  white-space: nowrap;
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

.for-admins-border {
  margin: 2.75rem 0 0.75rem 0;
  border: none;
  border-top: 1px solid var(--ct-slate-400);
}

.nav-footer-border {
  margin: 0 1rem;
  border: none;
  border-top: 1px solid var(--ct-slate-400);
}

.nav-footer {
  display: flex;
  justify-content: end;
  padding: 0 1rem;
  transition: padding 0.2s;
}

.nav-container.closed .nav-footer {
  padding: 0 0.5rem;
}

.nav-footer-button-icon {
  transition: transform 0.2s;
}

.nav-container.closed .nav-footer-button-icon {
  transform: rotate(180deg);
}

.nav-footer-open-close-button {
  background: rgba(var(--ct-slate-900-rgb), 0.05);
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 0.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ct-slate-600);
  transition: background-color 0.2s;
}

@media screen and (max-width: 560px) {
  .nav-container:not(.closed) {
    width: 100vw;
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    z-index: 1000;
  }
}
</style>
