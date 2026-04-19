<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="分类名称">
          <el-input v-model="searchInfo.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部状态" style="width: 140px">
            <el-option label="启用" :value="0" />
            <el-option label="禁用" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openCreateDialog">新增分类</el-button>
      </div>

      <el-table
        :data="tableData"
        row-key="id"
        tooltip-effect="dark"
        :tree-props="{ children: 'children' }"
        default-expand-all
      >
        <el-table-column align="left" label="ID" prop="id" width="90" />
        <el-table-column align="left" label="分类名称" prop="name" min-width="180" />
        <el-table-column align="left" label="父级分类" min-width="180">
          <template #default="scope">
            {{ scope.row.parentName || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="排序" prop="sortOrder" width="100" />
        <el-table-column align="left" label="图标" min-width="140">
          <template #default="scope">
            {{ scope.row.icon || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" min-width="120">
          <template #default="scope">
            <el-tag :type="scope.row.status === 0 ? 'success' : 'info'">
              {{ scope.row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="openEditDialog(scope.row)">
              编辑
            </el-button>
            <el-button type="warning" link class="table-button" @click="handleStatus(scope.row)">
              {{ scope.row.status === 0 ? '停用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="90px">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="父级分类">
          <el-select v-model="formData.parentId" clearable filterable placeholder="顶级分类" style="width: 100%">
            <el-option v-for="item in parentOptions" :key="item.id" :label="item.label" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.sortOrder" :min="0" :max="9999" style="width: 100%" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="formData.icon" placeholder="请输入图标名称" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :value="0">启用</el-radio>
            <el-radio :value="1">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="变更说明">
          <el-input
            v-model="formData.auditReason"
            type="textarea"
            :rows="4"
            maxlength="256"
            show-word-limit
            placeholder="可选，填写本次新增或编辑分类的说明"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button type="primary" :loading="submitLoading" @click="submitForm">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createCampusCategory,
  findCampusCategory,
  getCampusCategoryList,
  updateCampusCategory,
  updateCampusCategoryStatus
} from '@/api/campusCategory'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'CampusCategory'
})

const appStore = useAppStore()

const formRef = ref()
const tableData = ref([])
const dialogVisible = ref(false)
const submitLoading = ref(false)
const dialogMode = ref('create')

const createSearchInfo = () => ({
  name: '',
  status: undefined
})

const createFormData = () => ({
  id: 0,
  name: '',
  parentId: undefined,
  sortOrder: 0,
  icon: '',
  status: 0,
  auditReason: ''
})

const searchInfo = ref(createSearchInfo())
const formData = ref(createFormData())

const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const dialogTitle = computed(() => (dialogMode.value === 'create' ? '新增分类' : '编辑分类'))
const parentOptions = computed(() => flattenCategoryOptions(tableData.value, dialogMode.value === 'edit' ? formData.value.id : 0))

const getTableData = async () => {
  const params = { ...searchInfo.value }
  if (typeof params.status === 'undefined') {
    delete params.status
  }
  const res = await getCampusCategoryList(params)
  if (res.code === 0) {
    tableData.value = res.data || []
  }
}

const onSubmit = () => {
  getTableData()
}

const onReset = () => {
  searchInfo.value = createSearchInfo()
  getTableData()
}

const openCreateDialog = () => {
  dialogMode.value = 'create'
  formData.value = createFormData()
  dialogVisible.value = true
}

const openEditDialog = async (row) => {
  const res = await findCampusCategory({ id: row.id })
  if (res.code === 0) {
    dialogMode.value = 'edit'
    formData.value = {
      id: res.data.id,
      name: res.data.name,
      parentId: res.data.parentId,
      sortOrder: res.data.sortOrder,
      icon: res.data.icon || '',
      status: res.data.status,
      auditReason: ''
    }
    dialogVisible.value = true
  }
}

const submitForm = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) {
    return
  }

  submitLoading.value = true
  const payload = {
    ...formData.value,
    parentId: formData.value.parentId || null
  }
  const action = dialogMode.value === 'create' ? createCampusCategory : updateCampusCategory
  const res = await action(payload)
  submitLoading.value = false
  if (res.code === 0) {
    ElMessage.success(dialogMode.value === 'create' ? '创建成功' : '更新成功')
    dialogVisible.value = false
    getTableData()
  }
}

const handleStatus = async (row) => {
  const targetStatus = row.status === 0 ? 1 : 0
  const actionText = targetStatus === 0 ? '启用' : '停用'
  let auditReason = ''
  try {
    const promptResult = await ElMessageBox.prompt(`确定要${actionText}分类【${row.name}】吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
      inputType: 'textarea',
      inputPlaceholder: `请输入${actionText}原因`,
      inputValidator: (value) => {
        const trimmed = value?.trim?.() || ''
        if (!trimmed) {
          return `请输入${actionText}原因`
        }
        if (trimmed.length > 256) {
          return '原因最多 256 个字符'
        }
        return true
      }
    })
    auditReason = promptResult.value.trim()
  } catch (e) {
    return
  }
  const res = await updateCampusCategoryStatus({
    id: row.id,
    status: targetStatus,
    auditReason
  })
  if (res.code === 0) {
    ElMessage.success(`${actionText}成功`)
    getTableData()
  }
}

const flattenCategoryOptions = (tree, excludeId = 0, prefix = '') => {
  return tree.flatMap((item) => {
    if (item.id === excludeId) {
      return []
    }
    const label = prefix ? `${prefix} / ${item.name}` : item.name
    const current = [{ id: item.id, label }]
    if (!item.children?.length) {
      return current
    }
    return current.concat(flattenCategoryOptions(item.children, excludeId, label))
  })
}

getTableData()
</script>
