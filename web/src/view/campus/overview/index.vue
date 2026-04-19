<template>
  <div class="campus-overview" v-loading="loading">
    <div class="overview-header">
      <div>
        <h2 class="overview-title">校园数据总揽</h2>
        <p class="overview-subtitle">聚合展示校园模块总量、近 7 天趋势和状态分布</p>
      </div>
      <el-tag type="info" effect="plain">更新时间：{{ formatDateTime(overview.generatedAt) || '-' }}</el-tag>
    </div>

    <div class="kpi-grid">
      <el-card v-for="item in kpiCards" :key="item.key" shadow="hover" class="kpi-card">
        <div class="kpi-label">{{ item.label }}</div>
        <div class="kpi-value">{{ formatNumber(item.value) }}</div>
        <div class="kpi-desc">{{ item.desc }}</div>
      </el-card>
    </div>

    <div class="chart-grid">
      <el-card shadow="never" class="chart-card">
        <template #header>
          <div class="card-header">
            <span>近 7 天综合趋势</span>
          </div>
        </template>
        <base-chart :options="trendChartOptions" height="360px" />
      </el-card>

      <el-card shadow="never" class="chart-card">
        <template #header>
          <div class="card-header dist-header">
            <span>状态分布</span>
            <el-tabs v-model="activeDistribution" class="distribution-tabs">
              <el-tab-pane v-for="item in distributionTabs" :key="item.key" :label="item.label" :name="item.key" />
            </el-tabs>
          </div>
        </template>
        <base-chart :options="distributionChartOptions" height="360px" />
      </el-card>
    </div>

    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>模块总量表</span>
        </div>
      </template>
      <el-table :data="overview.moduleTotals" border>
        <el-table-column prop="module" label="模块" min-width="140" />
        <el-table-column label="总量" min-width="110" align="center">
          <template #default="{ row }">
            {{ formatNumber(row.total) }}
          </template>
        </el-table-column>
        <el-table-column label="核心状态1" min-width="160">
          <template #default="{ row }">
            {{ formatBucketCell(row.coreStatus1Label, row.coreStatus1Count) }}
          </template>
        </el-table-column>
        <el-table-column label="核心状态2" min-width="160">
          <template #default="{ row }">
            {{ formatBucketCell(row.coreStatus2Label, row.coreStatus2Count) }}
          </template>
        </el-table-column>
        <el-table-column label="核心状态3" min-width="160">
          <template #default="{ row }">
            {{ formatBucketCell(row.coreStatus3Label, row.coreStatus3Count) }}
          </template>
        </el-table-column>
        <el-table-column label="近7天新增" min-width="160">
          <template #default="{ row }">
            {{ formatBucketCell(row.last7DaysLabel, row.last7DaysCount) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>近 7 天明细表</span>
        </div>
      </template>
      <el-table :data="overview.dailyTrend" border>
        <el-table-column prop="date" label="日期" min-width="120" />
        <el-table-column label="新增用户" min-width="110" align="center">
          <template #default="{ row }">{{ formatNumber(row.newUsers) }}</template>
        </el-table-column>
        <el-table-column label="认证申请" min-width="110" align="center">
          <template #default="{ row }">{{ formatNumber(row.authApplications) }}</template>
        </el-table-column>
        <el-table-column label="认证通过" min-width="110" align="center">
          <template #default="{ row }">{{ formatNumber(row.authApproved) }}</template>
        </el-table-column>
        <el-table-column label="新增商品" min-width="110" align="center">
          <template #default="{ row }">{{ formatNumber(row.newProducts) }}</template>
        </el-table-column>
        <el-table-column label="新增订单" min-width="110" align="center">
          <template #default="{ row }">{{ formatNumber(row.newOrders) }}</template>
        </el-table-column>
        <el-table-column label="新增举报" min-width="110" align="center">
          <template #default="{ row }">{{ formatNumber(row.newReports) }}</template>
        </el-table-column>
        <el-table-column label="新增公告" min-width="110" align="center">
          <template #default="{ row }">{{ formatNumber(row.newAnnouncements) }}</template>
        </el-table-column>
        <el-table-column label="操作记录" min-width="110" align="center">
          <template #default="{ row }">{{ formatNumber(row.operationLogs) }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import BaseChart from '@/components/charts/index.vue'
import { getCampusOverview } from '@/api/campusOverview'

defineOptions({
  name: 'CampusOverview'
})

const loading = ref(false)
const activeDistribution = ref('productStatus')
const overview = ref({
  kpis: {
    userTotal: 0,
    verifiedUserTotal: 0,
    pendingAuthTotal: 0,
    productOnSaleTotal: 0,
    orderTotal: 0,
    orderCompletedTotal: 0,
    pendingReportTotal: 0,
    announcementOnlineTotal: 0
  },
  moduleTotals: [],
  dailyTrend: [],
  distributions: {
    userAuthStatus: [],
    productStatus: [],
    orderStatus: [],
    reportStatus: [],
    announcementStatus: [],
    staffStatus: []
  },
  generatedAt: ''
})

const distributionTabs = [
  { key: 'productStatus', label: '商品状态' },
  { key: 'userAuthStatus', label: '用户认证' },
  { key: 'orderStatus', label: '订单状态' },
  { key: 'reportStatus', label: '举报状态' },
  { key: 'announcementStatus', label: '公告状态' },
  { key: 'staffStatus', label: '管理员状态' }
]

const kpiCards = computed(() => [
  { key: 'userTotal', label: '用户总量', value: overview.value.kpis.userTotal, desc: '校园模块累计注册用户' },
  { key: 'verifiedUserTotal', label: '已认证用户', value: overview.value.kpis.verifiedUserTotal, desc: '已通过校园认证的用户' },
  { key: 'pendingAuthTotal', label: '待审核认证', value: overview.value.kpis.pendingAuthTotal, desc: '当前待处理的认证申请' },
  { key: 'productOnSaleTotal', label: '在售商品', value: overview.value.kpis.productOnSaleTotal, desc: '当前可交易的商品数量' },
  { key: 'orderTotal', label: '订单总量', value: overview.value.kpis.orderTotal, desc: '校园模块累计订单数' },
  { key: 'orderCompletedTotal', label: '已完成订单', value: overview.value.kpis.orderCompletedTotal, desc: '已完成成交的订单数' },
  { key: 'pendingReportTotal', label: '待处理举报', value: overview.value.kpis.pendingReportTotal, desc: '当前待处理的举报记录' },
  { key: 'announcementOnlineTotal', label: '上线公告', value: overview.value.kpis.announcementOnlineTotal, desc: '当前处于上线状态的公告' }
])

const trendChartOptions = computed(() => ({
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    top: 0
  },
  grid: {
    left: 24,
    right: 16,
    bottom: 24,
    top: 48,
    containLabel: true
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: overview.value.dailyTrend.map((item) => item.date)
  },
  yAxis: {
    type: 'value'
  },
  series: [
    { name: '新增用户', type: 'line', smooth: true, data: overview.value.dailyTrend.map((item) => item.newUsers) },
    { name: '认证申请', type: 'line', smooth: true, data: overview.value.dailyTrend.map((item) => item.authApplications) },
    { name: '新增商品', type: 'line', smooth: true, data: overview.value.dailyTrend.map((item) => item.newProducts) },
    { name: '新增订单', type: 'line', smooth: true, data: overview.value.dailyTrend.map((item) => item.newOrders) },
    { name: '新增举报', type: 'line', smooth: true, data: overview.value.dailyTrend.map((item) => item.newReports) },
    { name: '操作记录', type: 'line', smooth: true, data: overview.value.dailyTrend.map((item) => item.operationLogs) }
  ]
}))

const currentDistribution = computed(() => overview.value.distributions[activeDistribution.value] || [])

const distributionChartOptions = computed(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow'
    }
  },
  grid: {
    left: 24,
    right: 16,
    bottom: 24,
    top: 24,
    containLabel: true
  },
  xAxis: {
    type: 'category',
    data: currentDistribution.value.map((item) => item.label)
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      type: 'bar',
      barWidth: 36,
      data: currentDistribution.value.map((item) => item.count),
      itemStyle: {
        borderRadius: [8, 8, 0, 0]
      }
    }
  ]
}))

const formatNumber = (value) => Number(value || 0).toLocaleString('zh-CN')

const formatDateTime = (value) => {
  if (!value) {
    return ''
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  return date.toLocaleString('zh-CN', { hour12: false })
}

const formatBucketCell = (label, count) => {
  if (!label || label === '-') {
    return '-'
  }
  return `${label}：${formatNumber(count)}`
}

const fetchOverview = async () => {
  loading.value = true
  try {
    const res = await getCampusOverview()
    if (res.code === 0) {
      overview.value = {
        ...overview.value,
        ...res.data
      }
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchOverview()
})
</script>

<style scoped lang="scss">
.campus-overview {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.overview-header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
  flex-wrap: wrap;
}

.overview-title {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
}

.overview-subtitle {
  margin: 8px 0 0;
  color: #6b7280;
}

.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.kpi-card,
.chart-card,
.table-card {
  border: 1px solid #e5e7eb;
}

.kpi-label {
  color: #6b7280;
  font-size: 14px;
}

.kpi-value {
  margin-top: 12px;
  font-size: 32px;
  line-height: 1;
  font-weight: 700;
  color: #111827;
}

.kpi-desc {
  margin-top: 12px;
  color: #9ca3af;
  font-size: 13px;
}

.chart-grid {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(340px, 1fr);
  gap: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-weight: 600;
}

.dist-header {
  align-items: flex-start;
}

.distribution-tabs {
  margin: -8px 0 -12px;
}

@media (max-width: 1280px) {
  .kpi-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .chart-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .kpi-grid {
    grid-template-columns: 1fr;
  }
}
</style>
