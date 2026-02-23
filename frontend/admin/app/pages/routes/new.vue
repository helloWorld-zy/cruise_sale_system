<script setup lang="ts">
import { ref } from 'vue'
const form = ref({ code: '', name: '' })
</script>

<template>
  <div class="page">
    <h1>New Route</h1>
    <input v-model="form.code" placeholder="Code" />
    <input v-model="form.name" placeholder="Name" />
  </div>
</template>
