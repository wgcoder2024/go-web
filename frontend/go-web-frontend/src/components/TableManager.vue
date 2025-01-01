<template>
  <div class="container-fluid">
    <nav class="navbar navbar-dark bg-primary mb-4">
      <div class="container">
        <span class="navbar-brand mb-0 h1">数据库管理</span>
        <div class="d-flex gap-2">
          <button class="btn btn-outline-light" @click="showSqlExecutor = true">
            <i class="bi bi-terminal"></i> SQL 查询
          </button>
          <button class="btn btn-outline-light" @click="showImportModal = true">
            <i class="bi bi-upload"></i> 导入表
          </button>
        </div>
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
                placeholder="搜索表名..." 
                @input="filterTables"
              >
            </div>
            <button class="btn btn-primary" @click="showCreateTableForm = true">
              <i class="bi bi-plus-lg"></i> 创建新表
            </button>
          </div>

          <div class="table-responsive">
            <table class="table table-striped table-hover">
              <thead class="table-light">
                <tr>
                  <th>表名</th>
                  <th>字段数</th>
                  <th>记录数</th>
                  <th>创建时间</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="table in sortedTables" :key="table.name">
                  <td>{{ table.name }}</td>
                  <td>{{ table.columns }}</td>
                  <td>{{ table.rows }}</td>
                  <td>{{ formatDate(table.createTime) }}</td>
                  <td>
                    <div class="btn-group">
                      <button class="btn btn-sm btn-outline-primary" @click="viewTableData(table)">
                        <i class="bi bi-table"></i> 数据
                      </button>
                      <button class="btn btn-sm btn-outline-secondary" @click="viewTable(table)">
                        <i class="bi bi-pencil"></i> 结构
                      </button>
                      <button class="btn btn-sm btn-outline-danger" @click="deleteTable(table.name)">
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

      <!-- 导入表模态框 -->
      <div class="modal" v-if="showImportModal" tabindex="-1" style="display: block;">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">导入表</h5>
              <button type="button" class="btn-close" @click="showImportModal = false"></button>
            </div>
            <div class="modal-body">
              <div class="mb-3">
                <label class="form-label">SQL 文件</label>
                <input type="file" class="form-control" @change="handleFileUpload" accept=".sql">
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" @click="showImportModal = false">取消</button>
              <button type="button" class="btn btn-primary" @click="importTable">导入</button>
            </div>
          </div>
        </div>
        <div class="modal-backdrop" style="opacity: 0.5;"></div>
      </div>

      <!-- 其他模态框组件 -->
      <TableDataView 
        v-if="showTableDataView"
        :tableName="currentTable.name"
        @close="showTableDataView = false"
      />

      <SQLExecutor
        v-if="showSqlExecutor"
        @close="showSqlExecutor = false"
      />

      <!-- 表详情对话框 -->
      <div class="modal" v-if="showTableDetail" tabindex="-1" style="display: block;">
        <div class="modal-dialog modal-xl modal-dialog-centered">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">表详情: {{ currentTable.name }}</h5>
              <button type="button" class="btn-close" @click="showTableDetail = false"></button>
            </div>
            <div class="modal-body">
              <!-- 表格工具栏 -->
              <div class="d-flex justify-content-between mb-3">
                <div class="btn-group">
                  <button class="btn btn-outline-primary" @click="sortBy('name')">
                    按名称排序 <i class="bi" :class="getSortIcon('name')"></i>
                  </button>
                  <button class="btn btn-outline-primary" @click="sortBy('type')">
                    按类型排序 <i class="bi" :class="getSortIcon('type')"></i>
                  </button>
                </div>
                <button class="btn btn-primary" @click="showAddColumnForm = true">
                  <i class="bi bi-plus-lg"></i> 添加字段
                </button>
              </div>

              <!-- 字段列表 -->
              <table class="table table-bordered">
                <thead class="table-light">
                  <tr>
                    <th>字段名</th>
                    <th>类型</th>
                    <th>可空</th>
                    <th>默认值</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="column in sortedColumns" :key="column.name">
                    <td>{{ column.name }}</td>
                    <td>{{ column.type }}</td>
                    <td>{{ column.nullable ? '是' : '否' }}</td>
                    <td>{{ column.default || '-' }}</td>
                    <td>
                      <button class="btn btn-sm btn-outline-danger" @click="deleteColumn(column.name)">
                        <i class="bi bi-trash"></i>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
        <div class="modal-backdrop" style="opacity: 0.5;"></div>
      </div>

      <!-- 创建表对话框 -->
      <div class="modal" v-if="showCreateTableForm" tabindex="-1" style="display: block;">
        <div class="modal-dialog modal-lg modal-dialog-centered">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">创建新表</h5>
              <button type="button" class="btn-close" @click="showCreateTableForm = false"></button>
            </div>
            <form @submit.prevent="createTable">
              <div class="modal-body">
                <div class="mb-3">
                  <label class="form-label">表名</label>
                  <input type="text" class="form-control" v-model="newTable.name" required>
                </div>
                
                <div class="mb-3">
                  <label class="form-label d-flex justify-content-between">
                    <span>字段列表</span>
                    <button type="button" class="btn btn-sm btn-outline-primary" @click="addColumn">
                      <i class="bi bi-plus"></i> 添加字段
                    </button>
                  </label>
                  
                  <div class="table-responsive">
                    <table class="table table-bordered">
                      <thead class="table-light">
                        <tr>
                          <th>字段名</th>
                          <th>类型</th>
                          <th>可空</th>
                          <th>默认值</th>
                          <th>操作</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(column, index) in newTable.columns" :key="index">
                          <td>
                            <input type="text" class="form-control form-control-sm" v-model="column.name" required>
                          </td>
                          <td>
                            <select class="form-select form-select-sm" v-model="column.type" required>
                              <option value="INT">INT</option>
                              <option value="VARCHAR(255)">VARCHAR(255)</option>
                              <option value="TEXT">TEXT</option>
                              <option value="DATETIME">DATETIME</option>
                              <option value="BOOLEAN">BOOLEAN</option>
                            </select>
                          </td>
                          <td>
                            <div class="form-check">
                              <input class="form-check-input" type="checkbox" v-model="column.nullable">
                            </div>
                          </td>
                          <td>
                            <input type="text" class="form-control form-control-sm" v-model="column.default">
                          </td>
                          <td>
                            <button type="button" class="btn btn-sm btn-outline-danger" @click="removeColumn(index)">
                              <i class="bi bi-trash"></i>
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
              <div class="modal-footer">
                <button type="button" class="btn btn-secondary" @click="showCreateTableForm = false">取消</button>
                <button type="submit" class="btn btn-primary">创建</button>
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
import { ref, computed, onMounted } from 'vue'
import TableDataView from './TableDataView.vue'
import SQLExecutor from './SQLExecutor.vue'

const tables = ref([])
const searchQuery = ref('')
const showTableDetail = ref(false)
const currentTable = ref({})
const sortConfig = ref({ field: 'name', direction: 'asc' })

const API_URL = 'http://localhost:8080/api/v1/tables'

// 获取表列表
const fetchTables = async () => {
  try {
    const response = await fetch(API_URL)
    tables.value = await response.json()
  } catch (error) {
    console.error('Error fetching tables:', error)
  }
}

const sortedTables = computed(() => {
  return [...tables.value].sort((a, b) => {
    const direction = sortConfig.value.direction === 'asc' ? 1 : -1
    return a[sortConfig.value.field] > b[sortConfig.value.field] ? direction : -direction
  })
})

const filterTables = () => {
  // 实现表名搜索逻辑
}

// 获取表详情
const viewTable = async (table) => {
  try {
    const response = await fetch(`${API_URL}/${table.name}`)
    const data = await response.json()
    // 使用返回的完整数据
    currentTable.value = data
    showTableDetail.value = true
  } catch (error) {
    console.error('Error fetching table details:', error)
  }
}

// 删除表
const deleteTable = async (tableName) => {
  if (!confirm(`确定要删除表 ${tableName} 吗？`)) return
  
  try {
    const response = await fetch(`${API_URL}/${tableName}`, {
      method: 'DELETE'
    })
    
    if (response.ok) {
      fetchTables()
    }
  } catch (error) {
    console.error('Error deleting table:', error)
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

const sortBy = (field) => {
  if (sortConfig.value.field === field) {
    sortConfig.value.direction = sortConfig.value.direction === 'asc' ? 'desc' : 'asc'
  } else {
    sortConfig.value.field = field
    sortConfig.value.direction = 'asc'
  }
}

const getSortIcon = (field) => {
  if (sortConfig.value.field !== field) return 'bi-arrow-down-up'
  return sortConfig.value.direction === 'asc' ? 'bi-arrow-down' : 'bi-arrow-up'
}

// 新增的响应式变量
const showCreateTableForm = ref(false)
const newTable = ref({
  name: '',
  columns: []
})

// 添加字段
const addColumn = () => {
  newTable.value.columns.push({
    name: '',
    type: 'VARCHAR(255)',
    nullable: true,
    default: ''
  })
}

// 删除字段
const removeColumn = (index) => {
  newTable.value.columns.splice(index, 1)
}

// 创建表
const createTable = async () => {
  try {
    const response = await fetch(API_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(newTable.value)
    })

    if (response.ok) {
      showCreateTableForm.value = false
      newTable.value = { name: '', columns: [] }
      fetchTables()
    }
  } catch (error) {
    console.error('Error creating table:', error)
  }
}

// 新增的响应式变量
const showTableDataView = ref(false)
const showSqlExecutor = ref(false)
const showImportModal = ref(false)
const importFile = ref(null)

// 查看表数据
const viewTableData = (table) => {
  currentTable.value = table
  showTableDataView.value = true
}

// 处理文件上传
const handleFileUpload = (event) => {
  importFile.value = event.target.files[0]
}

// 导入表
const importTable = async () => {
  if (!importFile.value) return

  try {
    const reader = new FileReader()
    reader.onload = async (e) => {
      const sql = e.target.result
      const response = await fetch('http://localhost:8080/api/v1/tables/import', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          name: importFile.value.name.replace('.sql', ''),
          structure: sql.split('\n\n')[0],
          data: sql.split('\n\n')[1] || ''
        })
      })

      if (response.ok) {
        showImportModal.value = false
        fetchTables()
      }
    }
    reader.readAsText(importFile.value)
  } catch (error) {
    console.error('Error importing table:', error)
  }
}

// 添加计算属性用于排序列
const sortedColumns = computed(() => {
  if (!currentTable.value.columnDetails) return []
  return [...currentTable.value.columnDetails].sort((a, b) => {
    const direction = sortConfig.value.direction === 'asc' ? 1 : -1
    return a[sortConfig.value.field] > b[sortConfig.value.field] ? direction : -direction
  })
})

onMounted(() => {
  fetchTables()
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