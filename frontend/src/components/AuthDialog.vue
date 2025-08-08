<template>
  <div class="auth-overlay" v-if="visible">
    <div class="auth-dialog" :class="{ 'setup-mode': isSetupMode }">

      <div class="auth-header">
        <h2 class="auth-title">
          {{ isSetupMode ? '初始化 MoodStack' : 'MoodStack' }}
        </h2>
        <p class="auth-subtitle" v-if="isSetupMode">
          创建您的第一个账户来保护您的日记，账户仅保存于本地
        </p>

      </div>

      <div class="auth-dialog-content">
        <div class="auth-content">
          <!-- 设置模式 -->
        <div v-if="isSetupMode" class="setup-form">
          <div class="form-group">
            <label for="username">用户名</label>
            <input
              id="username"
              type="text"
              v-model="setupForm.username"
              placeholder="输入用户名"
              :disabled="loading"
              @keyup.enter="handleSetup"
            />
          </div>
          
          <div class="form-group">
            <label for="password">密码</label>
            <div class="password-input">
              <input
                id="password"
                :type="showPassword ? 'text' : 'password'"
                v-model="setupForm.password"
                placeholder="输入密码"
                :disabled="loading"
                @keyup.enter="handleSetup"
              />
              <button
                v-if="setupForm.password"
                type="button"
                class="password-toggle"
                @click="showPassword = !showPassword"
                :disabled="loading"
              >
                <svg v-if="showPassword" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label for="confirmPassword">确认密码</label>
            <div class="password-input">
              <input
                id="confirmPassword"
                :type="showPassword ? 'text' : 'password'"
                v-model="setupForm.confirmPassword"
                placeholder="再次输入密码"
                :disabled="loading"
                @keyup.enter="handleSetup"
              />
              <button
                v-if="setupForm.confirmPassword"
                type="button"
                class="password-toggle"
                @click="showPassword = !showPassword"
                :disabled="loading"
              >
                <svg v-if="showPassword" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
              </button>
            </div>
          </div>

          <div class="form-actions">
            <button
              class="auth-button primary"
              @click="handleSetup"
              :disabled="loading || !isSetupFormValid"
            >
              <div v-if="loading" class="loading-spinner"></div>
              <span v-else>创建账户</span>
            </button>
          </div>
        </div>

        <!-- 登录模式 -->
        <div v-else class="login-form">
          <div class="user-welcome" v-if="currentUser">
            <div class="user-avatar">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
            </div>
            <h3 class="user-name">欢迎回来，{{ currentUser.username }}！</h3>
          </div>
          
          <div class="form-group">
            <label for="loginPassword">密码</label>
            <div class="password-input">
              <input
                id="loginPassword"
                :type="showPassword ? 'text' : 'password'"
                v-model="loginForm.password"
                placeholder="输入密码"
                :disabled="loading"
                @keyup.enter="handlePasswordLogin"
              />
              <button
                v-if="loginForm.password"
                type="button"
                class="password-toggle"
                @click="showPassword = !showPassword"
                :disabled="loading"
              >
                <svg v-if="showPassword" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
              </button>
            </div>
          </div>

          <div class="form-actions">
            <button
              class="auth-button primary"
              @click="handlePasswordLogin"
              :disabled="loading || !isLoginFormValid"
            >
              <div v-if="loading" class="loading-spinner"></div>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="10" rx="2" ry="2"/>
                <circle cx="12" cy="16" r="1"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <span>密码登录</span>
            </button>
            
              <button
              v-if="biometricSupported && currentUser && currentUser.biometricEnabled"
              class="auth-button secondary biometric"
              @click="handleBiometricLogin"
              :disabled="loading"
            >
              <div v-if="biometricLoading" class="loading-spinner"></div>
              <svg v-else width="18" height="18" viewBox="0 0 24 24" aria-hidden="true" fill="none">
                <circle cx="8.5" cy="9" r="1.5" fill="currentColor"/>
                <circle cx="15.5" cy="9" r="1.5" fill="currentColor"/>
                <path d="M7 14c1.9 2.2 4 3.5 5 3.5s3.1-1.3 5-3.5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <span>{{ biometricLoading ? '请注意弹出窗口' : 'Windows Hello' }}</span>
            </button>
          </div>
        </div>

        <!-- 错误消息 -->
        <div v-if="errorMessage" class="error-message">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
          {{ errorMessage }}
        </div>

        <!-- 成功消息 -->
        <div v-if="successMessage" class="success-message">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22,4 12,14.01 9,11.01"/>
          </svg>
          {{ successMessage }}
        </div>

        <!-- 迁移提示 -->
        <div v-if="showMigrationNotice" class="migration-notice">
          <div class="migration-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7,10 12,15 17,10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
          </div>
          <div class="migration-content">
            <h4>检测到现有日记数据</h4>
            <p>发现 {{ migrationInfo.diaryCount }} 篇未加密的日记，建议迁移到加密数据库以提高安全性。</p>
            <button
              class="migration-button"
              @click="handleMigration"
              :disabled="migrationLoading"
            >
              <div v-if="migrationLoading" class="loading-spinner"></div>
              <span v-else>立即迁移</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 生物识别设置 -->
      <div v-if="showBiometricSetup" class="biometric-setup">
        <div class="biometric-header">
          <h3>启用生物识别</h3>
          <p>为了更便捷的登录体验，您可以启用Windows Hello生物识别认证</p>
        </div>
        
        <div class="form-group">
          <label for="biometricPassword">请输入当前密码以启用生物识别</label>
          <div class="password-input">
            <input
              id="biometricPassword"
              type="password"
              v-model="biometricPassword"
              placeholder="输入密码"
              :disabled="biometricLoading"
              @keyup.enter="handleEnableBiometric"
            />
          </div>
        </div>

        <div class="form-actions">
          <button
            class="auth-button secondary"
            @click="showBiometricSetup = false"
            :disabled="biometricLoading"
          >
            跳过
          </button>
          <button
            class="auth-button primary"
            @click="handleEnableBiometric"
            :disabled="biometricLoading || !biometricPassword"
          >
            <div v-if="biometricLoading" class="loading-spinner"></div>
            <span>{{ biometricLoading ? '请注意弹出窗口' : '启用' }}</span>
          </button>
        </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  CheckAuthStatus, 
  CreateFirstUser, 
  AuthenticateUser, 
  AuthenticateWithBiometric,
  CheckBiometricSupport,
  EnableBiometric,
  CheckMigrationStatus,
  MigrateData,
  GetFirstUser
} from '../../wailsjs/go/main/App'

const emit = defineEmits(['authenticated', 'close'])

const visible = ref(true)
const loading = ref(false)
const biometricLoading = ref(false)
const migrationLoading = ref(false)
const showPassword = ref(false)
const isSetupMode = ref(false)
const forceSetupMode = ref(false) // 强制设置模式，当没有用户时
const biometricSupported = ref(false)
const showBiometricSetup = ref(false)
const showMigrationNotice = ref(false)

const errorMessage = ref('')
const successMessage = ref('')
const biometricPassword = ref('')

const setupForm = ref({
  username: '',
  password: '',
  confirmPassword: ''
})

const loginForm = ref({
  password: ''
})

const currentUser = ref(null)

const migrationInfo = ref({
  diaryCount: 0
})

// Computed properties
const isSetupFormValid = computed(() => {
  return setupForm.value.username.trim() && 
         setupForm.value.password.length >= 6 && 
         setupForm.value.password === setupForm.value.confirmPassword
})

const isLoginFormValid = computed(() => {
  return loginForm.value.password.trim()
})

// Methods
const clearMessages = () => {
  errorMessage.value = ''
  successMessage.value = ''
}

const showError = (message) => {
  clearMessages()
  errorMessage.value = message
}

const showSuccess = (message) => {
  clearMessages()
  successMessage.value = message
}

const handleSetup = async () => {
  if (!isSetupFormValid.value) return
  
  loading.value = true
  clearMessages()
  
  try {
    const result = await CreateFirstUser(setupForm.value.username, setupForm.value.password)
    
    if (result.success) {
      showSuccess('账户创建成功！')
      
      // Check migration after successful setup
      await checkMigration()
      
      // Check biometric support
      if (biometricSupported.value) {
        showBiometricSetup.value = true
      } else {
        emit('authenticated', result.user)
      }
    } else {
      showError(result.message)
    }
  } catch (error) {
    showError('创建账户失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const handlePasswordLogin = async () => {
  if (!isLoginFormValid.value || !currentUser.value) return
  
  loading.value = true
  clearMessages()
  
  try {
    const result = await AuthenticateUser(currentUser.value.username, loginForm.value.password)
    
    if (result.success) {
      showSuccess('登录成功！')
      
      // Check migration after successful login
      await checkMigration()
      
      if (!showMigrationNotice.value) {
        emit('authenticated', result.user)
      }
    } else {
      showError(result.message)
    }
  } catch (error) {
    showError('登录失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const handleBiometricLogin = async () => {
  if (!currentUser.value) {
    showError('无法获取用户信息')
    return
  }
  
  biometricLoading.value = true
  clearMessages()
  
  try {
    const result = await AuthenticateWithBiometric(currentUser.value.username)
    
    if (result.success) {
      showSuccess('生物识别认证成功！')
      
      // Check migration after successful login
      await checkMigration()
      
      if (!showMigrationNotice.value) {
        emit('authenticated', result.user)
      }
    } else {
      showError(result.message)
    }
  } catch (error) {
    showError('生物识别认证失败: ' + error.message)
  } finally {
    biometricLoading.value = false
  }
}

const handleEnableBiometric = async () => {
  if (!biometricPassword.value) return
  
  biometricLoading.value = true
  clearMessages()
  
  try {
    await EnableBiometric(biometricPassword.value)
    showSuccess('生物识别已启用！')
    
    setTimeout(() => {
      showBiometricSetup.value = false
      if (!showMigrationNotice.value) {
        emit('authenticated')
      }
    }, 1500)
  } catch (error) {
    showError('启用生物识别失败: ' + error.message)
  } finally {
    biometricLoading.value = false
  }
}

const handleMigration = async () => {
  migrationLoading.value = true
  clearMessages()
  
  try {
    await MigrateData()
    showSuccess('数据迁移完成！')
    
    setTimeout(() => {
      showMigrationNotice.value = false
      emit('authenticated')
    }, 1500)
  } catch (error) {
    showError('数据迁移失败: ' + error.message)
  } finally {
    migrationLoading.value = false
  }
}

const checkMigration = async () => {
  try {
    const status = await CheckMigrationStatus()
    if (status.migrationNeeded && status.diaryCount > 0) {
      migrationInfo.value = status
      showMigrationNotice.value = true
    }
  } catch (error) {
    console.error('检查迁移状态失败:', error)
  }
}

const checkAuthStatus = async () => {
  try {
    const [authStatus, biometricInfo] = await Promise.all([
      CheckAuthStatus(),
      CheckBiometricSupport()
    ])
    
    // 检查返回值类型以兼容新旧版本
    let requireSetup = false
    let hasUsers = true
    
    if (typeof authStatus === 'boolean') {
      // 旧版本兼容 - 如果返回false，可能是需要初始化或需要登录
      requireSetup = !authStatus
      hasUsers = true // 假设有用户，让用户尝试登录
    } else if (authStatus && typeof authStatus === 'object') {
      // 新版本
      requireSetup = authStatus.requireSetup
      hasUsers = authStatus.hasUsers
    }
    
    isSetupMode.value = requireSetup
    forceSetupMode.value = requireSetup // 只有真正需要初始化时才强制设置模式
    biometricSupported.value = biometricInfo.supported
    
    if (requireSetup) {
      console.log('No users found, showing setup mode')
      currentUser.value = null
    } else {
      console.log('Users exist, showing login mode')
      // 获取第一个用户信息用于单用户模式
      try {
        const user = await GetFirstUser()
        currentUser.value = user
        console.log('Loaded user:', user.username)
      } catch (error) {
        console.error('获取用户信息失败:', error)
        currentUser.value = null
      }
    }
  } catch (error) {
    console.error('检查认证状态失败:', error)
    showError('初始化失败: ' + error.message)
  }
}

// 窗口控制方法已移除 - 登录界面不再需要窗口控制

// Lifecycle
onMounted(() => {
  checkAuthStatus()
})
</script>

<style scoped>
.auth-overlay {
  position: fixed;
  top: 48px; /* 顶栏高度 */
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--bg-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease-out;
}

.auth-dialog {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: 40px;
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 
    0 25px 80px rgba(0, 0, 0, 0.12), 
    0 10px 25px rgba(0, 0, 0, 0.08),
    0 0 0 1px rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-color);
  animation: slideUp 0.4s ease-out;
  display: flex;
  flex-direction: column;
  backdrop-filter: blur(20px);
  position: relative;
}

.auth-dialog-content {
  overflow-y: auto;
  flex: 1;
}

.auth-dialog.setup-mode {
  max-width: 520px;
}

/* 标题栏样式已移除 - 简化登录界面 */

.auth-header {
  text-align: center;
  margin-bottom: 32px;
  padding-top: 8px;
}

.mode-switch {
  margin-top: 16px;
}

.mode-switch-btn {
  background: none;
  border: none;
  color: var(--accent-primary);
  font-size: 13px;
  cursor: pointer;
  padding: 4px 0;
  transition: all 0.2s ease;
  text-decoration: underline;
  text-underline-offset: 2px;
}

.mode-switch-btn:hover:not(:disabled) {
  color: var(--accent-primary-hover);
  text-decoration-color: var(--accent-primary-hover);
}

.mode-switch-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.auth-title {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 12px 0;
  letter-spacing: -0.02em;
}

.auth-subtitle {
  font-size: 14px;
  color: var(--text-muted);
  margin: 0;
  line-height: 1.5;
}

.auth-content {
  margin-bottom: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.form-group input {
  width: 100%;
  height: 44px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 0 16px;
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--text-primary);
  transition: all 0.2s ease;
}

.password-input input {
  padding-right: 48px; /* 为眼睛图标留出空间 */
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 3px var(--accent-primary-alpha);
  z-index: 1;
  position: relative;
}

.form-group input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.password-input {
  position: relative;
  max-width: 380px; /* 使密码框更窄 */
  margin: 0 auto;
}

.input-wrapper {
  position: relative;
}

/* 移除浏览器原生的密码显示按钮 */
input[type="password"]::-ms-reveal {
  display: none;
}

input[type="password"]::-webkit-credentials-auto-fill-button {
  display: none;
}

.password-toggle {
  position: absolute;
  top: 50%;
  right: 12px;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 4px;
  border-radius: var(--radius-sm);
  transition: all 0.2s ease;
  z-index: 2;
}

.password-toggle:hover {
  color: var(--text-secondary);
  background: var(--bg-tertiary);
}

.password-toggle:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.auth-button {
  flex: 1;
  height: 44px;
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
  gap: 8px;
}

.auth-button.primary {
  background: var(--accent-primary);
  color: white;
}

.auth-button.primary:hover:not(:disabled),
.auth-button.primary:focus-visible:not(:disabled) {
  background: var(--accent-primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 16px var(--accent-primary-alpha);
}

.auth-button.secondary {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.auth-button.secondary:hover:not(:disabled) {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.auth-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.error-message,
.success-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  margin-top: 16px;
}

.error-message {
  background: var(--error-alpha);
  color: var(--error-primary);
  border: 1px solid var(--error-primary);
}

.success-message {
  background: var(--success-alpha);
  color: var(--success-primary);
  border: 1px solid var(--success-primary);
}

.migration-notice {
  background: var(--warning-alpha);
  border: 1px solid var(--warning-primary);
  border-radius: var(--radius-md);
  padding: 16px;
  margin-top: 20px;
  display: flex;
  gap: 12px;
}

.migration-icon {
  color: var(--warning-primary);
  flex-shrink: 0;
}

.migration-content h4 {
  margin: 0 0 4px 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--warning-primary);
}

.migration-content p {
  margin: 0 0 12px 0;
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.4;
}

.migration-button {
  background: var(--warning-primary);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  padding: 6px 12px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s ease;
}

.migration-button:hover:not(:disabled) {
  background: var(--warning-primary-hover);
}

.migration-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.biometric-setup {
  border-top: 1px solid var(--border-color);
  padding-top: 24px;
  margin-top: 24px;
}

.biometric-header {
  text-align: center;
  margin-bottom: 20px;
}

.biometric-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.biometric-header p {
  font-size: 13px;
  color: var(--text-muted);
  margin: 0;
  line-height: 1.4;
}

.user-welcome {
  text-align: center;
  margin-bottom: 24px;
  padding: 20px;
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
}

.user-avatar {
  margin: 0 auto 16px auto;
  width: 48px;
  height: 48px;
  background: var(--accent-primary-alpha);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-primary);
}

.user-name {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  letter-spacing: -0.01em;
}

.user-greeting {
  font-size: 14px;
  color: var(--text-muted);
  margin: 0;
  line-height: 1.4;
}

.auth-button.biometric {
  background: #eaaa29;
  color: var(--bg-primary);
  border: none;
  box-shadow: 0 2px 8px rgba(235, 203, 139, 0.2);
}

.auth-button.biometric:hover:not(:disabled),
.auth-button.biometric:focus-visible:not(:disabled) {
  background: #cb911b;
  color: var(--bg-primary);
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(235, 203, 139, 0.3);
}

.auth-button.biometric svg {
  color: currentColor;
  stroke: currentColor;
}

.auth-button.biometric:hover:not(:disabled) svg,
.auth-button.biometric:focus-visible:not(:disabled) svg {
  color: currentColor;
  stroke: currentColor;
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
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
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
