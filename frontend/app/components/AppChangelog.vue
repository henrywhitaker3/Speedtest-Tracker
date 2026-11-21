<template>
  <div class="inline-block">
    <div class="flex items-center gap-1">
      <NuxtLink
        :to="link"
        target="_blank"
        class="hover:text-primary inline-flex cursor-pointer items-center gap-1 text-green-400 dark:text-green-600"
      >
        {{ version }}
      </NuxtLink>
      <NuxtLink
        v-if="updateAvailable"
        :to="`https://github.com/henrywhitaker3/Speedtest-Tracker/releases/tag/` + latestVersion"
        rel="noopener noreferrer"
        target="_blank"
        :aria-label="`View version ${latestVersion} release`"
        class="hover:text-primary inline-flex cursor-pointer items-center text-green-400 dark:text-green-600"
      >
        <UIcon
          name="i-lucide-circle-arrow-up"
          :title="`Version ${latestVersion} is available`"
          aria-hidden="true"
        />
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
const version = useRuntimeConfig().public.version
const link = `https://github.com/henrywhitaker3/Speedtest-Tracker/releases/tag/${version}`
const releasesUrl = 'https://api.github.com/repos/henrywhitaker3/Speedtest-Tracker/releases/latest'

const latestVersion = ref(version)
const updateAvailable = computed(() => compareVersions(latestVersion.value ?? 'dev', version) > 0)

interface GitHubRelease {
  tag_name: string
}

function compareVersions(firstVersion: string, secondVersion: string) {
  const firstParts = parseVersion(firstVersion)
  const secondParts = parseVersion(secondVersion)

  if (!firstParts || !secondParts) {
    return 0
  }

  const partCount = Math.max(firstParts.length, secondParts.length)

  for (let index = 0; index < partCount; index++) {
    const firstPart = firstParts[index] ?? 0
    const secondPart = secondParts[index] ?? 0

    if (firstPart !== secondPart) {
      return firstPart - secondPart
    }
  }

  return 0
}

function parseVersion(value: string) {
  const match = /^v?(\d+(?:\.\d+)*)$/.exec(value)

  return match?.[1]?.split('.').map(Number)
}

onMounted(async () => {
  try {
    const response = await fetch(releasesUrl)

    if (!response.ok) {
      return
    }

    const release = await response.json() as GitHubRelease
    latestVersion.value = release.tag_name
  } catch {
    return
  }
})
</script>
