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
        <!-- 移除窗口控制按钮，保持空间平衡 -->
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
            <p class="loading-text">正在加载你的日记...</p>
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
  font-family: 'Inter', 'Nunito', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  height: 60px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(30px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  -webkit-app-region: drag;
  user-select: none;
  box-shadow: 0 1px 20px rgba(0, 0, 0, 0.1);
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
  font-weight: 700;
  color: rgba(255, 255, 255, 0.95);
  font-size: 20px;
  letter-spacing: -0.5px;
}

.titlebar-center {
  flex: 1;
  text-align: center;
}

.app-subtitle {
  color: rgba(255, 255, 255, 0.7);
  font-size: 15px;
  font-weight: 500;
  letter-spacing: 0.3px;
}

.titlebar-right {
  display: flex;
  align-items: center;
}

/* 窗口控制按钮样式已移除 */

/* 主布局容器 */
.main-container {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* 左侧边栏 */
.sidebar {
  width: 300px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(30px);
  border-right: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  flex-direction: column;
  padding: 32px 0;
  box-shadow: 2px 0 20px rgba(0, 0, 0, 0.1);
}

.sidebar-header {
  padding: 0 24px 24px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
}

.sidebar-title {
  font-size: 24px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.95);
  letter-spacing: -0.5px;
}

.sidebar-nav {
  flex: 1;
  padding: 24px 0;
}

.nav-item {
  width: 100%;
  padding: 18px 32px;
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: rgba(255, 255, 255, 0.8);
  font-size: 16px;
  font-weight: 600;
  border-radius: 0 25px 25px 0;
  margin: 4px 0;
  position: relative;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.15);
  color: rgba(255, 255, 255, 0.95);
  transform: translateX(8px);
}

.nav-item.active {
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
  box-shadow: 0 4px 15px rgba(255, 255, 255, 0.1);
  transform: translateX(8px);
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
  font-size: 28px;
  font-weight: 800;
  color: rgba(255, 255, 255, 0.95);
  background: linear-gradient(135deg, #ffffff 0%, #f0f0f0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.count-label {
  color: rgba(255, 255, 255, 0.7);
  font-size: 15px;
  font-weight: 500;
}

/* 右侧主内容区 */
.content-area {
  flex: 1;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  overflow-y: auto;
  border-radius: 20px 0 0 0;
}

.content-container {
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
}

/* 现代化加载状态 */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 500px;
  color: rgba(255, 255, 255, 0.9);
  animation: fadeInUp 0.6s ease-out;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top: 4px solid rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  animation: spin 1.2s ease-in-out infinite;
  margin-bottom: 24px;
  box-shadow: 0 0 20px rgba(255, 255, 255, 0.2);
}

.loading-text {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: 0.5px;
  animation: pulse 2s ease-in-out infinite;
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
