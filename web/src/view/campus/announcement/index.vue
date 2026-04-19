<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="标题">
          <el-input v-model="searchInfo.title" placeholder="请输入公告标题" />
        </el-form-item>
        <el-form-item label="发布人">
          <el-input v-model="searchInfo.publisherKeyword" placeholder="请输入昵称或手机号" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部状态" style="width: 140px">
            <el-option v-for="item in announcementStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openCreateDialog">新增公告</el-button>
      </div>

      <el-table :data="tableData" row-key="id" tooltip-effect="dark">
        <el-table-column align="left" label="ID" prop="id" width="90" />
        <el-table-column align="left" label="标题" prop="title" min-width="240" show-overflow-tooltip />
        <el-table-column align="left" label="发布人" min-width="180">
          <template #default="scope">
            {{ scope.row.publisherNickname || '-' }} / {{ scope.row.publisherPhone || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" min-width="120">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="更新时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" :min-width="operateColumnWidth">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">查看</el-button>
            <el-button type="primary" link class="table-button" @click="openEditDialog(scope.row)">编辑</el-button>
            <el-button type="warning" link class="table-button" @click="toggleStatus(scope.row)">
              {{ scope.row.status === 1 ? '下线' : '上线' }}
            </el-button>
            <el-button type="danger" link class="table-button" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      destroy-on-close
      :size="appStore.drawerSize"
      v-model="detailShow"
      :show-close="true"
      :before-close="closeDetailShow"
      title="公告详情"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="ID">{{ detailForm.id }}</el-descriptions-item>
        <el-descriptions-item label="标题">{{ detailForm.title || '-' }}</el-descriptions-item>
        <el-descriptions-item label="发布人">
          {{ detailForm.publisherNickname || '-' }} / {{ detailForm.publisherPhone || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">{{ detailForm.statusText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(detailForm.createdAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(detailForm.updatedAt) || '-' }}</el-descriptions-item>
      </el-descriptions>
      <div class="mt-4">
        <div class="text-base font-bold mb-3">公告内容</div>
        <div class="announcement-content" v-html="detailForm.content || '<p>-</p>'" />
      </div>
    </el-drawer>

    <el-dialog v-model="dialogVisible" :title="type === 'create' ? '新增公告' : '编辑公告'" width="960px">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="88px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入公告标题" />
        </el-form-item>
        <el-form-item label="发布人" prop="publisherId">
          <el-select v-model="formData.publisherId" filterable placeholder="请选择发布人" style="width: 100%">
            <el-option
              v-for="item in publisherOptions"
              :key="item.id"
              :label="`${item.nickname || '管理员'} / ${item.phone}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="formData.status" style="width: 100%">
            <el-option v-for="item in announcementStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <RichEdit v-model="formData.content" />
        </el-form-item>
        <el-form-item label="变更说明">
          <el-input
            v-model="formData.auditReason"
            type="textarea"
            :rows="4"
            maxlength="256"
            show-word-limit
            placeholder="可选，填写本次新增或编辑公告的说明"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" :loading="submitLoading" @click="submitDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import RichEdit from '@/components/richtext/rich-edit.vue'
import {
  createCampusAnnouncement,
  deleteCampusAnnouncement,
  findCampusAnnouncement,
  getCampusAnnouncementList,
  updateCampusAnnouncement,
  updateCampusAnnouncementStatus
} from '@/api/campusAnnouncement'
import { getCampusUserList } from '@/api/campusUser'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'CampusAnnouncement'
})

const appStore = useAppStore()

const announcementStatusOptions = [
  { value: 0, label: '下线' },
  { value: 1, label: '上线' }
]

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const publisherOptions = ref([])
const detailShow = ref(false)
const dialogVisible = ref(false)
const submitLoading = ref(false)
const type = ref('create')
const formRef = ref()

const createSearchInfo = () => ({
  title: '',
  publisherKeyword: '',
  status: undefined,
  createdAtRange: []
})

const createDetail = () => ({
  id: 0,
  title: '',
  content: '',
  publisherNickname: '',
  publisherPhone: '',
  statusText: '',
  createdAt: '',
  updatedAt: ''
})

const createFormData = () => ({
  id: 0,
  title: '',
  content: '',
  publisherId: undefined,
  status: 1,
  auditReason: ''
})

const searchInfo = ref(createSearchInfo())
const detailForm = ref(createDetail())
const formData = ref(createFormData())
const operateColumnWidth = computed(() => Number(appStore.operateMinWith) + 80)
const rules = reactive({
  title: [{ required: true, message: '请输入公告标题', trigger: 'blur' }],
  publisherId: [{ required: true, message: '请选择发布人', trigger: 'change' }],
  status: [{ required: true, message: '请选择公告状态', trigger: 'change' }]
})

const loadPublishers = async () => {
  const res = await getCampusUserList({
    page: 1,
    pageSize: 200,
    role: 1
  })
  if (res.code === 0) {
    publisherOptions.value = res.data.list || []
  }
}

const getTableData = async () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  }
  if (!params.createdAtRange?.length) {
    delete params.createdAtRange
  }
  if (typeof params.status === 'undefined') {
    delete params.status
  }

  const table = await getCampusAnnouncementList(params)
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = createSearchInfo()
  page.value = 1
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getDetails = async (row) => {
  const res = await findCampusAnnouncement({ id: row.id })
  if (res.code === 0) {
    detailForm.value = {
      ...createDetail(),
      ...res.data
    }
    detailShow.value = true
  }
}

const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = createDetail()
}

const openCreateDialog = () => {
  type.value = 'create'
  formData.value = createFormData()
  dialogVisible.value = true
}

const openEditDialog = async (row) => {
  const res = await findCampusAnnouncement({ id: row.id })
  if (res.code === 0) {
    type.value = 'update'
    formData.value = {
      id: res.data.id,
      title: res.data.title,
      content: res.data.content,
      publisherId: res.data.publisherId,
      status: res.data.status,
      auditReason: ''
    }
    dialogVisible.value = true
  }
}

const closeDialog = () => {
  dialogVisible.value = false
  formData.value = createFormData()
  formRef.value?.clearValidate()
}

const submitDialog = async () => {
  await formRef.value?.validate()
  submitLoading.value = true
  const request = type.value === 'create' ? createCampusAnnouncement : updateCampusAnnouncement
  const res = await request(formData.value)
  submitLoading.value = false
  if (res.code === 0) {
    ElMessage.success(type.value === 'create' ? '创建成功' : '更新成功')
    closeDialog()
    getTableData()
  }
}

const toggleStatus = async (row) => {
  const targetStatus = row.status === 1 ? 0 : 1
  const actionText = targetStatus === 1 ? '上线' : '下线'
  let auditReason = ''
  try {
    const promptResult = await ElMessageBox.prompt(`确定要${actionText}公告【${row.title}】吗？`, '提示', {
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
  const res = await updateCampusAnnouncementStatus({
    id: row.id,
    status: targetStatus,
    auditReason
  })
  if (res.code === 0) {
    ElMessage.success(`${actionText}成功`)
    getTableData()
    if (detailShow.value && detailForm.value.id === row.id) {
      getDetails(row)
    }
  }
}

const deleteRow = async (row) => {
  let auditReason = ''
  try {
    const promptResult = await ElMessageBox.prompt(`确定要删除公告【${row.title}】吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
      inputType: 'textarea',
      inputPlaceholder: '请输入删除原因',
      inputValidator: (value) => {
        const trimmed = value?.trim?.() || ''
        if (!trimmed) {
          return '请输入删除原因'
        }
        if (trimmed.length > 256) {
          return '删除原因最多 256 个字符'
        }
        return true
      }
    })
    auditReason = promptResult.value.trim()
  } catch (e) {
    return
  }
  const res = await deleteCampusAnnouncement({ id: row.id, auditReason })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

loadPublishers()
getTableData()
</script>

<style scoped>
.announcement-content {
  min-height: 180px;
  padding: 16px;
  border: 1px solid var(--el-border-color);
  border-radius: 8px;
  background: var(--el-bg-color-page);
}
</style>
