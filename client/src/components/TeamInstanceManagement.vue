<script setup lang="ts">
import type InstanceCardList from '@/components/InstanceCardList.vue'
import { useTeamInstances } from '@/lib/useServerData'
import { computed, ref } from 'vue'

const { teamId } = defineProps<{ teamId: string }>()

const showDeleted = ref(false)

const { data: instances, refetch } = useTeamInstances(teamId)

const visibleInstances = computed(() =>
  instances.value?.filter((i) => showDeleted.value || i.status !== 'deleted'),
)

setInterval(() => {
  const loadingStatuses = ['building', 'starting', 'stopping', 'deleting']
  if (instances.value?.some((i) => loadingStatuses.includes(i.status))) {
    refetch()
  }
}, 500)
</script>

<template>
  <div class="team-instance-management-container">
    <div class="instance-card-list-header">
      <MainSwitch v-model="showDeleted">削除済みのインスタンスも表示する</MainSwitch>
    </div>
    <InstanceCardList :teamId="teamId" :instances="visibleInstances" />
  </div>
</template>

<style scoped>
.team-instance-management-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.instance-card-list-header {
  display: flex;
  justify-content: flex-end;
}
</style>
