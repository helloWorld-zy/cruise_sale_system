<script setup lang="ts">
import BookingRow from '../../components/BookingRow.vue'

const items = [{ id: 1, status: 'created', total: 19900 }]
</script>

<template>
  <div class="page">
    <h1>Bookings</h1>
    <BookingRow v-for="b in items" :key="b.id" :booking="b" />
  </div>
</template>
