<template>
  <div class="user-list">
    <div class="controls">
      <input 
        v-model="searchQuery" 
        placeholder="搜索用户名..." 
        @input="searchUsers"
      >
      <button @click="showCreateForm = true">添加用户</button>
    </div>

    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>姓名</th>
          <th>邮箱</th>
          <th>年龄</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.id">
          <td>{{ user.id }}</td>
          <td>{{ user.name }}</td>
          <td>{{ user.email }}</td>
          <td>{{ user.age }}</td>
          <td>
            <button @click="editUser(user)">编辑</button>
            <button @click="deleteUser(user.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- 创建/编辑用户对话框 -->
    <div v-if="showCreateForm" class="modal">
      <div class="modal-content">
        <h3>{{ isEditing ? '编辑用户' : '创建用户' }}</h3>
        <form @submit.prevent="submitUser">
          <div>
            <label>姓名:</label>
            <input v-model="currentUser.name" required>
          </div>
          <div>
            <label>邮箱:</label>
            <input v-model="currentUser.email" type="email" required>
          </div>
          <div>
            <label>年龄:</label>
            <input v-model="currentUser.age" type="number" required>
          </div>
          <div class="form-buttons">
            <button type="submit">{{ isEditing ? '更新' : '创建' }}</button>
            <button type="button" @click="showCreateForm = false">取消</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const users = ref([])
const showCreateForm = ref(false)
const currentUser = ref({})
const isEditing = ref(false)
const searchQuery = ref('')

const API_URL = 'http://localhost:8080/api/v1/users'

// 获取用户列表
const fetchUsers = async () => {
  try {
    const response = await fetch(API_URL)
    users.value = await response.json()
  } catch (error) {
    console.error('Error fetching users:', error)
  }
}

// 搜索用户
const searchUsers = async () => {
  try {
    const response = await fetch(`${API_URL}?name=${searchQuery.value}`)
    users.value = await response.json()
  } catch (error) {
    console.error('Error searching users:', error)
  }
}

// 创建或更新用户
const submitUser = async () => {
  try {
    const url = isEditing.value ? `${API_URL}/${currentUser.value.id}` : API_URL
    const method = isEditing.value ? 'PUT' : 'POST'
    
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(currentUser.value),
    })

    if (response.ok) {
      showCreateForm.value = false
      fetchUsers()
      currentUser.value = {}
    }
  } catch (error) {
    console.error('Error submitting user:', error)
  }
}

// 编辑用户
const editUser = (user) => {
  currentUser.value = { ...user }
  isEditing.value = true
  showCreateForm.value = true
}

// 删除用户
const deleteUser = async (id) => {
  if (!confirm('确定要删除这个用户吗？')) return
  
  try {
    const response = await fetch(`${API_URL}/${id}`, {
      method: 'DELETE',
    })
    
    if (response.ok) {
      fetchUsers()
    }
  } catch (error) {
    console.error('Error deleting user:', error)
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.user-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.controls {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

button {
  padding: 8px 16px;
  margin: 0 4px;
  cursor: pointer;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  min-width: 300px;
}

form div {
  margin: 10px 0;
}

label {
  display: inline-block;
  width: 60px;
}

.form-buttons {
  margin-top: 20px;
  text-align: right;
}
</style> 