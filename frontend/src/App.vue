<script setup>
import { ref, onMounted } from 'vue'
import { GetDiariesList, GetDiaryByID } from '../wailsjs/go/main/App'
import DiaryList from './components/DiaryList.vue'
import DiaryViewer from './components/DiaryViewer.vue'
import DiaryEditor from './components/DiaryEditor.vue'

const currentView = ref('list')
const diaries = ref([])
const selectedDiary = ref(null)
const loading = ref(false)
const editingDiary = ref(null)

onMounted(async () => {
  await loadDiaries()
})

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
          <h1 class="app-title">MoodStack</h1>
          <span class="app-subtitle">心栈</span>
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
      <aside class="sidebar">
        <nav class="navigation">
          <button 
            @click="showListView" 
            :class="['nav-item', { active: currentView === 'list' || currentView === 'viewer' || currentView === 'editor' }]"
          >
            <span class="nav-text">我的日记</span>
          </button>
        </nav>
        
        <div class="sidebar-footer">
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
        
        <DiaryList 
          v-if="currentView === 'list'" 
          :diaries="diaries"
          @view-diary="handleViewDiary"
          @upload-success="handleUploadSuccess"
          @create-diary="handleCreateDiary"
        />
        
        <DiaryViewer 
          v-else-if="currentView === 'viewer' && selectedDiary" 
          :diary="selectedDiary"
          @back="showListView"
          @diary-updated="handleDiaryUpdated"
        />
        
        <DiaryEditor
          v-else-if="currentView === 'editor' && editingDiary"
          :diary="editingDiary"
          @back="showListView"
          @diary-saved="handleDiarySaved"
        />
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
  -webkit-app-region: drag;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--nord4);
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
  -webkit-app-region: drag;
}

.app-info {
  display: flex;
  align-items: center;
  gap: 12px;
  -webkit-app-region: drag;
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
  -webkit-app-region: no-drag;
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
  background: var(--nord4);
  color: var(--text-secondary);
}

.control-btn.close:hover {
  background: var(--accent-error);
  color: white;
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
  border-right: 1px solid var(--nord4);
  display: flex;
  flex-direction: column;
  padding: 24px 16px;
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
}

.nav-item:hover {
  background: var(--nord5);
  color: var(--text-primary);
}

.nav-item.active {
  background: var(--accent-primary);
  color: var(--nord0);
  font-weight: 600;
}

.nav-text {
  flex: 1;
}

.sidebar-footer {
  margin-top: auto;
}

.stats {
  background: var(--nord5);
  border-radius: var(--radius-md);
  padding: 16px;
  border: 1px solid var(--nord4);
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
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 16px;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--nord4);
  border-top: 3px solid var(--accent-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-text {
  color: var(--text-muted);
  font-weight: 500;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
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
    border-bottom: 1px solid var(--nord4);
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