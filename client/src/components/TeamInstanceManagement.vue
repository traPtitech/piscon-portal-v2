<script setup lang="ts">
import { useTeamInstances } from '@/lib/useServerData'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const { teamId } = defineProps<{ teamId: string }>()

const showDeleted = ref(false)

const { data: instances, refetch } = useTeamInstances(teamId)

const visibleInstances = computed(() =>
  instances.value ? instances.value.filter((i) => showDeleted.value || i.status !== 'deleted') : [],
)

let interval: number
onMounted(() => {
  interval = setInterval(() => {
    const loadingStatuses = ['building', 'starting', 'stopping', 'deleting']
    if (instances.value?.some((i) => loadingStatuses.includes(i.status))) {
      void refetch()
    }
  }, 500)
})

onUnmounted(() => {
  clearInterval(interval)
})
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
