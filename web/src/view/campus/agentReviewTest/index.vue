<template>
  <div>
    <div class="gva-search-box">
      <el-alert
        title="这是一个 B 端模拟 C 端校园认证申请的测试入口，提交后仍会走真实 Agent 审核链路。"
        type="info"
        :closable="false"
        show-icon
      />
    </div>

    <div class="agent-review-test-grid">
      <div class="gva-table-box">
        <div class="panel-header">
          <div>
            <div class="panel-title">测试人物选择</div>
            <div class="panel-desc">从现有校园用户中选择一个真实用户，模拟 C 端提交认证。</div>
          </div>
          <div class="panel-actions">
            <el-button type="primary" @click="openUserSelector">选择测试人物</el-button>
            <el-button :disabled="!selectedUser.id" @click="clearSelectedUser">清空</el-button>
          </div>
        </div>

        <el-empty v-if="!selectedUser.id" description="暂未选择测试人物" />

        <el-descriptions v-else :column="1" border>
          <el-descriptions-item label="用户ID">{{ selectedUser.id }}</el-descriptions-item>
          <el-descriptions-item label="手机号">{{ selectedUser.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item label="昵称">{{ selectedUser.nickname || '-' }}</el-descriptions-item>
          <el-descriptions-item label="认证状态">
            <el-tag :type="getAuthStatusTagType(selectedUser.authStatus)">
              {{ selectedUser.authStatusText || '-' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="当前学号">{{ selectedUser.studentId || '-' }}</el-descriptions-item>
          <el-descriptions-item label="当前姓名">{{ selectedUser.realName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="当前学院">{{ selectedUser.college || '-' }}</el-descriptions-item>
        </el-descriptions>

        <el-alert
          v-if="selectedUser.id && submitDisabledReason"
          class="mt-4"
          type="warning"
          show-icon
          :closable="false"
          :title="submitDisabledReason"
        />
      </div>

      <div class="gva-table-box">
        <div class="panel-title panel-title--compact">模拟提交表单</div>
        <div class="panel-desc panel-desc--spaced">你可以覆盖学号、姓名、学院，来模拟 C 端重新填写并发起认证申请。</div>

        <el-form ref="submitFormRef" :model="submitForm" :rules="submitRules" label-position="top">
          <el-form-item label="学号" prop="studentId">
            <el-input v-model="submitForm.studentId" maxlength="32" placeholder="请输入学号" />
          </el-form-item>
          <el-form-item label="姓名" prop="realName">
            <el-input v-model="submitForm.realName" maxlength="32" placeholder="请输入姓名" />
          </el-form-item>
          <el-form-item label="学院" prop="college">
            <el-input v-model="submitForm.college" maxlength="64" placeholder="请输入学院" />
          </el-form-item>
          <el-form-item>
            <div class="panel-actions">
              <el-button type="primary" :loading="submitLoading" :disabled="Boolean(submitDisabledReason)" @click="handleSubmit">
                发起 Agent 审核测试
              </el-button>
              <el-button :disabled="!selectedUser.id" @click="fillFromSelectedUser">恢复用户现有信息</el-button>
            </div>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <div v-if="submitResult" class="gva-table-box result-box">
      <div class="panel-header">
        <div>
          <div class="panel-title">本次提交结果</div>
          <div class="panel-desc">这里只展示当前这次提交的即时返回，最终审核结果以审核列表/回调回写为准。</div>
        </div>
        <el-button type="success" @click="goToCampusAuth">去校园身份审核页</el-button>
      </div>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="authRecordId">{{ submitResult.authRecordId }}</el-descriptions-item>
        <el-descriptions-item label="taskId">{{ submitResult.taskId }}</el-descriptions-item>
        <el-descriptions-item label="审核状态">
          <el-tag :type="getReviewStatusTagType(submitResult.reviewStatus)">
            {{ getReviewStatusText(submitResult.reviewStatus) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="用户认证状态">
          <el-tag :type="getAuthStatusTagType(submitResult.authStatus)">
            {{ getAuthStatusText(submitResult.authStatus) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="提交时间">{{ formatDate(submitResult.submittedAt) }}</el-descriptions-item>
        <el-descriptions-item label="返回文案">{{ submitResult.message }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <el-dialog v-model="userSelectorVisible" title="选择测试人物" width="1080px" destroy-on-close>
      <div class="gva-search-box">
        <el-form :inline="true" :model="userSearch" @keyup.enter="fetchUserOptions">
          <el-form-item label="用户ID">
            <el-input v-model="userSearch.id" placeholder="请输入用户ID" />
          </el-form-item>
          <el-form-item label="手机号">
            <el-input v-model="userSearch.phone" placeholder="请输入手机号" />
          </el-form-item>
          <el-form-item label="昵称">
            <el-input v-model="userSearch.nickname" placeholder="请输入昵称" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="search" @click="fetchUserOptions">查询</el-button>
            <el-button icon="refresh" @click="resetUserSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table :data="userOptions" row-key="id" tooltip-effect="dark">
        <el-table-column label="用户ID" prop="id" width="100" />
        <el-table-column label="手机号" prop="phone" min-width="150" />
        <el-table-column label="昵称" min-width="140">
          <template #default="scope">
            {{ scope.row.nickname || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="认证状态" min-width="120">
          <template #default="scope">
            <el-tag :type="getAuthStatusTagType(scope.row.authStatus)">
              {{ scope.row.authStatusText || '-' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="学号" min-width="140">
          <template #default="scope">
            {{ scope.row.studentId || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="姓名" min-width="120">
          <template #default="scope">
            {{ scope.row.realName || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="学院" min-width="180">
          <template #default="scope">
            {{ scope.row.college || '-' }}
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="120">
          <template #default="scope">
            <el-button type="primary" link @click="selectUser(scope.row)">选择</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="userPage"
          :page-size="userPageSize"
          :page-sizes="[10, 20, 50]"
          :total="userTotal"
          @current-change="handleUserPageChange"
          @size-change="handleUserPageSizeChange"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCampusUserList } from '@/api/campusUser'
import { submitCampusAuthTest } from '@/api/campusAuth'
import { formatDate } from '@/utils/format'

defineOptions({
  name: 'CampusAgentReviewTest'
})

const router = useRouter()
const submitFormRef = ref()
const submitLoading = ref(false)
const userSelectorVisible = ref(false)
const userOptions = ref([])
const userPage = ref(1)
const userPageSize = ref(10)
const userTotal = ref(0)
const submitResult = ref(null)

const createSelectedUser = () => ({
  id: 0,
  phone: '',
  nickname: '',
  authStatus: 0,
  authStatusText: '',
  studentId: '',
  realName: '',
  college: ''
})

const createSubmitForm = () => ({
  studentId: '',
  realName: '',
  college: ''
})

const createUserSearch = () => ({
  id: '',
  phone: '',
  nickname: ''
})

const selectedUser = ref(createSelectedUser())
const submitForm = ref(createSubmitForm())
const userSearch = ref(createUserSearch())

const submitRules = {
  studentId: [
    { required: true, message: '请输入学号', trigger: 'blur' },
    { pattern: /^[A-Za-z0-9_-]{4,32}$/, message: '学号仅支持 4-32 位字母、数字、下划线或中划线', trigger: 'blur' }
  ],
  realName: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { max: 32, message: '姓名最多 32 个字符', trigger: 'blur' }
  ],
  college: [
    { required: true, message: '请输入学院', trigger: 'blur' },
    { max: 64, message: '学院最多 64 个字符', trigger: 'blur' }
  ]
}

const submitDisabledReason = computed(() => {
  if (!selectedUser.value.id) {
    return '请先选择一个测试人物'
  }
  if (selectedUser.value.authStatus === 3) {
    return '当前用户已认证，无需重复提交'
  }
  if (selectedUser.value.authStatus === 2) {
    return '当前用户已有审核中的申请或人工复核中的申请'
  }
  return ''
})

const getAuthStatusText = (status) => {
  switch (status) {
    case 1:
      return '已拒绝'
    case 2:
      return '审核中'
    case 3:
      return '已认证'
    default:
      return '未认证'
  }
}

const getReviewStatusText = (status) => {
  switch (status) {
    case 'approved':
      return '已通过'
    case 'rejected':
      return '已拒绝'
    default:
      return '审核中'
  }
}

const getAuthStatusTagType = (status) => {
  switch (status) {
    case 1:
      return 'danger'
    case 2:
      return 'warning'
    case 3:
      return 'success'
    default:
      return 'info'
  }
}

const getReviewStatusTagType = (status) => {
  switch (status) {
    case 'approved':
      return 'success'
    case 'rejected':
      return 'danger'
    default:
      return 'warning'
  }
}

const fillFromSelectedUser = () => {
  submitForm.value = {
    studentId: selectedUser.value.studentId || '',
    realName: selectedUser.value.realName || '',
    college: selectedUser.value.college || ''
  }
}

const openUserSelector = async () => {
  userSelectorVisible.value = true
  userPage.value = 1
  await fetchUserOptions()
}

const fetchUserOptions = async () => {
  const params = {
    page: userPage.value,
    pageSize: userPageSize.value,
    ...userSearch.value,
    role: 0
  }
  if (!params.id) {
    delete params.id
  }
  if (!params.phone) {
    delete params.phone
  }
  if (!params.nickname) {
    delete params.nickname
  }

  const res = await getCampusUserList(params)
  if (res.code === 0) {
    userOptions.value = res.data.list || []
    userTotal.value = res.data.total || 0
    userPage.value = res.data.page || userPage.value
    userPageSize.value = res.data.pageSize || userPageSize.value
  }
}

const resetUserSearch = () => {
  userSearch.value = createUserSearch()
  userPage.value = 1
  fetchUserOptions()
}

const handleUserPageChange = (val) => {
  userPage.value = val
  fetchUserOptions()
}

const handleUserPageSizeChange = (val) => {
  userPageSize.value = val
  userPage.value = 1
  fetchUserOptions()
}

const selectUser = (row) => {
  selectedUser.value = {
    id: row.id,
    phone: row.phone || '',
    nickname: row.nickname || '',
    authStatus: row.authStatus ?? 0,
    authStatusText: row.authStatusText || getAuthStatusText(row.authStatus),
    studentId: row.studentId || '',
    realName: row.realName || '',
    college: row.college || ''
  }
  fillFromSelectedUser()
  submitResult.value = null
  userSelectorVisible.value = false
}

const clearSelectedUser = () => {
  selectedUser.value = createSelectedUser()
  submitForm.value = createSubmitForm()
  submitResult.value = null
}

const handleSubmit = async () => {
  if (submitDisabledReason.value) {
    ElMessage.warning(submitDisabledReason.value)
    return
  }
  const valid = await submitFormRef.value?.validate().catch(() => false)
  if (!valid) {
    return
  }

  submitLoading.value = true
  const res = await submitCampusAuthTest({
    userId: selectedUser.value.id,
    studentId: submitForm.value.studentId.trim(),
    realName: submitForm.value.realName.trim(),
    college: submitForm.value.college.trim()
  })
  submitLoading.value = false

  if (res.code !== 0) {
    ElMessage.error(res.msg || '提交失败')
    return
  }

  selectedUser.value = {
    ...selectedUser.value,
    authStatus: res.data.authStatus,
    authStatusText: getAuthStatusText(res.data.authStatus),
    studentId: submitForm.value.studentId.trim(),
    realName: submitForm.value.realName.trim(),
    college: submitForm.value.college.trim()
  }
  submitResult.value = {
    ...res.data,
    message: res.msg || '申请已提交，审核中',
    submittedAt: new Date().toISOString()
  }
  ElMessage.success(res.msg || '申请已提交，审核中')
}

const goToCampusAuth = () => {
  if (!selectedUser.value.id) {
    router.push({ name: 'campusAuth' })
    return
  }
  router.push({
    name: 'campusAuth',
    query: {
      studentId: submitForm.value.studentId?.trim() || selectedUser.value.studentId || '',
      realName: submitForm.value.realName?.trim() || selectedUser.value.realName || '',
      college: submitForm.value.college?.trim() || selectedUser.value.college || ''
    }
  })
}
</script>

<style scoped>
.agent-review-test-grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: 1.15fr 1fr;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.panel-title {
  font-size: 1.125rem;
  font-weight: 600;
}

.panel-title--compact {
  margin-bottom: 0.25rem;
}

.panel-desc {
  margin-top: 0.25rem;
  color: var(--el-text-color-secondary);
  font-size: 0.875rem;
}

.panel-desc--spaced {
  margin-bottom: 1rem;
}

.panel-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.result-box {
  margin-top: 1rem;
}

@media (max-width: 960px) {
  .agent-review-test-grid {
    grid-template-columns: 1fr;
  }

  .panel-header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
