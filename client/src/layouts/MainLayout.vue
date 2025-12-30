<script setup lang="ts"></script>

<template>
  <div class="main-layout-container">
    <SideNavigation />
    <div class="main-content-wrapper">
      <div class="main-content">
        <slot />
      </div>
    </div>
  </div>
</template>

<style scoped>
.main-layout-container {
  display: flex;
  height: 100vh;
}

.main-content-wrapper {
  padding-top: 1rem;
  padding-bottom: 10rem;
  overflow-y: auto;
  min-width: 0;
  flex: 1;
}

.main-content {
  max-width: 1400px;
  margin: 0 auto;
}
</style>
