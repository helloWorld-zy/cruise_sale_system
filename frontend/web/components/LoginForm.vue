<script setup lang="ts">
import { ref } from 'vue'

const phone = ref('')
</script>

<template>
  <div class="login-form">
    <input v-model="phone" placeholder="Phone" />
  </div>
</template>
