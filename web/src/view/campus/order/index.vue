<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="订单号">
          <el-input v-model="searchInfo.orderNo" placeholder="请输入订单号" />
        </el-form-item>
        <el-form-item label="买家ID">
          <el-input v-model="searchInfo.buyerId" placeholder="请输入买家ID" />
        </el-form-item>
        <el-form-item label="卖家ID">
          <el-input v-model="searchInfo.sellerId" placeholder="请输入卖家ID" />
        </el-form-item>
        <el-form-item label="商品ID">
          <el-input v-model="searchInfo.productId" placeholder="请输入商品ID" />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部状态" style="width: 140px">
            <el-option v-for="item in orderStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
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
        <el-table-column align="left" label="订单号" prop="orderNo" min-width="220" show-overflow-tooltip />
        <el-table-column align="left" label="商品标题" prop="productTitle" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="买家" min-width="170">
          <template #default="scope">
            {{ scope.row.buyerNickname || '-' }} / {{ scope.row.buyerPhone || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="卖家" min-width="170">
          <template #default="scope">
            {{ scope.row.sellerNickname || '-' }} / {{ scope.row.sellerPhone || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="价格" min-width="110">
          <template #default="scope">
            {{ formatPrice(scope.row.price) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" min-width="120">
          <template #default="scope">
            <el-tag :type="orderStatusTagMap[scope.row.status] || 'info'">
              {{ scope.row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="关闭原因" min-width="220" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.closeReason || '-' }}
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
      title="交易信息详情"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="订单ID">{{ detailForm.id }}</el-descriptions-item>
        <el-descriptions-item label="订单号">{{ detailForm.orderNo || '-' }}</el-descriptions-item>
        <el-descriptions-item label="商品信息">
          <div>{{ detailForm.productTitle || '-' }}</div>
          <div class="mt-2">
            <el-image
              v-if="detailForm.productImage"
              :src="detailForm.productImage"
              fit="cover"
              style="width: 96px; height: 96px; border-radius: 8px"
              :preview-src-list="[detailForm.productImage]"
              preview-teleported
            />
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="买家">
          {{ detailForm.buyerNickname || '-' }} / {{ detailForm.buyerPhone || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="卖家">
          {{ detailForm.sellerNickname || '-' }} / {{ detailForm.sellerPhone || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="成交价格">{{ formatPrice(detailForm.price) }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">{{ detailForm.statusText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="订单备注">{{ detailForm.remark || '-' }}</el-descriptions-item>
        <el-descriptions-item label="关闭原因">{{ detailForm.closeReason || '-' }}</el-descriptions-item>
        <el-descriptions-item label="关闭方">{{ detailForm.closeByText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="关闭确认">{{ detailForm.closeConfirmedText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="确认时间">{{ formatDate(detailForm.confirmedAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ formatDate(detailForm.completedAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="取消时间">{{ formatDate(detailForm.cancelledAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(detailForm.createdAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(detailForm.updatedAt) || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { findCampusOrder, getCampusOrderList } from '@/api/campusOrder'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'CampusOrder'
})

const appStore = useAppStore()

const orderStatusOptions = [
  { value: 1, label: '待付款' },
  { value: 2, label: '待确认' },
  { value: 3, label: '已完成' },
  { value: 4, label: '已取消' },
  { value: 5, label: '已关闭' }
]

const orderStatusTagMap = {
  1: 'warning',
  2: 'warning',
  3: 'success',
  4: 'info',
  5: 'danger'
}

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const detailShow = ref(false)

const createSearchInfo = () => ({
  orderNo: '',
  buyerId: '',
  sellerId: '',
  productId: '',
  status: undefined,
  createdAtRange: []
})

const createDetail = () => ({
  id: 0,
  orderNo: '',
  productTitle: '',
  productImage: '',
  buyerNickname: '',
  buyerPhone: '',
  sellerNickname: '',
  sellerPhone: '',
  price: 0,
  statusText: '',
  remark: '',
  closeReason: '',
  closeByText: '',
  closeConfirmedText: '',
  confirmedAt: '',
  completedAt: '',
  cancelledAt: '',
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
  if (typeof params.status === 'undefined') {
    delete params.status
  }

  const table = await getCampusOrderList(params)
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
  const res = await findCampusOrder({ id: row.id })
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

const formatPrice = (price) => {
  if (price === null || typeof price === 'undefined') {
    return '-'
  }
  return `￥${Number(price).toFixed(2)}`
}

getTableData()
</script>
