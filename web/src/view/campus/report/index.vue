<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="举报ID">
          <el-input v-model="searchInfo.id" placeholder="请输入举报ID" />
        </el-form-item>
        <el-form-item label="举报人ID">
          <el-input v-model="searchInfo.reporterId" placeholder="请输入举报人ID" />
        </el-form-item>
        <el-form-item label="目标类型">
          <el-select v-model="searchInfo.targetType" clearable placeholder="全部类型" style="width: 140px">
            <el-option v-for="item in reportTargetTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标ID">
          <el-input v-model="searchInfo.targetId" placeholder="请输入目标ID" />
        </el-form-item>
        <el-form-item label="举报原因">
          <el-select v-model="searchInfo.reason" clearable placeholder="全部原因" style="width: 160px">
            <el-option v-for="item in reportReasonOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部状态" style="width: 140px">
            <el-option v-for="item in reportStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
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
      <el-table :data="tableData" row-key="id" tooltip-effect="dark">
        <el-table-column align="left" label="ID" prop="id" width="90" />
        <el-table-column align="left" label="举报人" min-width="180">
          <template #default="scope">
            {{ scope.row.reporterNickname || '-' }} / {{ scope.row.reporterPhone || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="目标类型" prop="targetTypeText" min-width="120" />
        <el-table-column align="left" label="目标ID" prop="targetId" width="100" />
        <el-table-column align="left" label="原因" prop="reasonText" min-width="140" />
        <el-table-column align="left" label="状态" min-width="120">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'warning'">
              {{ scope.row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="处理人" min-width="180">
          <template #default="scope">
            <span v-if="scope.row.handledByNickname || scope.row.handledByPhone">
              {{ scope.row.handledByNickname || '-' }} / {{ scope.row.handledByPhone || '-' }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" :min-width="operateColumnWidth">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">查看</el-button>
            <el-button type="warning" link class="table-button" @click="openHandleDialog(scope.row)">处理</el-button>
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
      title="举报详情"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="举报ID">{{ detailForm.id }}</el-descriptions-item>
        <el-descriptions-item label="举报人">
          {{ detailForm.reporterNickname || '-' }} / {{ detailForm.reporterPhone || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="目标类型">{{ detailForm.targetTypeText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="目标ID">{{ detailForm.targetId || '-' }}</el-descriptions-item>
        <el-descriptions-item label="目标商品">{{ detailForm.targetProductTitle || '-' }}</el-descriptions-item>
        <el-descriptions-item label="举报原因">{{ detailForm.reasonText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="举报描述">{{ detailForm.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="处理状态">{{ detailForm.statusText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="处理人">
          <span v-if="detailForm.handledByNickname || detailForm.handledByPhone">
            {{ detailForm.handledByNickname || '-' }} / {{ detailForm.handledByPhone || '-' }}
          </span>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="处理结果">{{ detailForm.handleResult || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(detailForm.createdAt) || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>

    <el-dialog v-model="handleDialogVisible" title="处理举报" width="460px">
      <el-form ref="handleFormRef" :model="handleForm" :rules="handleRules" label-width="88px">
        <el-form-item label="举报ID">
          <span>{{ handleForm.id }}</span>
        </el-form-item>
        <el-form-item label="处理状态" prop="status">
          <el-select v-model="handleForm.status" style="width: 100%">
            <el-option v-for="item in reportStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="处理结果" prop="handleResult">
          <el-input
            v-model="handleForm.handleResult"
            type="textarea"
            :rows="4"
            placeholder="请输入处理结果"
          />
        </el-form-item>
        <el-form-item label="执行原因" prop="auditReason">
          <el-input
            v-model="handleForm.auditReason"
            type="textarea"
            :rows="4"
            maxlength="256"
            show-word-limit
            placeholder="请输入处理原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeHandleDialog">取 消</el-button>
          <el-button type="primary" :loading="handleLoading" @click="submitHandle">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { findCampusReport, getCampusReportList, handleCampusReport } from '@/api/campusReport'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'CampusReport'
})

const appStore = useAppStore()

const reportTargetTypeOptions = [
  { value: 1, label: '商品' }
]

const reportReasonOptions = [
  { value: 1, label: '虚假信息' },
  { value: 2, label: '疑似违规' },
  { value: 3, label: '疑似诈骗' }
]

const reportStatusOptions = [
  { value: 0, label: '待处理' },
  { value: 1, label: '已处理' }
]

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const detailShow = ref(false)
const handleDialogVisible = ref(false)
const handleLoading = ref(false)
const handleFormRef = ref()

const createSearchInfo = () => ({
  id: '',
  reporterId: '',
  targetType: undefined,
  targetId: '',
  reason: undefined,
  status: undefined,
  createdAtRange: []
})

const createDetail = () => ({
  id: 0,
  reporterNickname: '',
  reporterPhone: '',
  targetTypeText: '',
  targetId: 0,
  targetProductTitle: '',
  reasonText: '',
  description: '',
  statusText: '',
  handledByNickname: '',
  handledByPhone: '',
  handleResult: '',
  createdAt: ''
})

const createHandleForm = () => ({
  id: 0,
  status: 1,
  handleResult: '',
  auditReason: ''
})

const searchInfo = ref(createSearchInfo())
const detailForm = ref(createDetail())
const handleForm = ref(createHandleForm())
const operateColumnWidth = computed(() => Number(appStore.operateMinWith) + 40)
const handleRules = reactive({
  status: [{ required: true, message: '请选择处理状态', trigger: 'change' }],
  handleResult: [{ required: true, message: '请输入处理结果', trigger: 'blur' }],
  auditReason: [{ required: true, message: '请输入处理原因', trigger: 'blur' }]
})

const getTableData = async () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  }
  if (!params.createdAtRange?.length) {
    delete params.createdAtRange
  }
  if (typeof params.targetType === 'undefined') {
    delete params.targetType
  }
  if (typeof params.reason === 'undefined') {
    delete params.reason
  }
  if (typeof params.status === 'undefined') {
    delete params.status
  }

  const table = await getCampusReportList(params)
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
  const res = await findCampusReport({ id: row.id })
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

const openHandleDialog = (row) => {
  handleForm.value = {
    id: row.id,
    status: 1,
    handleResult: row.handleResult || '',
    auditReason: ''
  }
  handleDialogVisible.value = true
}

const closeHandleDialog = () => {
  handleDialogVisible.value = false
  handleForm.value = createHandleForm()
  handleFormRef.value?.clearValidate()
}

const submitHandle = async () => {
  await handleFormRef.value?.validate()
  handleLoading.value = true
  const currentID = handleForm.value.id
  const res = await handleCampusReport(handleForm.value)
  handleLoading.value = false
  if (res.code === 0) {
    ElMessage.success('处理成功')
    closeHandleDialog()
    getTableData()
    if (detailShow.value && detailForm.value.id === currentID) {
      getDetails({ id: currentID })
    }
  }
}

getTableData()
</script>
