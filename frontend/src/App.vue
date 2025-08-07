<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue'
import { GetDiariesList, GetDiaryByID, SearchDiariesWithContext } from '../wailsjs/go/main/App'
import DiaryList from './components/DiaryList.vue'
import DiaryViewer from './components/DiaryViewer.vue'
import DiaryEditor from './components/DiaryEditor.vue'

const currentView = ref('list')
const diaries = ref([])
const selectedDiary = ref(null)
const loading = ref(false)
const editingDiary = ref(null)
const searchQuery = ref('')
const searchMode = ref(false)
const searchResults = ref([])
const searchLoading = ref(false)
const showSearchResults = ref(false)
let searchTimeout = null

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

const focusSearch = () => {
  searchMode.value = true
  // 使用 nextTick 确保 input 在 DOM 中渲染完毕
  nextTick(() => {
    document.querySelector('.search-input')?.focus()
  })
}

const handleKeydown = (event) => {
  if (event.ctrlKey && event.key === 'k') {
    event.preventDefault()
    focusSearch()
  }
  if (event.key === 'Escape' && searchMode.value) {
    searchMode.value = false
    searchQuery.value = ''
    showSearchResults.value = false
  }
}

// 搜索相关方法
const onSearchFocus = () => {
  searchMode.value = true
  if (searchQuery.value.trim()) {
    showSearchResults.value = true
  }
}

const onSearchBlur = () => {
  // 延迟关闭搜索结果，以便处理点击事件
  setTimeout(() => {
    searchMode.value = false
    showSearchResults.value = false
  }, 200)
}

const onSearchInput = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }

  if (!searchQuery.value.trim()) {
    searchResults.value = []
    showSearchResults.value = false
    return
  }

  showSearchResults.value = true
  searchLoading.value = true

  // 防抖搜索
  searchTimeout = setTimeout(async () => {
    try {
      const results = await SearchDiariesWithContext(searchQuery.value.trim())
      searchResults.value = results
    } catch (error) {
      console.error('搜索失败:', error)
      searchResults.value = []
    } finally {
      searchLoading.value = false
    }
  }, 300)
}

const handleSearchResultClick = async (result) => {
  try {
    // 加载完整的日记
    const fullDiary = await GetDiaryByID(result.diary.id)
    selectedDiary.value = fullDiary
    currentView.value = 'viewer'
    
    // 关闭搜索
    searchMode.value = false
    showSearchResults.value = false
    
    // 清空搜索查询
    searchQuery.value = ''
    
    // TODO: 如果有匹配的内容片段，可以在这里实现滚动到相应位置的逻辑
  } catch (error) {
    console.error('打开日记失败:', error)
  }
}

const highlightSearchText = (text, query) => {
  if (!query.trim()) return text
  
  const regex = new RegExp(`(${query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
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
        
        <!-- 全局搜索 -->
        <div class="global-search">
          <div class="search-container" :class="{ 'search-mode': searchMode }">
            <div class="search-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
            </div>
            <input 
              type="text" 
              class="search-input"
              v-model="searchQuery"
              placeholder="全局搜索... (Ctrl+K)"
              @focus="onSearchFocus"
              @blur="onSearchBlur"
              @input="onSearchInput"
            />
            <div class="search-shortcut" v-show="!searchMode && !searchQuery">
              <kbd>Ctrl</kbd>
              <kbd>K</kbd>
            </div>
            <!-- 搜索加载状态 -->
            <div v-if="searchLoading" class="search-loading">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="search-spinner">
                <path d="M21 12a9 9 0 11-6.219-8.56"/>
              </svg>
            </div>
            
            <!-- 搜索结果下拉框 -->
            <div v-if="showSearchResults" class="search-dropdown" @mousedown.prevent>
              <div v-if="searchLoading" class="search-loading-state">
                <div class="loading-spinner"></div>
                <span>搜索中...</span>
              </div>
              
              <div v-else-if="searchResults.length === 0 && searchQuery.trim()" class="search-empty-state">
                <div class="empty-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"/>
                    <line x1="21" y1="21" x2="16.65" y2="16.65"/>
                    <line x1="8" y1="11" x2="14" y2="11"/>
                  </svg>
                </div>
                <span>未找到相关日记</span>
              </div>
              
              <div v-else-if="searchResults.length > 0" class="search-results">
                <div 
                  v-for="result in searchResults" 
                  :key="result.diary.id"
                  class="search-result-item"
                  @click="handleSearchResultClick(result)"
                >
                  <div class="result-header">
                    <h3 class="result-title">{{ result.diary.title }}</h3>
                    <span class="result-date">{{ formatDate(result.diary.createdAt) }}</span>
                  </div>
                  
                  <div v-if="result.matchType === 'title'" class="result-match-type">
                    <span class="match-tag title-match">标题匹配</span>
                  </div>
                  
                  <div v-if="result.matchedSnippets.length > 0" class="result-snippets">
                    <div 
                      v-for="(snippet, index) in result.matchedSnippets.slice(0, 2)" 
                      :key="index"
                      class="snippet"
                      v-html="highlightSearchText(snippet, searchQuery)"
                    ></div>
                  </div>
                  
                  <div class="result-footer">
                    <span class="result-match-type">
                      <span v-if="result.matchType === 'content'" class="match-tag content-match">内容匹配</span>
                      <span v-if="result.matchType === 'both'" class="match-tag both-match">标题+内容匹配</span>
                    </span>
                  </div>
                </div>
              </div>
              
              <div v-if="searchQuery.trim() && !searchLoading" class="search-footer">
                <div class="search-tip">
                  <kbd>Enter</kbd> 选择 • <kbd>Esc</kbd> 关闭
                </div>
              </div>
            </div>
          </div>
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
  -webkit-app-region: drag;
  background: var(--bg-header);
  border-bottom: 1px solid var(--bg-tertiary);
  display: flex;
  align-items: center;
  position: relative;
  z-index: 100;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.titlebar-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  gap: 16px;
  -webkit-app-region: drag;
}

.global-search {
  flex-grow: 1;
  display: flex;
  justify-content: center;
  -webkit-app-region: drag;
}

.search-container {
  position: relative;
  width: 100%;
  max-width: 400px;
  transition: all 0.3s var(--ease-spring);
  -webkit-app-region: no-drag;
}

.search-container.search-mode {
  max-width: 600px;
}

.search-icon {
  position: absolute;
  top: 50%;
  left: 14px;
  transform: translateY(-50%);
  color: var(--text-muted);
  pointer-events: none;
  transition: color 0.2s ease;
}

.search-input {
  width: 100%;
  height: 36px;
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  padding: 0 40px 0 40px;
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--text-primary);
  transition: all 0.2s ease;
  outline: none;
  -webkit-app-region: no-drag;
}

.search-input:focus {
  border-color: var(--accent-primary);
  background: var(--bg-primary);
  box-shadow: 0 0 0 3px var(--accent-primary-alpha);
}

.search-input:focus + .search-icon {
  color: var(--accent-primary);
}

.search-shortcut {
  position: absolute;
  top: 50%;
  right: 12px;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  gap: 4px;
  pointer-events: none;
}

kbd {
  background-color: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  padding: 2px 6px;
  font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
  font-size: 11px;
  color: var(--text-secondary);
}

/* 搜索加载状态 */
.search-loading {
  position: absolute;
  top: 50%;
  right: 12px;
  transform: translateY(-50%);
  color: var(--accent-primary);
}

.search-spinner {
  animation: search-spin 1s linear infinite;
}

@keyframes search-spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* 搜索下拉框样式 */
.search-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  box-shadow: 0 8px 32px var(--shadow-color);
  backdrop-filter: blur(16px);
  z-index: 1000;
  margin-top: 4px;
  max-height: 400px;
  overflow-y: auto;
}

.search-loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 24px;
  color: var(--text-muted);
  font-size: 14px;
}

.search-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
  color: var(--text-muted);
  font-size: 14px;
}

.empty-icon {
  color: var(--text-disabled);
}

.search-results {
  padding: 8px 0;
}

.search-result-item {
  padding: 14px 18px;
  cursor: pointer;
  border-bottom: 1px solid var(--bg-tertiary);
  transition: background-color 0.2s ease;
}

.search-result-item:hover {
  background: var(--surface-hover);
}

.search-result-item:last-child {
  border-bottom: none;
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.result-title {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.3;
  flex: 1;
  margin-right: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-date {
  font-size: 12px;
  color: var(--text-muted);
  white-space: nowrap;
}

.result-snippets {
  margin: 8px 0;
}

.snippet {
  font-size: 13px;
  line-height: 1.4;
  color: var(--text-secondary);
  margin-bottom: 4px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.snippet:last-child {
  margin-bottom: 0;
}

.snippet :global(mark) {
  background: var(--accent-primary-alpha);
  color: var(--accent-primary);
  border-radius: 2px;
  padding: 1px 2px;
  font-weight: 500;
}

.result-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
}

.match-tag {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: var(--radius-sm);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.title-match {
  background: var(--accent-primary-alpha);
  color: var(--accent-primary);
}

.content-match {
  background: var(--success-alpha);
  color: var(--success-primary);
}

.both-match {
  background: var(--warning-alpha);
  color: var(--warning-primary);
}

.search-footer {
  padding: 8px 16px;
  border-top: 1px solid var(--bg-tertiary);
  background: var(--bg-secondary);
}

.search-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text-muted);
}

.search-tip kbd {
  font-size: 10px;
  padding: 1px 4px;
}

.loading-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid var(--bg-tertiary);
  border-top-color: var(--accent-primary);
  border-radius: 50%;
  animation: search-spin 1s linear infinite;
}


.app-info {
  display: flex;
  align-items: center;
  gap: 12px;
  -webkit-app-region: drag;
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
  -webkit-app-region: no-drag;
}

.sidebar-toggle:hover, .theme-toggle:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.titlebar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  -webkit-app-region: no-drag;
}

.app-title {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 16px;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.025em;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  -webkit-app-region: drag;
}

.app-subtitle {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  -webkit-app-region: drag;
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
  -webkit-app-region: no-drag;
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

