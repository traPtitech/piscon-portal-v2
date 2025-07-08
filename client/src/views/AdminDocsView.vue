<script setup lang="ts">
import DocsMarkdown from '@/components/DocsMarkdown.vue'
import MainButton from '@/components/MainButton.vue'
import PageTitle from '@/components/PageTitle.vue'
import { useDocs, useUpdateDocs } from '@/lib/useServerData'
import { Icon } from '@iconify/vue'
import { useConfirm } from 'primevue'
import { ref, watch } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'

const { data: docs } = useDocs()
const { mutate: updateDocs } = useUpdateDocs()
const confirm = useConfirm()

const docsValue = ref('')
watch(
  docs,
  () => {
    if (docs.value?.body) docsValue.value = docs.value?.body
  },
  { immediate: true },
)

onBeforeRouteLeave((_to, _from, next) => {
  if (docsValue.value == docs?.value?.body) {
    next()
    return
  }

  confirm.require({
    header: '変更の確認',
    message: 'ドキュメントが変更されています。このまま保存せずに移動しますか？',
    accept: () => {
      next()
    },
  })
})
</script>

<template>
  <main class="admin-docs-container">
    <PageTitle icon="mdi:text-box-edit">ドキュメント (管理者)</PageTitle>

    <div class="docs-container">
      <div class="docs-two-column-container">
        <Textarea autoResize rows="5" v-model="docsValue"></Textarea>
        <DocsMarkdown :markdown="docsValue" />
      </div>
      <MainButton @click="updateDocs({ body: docsValue })" :disabled="docsValue === docs?.body">
        <Icon icon="mdi:content-save" width="20" height="20" />
        <span>更新</span>
      </MainButton>
    </div>
  </main>

  <ConfirmDialog>
    <template #container="{ message, acceptCallback, rejectCallback }">
      <div class="confirm-dialog-container">
        <div class="confirm-dialog-header">
          {{ message.header }}
        </div>
        <div class="confirm-dialog-message">{{ message.message }}</div>
        <div class="confirm-dialog-actions">
          <MainButton
            variant="primary"
            outlined
            @click="rejectCallback"
            class="confirm-dialog-reject"
          >
            移動しない
          </MainButton>
          <MainButton variant="destructive" @click="acceptCallback" class="confirm-dialog-accept">
            保存せずに移動
          </MainButton>
        </div>
      </div>
    </template>
  </ConfirmDialog>
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

.docs-two-column-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.confirm-dialog-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;
}

.confirm-dialog-header {
  font-size: 1.25rem;
  font-weight: 600;
}

.confirm-dialog-message {
  font-size: 1rem;
  color: var(--ct-slate-700);
}

.confirm-dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}
</style>
