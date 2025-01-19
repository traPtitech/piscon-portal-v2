<script setup lang="ts">
import MainButton from '@/components/MainButton.vue'
import PageTitle from '@/components/PageTitle.vue'
import { useDocs, useUpdateDocs } from '@/lib/useServerData'
import { Icon } from '@iconify/vue'
import { ref, watch } from 'vue'

const { data: docs } = useDocs()
const { mutate: updateDocs } = useUpdateDocs()

const docsValue = ref('')
watch(docs, () => {
  if (docs.value?.body) docsValue.value = docs.value?.body
})
</script>

<template>
  <main class="admin-docs-container">
    <PageTitle icon="mdi:text-box-edit">ドキュメント (管理者)</PageTitle>

    <div class="docs-container">
      <Textarea autoResize rows="5" v-model="docsValue"></Textarea>
      <MainButton @click="updateDocs({ body: docsValue })" :disabled="docsValue === docs?.body">
        <Icon icon="mdi:content-save" width="20" height="20" />
        <span>更新</span>
      </MainButton>
    </div>
  </main>
</template>

<style scoped>
.admin-docs-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.docs-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
}
</style>
