<template>
  <div class="settings-overlay" v-if="visible">
    <div class="settings-dialog">
      <div class="settings-header">
        <h2 class="settings-title">设置</h2>
        <button class="close-btn" @click="$emit('close')" title="关闭">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>

      <div class="settings-content">
        <!-- 账户信息 -->
        <section class="settings-section">
          <h3 class="section-title">账户信息</h3>
          
          <div class="setting-item account-info">
            <div class="setting-info">
              <div class="setting-label">当前用户</div>
              <div class="setting-description">
                {{ currentUser?.username || '未知用户' }}
              </div>
            </div>
            <div class="setting-control">
              <button
                class="setting-btn danger"
                @click="handleLogout"
                title="登出"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                  <polyline points="16,17 21,12 16,7"/>
                  <line x1="21" y1="12" x2="9" y2="12"/>
                </svg>
                登出
              </button>
            </div>
          </div>
        </section>

        <!-- 安全设置 -->
        <section class="settings-section">
          <h3 class="section-title">安全设置</h3>
          
          <!-- 生物识别设置 -->
          <div class="setting-item">
            <div class="setting-info">
              <div class="setting-label">Windows Hello 生物识别</div>
              <div class="setting-description">
                使用指纹、面部识别或PIN码快速登录
              </div>
              <div v-if="!biometricSupported" class="setting-warning">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                  <line x1="12" y1="9" x2="12" y2="13"/>
                  <line x1="12" y1="17" x2="12.01" y2="17"/>
                </svg>
                {{ biometricMessage }}
              </div>
            </div>
            <div class="setting-control">
              <button
                v-if="biometricSupported && !currentUser?.biometricEnabled"
                class="setting-btn primary"
                @click="showEnableBiometric = true"
                :disabled="loading"
              >
                启用
              </button>
              <button
                v-if="biometricSupported && currentUser?.biometricEnabled"
                class="setting-btn secondary"
                @click="handleDisableBiometric"
                :disabled="loading"
              >
                <div v-if="loading" class="loading-spinner"></div>
                <span v-else>禁用</span>
              </button>
              <div v-if="!biometricSupported" class="setting-status disabled">
                不支持
              </div>
            </div>
          </div>

          <!-- 更改密码 -->
          <div class="setting-item">
            <div class="setting-info">
              <div class="setting-label">更改密码</div>
              <div class="setting-description">
                定期更改密码以确保账户安全
              </div>
            </div>
            <div class="setting-control">
              <button
                class="setting-btn secondary"
                @click="showChangePassword = true"
              >
                更改密码
              </button>
            </div>
          </div>
        </section>

        <!-- 数据管理 -->
        <section class="settings-section">
          <h3 class="section-title">数据管理</h3>
          
          <!-- 数据迁移状态 -->
          <div v-if="migrationStatus?.migrationNeeded" class="setting-item">
            <div class="setting-info">
              <div class="setting-label">数据迁移</div>
              <div class="setting-description">
                发现 {{ migrationStatus.diaryCount }} 篇未加密的日记，建议迁移到加密数据库
              </div>
            </div>
            <div class="setting-control">
              <button
                class="setting-btn warning"
                @click="handleMigration"
                :disabled="migrationLoading"
              >
                <div v-if="migrationLoading" class="loading-spinner"></div>
                <span v-else>立即迁移</span>
              </button>
            </div>
          </div>

          <!-- 导出数据 -->
          <div class="setting-item">
            <div class="setting-info">
              <div class="setting-label">数据导出</div>
              <div class="setting-description">
                导出您的日记数据作为备份
              </div>
            </div>
            <div class="setting-control">
              <button class="setting-btn secondary" disabled>
                导出 (即将推出)
              </button>
            </div>
          </div>
        </section>
      </div>

      <!-- 启用生物识别对话框 -->
      <div v-if="showEnableBiometric" class="modal-overlay">
        <div class="modal-dialog">
          <h3>启用生物识别</h3>
          <p>请输入当前密码以启用Windows Hello生物识别认证</p>
          
          <div class="form-group">
            <input
              type="password"
              v-model="biometricPassword"
              placeholder="输入当前密码"
              :disabled="biometricLoading"
              @keyup.enter="handleEnableBiometric"
            />
          </div>

          <div class="modal-actions">
            <button
              class="modal-btn secondary"
              @click="showEnableBiometric = false"
              :disabled="biometricLoading"
            >
              取消
            </button>
            <button
              class="modal-btn primary"
              @click="handleEnableBiometric"
              :disabled="biometricLoading || !biometricPassword"
            >
              <div v-if="biometricLoading" class="loading-spinner"></div>
              <span v-else>启用</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 更改密码对话框 -->
      <div v-if="showChangePassword" class="modal-overlay">
        <div class="modal-dialog">
          <h3>更改密码</h3>
          
          <div class="form-group">
            <label>当前密码</label>
            <input
              type="password"
              v-model="passwordForm.currentPassword"
              placeholder="输入当前密码"
              :disabled="passwordLoading"
            />
          </div>

          <div class="form-group">
            <label>新密码</label>
            <input
              type="password"
              v-model="passwordForm.newPassword"
              placeholder="输入新密码"
              :disabled="passwordLoading"
            />
          </div>

          <div class="form-group">
            <label>确认新密码</label>
            <input
              type="password"
              v-model="passwordForm.confirmPassword"
              placeholder="再次输入新密码"
              :disabled="passwordLoading"
            />
          </div>

          <div class="modal-actions">
            <button
              class="modal-btn secondary"
              @click="showChangePassword = false"
              :disabled="passwordLoading"
            >
              取消
            </button>
            <button
              class="modal-btn primary"
              @click="handleChangePassword"
              :disabled="passwordLoading || !isPasswordFormValid"
            >
              <div v-if="passwordLoading" class="loading-spinner"></div>
              <span v-else>更改</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 消息提示 -->
      <div v-if="message" :class="['message', messageType]">
        <svg v-if="messageType === 'error'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="15" y1="9" x2="9" y2="15"/>
          <line x1="9" y1="9" x2="15" y2="15"/>
        </svg>
        <svg v-if="messageType === 'success'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
          <polyline points="22,4 12,14.01 9,11.01"/>
        </svg>
        {{ message }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  CheckBiometricSupport,
  EnableBiometric,
  DisableBiometric,
  CheckMigrationStatus,
  MigrateData,
  GetCurrentUser,
  Logout
} from '../../wailsjs/go/main/App'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'userUpdated', 'logout'])

const loading = ref(false)
const biometricLoading = ref(false)
const migrationLoading = ref(false)
const passwordLoading = ref(false)

const showEnableBiometric = ref(false)
const showChangePassword = ref(false)

const biometricSupported = ref(false)
const biometricMessage = ref('')
const currentUser = ref(null)
const migrationStatus = ref(null)

const biometricPassword = ref('')
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const message = ref('')
const messageType = ref('info')

// Computed
const isPasswordFormValid = computed(() => {
  return passwordForm.value.currentPassword &&
         passwordForm.value.newPassword.length >= 6 &&
         passwordForm.value.newPassword === passwordForm.value.confirmPassword
})

// Methods
const showMessage = (text, type = 'info') => {
  message.value = text
  messageType.value = type
  setTimeout(() => {
    message.value = ''
  }, 3000)
}

const handleEnableBiometric = async () => {
  if (!biometricPassword.value) return
  
  biometricLoading.value = true
  
  try {
    await EnableBiometric(biometricPassword.value)
    showMessage('生物识别已启用', 'success')
    showEnableBiometric.value = false
    biometricPassword.value = ''
    
    // 立即更新本地用户状态
    if (currentUser.value) {
      currentUser.value.biometricEnabled = true
    }
    
    emit('userUpdated', currentUser.value)
  } catch (error) {
    showMessage('启用生物识别失败: ' + error.message, 'error')
  } finally {
    biometricLoading.value = false
  }
}

const handleDisableBiometric = async () => {
  loading.value = true
  
  try {
    await DisableBiometric()
    showMessage('生物识别已禁用', 'success')
    
    // 立即更新本地用户状态
    if (currentUser.value) {
      currentUser.value.biometricEnabled = false
    }
    
    emit('userUpdated', currentUser.value)
  } catch (error) {
    showMessage('禁用生物识别失败: ' + error.message, 'error')
  } finally {
    loading.value = false
  }
}

const handleChangePassword = async () => {
  if (!isPasswordFormValid.value) return
  
  passwordLoading.value = true
  
  try {
    // TODO: Implement change password API
    showMessage('密码更改功能即将推出', 'info')
    showChangePassword.value = false
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  } catch (error) {
    showMessage('更改密码失败: ' + error.message, 'error')
  } finally {
    passwordLoading.value = false
  }
}

const handleLogout = () => {
  // 交由父组件处理实际登出逻辑与界面切换
  emit('logout')
  emit('close')
}

const handleMigration = async () => {
  migrationLoading.value = true
  
  try {
    await MigrateData()
    showMessage('数据迁移完成', 'success')
    await loadMigrationStatus()
  } catch (error) {
    showMessage('数据迁移失败: ' + error.message, 'error')
  } finally {
    migrationLoading.value = false
  }
}

const loadBiometricInfo = async () => {
  try {
    const info = await CheckBiometricSupport()
    biometricSupported.value = info.supported
    biometricMessage.value = info.message
  } catch (error) {
    console.error('检查生物识别支持失败:', error)
  }
}

const loadCurrentUser = async () => {
  try {
    currentUser.value = await GetCurrentUser()
  } catch (error) {
    console.error('获取当前用户失败:', error)
  }
}

const loadMigrationStatus = async () => {
  try {
    migrationStatus.value = await CheckMigrationStatus()
  } catch (error) {
    console.error('检查迁移状态失败:', error)
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    loadBiometricInfo(),
    loadCurrentUser(),
    loadMigrationStatus()
  ])
})
</script>

<style scoped>
.settings-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease-out;
}

.settings-dialog {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border-color);
  animation: slideUp 0.4s ease-out;
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 32px;
  border-bottom: 1px solid var(--border-color);
}

.settings-title {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

.settings-content {
  padding: 0 32px 32px;
}

.settings-section {
  margin-top: 32px;
}

.settings-section:first-child {
  margin-top: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color);
}

.setting-item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid var(--bg-tertiary);
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-info {
  flex: 1;
  margin-right: 16px;
}

.setting-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.setting-description {
  font-size: 13px;
  color: var(--text-muted);
  line-height: 1.4;
}

.setting-warning {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--warning-primary);
  margin-top: 6px;
}

.setting-control {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.setting-btn {
  height: 36px;
  padding: 0 16px;
  border: none;
  border-radius: var(--radius-md);
  font-family: var(--font-body);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.setting-btn.primary {
  background: var(--accent-primary);
  color: white;
}

.setting-btn.primary:hover:not(:disabled) {
  background: var(--accent-primary-hover);
}

.setting-btn.secondary {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.setting-btn.secondary:hover:not(:disabled) {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.setting-btn.warning {
  background: var(--warning-primary);
  color: white;
}

.setting-btn.warning:hover:not(:disabled) {
  background: var(--warning-primary-hover);
}

.setting-btn.danger {
  background: var(--accent-error);
  color: white;
}

.setting-btn.danger:hover:not(:disabled) {
  background: var(--accent-error-hover);
}

.setting-btn.danger:focus:not(:disabled),
.setting-btn.danger:focus-visible:not(:disabled) {
  background: var(--accent-error-focus);
  outline: 2px solid var(--accent-error);
  outline-offset: 2px;
  box-shadow: 0 0 0 4px var(--accent-error-alpha);
}

.setting-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.setting-status.disabled {
  font-size: 13px;
  color: var(--text-disabled);
  padding: 8px 0;
}

.account-info .setting-label {
  color: var(--accent-primary);
  font-weight: 600;
}

.account-info .setting-description {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.modal-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-xl);
}

.modal-dialog {
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  padding: 24px;
  width: 90%;
  max-width: 400px;
  border: 1px solid var(--border-color);
}

.modal-dialog h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.modal-dialog p {
  margin: 0 0 20px 0;
  font-size: 14px;
  color: var(--text-muted);
  line-height: 1.4;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.form-group input {
  width: 100%;
  height: 40px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 0 12px;
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--text-primary);
  transition: all 0.2s ease;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 3px var(--accent-primary-alpha);
}

.form-group input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.modal-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}

.modal-btn {
  flex: 1;
  height: 40px;
  border: none;
  border-radius: var(--radius-md);
  font-family: var(--font-body);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.modal-btn.primary {
  background: var(--accent-primary);
  color: white;
}

.modal-btn.primary:hover:not(:disabled) {
  background: var(--accent-primary-hover);
}

.modal-btn.secondary {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.modal-btn.secondary:hover:not(:disabled) {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.modal-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.message {
  position: absolute;
  bottom: 24px;
  left: 32px;
  right: 32px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  animation: slideUp 0.3s ease-out;
}

.message.error {
  background: var(--error-alpha);
  color: var(--error-primary);
  border: 1px solid var(--error-primary);
}

.message.success {
  background: var(--success-alpha);
  color: var(--success-primary);
  border: 1px solid var(--success-primary);
}

.message.info {
  background: var(--accent-primary-alpha);
  color: var(--accent-primary);
  border: 1px solid var(--accent-primary);
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
