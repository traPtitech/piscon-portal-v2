<script setup lang="ts">
import DocsMarkdown from '@/components/DocsMarkdown.vue'
import { useDocs } from '@/lib/useServerData'

const { data: docs } = useDocs()
</script>

<template>
  <main class="docs-container">
    <PageTitle icon="mdi:drive-document">ドキュメント</PageTitle>

    <DocsMarkdown :markdown="docs?.body ?? ''" />
  </main>
</template>

<style scoped>
.docs-container {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
</style>
