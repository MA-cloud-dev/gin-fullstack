<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="用户ID">
          <el-input v-model="searchInfo.id" placeholder="请输入用户ID" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="searchInfo.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchInfo.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="searchInfo.role" clearable placeholder="全部角色" style="width: 140px">
            <el-option label="普通用户" :value="0" />
            <el-option label="管理员" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部状态" style="width: 140px">
            <el-option label="启用" :value="0" />
            <el-option label="禁用" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="认证状态">
          <el-select v-model="searchInfo.authStatus" clearable placeholder="全部状态" style="width: 140px">
            <el-option label="未认证" :value="0" />
            <el-option label="已拒绝" :value="1" />
            <el-option label="审核中" :value="2" />
            <el-option label="已认证" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="学号">
          <el-input v-model="searchInfo.studentId" placeholder="请输入学号" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="searchInfo.realName" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="学院">
          <el-input v-model="searchInfo.college" placeholder="请输入学院" />
        </el-form-item>
        <el-form-item label="注册时间">
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
        <el-table-column align="left" label="手机号" prop="phone" min-width="150" />
        <el-table-column align="left" label="昵称" min-width="140">
          <template #default="scope">
            {{ scope.row.nickname || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="角色" min-width="120" prop="roleText" />
        <el-table-column align="left" label="状态" min-width="120">
          <template #default="scope">
            <el-tag :type="scope.row.status === 0 ? 'success' : 'info'">
              {{ scope.row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="认证状态" min-width="120">
          <template #default="scope">
            <el-tag :type="getAuthStatusTagType(scope.row.authStatusText)">
              {{ scope.row.authStatusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="学号" min-width="150">
          <template #default="scope">
            {{ scope.row.studentId || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="姓名" min-width="120">
          <template #default="scope">
            {{ scope.row.realName || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="学院" min-width="180">
          <template #default="scope">
            {{ scope.row.college || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="注册时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">查看</el-button>
            <el-button v-if="canReviewAuth(scope.row)" type="success" link class="table-button" @click="handleApproveAuth(scope.row)">
              审核通过
            </el-button>
            <el-button v-if="canReviewAuth(scope.row)" type="danger" link class="table-button" @click="handleRejectAuth(scope.row)">
              审核拒绝
            </el-button>
            <el-button type="warning" link class="table-button" @click="handleStatus(scope.row)">
              {{ scope.row.status === 0 ? '禁用' : '启用' }}
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
      title="校园用户详情"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="用户ID">{{ detailForm.id }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ detailForm.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="昵称">{{ detailForm.nickname || '-' }}</el-descriptions-item>
        <el-descriptions-item label="角色">{{ detailForm.roleText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ detailForm.statusText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="认证状态">{{ detailForm.authStatusText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="学号">{{ detailForm.studentId || '-' }}</el-descriptions-item>
        <el-descriptions-item label="姓名">{{ detailForm.realName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="学院">{{ detailForm.college || '-' }}</el-descriptions-item>
        <el-descriptions-item label="年级">{{ detailForm.grade || '-' }}</el-descriptions-item>
        <el-descriptions-item label="宿舍">{{ detailForm.dormitory || '-' }}</el-descriptions-item>
        <el-descriptions-item label="微信号">{{ detailForm.wechatId || '-' }}</el-descriptions-item>
        <el-descriptions-item label="认证备注">{{ detailForm.reviewRemark || '-' }}</el-descriptions-item>
        <el-descriptions-item label="审核人">{{ detailForm.reviewedByName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="审核时间">{{ formatDate(detailForm.reviewedAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ formatDate(detailForm.createdAt) || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { findCampusUser, getCampusUserList, updateCampusUserStatus } from '@/api/campusUser'
import { rejectCampusAuth, reviewCampusAuth } from '@/api/campusAuth'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'CampusUser'
})

const appStore = useAppStore()

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const detailShow = ref(false)

const createSearchInfo = () => ({
  id: '',
  phone: '',
  nickname: '',
  role: undefined,
  status: undefined,
  authStatus: undefined,
  studentId: '',
  realName: '',
  college: '',
  createdAtRange: []
})

const createDetail = () => ({
  id: 0,
  phone: '',
  nickname: '',
  roleText: '',
  statusText: '',
  authStatusText: '',
  authRecordId: undefined,
  studentId: '',
  realName: '',
  college: '',
  grade: '',
  dormitory: '',
  wechatId: '',
  reviewRemark: '',
  reviewedByName: '',
  reviewedAt: '',
  createdAt: ''
})

const searchInfo = ref(createSearchInfo())
const detailForm = ref(createDetail())

const getTableData = async () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  }
  if (!params.createdAtRange?.length) {
    delete params.createdAtRange
  }
  if (typeof params.role === 'undefined') {
    delete params.role
  }
  if (typeof params.status === 'undefined') {
    delete params.status
  }
  if (typeof params.authStatus === 'undefined') {
    delete params.authStatus
  }

  const table = await getCampusUserList(params)
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
  const res = await findCampusUser({ id: row.id })
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

const getAuthStatusTagType = (authStatusText) => {
  switch (authStatusText) {
    case '已认证':
      return 'success'
    case '已拒绝':
      return 'danger'
    case '审核中':
      return 'warning'
    default:
      return 'info'
  }
}

const canReviewAuth = (row) => {
  return Boolean(row.authRecordId) && row.authStatusText === '审核中'
}

const handleApproveAuth = async (row) => {
  const displayName = row.nickname || row.phone || `ID:${row.id}`
  let auditReason = ''

  try {
    const promptResult = await ElMessageBox.prompt(`确定通过用户【${displayName}】的校园认证吗？`, '审核通过', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputType: 'textarea',
      inputPlaceholder: '请输入审核备注',
      inputValidator: (value) => {
        const trimmed = value?.trim?.() || ''
        if (!trimmed) {
          return '请输入审核备注'
        }
        if (trimmed.length > 256) {
          return '审核备注最多 256 个字符'
        }
        return true
      }
    })
    auditReason = promptResult.value?.trim?.() || ''
  } catch (e) {
    return
  }

  const res = await reviewCampusAuth({
    id: row.authRecordId,
    reviewRemark: auditReason,
    auditReason
  })
  if (res.code === 0) {
    ElMessage.success('审核成功')
    getTableData()
    if (detailShow.value && detailForm.value.id === row.id) {
      getDetails(row)
    }
  }
}

const handleRejectAuth = async (row) => {
  const displayName = row.nickname || row.phone || `ID:${row.id}`
  let auditReason = ''

  try {
    const promptResult = await ElMessageBox.prompt(`确定拒绝用户【${displayName}】的校园认证吗？`, '审核拒绝', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputType: 'textarea',
      inputPlaceholder: '请输入拒绝原因',
      inputValidator: (value) => {
        const trimmed = value?.trim?.() || ''
        if (!trimmed) {
          return '请输入拒绝原因'
        }
        if (trimmed.length > 256) {
          return '拒绝原因最多 256 个字符'
        }
        return true
      }
    })
    auditReason = promptResult.value?.trim?.() || ''
  } catch (e) {
    return
  }

  const res = await rejectCampusAuth({
    id: row.authRecordId,
    reviewRemark: auditReason,
    auditReason
  })
  if (res.code === 0) {
    ElMessage.success('拒绝成功')
    getTableData()
    if (detailShow.value && detailForm.value.id === row.id) {
      getDetails(row)
    }
  }
}

const handleStatus = async (row) => {
  const targetStatus = row.status === 0 ? 1 : 0
  const actionText = targetStatus === 0 ? '启用' : '禁用'
  let auditReason = ''
  try {
    const promptResult = await ElMessageBox.prompt(`确定要${actionText}用户【${row.nickname || row.phone}】吗？`, '提示', {
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
  const res = await updateCampusUserStatus({
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

getTableData()
</script>
