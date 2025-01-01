<template>
  <div class="modal" tabindex="-1" style="display: block;">
    <div class="modal-dialog modal-xl modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">SQL 查询执行器</h5>
          <button type="button" class="btn-close" @click="$emit('close')"></button>
        </div>
        <div class="modal-body">
          <div class="mb-3">
            <label class="form-label">SQL 查询语句 (仅支持 SELECT)</label>
            <textarea 
              class="form-control font-monospace" 
              v-model="sql" 
              rows="5"
              placeholder="SELECT * FROM users WHERE age > 18"
            ></textarea>
          </div>
          
          <div class="d-flex justify-content-between mb-3">
            <button class="btn btn-primary" @click="executeQuery">
              <i class="bi bi-play-fill"></i> 执行查询
            </button>
            <button class="btn btn-outline-secondary" @click="clearResults">
              <i class="bi bi-trash"></i> 清除结果
            </button>
          </div>

          <div v-if="results.length" class="table-responsive">
            <table class="table table-striped table-hover">
              <thead class="table-light">
                <tr>
                  <th v-for="col in columns" :key="col">{{ col }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, index) in results" :key="index">
                  <td v-for="col in columns" :key="col">{{ row[col] }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
    <div class="modal-backdrop" style="opacity: 0.5;"></div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['close'])

const sql = ref('')
const results = ref([])

const columns = computed(() => {
  if (!results.value.length) return []
  return Object.keys(results.value[0])
})

const executeQuery = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/v1/tables/query', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ sql: sql.value })
    })
    results.value = await response.json()
  } catch (error) {
    console.error('Error executing query:', error)
  }
}

const clearResults = () => {
  results.value = []
}
</script>

<style scoped>
.modal-dialog {
  max-width: 90%;
}
</style> 