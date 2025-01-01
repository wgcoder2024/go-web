<template>
  <div class="container-fluid">
    <nav class="navbar navbar-dark bg-primary mb-4">
      <div class="container">
        <span class="navbar-brand mb-0 h1">用户管理系统</span>
      </div>
    </nav>
    
    <div class="container">
      <div class="card">
        <div class="card-body">
          <div class="d-flex justify-content-between align-items-center mb-4">
            <div class="input-group" style="max-width: 300px;">
              <span class="input-group-text">
                <i class="bi bi-search"></i>
              </span>
              <input 
                type="text"
                class="form-control"
                v-model="searchQuery" 
                placeholder="搜索用户名..." 
                @input="searchUsers"
              >
            </div>
            <button class="btn btn-primary" @click="showCreateForm = true">
              <i class="bi bi-plus-lg"></i> 添加用户
            </button>
          </div>

          <div class="table-responsive">
            <table class="table table-striped table-hover">
              <thead class="table-light">
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
                    <button class="btn btn-sm btn-outline-primary me-2" @click="editUser(user)">
                      <i class="bi bi-pencil"></i> 编辑
                    </button>
                    <button class="btn btn-sm btn-outline-danger" @click="deleteUser(user.id)">
                      <i class="bi bi-trash"></i> 删除
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- 创建/编辑用户对话框 -->
      <div class="modal" v-if="showCreateForm" tabindex="-1" style="display: block;">
        <div class="modal-dialog modal-dialog-centered">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">{{ isEditing ? '编辑用户' : '创建用户' }}</h5>
              <button type="button" class="btn-close" @click="showCreateForm = false"></button>
            </div>
            <form @submit.prevent="submitUser">
              <div class="modal-body">
                <div class="mb-3">
                  <label class="form-label">姓名</label>
                  <input type="text" class="form-control" v-model="currentUser.name" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">邮箱</label>
                  <input type="email" class="form-control" v-model="currentUser.email" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">年龄</label>
                  <input type="number" class="form-control" v-model="currentUser.age" required>
                </div>
              </div>
              <div class="modal-footer">
                <button type="button" class="btn btn-secondary" @click="showCreateForm = false">取消</button>
                <button type="submit" class="btn btn-primary">{{ isEditing ? '更新' : '创建' }}</button>
              </div>
            </form>
          </div>
        </div>
        <div class="modal-backdrop" style="opacity: 0.5;"></div>
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

<style>
.modal {
  display: block;
  background-color: rgba(0, 0, 0, 0.5);
}

.modal-dialog {
  z-index: 1056;
}

.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #000;
  z-index: 1055;
}
</style> 