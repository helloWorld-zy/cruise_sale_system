<script setup lang="ts">
defineProps<{ items: { id: number; code: string; name: string }[] }>()
</script>

<template>
  <table>
    <tbody>
      <tr v-for="r in items" :key="r.id">
        <td>{{ r.code }}</td>
        <td>{{ r.name }}</td>
      </tr>
    </tbody>
  </table>
</template>
