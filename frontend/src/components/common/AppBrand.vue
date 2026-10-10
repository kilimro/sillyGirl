<script setup lang="ts">
withDefaults(
  defineProps<{
    href?: string;
    label?: string;
  }>(),
  {
    href: "",
    label: "MianPro",
  },
);
</script>

<template>
  <a v-if="href" :href="href" class="app-brand">
    <span class="app-brand-mark" role="img" aria-label="MianPro Logo"></span>
    <span>{{ label }}</span>
  </a>
  <div v-else class="app-brand">
    <span class="app-brand-mark" role="img" aria-label="MianPro Logo"></span>
    <span>{{ label }}</span>
  </div>
</template>

<style scoped>
.app-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #1f2937;
  font-weight: 700;
  text-decoration: none;
}

.app-brand-mark {
  display: grid;
  flex: 0 0 auto;
  place-items: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: url("../../../logo.png") center center / contain no-repeat;
}
</style>
