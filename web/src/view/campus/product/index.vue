<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="商品ID">
          <el-input v-model="searchInfo.id" placeholder="请输入商品ID" />
        </el-form-item>
        <el-form-item label="商品标题">
          <el-input v-model="searchInfo.title" placeholder="请输入商品标题" />
        </el-form-item>
        <el-form-item label="发布人ID">
          <el-input v-model="searchInfo.userId" placeholder="请输入发布人ID" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchInfo.categoryId" clearable filterable placeholder="全部分类" style="width: 180px">
            <el-option
              v-for="item in categoryOptions"
              :key="item.id"
              :label="item.label"
              :value="String(item.id)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="商品状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部状态" style="width: 140px">
            <el-option v-for="item in productStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="交易方式">
          <el-select v-model="searchInfo.tradeMode" clearable placeholder="全部方式" style="width: 140px">
            <el-option v-for="item in tradeModeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="发布时间">
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
        <el-table-column align="left" label="封面" width="100">
          <template #default="scope">
            <el-image
              v-if="scope.row.coverUrl"
              :src="scope.row.coverUrl"
              fit="cover"
              style="width: 64px; height: 64px; border-radius: 8px"
              :preview-src-list="[scope.row.coverUrl]"
              preview-teleported
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="标题" prop="title" min-width="220" show-overflow-tooltip />
        <el-table-column align="left" label="价格" min-width="120">
          <template #default="scope">
            {{ formatPrice(scope.row.price) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="分类" min-width="140">
          <template #default="scope">
            {{ scope.row.categoryName || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="发布人" min-width="180">
          <template #default="scope">
            {{ scope.row.publisherNickname || '-' }} / {{ scope.row.publisherPhone || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" min-width="120">
          <template #default="scope">
            <el-tag :type="productStatusTagMap[scope.row.status] || 'info'">
              {{ scope.row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="交易方式" min-width="120" prop="tradeModeText" />
        <el-table-column align="left" label="浏览量" prop="viewCount" width="100" />
        <el-table-column align="left" label="想要数" prop="wantCount" width="100" />
        <el-table-column align="left" label="过期时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.expireAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="发布时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              查看
            </el-button>
            <el-button type="warning" link class="table-button" @click="openStatusDialog(scope.row)">
              状态调整
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
      title="商品详情"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="商品ID">{{ detailForm.id }}</el-descriptions-item>
        <el-descriptions-item label="标题">{{ detailForm.title || '-' }}</el-descriptions-item>
        <el-descriptions-item label="分类">{{ detailForm.categoryName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="发布人">
          {{ detailForm.publisherNickname || '-' }} / {{ detailForm.publisherPhone || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="价格">{{ formatPrice(detailForm.price) }}</el-descriptions-item>
        <el-descriptions-item label="原价">
          {{ detailForm.originalPrice ? formatPrice(detailForm.originalPrice) : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="productStatusTagMap[detailForm.status] || 'info'">
            {{ detailForm.statusText || '-' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="交易方式">{{ detailForm.tradeModeText || '-' }}</el-descriptions-item>
        <el-descriptions-item label="联系方式">{{ detailForm.contactInfo || '-' }}</el-descriptions-item>
        <el-descriptions-item label="商品描述">{{ detailForm.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="浏览量">{{ detailForm.viewCount ?? '-' }}</el-descriptions-item>
        <el-descriptions-item label="想要数">{{ detailForm.wantCount ?? '-' }}</el-descriptions-item>
        <el-descriptions-item label="过期时间">{{ formatDate(detailForm.expireAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="发布时间">{{ formatDate(detailForm.createdAt) || '-' }}</el-descriptions-item>
      </el-descriptions>

      <div class="mt-4">
        <div class="text-base font-bold mb-3">商品图片</div>
        <div v-if="detailForm.images?.length" class="flex flex-wrap gap-3">
          <el-image
            v-for="item in detailForm.images"
            :key="item.id"
            :src="item.imageUrl"
            fit="cover"
            style="width: 120px; height: 120px; border-radius: 8px"
            :preview-src-list="detailImageList"
            preview-teleported
          />
        </div>
        <el-empty v-else description="暂无图片" />
      </div>
    </el-drawer>

    <el-dialog v-model="statusDialogVisible" title="调整商品状态" width="420px">
      <el-form ref="statusFormRef" :model="statusForm" :rules="statusRules" label-width="80px">
        <el-form-item label="商品ID">
          <span>{{ statusForm.id }}</span>
        </el-form-item>
        <el-form-item label="新状态" prop="status">
          <el-select v-model="statusForm.status" style="width: 100%">
            <el-option v-for="item in productStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="执行原因" prop="auditReason">
          <el-input
            v-model="statusForm.auditReason"
            type="textarea"
            :rows="4"
            maxlength="256"
            show-word-limit
            placeholder="请输入状态调整原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="statusDialogVisible = false">取 消</el-button>
          <el-button type="primary" :loading="statusLoading" @click="submitStatus">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { findCampusProduct, getCampusProductList, updateCampusProductStatus } from '@/api/campusProduct'
import { getCampusCategoryList } from '@/api/campusCategory'
import { formatDate } from '@/utils/format'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'CampusProduct'
})

const appStore = useAppStore()

const productStatusOptions = [
  { value: 0, label: '在售' },
  { value: 2, label: '交易中' },
  { value: 3, label: '已下架' }
]

const tradeModeOptions = [
  { value: 0, label: '线下自提' },
  { value: 1, label: '快递邮寄' }
]

const productStatusTagMap = {
  0: 'success',
  2: 'warning',
  3: 'info'
}

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const categoryTree = ref([])
const detailShow = ref(false)
const statusDialogVisible = ref(false)
const statusLoading = ref(false)
const statusFormRef = ref()

const createSearchInfo = () => ({
  id: '',
  title: '',
  userId: '',
  categoryId: '',
  status: undefined,
  tradeMode: undefined,
  createdAtRange: []
})

const searchInfo = ref(createSearchInfo())

const createEmptyDetail = () => ({
  id: 0,
  title: '',
  categoryName: '',
  publisherNickname: '',
  publisherPhone: '',
  price: 0,
  originalPrice: null,
  status: 0,
  statusText: '',
  tradeModeText: '',
  contactInfo: '',
  description: '',
  viewCount: 0,
  wantCount: 0,
  expireAt: '',
  createdAt: '',
  images: []
})

const detailForm = ref(createEmptyDetail())
const statusForm = ref({
  id: 0,
  status: 0,
  auditReason: ''
})

const statusRules = {
  status: [{ required: true, message: '请选择新状态', trigger: 'change' }],
  auditReason: [{ required: true, message: '请输入状态调整原因', trigger: 'blur' }]
}

const categoryOptions = computed(() => flattenCategoryTree(categoryTree.value))
const detailImageList = computed(() => (detailForm.value.images || []).map((item) => item.imageUrl))

const loadCategories = async () => {
  const res = await getCampusCategoryList()
  if (res.code === 0) {
    categoryTree.value = res.data || []
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
  if (!params.categoryId) {
    delete params.categoryId
  }
  if (typeof params.status === 'undefined') {
    delete params.status
  }
  if (typeof params.tradeMode === 'undefined') {
    delete params.tradeMode
  }

  const table = await getCampusProductList(params)
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
  const res = await findCampusProduct({ id: row.id })
  if (res.code === 0) {
    detailForm.value = {
      ...createEmptyDetail(),
      ...res.data
    }
    detailShow.value = true
  }
}

const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = createEmptyDetail()
}

const openStatusDialog = (row) => {
  statusForm.value = {
    id: row.id,
    status: row.status,
    auditReason: ''
  }
  statusDialogVisible.value = true
}

const submitStatus = async () => {
  const valid = await statusFormRef.value?.validate().catch(() => false)
  if (!valid) {
    return
  }
  const currentID = statusForm.value.id
  statusLoading.value = true
  const res = await updateCampusProductStatus(statusForm.value)
  statusLoading.value = false
  if (res.code === 0) {
    ElMessage.success('状态更新成功')
    statusDialogVisible.value = false
    statusForm.value = {
      id: 0,
      status: 0,
      auditReason: ''
    }
    getTableData()
    if (detailShow.value && detailForm.value.id === currentID) {
      getDetails({ id: currentID })
    }
  }
}

const flattenCategoryTree = (tree, prefix = '') => {
  return tree.flatMap((item) => {
    const label = prefix ? `${prefix} / ${item.name}` : item.name
    const current = [{ id: item.id, label }]
    if (!item.children?.length) {
      return current
    }
    return current.concat(flattenCategoryTree(item.children, label))
  })
}

const formatPrice = (price) => {
  if (price === null || typeof price === 'undefined') {
    return '-'
  }
  return `￥${Number(price).toFixed(2)}`
}

loadCategories()
getTableData()
</script>
