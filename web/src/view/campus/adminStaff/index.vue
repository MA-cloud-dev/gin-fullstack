<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="ID">
          <el-input v-model="searchInfo.id" placeholder="请输入ID" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchInfo.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="显示名">
          <el-input v-model="searchInfo.displayName" placeholder="请输入显示名" />
        </el-form-item>
        <el-form-item label="角色类型">
          <el-select v-model="searchInfo.roleType" clearable placeholder="全部角色" style="width: 140px">
            <el-option label="超级管理员" :value="1" />
            <el-option label="运营管理员" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部状态" style="width: 140px">
            <el-option label="启用" :value="0" />
            <el-option label="禁用" :value="1" />
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
        <el-table-column align="left" label="用户名" prop="username" min-width="150" />
        <el-table-column align="left" label="显示名" prop="displayName" min-width="140" />
        <el-table-column align="left" label="角色类型" min-width="140" prop="roleTypeText" />
        <el-table-column align="left" label="状态" min-width="120">
          <template #default="scope">
            <el-tag :type="scope.row.status === 0 ? 'success' : 'info'">
              {{ scope.row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="最近登录时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.lastLoginAt) || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">查看</el-button>
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
      title="B端管理员详情"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="ID">{{ detailForm.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ detailForm.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="显示名">{{ detailForm.displayName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="角色类型">{{ detailForm.roleTypeText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ detailForm.statusText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="最近登录时间">{{ formatDate(detailForm.lastLoginAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(detailForm.createdAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(detailForm.updatedAt) || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { findCampusAdminStaff, getCampusAdminStaffList, updateCampusAdminStaffStatus } from '@/api/campusAdminStaff'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'CampusAdminStaff'
})

const appStore = useAppStore()

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const detailShow = ref(false)

const createSearchInfo = () => ({
  id: '',
  username: '',
  displayName: '',
  roleType: undefined,
  status: undefined,
  createdAtRange: []
})

const createDetail = () => ({
  id: 0,
  username: '',
  displayName: '',
  roleTypeText: '',
  statusText: '',
  lastLoginAt: '',
  createdAt: '',
  updatedAt: ''
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
  if (typeof params.roleType === 'undefined') {
    delete params.roleType
  }
  if (typeof params.status === 'undefined') {
    delete params.status
  }

  const table = await getCampusAdminStaffList(params)
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
  const res = await findCampusAdminStaff({ id: row.id })
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

const handleStatus = async (row) => {
  const targetStatus = row.status === 0 ? 1 : 0
  const actionText = targetStatus === 0 ? '启用' : '禁用'
  await ElMessageBox.confirm(`确定要${actionText}管理员【${row.displayName || row.username}】吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  const res = await updateCampusAdminStaffStatus({
    id: row.id,
    status: targetStatus
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
