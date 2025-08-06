<script setup>
import { ref, onMounted } from 'vue'
import { UploadDiary, GetDiariesList, GetDiaryByID } from '../wailsjs/go/main/App'
import DiaryUpload from './components/DiaryUpload.vue'
import DiaryList from './components/DiaryList.vue'
import DiaryViewer from './components/DiaryViewer.vue'

const currentView = ref('list')
const diaries = ref([])
const selectedDiary = ref(null)
const loading = ref(false)

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
  currentView.value = 'list'
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

const showUploadView = () => {
  currentView.value = 'upload'
}

const showListView = () => {
  currentView.value = 'list'
  selectedDiary.value = null
}
</script>

<template>
  <div id="app">
    <!-- 自定义状态栏 -->
    <div class="custom-titlebar">
      <div class="titlebar-left">
        <div class="app-icon">📔</div>
        <span class="app-name">MoodStack</span>
      </div>
      <div class="titlebar-center">
        <span class="app-subtitle">你的心事，只让你和本地大模型知道</span>
      </div>
      <div class="titlebar-right">
        <div class="window-controls">
          <div class="control-btn minimize"></div>
          <div class="control-btn maximize"></div>
          <div class="control-btn close"></div>
        </div>
      </div>
    </div>

    <!-- 主布局容器 -->
    <div class="main-container">
      <!-- 左侧边栏 -->
      <aside class="sidebar">
        <div class="sidebar-header">
          <h1 class="sidebar-title">📚 日记管理</h1>
        </div>
        
        <nav class="sidebar-nav">
          <button 
            @click="showListView" 
            :class="{ active: currentView === 'list' }"
            class="nav-item"
          >
            <span class="nav-icon">📖</span>
            <span class="nav-text">日记列表</span>
          </button>
          <button 
            @click="showUploadView" 
            :class="{ active: currentView === 'upload' }"
            class="nav-item"
          >
            <span class="nav-icon">📝</span>
            <span class="nav-text">上传日记</span>
          </button>
        </nav>

        <div class="sidebar-footer">
          <div class="diary-count">
            <span class="count-number">{{ diaries.length }}</span>
            <span class="count-label">篇日记</span>
          </div>
        </div>
      </aside>

      <!-- 右侧主内容区 -->
      <main class="content-area">
        <div class="content-container">
          <div v-if="loading" class="loading">
            <div class="loading-spinner"></div>
            <p>加载中...</p>
          </div>
          
          <DiaryUpload 
            v-else-if="currentView === 'upload'" 
            @upload-success="handleUploadSuccess"
          />
          
          <DiaryList 
            v-else-if="currentView === 'list'" 
            :diaries="diaries"
            @view-diary="handleViewDiary"
            @refresh="loadDiaries"
          />
          
          <DiaryViewer 
            v-else-if="currentView === 'viewer' && selectedDiary" 
            :diary="selectedDiary"
            @back="showListView"
          />
        </div>
      </main>
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Nunito', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
  overflow: hidden;
}

#app {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

/* 自定义状态栏 */
.custom-titlebar {
  height: 40px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  -webkit-app-region: drag;
  user-select: none;
}

.titlebar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-icon {
  font-size: 20px;
}

.app-name {
  font-weight: 600;
  color: #4a5568;
  font-size: 16px;
}

.titlebar-center {
  flex: 1;
  text-align: center;
}

.app-subtitle {
  color: #718096;
  font-size: 14px;
  font-weight: 400;
}

.titlebar-right {
  display: flex;
  align-items: center;
}

.window-controls {
  display: flex;
  gap: 8px;
  -webkit-app-region: no-drag;
}

.control-btn {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s ease;
}

.control-btn.minimize {
  background: #ffbd2e;
}

.control-btn.maximize {
  background: #28ca42;
}

.control-btn.close {
  background: #ff5f57;
}

.control-btn:hover {
  transform: scale(1.1);
}

/* 主布局容器 */
.main-container {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* 左侧边栏 */
.sidebar {
  width: 280px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  flex-direction: column;
  padding: 24px 0;
}

.sidebar-header {
  padding: 0 24px 24px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
}

.sidebar-title {
  font-size: 20px;
  font-weight: 600;
  color: #2d3748;
}

.sidebar-nav {
  flex: 1;
  padding: 24px 0;
}

.nav-item {
  width: 100%;
  padding: 16px 24px;
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #4a5568;
  font-size: 16px;
  font-weight: 500;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.5);
  color: #2d3748;
}

.nav-item.active {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border-right: 3px solid #667eea;
}

.nav-icon {
  font-size: 18px;
}

.nav-text {
  font-weight: 500;
}

.sidebar-footer {
  padding: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.3);
}

.diary-count {
  display: flex;
  align-items: center;
  gap: 8px;
}

.count-number {
  font-size: 24px;
  font-weight: 700;
  color: #667eea;
}

.count-label {
  color: #718096;
  font-size: 14px;
}

/* 右侧主内容区 */
.content-area {
  flex: 1;
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  overflow-y: auto;
}

.content-container {
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
}

/* 加载状态 */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 400px;
  color: #4a5568;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(102, 126, 234, 0.2);
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 240px;
  }
  
  .content-container {
    padding: 20px;
  }
  
  .sidebar-title {
    font-size: 18px;
  }
  
  .nav-item {
    padding: 14px 20px;
    font-size: 14px;
  }
}
</style>
