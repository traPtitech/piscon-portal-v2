<script setup lang="ts">
import { useId } from 'vue'

const model = defineModel()
const id = useId()
</script>

<template>
  <div class="checkbox-wrapper">
    <input :id="id" type="checkbox" class="switch" v-model="model" />
    <label :for="id">
      <slot />
    </label>
  </div>
</template>

<style>
.checkbox-wrapper {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--ct-slate-600);
}

.checkbox-wrapper .switch {
  appearance: none;
  background-color: var(--ct-slate-300);
  border-radius: 72px;
  border-style: none;
  flex-shrink: 0;
  height: 20px;
  margin: 0;
  position: relative;
  width: 30px;
  cursor: default;
}

.checkbox-wrapper .switch::before {
  bottom: -6px;
  content: '';
  left: -6px;
  position: absolute;
  right: -6px;
  top: -6px;
}

.checkbox-wrapper .switch,
.checkbox-wrapper .switch::after {
  transition: all 200ms;
}

.checkbox-wrapper .switch::after {
  background-color: #fff;
  border-radius: 50%;
  content: '';
  height: 14px;
  left: 3px;
  position: absolute;
  top: 3px;
  width: 14px;
}

.checkbox-wrapper .switch:hover {
  background-color: var(--ct-slate-400);
}

.checkbox-wrapper .switch:checked {
  background-color: var(--color-primary);
}

.checkbox-wrapper .switch:checked::after {
  background-color: #fff;
  left: 13px;
}

.checkbox-wrapper :focus:not(.focus-visible) {
  outline: 0;
}

.checkbox-wrapper .switch:checked:hover {
  background-color: var(--color-primary-hover-bg);
}
</style>
