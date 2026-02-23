<script setup lang="ts">
const items = [{ id: 1, code: 'V1', name: 'Voyage 1' }]
</script>

<template>
  <div class="page">
    <h1>Voyages</h1>
    <table>
      <tbody>
        <tr v-for="r in items" :key="r.id">
          <td>{{ r.code }}</td>
          <td>{{ r.name }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
