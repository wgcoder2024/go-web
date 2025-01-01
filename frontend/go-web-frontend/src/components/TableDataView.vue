<template>
  <div class="modal" tabindex="-1" style="display: block;">
    <div class="modal-dialog modal-xl modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">数据预览: {{ tableName }}</h5>
          <button type="button" class="btn-close" @click="$emit('close')"></button>
        </div>
        <div class="modal-body">
          <!-- 工具栏 -->
          <div class="d-flex justify-content-between mb-3">
            <div class="d-flex gap-2">
              <select class="form-select" v-model="pageSize" @change="fetchData">
                <option value="10">10条/页</option>
                <option value="20">20条/页</option>
                <option value="50">50条/页</option>
              </select>
              <select class="form-select" v-model="sortField" @change="fetchData">
                <option v-for="col in columns" :key="col" :value="col">
                  按{{ col }}排序
                </option>
              </select>
              <select class="form-select" v-model="sortOrder" @change="fetchData">
                <option value="asc">升序</option>
                <option value="desc">降序</option>
              </select>
            </div>
            <button class="btn btn-primary" @click="exportData">
              <i class="bi bi-download"></i> 导出数据
            </button>
          </div>

          <!-- 数据表格 -->
          <div class="table-responsive">
            <table class="table table-striped table-hover">
              <thead class="table-light">
                <tr>
                  <th v-for="col in columns" :key="col">{{ col }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, index) in rows" :key="index">
                  <td v-for="(value, colIndex) in row" :key="colIndex">
                    {{ value }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- 分页 -->
          <nav class="mt-3">
            <ul class="pagination justify-content-center">
              <li class="page-item" :class="{ disabled: currentPage === 1 }">
                <a class="page-link" href="#" @click.prevent="changePage(currentPage - 1)">
                  上一页
                </a>
              </li>
              <li v-for="p in pages" :key="p" class="page-item" :class="{ active: currentPage === p }">
                <a class="page-link" href="#" @click.prevent="changePage(p)">{{ p }}</a>
              </li>
              <li class="page-item" :class="{ disabled: currentPage === totalPages }">
                <a class="page-link" href="#" @click.prevent="changePage(currentPage + 1)">
                  下一页
                </a>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </div>
    <div class="modal-backdrop" style="opacity: 0.5;"></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const props = defineProps({
  tableName: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close'])

const columns = ref([])
const rows = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const sortField = ref('id')
const sortOrder = ref('asc')

const API_URL = `http://localhost:8080/api/v1/tables/${props.tableName}/data`

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))
const pages = computed(() => {
  const arr = []
  for (let i = 1; i <= totalPages.value; i++) {
    arr.push(i)
  }
  return arr
})

const fetchData = async () => {
  try {
    const response = await fetch(
      `${API_URL}?page=${currentPage.value}&pageSize=${pageSize.value}&sortField=${sortField.value}&sortOrder=${sortOrder.value}`
    )
    const data = await response.json()
    columns.value = data.columns
    rows.value = data.rows
    total.value = data.total
  } catch (error) {
    console.error('Error fetching table data:', error)
  }
}

const changePage = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchData()
}

const exportData = async () => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/tables/${props.tableName}/export`)
    const data = await response.json()
    
    // 创建下载链接
    const blob = new Blob([data.structure + '\n\n' + data.data], { type: 'text/plain' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${props.tableName}_backup.sql`
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
  } catch (error) {
    console.error('Error exporting table:', error)
  }
}

onMounted(fetchData)
</script>

<style scoped>
.modal-dialog {
  max-width: 90%;
}

.pagination {
  margin-bottom: 0;
}
</style> 