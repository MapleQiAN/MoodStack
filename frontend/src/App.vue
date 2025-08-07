<script setup>
import { ref, onMounted, computed } from 'vue'
import { GetDiariesList, GetDiaryByID } from '../wailsjs/go/main/App'
import DiaryList from './components/DiaryList.vue'
import DiaryViewer from './components/DiaryViewer.vue'
import DiaryEditor from './components/DiaryEditor.vue'

const currentView = ref('list')
const diaries = ref([])
const selectedDiary = ref(null)
const loading = ref(false)
const editingDiary = ref(null)

// 新增状态管理
const sidebarCollapsed = ref(false)
const isDarkMode = ref(false)

// 从本地存储恢复主题设置
onMounted(async () => {
  const savedTheme = localStorage.getItem('moodstack-theme')
  if (savedTheme === 'dark') {
    isDarkMode.value = true
    document.documentElement.setAttribute('data-theme', 'dark')
  }
  
  const savedSidebarState = localStorage.getItem('moodstack-sidebar-collapsed')
  if (savedSidebarState === 'true') {
    sidebarCollapsed.value = true
  }
  
  await loadDiaries()
})

// 切换侧边栏状态
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('moodstack-sidebar-collapsed', sidebarCollapsed.value.toString())
}

// 切换主题模式
const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value
  const theme = isDarkMode.value ? 'dark' : 'light'
  document.documentElement.setAttribute('data-theme', theme)
  localStorage.setItem('moodstack-theme', theme)
}

const loadDiaries = async () => {
  try {
    loading.value = true
    const result = await GetDiariesList()
    diaries.value = result || []
  } catch (error) {
    console.error('加载日记列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleUploadSuccess = async () => {
  await loadDiaries()
}

const handleViewDiary = async (diary) => {
  try {
    const fullDiary = await GetDiaryByID(diary.id)
    selectedDiary.value = fullDiary
    currentView.value = 'viewer'
  } catch (error) {
    console.error('加载日记详情失败:', error)
  }
}



const showListView = () => {
  currentView.value = 'list'
  selectedDiary.value = null
  editingDiary.value = null
}

const handleCreateDiary = () => {
  editingDiary.value = {
    id: '',
    title: '',
    content: '',
    tags: []
  }
  currentView.value = 'editor'
}

// handleEditDiary 函数已移除，因为编辑功能现在集成在 DiaryViewer 中

const handleDiarySaved = async () => {
  await loadDiaries()
  showListView()
}

const handleDiaryUpdated = async (updatedDiary) => {
  // 更新选中的日记数据
  selectedDiary.value = updatedDiary
  // 重新加载日记列表以保持同步
  await loadDiaries()
}

// 窗口控制函数
const minimizeWindow = () => {
  window.runtime.WindowMinimise()
}

const maximizeWindow = () => {
  window.runtime.WindowToggleMaximise()
}

const closeWindow = () => {
  window.runtime.Quit()
}
</script>

<template>
  <div id="app">
    <!-- 标题栏 -->
    <div class="titlebar">
      <div class="titlebar-content">
        <div class="app-info">
          <button class="sidebar-toggle" @click="toggleSidebar" title="切换侧边栏">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="3" y1="6" x2="21" y2="6"/>
              <line x1="3" y1="12" x2="21" y2="12"/>
              <line x1="3" y1="18" x2="21" y2="18"/>
            </svg>
          </button>
          <h1 class="app-title">MoodStack</h1>
          <span class="app-subtitle">心栈</span>
        </div>
        
        <div class="titlebar-actions">
          <button class="theme-toggle" @click="toggleTheme" title="切换主题">
            <svg v-if="!isDarkMode" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5"/>
              <line x1="12" y1="1" x2="12" y2="3"/>
              <line x1="12" y1="21" x2="12" y2="23"/>
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
              <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
              <line x1="1" y1="12" x2="3" y2="12"/>
              <line x1="21" y1="12" x2="23" y2="12"/>
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
              <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
            </svg>
            <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
            </svg>
          </button>
        </div>
        
        <div class="window-controls">
          <button class="control-btn minimize" @click="minimizeWindow">
            <svg width="12" height="1" viewBox="0 0 12 1" fill="currentColor">
              <rect width="12" height="1" />
            </svg>
          </button>
          <button class="control-btn maximize" @click="maximizeWindow">
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1">
              <rect x="1" y="1" width="10" height="10" />
            </svg>
          </button>
          <button class="control-btn close" @click="closeWindow">
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M2 2l8 8M10 2l-8 8" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 主布局 -->
    <div class="main-layout">
      <!-- 左侧边栏 -->
      <aside :class="['sidebar', { collapsed: sidebarCollapsed }]">
        <nav class="navigation">
          <button 
            @click="showListView" 
            :class="['nav-item', { active: currentView === 'list' || currentView === 'viewer' || currentView === 'editor' }]"
            :title="sidebarCollapsed ? '我的日记' : ''"
          >
            <span class="nav-icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
              </svg>
            </span>
            <span class="nav-text" v-show="!sidebarCollapsed">我的日记</span>
          </button>
        </nav>
        
        <div class="sidebar-footer" v-show="!sidebarCollapsed">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-label">记录总数</span>
              <span class="stat-value">{{ diaries.length }}</span>
            </div>
          </div>
        </div>
      </aside>

      <!-- 右侧主内容区域 -->
      <main class="content-area">
        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <p class="loading-text">正在加载...</p>
        </div>
        
        <Transition name="page" mode="out-in" appear>
          <DiaryList 
            v-if="currentView === 'list'" 
            key="list"
            :diaries="diaries"
            @view-diary="handleViewDiary"
            @upload-success="handleUploadSuccess"
            @create-diary="handleCreateDiary"
          />
          
          <DiaryViewer 
            v-else-if="currentView === 'viewer' && selectedDiary" 
            key="viewer"
            :diary="selectedDiary"
            @back="showListView"
            @diary-updated="handleDiaryUpdated"
          />
          
          <DiaryEditor
            v-else-if="currentView === 'editor' && editingDiary"
            key="editor"
            :diary="editingDiary"
            @back="showListView"
            @diary-saved="handleDiarySaved"
          />
        </Transition>
      </main>
    </div>
  </div>
</template>

<style scoped>
#app {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-gradient);
}

/* 标题栏 */
.titlebar {
  height: 48px;
  background: var(--bg-header);
  border-bottom: 1px solid var(--bg-tertiary);
  display: flex;
  align-items: center;
  position: relative;
  z-index: 100;
}

.titlebar-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
}

.app-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.sidebar-toggle, .theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s ease;
}

.sidebar-toggle:hover, .theme-toggle:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.titlebar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-title {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 16px;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.025em;
}

.app-subtitle {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
}

.window-controls {
  display: flex;
  gap: 8px;
}

.control-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: var(--radius-sm);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  transition: all 0.2s ease;
}

.control-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  transform: scale(1.1);
}

.control-btn.close:hover {
  background: var(--accent-error);
  color: white;
  transform: scale(1.1);
}

.control-btn:active {
  transform: scale(0.95);
}

/* 主布局 */
.main-layout {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* 左侧边栏 */
.sidebar {
  width: 280px;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-right: 1px solid var(--bg-tertiary);
  display: flex;
  flex-direction: column;
  padding: 24px 16px;
  transition: width 0.3s var(--ease-spring), padding 0.3s var(--ease-spring);
  overflow: hidden;
}

.sidebar.collapsed {
  width: 72px;
  padding: 24px 12px;
}

.navigation {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 32px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: var(--font-body);
  font-weight: 500;
  font-size: 14px;
  color: var(--text-secondary);
  text-align: left;
  width: 100%;
  justify-content: flex-start;
  min-height: 44px;
}

.sidebar.collapsed .nav-item {
  justify-content: center;
  padding: 12px;
}

.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nav-item:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  transform: translateX(4px);
}

.nav-item.active {
  background: var(--accent-primary);
  color: var(--bg-primary);
  font-weight: 600;
  transform: translateX(6px);
  box-shadow: var(--shadow-md);
}

.nav-text {
  flex: 1;
}

.sidebar-footer {
  margin-top: auto;
}

.stats {
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  padding: 16px;
  border: 1px solid var(--bg-tertiary);
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-label {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 500;
}

.stat-value {
  font-size: 15px;
  color: var(--text-primary);
  font-weight: 600;
}

/* 主内容区域 */
.content-area {
  flex: 1;
  overflow: auto;
  background: var(--bg-primary);
  position: relative;
  perspective: 1000px;
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 20px;
  animation: fadeInScale 0.4s ease-out;
}

.loading-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid var(--bg-tertiary);
  border-top: 3px solid var(--accent-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  filter: drop-shadow(0 2px 4px rgba(136, 192, 208, 0.2));
}

.loading-text {
  color: var(--text-muted);
  font-weight: 500;
  font-size: 14px;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 页面切换动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.45s var(--ease-spring);
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  will-change: transform, opacity;
  backface-visibility: hidden;
}

.page-enter-active {
  transition-delay: 0.05s;
  z-index: 2;
}

.page-leave-active {
  z-index: 1;
}

.page-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

.page-enter-to,
.page-leave-from {
  opacity: 1;
  transform: translateX(0);
}

/* 移动端动画优化 */
@media (max-width: 768px) {
  .page-enter-active,
  .page-leave-active {
    transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  }
  
  .page-enter-from {
    transform: translateX(20px);
  }
  
  .page-leave-to {
    transform: translateX(-20px);
  }
}



/* 响应式设计 */
@media (max-width: 1024px) {
  .sidebar {
    width: 240px;
    padding: 20px 12px;
  }
}

@media (max-width: 768px) {
  .main-layout {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
    padding: 16px;
    border-right: none;
    border-bottom: 1px solid var(--bg-tertiary);
  }
  
  .navigation {
    flex-direction: row;
    gap: 4px;
    margin-bottom: 16px;
  }
  
  .nav-item {
    flex: 1;
    justify-content: center;
    padding: 10px 12px;
  }
  
  .sidebar-footer {
    margin-top: 0;
  }
  
  .app-title {
    font-size: 15px;
  }
  
  .app-subtitle {
    font-size: 11px;
  }
}

@media (max-width: 480px) {
  .titlebar-content {
    padding: 0 12px;
  }
  
  .sidebar {
    padding: 12px;
  }
  
  .navigation {
    gap: 2px;
  }
  
  .nav-item {
    padding: 8px 10px;
    font-size: 13px;
  }
}
</style>

