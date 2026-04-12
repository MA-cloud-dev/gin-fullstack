<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="学号">
          <el-input v-model="searchInfo.studentId" placeholder="请输入学号" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="searchInfo.realName" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="学院">
          <el-input v-model="searchInfo.college" placeholder="请输入学院" />
        </el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model="searchInfo.reviewed" clearable placeholder="全部状态" style="width: 120px">
            <el-option :value="false" label="待审" />
            <el-option :value="true" label="已审" />
          </el-select>
        </el-form-item>
        <el-form-item label="申请时间">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>
        <el-form-item label="审核时间">
          <el-date-picker
            v-model="searchInfo.reviewedAtRange"
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
        <el-table-column align="left" label="用户ID" prop="userId" width="100" />
        <el-table-column align="left" label="学号" prop="studentId" min-width="140" />
        <el-table-column align="left" label="姓名" prop="realName" min-width="120" />
        <el-table-column align="left" label="学院" prop="college" min-width="180" />
        <el-table-column align="left" label="审核状态" min-width="100">
          <template #default="scope">
            <el-tag :type="scope.row.reviewedAt ? 'success' : 'warning'">
              {{ scope.row.reviewedAt ? '已审' : '待审' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="审核备注" prop="reviewRemark" min-width="220" show-overflow-tooltip />
        <el-table-column align="left" label="审核人" min-width="140">
          <template #default="scope">
            {{ scope.row.reviewedByName || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="审核时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.reviewedAt) || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="申请时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看
            </el-button>
            <el-button
              v-if="!scope.row.reviewedAt"
              type="primary"
              link
              icon="check"
              class="table-button"
              @click="openReviewDialog(scope.row)"
            >
              审核通过
            </el-button>
            <el-button
              v-else
              type="warning"
              link
              icon="refresh-left"
              class="table-button"
              @click="handleRevoke(scope.row)"
            >
              审核撤回
            </el-button>
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
      title="审核详情"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="ID">{{ detailForm.id }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ detailForm.userId }}</el-descriptions-item>
        <el-descriptions-item label="学号">{{ detailForm.studentId }}</el-descriptions-item>
        <el-descriptions-item label="姓名">{{ detailForm.realName }}</el-descriptions-item>
        <el-descriptions-item label="学院">{{ detailForm.college }}</el-descriptions-item>
        <el-descriptions-item label="审核状态">
          <el-tag :type="detailForm.reviewedAt ? 'success' : 'warning'">
            {{ detailForm.reviewedAt ? '已审' : '待审' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="审核备注">{{ detailForm.reviewRemark || '-' }}</el-descriptions-item>
        <el-descriptions-item label="审核人">{{ detailForm.reviewedByName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="审核时间">{{ formatDate(detailForm.reviewedAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ formatDate(detailForm.createdAt) }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>

    <el-drawer
      destroy-on-close
      :size="appStore.drawerSize"
      v-model="reviewDialogVisible"
      :show-close="false"
      :before-close="closeReviewDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">审核校园身份</span>
          <div>
            <el-button :loading="reviewLoading" type="primary" @click="submitReview">确 定</el-button>
            <el-button @click="closeReviewDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-descriptions :column="1" border class="mb-4">
        <el-descriptions-item label="用户ID">{{ reviewTarget.userId }}</el-descriptions-item>
        <el-descriptions-item label="学号">{{ reviewTarget.studentId }}</el-descriptions-item>
        <el-descriptions-item label="姓名">{{ reviewTarget.realName }}</el-descriptions-item>
        <el-descriptions-item label="学院">{{ reviewTarget.college }}</el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ formatDate(reviewTarget.createdAt) }}</el-descriptions-item>
      </el-descriptions>

      <el-form ref="reviewFormRef" :model="reviewForm" :rules="reviewRules" label-position="top">
        <el-form-item label="审核备注" prop="reviewRemark">
          <el-input
            v-model="reviewForm.reviewRemark"
            type="textarea"
            :rows="5"
            maxlength="256"
            show-word-limit
            placeholder="请输入审核备注，可为空"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { InfoFilled } from '@element-plus/icons-vue'
import { findCampusAuth, getCampusAuthList, reviewCampusAuth, revokeCampusAuth } from '@/api/campusAuth'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'CampusAuth'
})

const appStore = useAppStore()

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const detailShow = ref(false)
const reviewDialogVisible = ref(false)
const reviewLoading = ref(false)
const reviewFormRef = ref()

const createSearchInfo = () => ({
  studentId: '',
  realName: '',
  college: '',
  reviewed: undefined,
  createdAtRange: [],
  reviewedAtRange: []
})

const searchInfo = ref(createSearchInfo())

const createEmptyDetail = () => ({
  id: 0,
  userId: '',
  studentId: '',
  realName: '',
  college: '',
  reviewRemark: '',
  reviewedAt: '',
  reviewedByName: '',
  createdAt: ''
})

const detailForm = ref(createEmptyDetail())
const reviewTarget = ref(createEmptyDetail())
const reviewForm = ref({
  id: 0,
  reviewRemark: ''
})

const reviewRules = {
  reviewRemark: [{ max: 256, message: '审核备注最多 256 个字符', trigger: 'blur' }]
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
  if (!params.reviewedAtRange?.length) {
    delete params.reviewedAtRange
  }
  if (typeof params.reviewed === 'undefined') {
    delete params.reviewed
  }

  const table = await getCampusAuthList(params)
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
  const res = await findCampusAuth({ id: row.id })
  if (res.code === 0) {
    detailForm.value = res.data
    detailShow.value = true
  }
}

const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = createEmptyDetail()
}

const openReviewDialog = (row) => {
  reviewTarget.value = { ...row }
  reviewForm.value = {
    id: row.id,
    reviewRemark: row.reviewRemark || ''
  }
  reviewDialogVisible.value = true
}

const closeReviewDialog = () => {
  reviewDialogVisible.value = false
  reviewLoading.value = false
  reviewTarget.value = createEmptyDetail()
  reviewForm.value = {
    id: 0,
    reviewRemark: ''
  }
  reviewFormRef.value?.clearValidate()
}

const submitReview = async () => {
  const valid = await reviewFormRef.value.validate().catch(() => false)
  if (!valid) {
    return
  }
  reviewLoading.value = true
  const res = await reviewCampusAuth(reviewForm.value).finally(() => {
    reviewLoading.value = false
  })
  if (res.code === 0) {
    ElMessage.success('审核成功')
    closeReviewDialog()
    getTableData()
  }
}

const handleRevoke = async (row) => {
  try {
    await ElMessageBox.confirm(`确定撤回 ${row.realName} 的校园身份审核结果吗？`, '撤回确认', {
      type: 'warning'
    })
  } catch (e) {
    return
  }

  const res = await revokeCampusAuth({ id: row.id })
  if (res.code === 0) {
    ElMessage.success('撤回成功')
    getTableData()
  }
}

getTableData()
</script>
