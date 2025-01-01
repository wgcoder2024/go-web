<template>
  <div>
    <header>
      <h1>用户管理系统</h1>
    </header>
    <UserList />
  </div>
</template>

<script setup>
import UserList from './components/UserList.vue'
</script>

<style>
body {
  margin: 0;
  font-family: Arial, sans-serif;
}

header {
  background: #2c3e50;
  color: white;
  padding: 1rem;
  text-align: center;
}
</style> 