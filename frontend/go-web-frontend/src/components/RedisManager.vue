<template>
  <div class="container-fluid">
    <nav class="navbar navbar-dark bg-primary mb-4">
      <div class="container">
        <span class="navbar-brand mb-0 h1">Redis 管理</span>
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
                v-model="pattern" 
                placeholder="搜索键名..."
                @keyup.enter="fetchKeys"
              >
            </div>
            <button class="btn btn-primary" @click="showCreateForm = true">
              <i class="bi bi-plus-lg"></i> 添加键值对
            </button>
          </div>

          <div class="table-responsive">
            <table class="table table-striped table-hover">
              <thead class="table-light">
                <tr>
                  <th>键名</th>
                  <th>类型</th>
                  <th>值</th>
                  <th>过期时间</th>
                  <th>节点</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="entry in entries" :key="entry.key">
                  <td>{{ entry.key }}</td>
                  <td>{{ entry.type }}</td>
                  <td>{{ entry.value }}</td>
                  <td>{{ formatTTL(entry.ttl) }}</td>
                  <td>{{ entry.node || '单机' }}</td>
                  <td>
                    <div class="btn-group">
                      <button class="btn btn-sm btn-outline-primary" @click="editEntry(entry)">
                        <i class="bi bi-pencil"></i> 编辑
                      </button>
                      <button class="btn btn-sm btn-outline-danger" @click="deleteEntry(entry.key)">
                        <i class="bi bi-trash"></i> 删除
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- 创建/编辑表单 -->
      <div class="modal" v-if="showCreateForm" tabindex="-1" style="display: block;">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">{{ isEditing ? '编辑键值对' : '创建键值对' }}</h5>
              <button type="button" class="btn-close" @click="showCreateForm = false"></button>
            </div>
            <form @submit.prevent="submitEntry">
              <div class="modal-body">
                <div class="mb-3">
                  <label class="form-label">键名</label>
                  <input 
                    type="text" 
                    class="form-control" 
                    v-model="currentEntry.key"
                    :readonly="isEditing"
                    required
                  >
                </div>
                <div class="mb-3">
                  <label class="form-label">值</label>
                  <input type="text" class="form-control" v-model="currentEntry.value" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">过期时间（秒）</label>
                  <input type="number" class="form-control" v-model="currentEntry.ttl">
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

const API_URL = 'http://localhost:8080/api/v1/redis/keys'
const entries = ref([])
const pattern = ref('*')
const showCreateForm = ref(false)
const currentEntry = ref({
  key: '',
  value: '',
  ttl: 0
})
const isEditing = ref(false)

const fetchKeys = async () => {
  try {
    const response = await fetch(`${API_URL}?pattern=${pattern.value}`)
    entries.value = await response.json()
  } catch (error) {
    console.error('Error fetching redis keys:', error)
  }
}

const editEntry = (entry) => {
  currentEntry.value = { ...entry }
  isEditing.value = true
  showCreateForm.value = true
}

const deleteEntry = async (key) => {
  if (!confirm(`确定要删除键 ${key} 吗？`)) return

  try {
    const response = await fetch(`${API_URL}/${key}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      fetchKeys()
    }
  } catch (error) {
    console.error('Error deleting redis key:', error)
  }
}

const submitEntry = async () => {
  try {
    const response = await fetch(API_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(currentEntry.value)
    })

    if (response.ok) {
      showCreateForm.value = false
      currentEntry.value = { key: '', value: '', ttl: 0 }
      isEditing.value = false
      fetchKeys()
    }
  } catch (error) {
    console.error('Error submitting redis entry:', error)
  }
}

const formatTTL = (ttl) => {
  if (ttl < 0) return '永久'
  return `${ttl}秒`
}

onMounted(() => {
  fetchKeys()
})
</script>

<style scoped>
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