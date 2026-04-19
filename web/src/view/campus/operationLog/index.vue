<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="操作人">
          <el-input v-model="searchInfo.operatorKeyword" placeholder="请输入账号或昵称" />
        </el-form-item>
        <el-form-item label="操作来源">
          <el-select v-model="searchInfo.operatorSource" clearable placeholder="全部来源" style="width: 140px">
            <el-option label="B端页面" value="web" />
            <el-option label="CLI工具" value="cli" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作模块">
          <el-select v-model="searchInfo.module" clearable placeholder="全部模块" style="width: 160px">
            <el-option v-for="item in moduleOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="动作">
          <el-select v-model="searchInfo.action" clearable filterable placeholder="全部动作" style="width: 180px">
            <el-option v-for="item in actionOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标ID">
          <el-input v-model="searchInfo.targetId" placeholder="请输入目标ID" />
        </el-form-item>
        <el-form-item label="操作时间">
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
        <el-table-column align="left" label="操作时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作人" min-width="180">
          <template #default="scope">
            {{ scope.row.operatorNickname || scope.row.operatorUsername || '-' }}
            <span v-if="scope.row.operatorUsername">({{ scope.row.operatorUsername }})</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="来源" prop="operatorSourceText" min-width="120" />
        <el-table-column align="left" label="模块" prop="moduleText" min-width="140" />
        <el-table-column align="left" label="动作" prop="actionText" min-width="180" />
        <el-table-column align="left" label="目标" min-width="180" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.targetLabel || `ID:${scope.row.targetId}` }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="执行原因" prop="reason" min-width="220" show-overflow-tooltip />
        <el-table-column align="left" label="执行结果" prop="result" min-width="220" show-overflow-tooltip />
        <el-table-column align="left" fixed="right" label="操作" width="100">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">查看</el-button>
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
      title="审核操作记录详情"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="记录ID">{{ detailForm.id }}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ formatDate(detailForm.createdAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="操作人">
          {{ detailForm.operatorNickname || detailForm.operatorUsername || '-' }}
          <span v-if="detailForm.operatorUsername">({{ detailForm.operatorUsername }})</span>
        </el-descriptions-item>
        <el-descriptions-item label="操作来源">{{ detailForm.operatorSourceText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="操作模块">{{ detailForm.moduleText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="执行动作">{{ detailForm.actionText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="目标对象">
          {{ detailForm.targetLabel || `ID:${detailForm.targetId || '-'}` }}
        </el-descriptions-item>
        <el-descriptions-item label="执行原因">{{ detailForm.reason || '-' }}</el-descriptions-item>
        <el-descriptions-item label="执行结果">{{ detailForm.result || '-' }}</el-descriptions-item>
        <el-descriptions-item label="请求路径">{{ detailForm.requestPath || '-' }}</el-descriptions-item>
        <el-descriptions-item label="请求方法">{{ detailForm.requestMethod || '-' }}</el-descriptions-item>
        <el-descriptions-item label="操作IP">{{ detailForm.operatorIp || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getCampusOperationLogList, findCampusOperationLog } from '@/api/campusOperationLog'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'CampusOperationLog'
})

const appStore = useAppStore()

const moduleOptions = [
  { value: 'auth', label: '身份审核' },
  { value: 'user', label: '用户' },
  { value: 'staff', label: 'B端管理' },
  { value: 'product', label: '商品' },
  { value: 'report', label: '举报' },
  { value: 'category', label: '分类' },
  { value: 'announcement', label: '公告' }
]

const actionOptions = [
  { value: 'approve_auth', label: '通过校园审核' },
  { value: 'revoke_auth', label: '撤回校园审核' },
  { value: 'enable_user', label: '启用用户' },
  { value: 'disable_user', label: '禁用用户' },
  { value: 'enable_staff', label: '启用管理员' },
  { value: 'disable_staff', label: '禁用管理员' },
  { value: 'set_product_status', label: '调整商品状态' },
  { value: 'handle_report', label: '处理举报' },
  { value: 'create_category', label: '新增分类' },
  { value: 'update_category', label: '编辑分类' },
  { value: 'enable_category', label: '启用分类' },
  { value: 'disable_category', label: '停用分类' },
  { value: 'create_announcement', label: '新增公告' },
  { value: 'update_announcement', label: '编辑公告' },
  { value: 'publish_announcement', label: '上线公告' },
  { value: 'unpublish_announcement', label: '下线公告' },
  { value: 'delete_announcement', label: '删除公告' }
]

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const detailShow = ref(false)

const createSearchInfo = () => ({
  operatorKeyword: '',
  operatorSource: '',
  module: '',
  action: '',
  targetId: '',
  createdAtRange: []
})

const createDetail = () => ({
  id: 0,
  operatorUsername: '',
  operatorNickname: '',
  operatorSourceText: '',
  moduleText: '',
  actionText: '',
  targetId: 0,
  targetLabel: '',
  reason: '',
  result: '',
  requestPath: '',
  requestMethod: '',
  operatorIp: '',
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
  if (!params.operatorSource) {
    delete params.operatorSource
  }
  if (!params.module) {
    delete params.module
  }
  if (!params.action) {
    delete params.action
  }
  if (!params.operatorKeyword) {
    delete params.operatorKeyword
  }
  if (!params.targetId) {
    delete params.targetId
  }

  const table = await getCampusOperationLogList(params)
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
  const res = await findCampusOperationLog({ id: row.id })
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

getTableData()
</script>
