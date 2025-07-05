<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { ref } from 'vue'

const { text } = defineProps<{ text: string }>()

const clicked = ref(false)

const copyToClipboard = async () => {
  try {
    if (!navigator.clipboard) {
      throw new Error('Clipboard API not supported')
    }
    await navigator.clipboard.writeText(text)
    clicked.value = true
    setTimeout(() => {
      clicked.value = false
    }, 3000)
  } catch (err) {
    console.error('Failed to copy text to clipboard:', err)
    // TODO: もっとリッチな表示にする
    alert('クリップボードへのコピーに失敗しました。')
  }
}
</script>

<template>
  <button @click="copyToClipboard" class="copy-to-clipboard-button" :class="{ clicked }">
    <Icon class="copy-icon" icon="mdi:content-copy" width="16" height="16" />
    <Icon class="copied-icon" icon="mdi:check-bold" width="16" height="16" />
  </button>
</template>

<style scoped>
.copy-to-clipboard-button {
  background-color: transparent;
  border: none;
  cursor: pointer;
  padding: 0;
  margin: 0;
  position: relative;
  width: 16px;
  height: 16px;
}
.copy-icon,
.copied-icon {
  position: absolute;
  inset: 0;
}

.copy-icon {
  opacity: 1;
}

.copied-icon {
  opacity: 0;
}

.copy-to-clipboard-button.clicked .copy-icon {
  animation: copy-icon 3s;
}

@keyframes copy-icon {
  0% {
    opacity: 1;
    transform: scale(1);
  }
  10% {
    opacity: 0;
    transform: scale(0);
  }
  90% {
    opacity: 0;
    transform: scale(0);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

.copy-to-clipboard-button.clicked .copied-icon {
  animation: copied-icon 3s;
}

@keyframes copied-icon {
  0% {
    opacity: 0;
    transform: scale(0);
  }
  10% {
    opacity: 1;
    transform: scale(1);
  }
  90% {
    opacity: 1;
    transform: scale(1);
  }
  100% {
    opacity: 0;
    transform: scale(0);
  }
}
</style>
