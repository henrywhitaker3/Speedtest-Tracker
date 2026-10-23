<template>
  <NuxtLink
    :to="link"
    target="_blank"
    class="hover:text-primary cursor-pointer text-green-400 dark:text-green-600"
  >
    {{ version }}
  </NuxtLink>
</template>

<script setup lang="ts">
const version = useRuntimeConfig().public.version
const link = 'https://github.com/henrywhitaker3/Speedtest-Tracker/blob/v2/CHANGELOG.md'
</script>
