<script setup lang="ts">
import { renderMarkdown } from '@/lib/renderMarkdown'
import { computed } from 'vue'

const { markdown } = defineProps<{ markdown: string }>()

const rendered = computed(() => renderMarkdown(markdown))
</script>

<template>
  <div v-html="rendered" class="docs-markdown-root"></div>
</template>

<style scoped>
.docs-markdown-root {
  border: 1px solid var(--ct-slate-300);
  border-radius: 4px;
  padding: 1rem;
}
</style>

<style>
.docs-markdown-root h1 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
}
.docs-markdown-root h2 {
  font-size: 1.2rem;
  margin-block: 0.5rem;
}
.docs-markdown-root p + p {
  margin-top: 1rem;
}
.docs-markdown-root a {
  color: var(--color-primary);
}
.docs-markdown-root a:hover {
  color: var(--color-primary-hover-text);
}
.docs-markdown-root ul,
.docs-markdown-root ol {
  margin-bottom: 1rem;
}
.docs-markdown-root code {
  background-color: var(--ct-slate-200);
  padding: 0.1rem 0.25rem;
  border-radius: 4px;
  font-size: 0.9rem;
}

.docs-markdown-root blockquote {
  border-left: 4px solid var(--ct-slate-300);
  padding: 0.25rem 0 0.25rem 0.5rem;
  margin: 1rem 0;
  color: var(--ct-slate-500);
}

.docs-markdown-root pre {
  background-color: var(--ct-slate-700);
  color: var(--ct-slate-100);
  border-radius: 4px;
  padding: 1rem;
  overflow-x: auto;
}
.docs-markdown-root pre code {
  background-color: inherit;
  padding: 0;
}
</style>
