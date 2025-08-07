<template>
  <div class="list-container">
    <div class="list-header">
            <div class="header-content">
        <h1 class="page-title">我的日记</h1>
        <p class="page-description">共 {{ filteredDiaries.length }} 篇日记记录</p>
      </div>
    
      <div class="header-actions">
        <div class="add-diary-dropdown" :class="{ open: showDropdown }">
          <button @click="toggleDropdown" class="dropdown-trigger" :disabled="isUploading">
          <span class="add-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
          </span>
            <span v-if="!isUploading">新增日记</span>
            <span v-else>处理中...</span>
            <span class="dropdown-arrow">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="6,9 12,15 18,9"/>
            </svg>
          </span>
          </button>

          <div v-if="showDropdown" class="dropdown-menu">
            <button @click="createNewDiary" class="dropdown-item">
            <span class="item-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
                <line x1="12" y1="11" x2="12" y2="17"/>
                <line x1="9" y1="14" x2="15" y2="14"/>
              </svg>
            </span>
              <div class="item-content">
                <span class="item-title">新建日记</span>
                <span class="item-description">创建并直接编辑</span>
              </div>
            </button>

            <button @click="selectFile" class="dropdown-item">
            <span class="item-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
                <line x1="12" y1="11" x2="12" y2="17"/>
                <polyline points="9,14 12,11 15,14"/>
              </svg>
            </span>
              <div class="item-content">
                <span class="item-title">文件导入</span>
                <span class="item-description">上传现有文件</span>
              </div>
            </button>
          </div>
        </div>

                 <input
             ref="fileInput"
             type="file"
             accept=".txt,.md,.docx,.pdf"
             @change="handleFileSelect"
             style="display: none"
         />
         
         <!-- 排序和筛选控件 -->
         <div class="controls-section">
           <div class="control-item sort-control" @click="toggleSortOrder">
             <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                  :style="{ transform: sortOrder === 'desc' ? 'rotate(0deg)' : 'rotate(180deg)' }">
               <polyline points="6,9 12,15 18,9"/>
             </svg>
             <span class="control-text">{{ sortOrder === 'desc' ? '最新在前' : '最旧在前' }}</span>
           </div>

           <div class="control-item date-control">
             <div @click="showDatePicker = true" class="date-trigger">
               <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                 <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                 <line x1="16" y1="2" x2="16" y2="6"/>
                 <line x1="8" y1="2" x2="8" y2="6"/>
                 <line x1="3" y1="10" x2="21" y2="10"/>
               </svg>
               <span class="control-text">{{ selectedDate ? formatSelectedDateShort(selectedDate) : '选择日期' }}</span>
             </div>

             <!-- 自定义日期选择器 -->
             <div v-if="showDatePicker" class="custom-date-picker">
               <div class="date-picker-header">
                 <button @click="previousMonth" class="nav-btn">
                   <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                     <polyline points="15,18 9,12 15,6"/>
                   </svg>
                 </button>
                 <span class="month-year">{{ currentMonth }}月 {{ currentYear }}</span>
                 <button @click="nextMonth" class="nav-btn">
                   <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                     <polyline points="9,18 15,12 9,6"/>
                   </svg>
                 </button>
               </div>
               
               <div class="weekdays">
                 <span v-for="day in weekdays" :key="day" class="weekday">{{ day }}</span>
               </div>
               
               <div class="calendar-grid">
                 <button
                   v-for="date in calendarDays"
                   :key="date.key"
                   @click="selectDate(date)"
                   :class="[
                     'calendar-day',
                     { 
                       'other-month': !date.isCurrentMonth,
                       'selected': isSelectedDate(date),
                       'today': isToday(date)
                     }
                   ]"
                   :disabled="!date.isCurrentMonth"
                 >
                   {{ date.day }}
                 </button>
               </div>
               
               <div class="date-picker-footer">
                 <button @click="clearDateFilter" class="clear-btn" v-if="selectedDate">清除</button>
                 <button @click="selectToday" class="today-btn">今天</button>
               </div>
             </div>
           </div>
         </div>
       </div>

      <div v-if="diaries.length === 0" class="empty-state">
        <div class="empty-visual">
          <div class="empty-icon">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14,2 14,8 20,8"/>
            </svg>
          </div>
          <div class="empty-decoration"></div>
        </div>
        <div class="empty-content">
          <h3>还没有日记</h3>
          <p>开始记录你的第一个心情故事吧</p>
          <div class="empty-actions">
            <button @click="createNewDiary" class="empty-action primary">
              新建日记
            </button>
            <button @click="selectFile" class="empty-action secondary">
              导入文件
            </button>
          </div>
        </div>
      </div>

             <div v-else class="diary-grid">
         <div
             v-for="diary in filteredDiaries"
             :key="diary.id"
             class="diary-card"
             @click="$emit('view-diary', diary)"
         >
          <div class="card-header">
          </div>

          <div class="card-body">
            <h3 class="diary-title">{{ diary.title }}</h3>
            <p class="diary-preview">{{ getPreview(diary.content) }}</p>
          </div>

          <div class="card-footer">
            <div class="file-info">
              <span class="file-name">{{ diary.fileName }}</span>
            </div>
            <div class="footer-right">
              <span class="diary-date">{{ formatSpecificDate(diary.createdAt) }}</span>
              <div class="read-more">
                <span>阅读</span>
                <span class="arrow">→</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {computed, onMounted, onUnmounted, ref} from 'vue'
import {UploadDiary} from '../../wailsjs/go/main/App'

const props = defineProps({
  diaries: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['view-diary', 'refresh', 'upload-success', 'create-diary'])

const fileInput = ref(null)
const selectedFile = ref(null)
const isUploading = ref(false)
const uploadError = ref('')
const showDropdown = ref(false)
const sortOrder = ref('desc') // 'desc' for newest first, 'asc' for oldest first
const selectedDate = ref('')
const showDatePicker = ref(false)
const tempSelectedDate = ref('')
// 本地时间处理函数，避免时区问题
const getLocalDateString = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth() + 1)
const weekdays = ['日', '一', '二', '三', '四', '五', '六']

const filteredDiaries = computed(() => {
  if (!props.diaries) return []

  let filtered = [...props.diaries]

  // 按日期筛选
  if (selectedDate.value) {
    filtered = filtered.filter(diary => {
      const diaryDate = getLocalDateString(new Date(diary.createdAt))
      return diaryDate === selectedDate.value
    })
  }

  // 排序
  return filtered.sort((a, b) => {
    const dateA = new Date(a.createdAt)
    const dateB = new Date(b.createdAt)
    return sortOrder.value === 'desc' ? dateB - dateA : dateA - dateB
  })
})

// 生成日历天数
const calendarDays = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const firstDay = new Date(year, month - 1, 1)
  const lastDay = new Date(year, month, 0)
  const daysInMonth = lastDay.getDate()
  const startDate = new Date(firstDay)
  startDate.setDate(startDate.getDate() - firstDay.getDay())
  
  const days = []
  const today = new Date()
  
  for (let i = 0; i < 42; i++) {
    const date = new Date(startDate)
    date.setDate(startDate.getDate() + i)
    
    const isCurrentMonth = date.getMonth() === month - 1
    
    days.push({
      day: date.getDate(),
      date: date,
      isCurrentMonth,
      key: `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}`
    })
  }
  
  return days
})

const formatDate = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = now - date
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays === 0) {
    return '今天'
  } else if (diffDays === 1) {
    return '昨天'
  } else if (diffDays < 7) {
    return `${diffDays} 天前`
  } else {
    return date.toLocaleDateString('zh-CN', {
      month: 'long',
      day: 'numeric'
    })
  }
}

const formatSpecificDate = (dateString) => {
  const date = new Date(dateString)
  return getLocalDateString(date)
}

// 排序和筛选相关函数
const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
}

const clearDateFilter = () => {
  selectedDate.value = ''
}

const formatSelectedDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 格式化选中日期为简短格式
const formatSelectedDateShort = (dateString) => {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${year}年${month}月${day}日`
}

// 日期选择器相关函数
const previousMonth = () => {
  if (currentMonth.value === 1) {
    currentMonth.value = 12
    currentYear.value--
  } else {
    currentMonth.value--
  }
}

const nextMonth = () => {
  if (currentMonth.value === 12) {
    currentMonth.value = 1
    currentYear.value++
  } else {
    currentMonth.value++
  }
}

const selectDate = (dateObj) => {
  if (!dateObj.isCurrentMonth) return
  
  const dateString = getLocalDateString(dateObj.date)
  selectedDate.value = dateString
  showDatePicker.value = false
}

const isSelectedDate = (dateObj) => {
  if (!selectedDate.value) return false
  const dateString = getLocalDateString(dateObj.date)
  return dateString === selectedDate.value
}

const isToday = (dateObj) => {
  const today = new Date()
  const dateString = getLocalDateString(dateObj.date)
  const todayString = getLocalDateString(today)
  return dateString === todayString
}

const selectToday = () => {
  const today = new Date()
  selectedDate.value = getLocalDateString(today)
  showDatePicker.value = false
  
  // 跳转到今天所在的月份
  currentYear.value = today.getFullYear()
  currentMonth.value = today.getMonth() + 1
}

// 移除emoji图标函数，改用SVG图标

const toggleDropdown = () => {
  showDropdown.value = !showDropdown.value
}

const createNewDiary = () => {
  showDropdown.value = false
  emit('create-diary')
}

const selectFile = () => {
  showDropdown.value = false
  fileInput.value.click()
}

const handleFileSelect = async (event) => {
  const file = event.target.files[0]
  if (file) {
    selectedFile.value = file
    uploadError.value = ''
    await uploadFile()
  }
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const uploadFile = async () => {
  if (!selectedFile.value) return

  // Check file size (10MB limit)
  if (selectedFile.value.size > 10 * 1024 * 1024) {
    uploadError.value = '文件大小不能超过 10MB'
    alert('文件大小不能超过 10MB')
    return
  }

  // Check file type
  const allowedTypes = ['.txt', '.md', '.docx', '.pdf']
  const fileExt = '.' + selectedFile.value.name.split('.').pop().toLowerCase()
  if (!allowedTypes.includes(fileExt)) {
    uploadError.value = '不支持的文件类型，请选择 .txt、.md、.docx 或 .pdf 文件'
    alert('不支持的文件类型，请选择 .txt、.md、.docx 或 .pdf 文件')
    return
  }

  try {
    isUploading.value = true
    uploadError.value = ''

    // Read file content
    const arrayBuffer = await selectedFile.value.arrayBuffer()
    const uint8Array = new Uint8Array(arrayBuffer)

    // Convert to regular array for Go
    const content = Array.from(uint8Array)

    // Upload file
    await UploadDiary(selectedFile.value.name, content)

    // Reset form
    selectedFile.value = null
    fileInput.value.value = ''

    // 发送成功事件
    emit('upload-success')

  } catch (error) {
    console.error('Upload error:', error)
    uploadError.value = error.message || '上传失败，请重试'
    alert(error.message || '上传失败，请重试')
  } finally {
    isUploading.value = false
  }
}

const getPreview = (content) => {
  if (!content) return '暂无内容预览'

  // Remove markdown formatting for preview
  let preview = content
      .replace(/^#+\s+/gm, '') // Remove headers
      .replace(/\*\*(.*?)\*\*/g, '$1') // Remove bold
      .replace(/\*(.*?)\*/g, '$1') // Remove italic
      .replace(/`(.*?)`/g, '$1') // Remove code
      .replace(/\n+/g, ' ') // Replace newlines with spaces
      .trim()

  return preview.length > 100 ? preview.substring(0, 100) + '...' : preview
}


// 点击外部关闭下拉菜单和日期选择器
const handleClickOutside = (event) => {
  if (!event.target.closest('.add-diary-dropdown')) {
    showDropdown.value = false
  }
  if (!event.target.closest('.date-control')) {
    showDatePicker.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.list-container {
  max-width: 100%;
  padding: 32px;
  box-sizing: border-box;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
  flex-wrap: wrap;
  gap: 20px;
}

.header-actions {
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: flex-end;
}

.controls-section {
  display: flex;
  gap: 16px;
  align-items: center;
}

.header-content {
  flex: 1;
}

.page-title {
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
  letter-spacing: -0.025em;
}

.page-description {
  font-size: 16px;
  color: var(--text-muted);
  font-weight: 500;
}

/* 简单的图标+文字控件样式 */
.control-item {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  color: var(--text-muted);
  transition: all 0.2s ease;
  position: relative;
}

.control-item:hover {
  color: var(--accent-primary);
}

.control-text {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}

.date-control {
  position: relative;
}

.date-trigger {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 自定义日期选择器样式 */
.custom-date-picker {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  padding: 20px;
  z-index: 1000;
  margin-top: 8px;
  min-width: 300px;
}

.date-picker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.nav-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--text-muted);
  transition: all 0.2s ease;
}

.nav-btn:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border-color: var(--accent-primary);
}

.month-year {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  margin-bottom: 8px;
}

.weekday {
  text-align: center;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  padding: 8px 4px;
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  margin-bottom: 16px;
}

.calendar-day {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  transition: all 0.2s ease;
}

.calendar-day:hover:not(:disabled) {
  background: var(--bg-secondary);
  color: var(--accent-primary);
}

.calendar-day.other-month {
  color: var(--text-muted);
  opacity: 0.5;
}

.calendar-day.selected {
  background: var(--accent-primary);
  color: white;
}

.calendar-day.today {
  background: var(--bg-tertiary);
  font-weight: 600;
}

.calendar-day.today.selected {
  background: var(--accent-primary);
}

.calendar-day:disabled {
  cursor: not-allowed;
  opacity: 0.3;
}

.date-picker-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.clear-btn, .today-btn {
  padding: 8px 16px;
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.clear-btn {
  background: transparent;
  color: var(--text-muted);
}

.clear-btn:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border-color: var(--accent-primary);
}

.today-btn {
  background: var(--accent-primary);
  color: white;
  border-color: var(--accent-primary);
}

.today-btn:hover {
  background: var(--accent-secondary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}


/* 下拉菜单样式 */
.add-diary-dropdown {
  position: relative;
}

.dropdown-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: var(--accent-primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s ease;
  box-shadow: var(--shadow-sm);
}

.dropdown-trigger:hover:not(:disabled) {
  background: var(--accent-secondary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.dropdown-trigger:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.add-icon {
  font-size: 16px;
}

.dropdown-arrow {
  font-size: 14px;
  transition: transform 0.2s ease;
}

.add-diary-dropdown.open .dropdown-arrow {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
  z-index: 1000;
  min-width: 240px;
  margin-top: 8px;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 16px 20px;
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}

.dropdown-item:hover {
  background: var(--bg-secondary);
}

.dropdown-item:not(:last-child) {
  border-bottom: 1px solid var(--bg-tertiary);
}

.item-icon {
  color: var(--accent-primary);
  flex-shrink: 0;
}

.item-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.item-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.item-description {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 40px;
  text-align: center;
}

.empty-visual {
  position: relative;
  margin-bottom: 32px;
}

.empty-icon {
  color: var(--text-muted);
  opacity: 0.6;
  position: relative;
  z-index: 2;
}

.empty-decoration {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 120px;
  height: 120px;
  border: 2px solid var(--bg-tertiary);
  border-radius: 50%;
  z-index: 1;
  opacity: 0.3;
}

.empty-content h3 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.empty-content p {
  font-size: 16px;
  color: var(--text-muted);
  margin-bottom: 32px;
  font-weight: 500;
}

.empty-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
}

.empty-action {
  padding: 12px 24px;
  border: none;
  border-radius: var(--radius-md);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
}

.empty-action.primary {
  background: var(--accent-primary);
  color: white;
}

.empty-action.primary:hover {
  background: var(--accent-secondary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.empty-action.secondary {
  background: transparent;
  color: var(--accent-primary);
  border: 1px solid var(--accent-primary);
}

.empty-action.secondary:hover {
  background: var(--accent-primary);
  color: white;
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.diary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 24px;
}

.diary-card {
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-xl);
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid var(--bg-tertiary);
  box-shadow: var(--shadow-sm);
  position: relative;
  overflow: hidden;
}

.diary-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--accent-primary), var(--nord15));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.diary-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
  border-color: var(--nord3);
}

.diary-card:hover::before {
  opacity: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  position: relative;
}


.edit-btn {
  position: absolute;
  top: 0;
  right: 0;
  width: 32px;
  height: 32px;
  background: var(--accent-primary);
  border: none;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  opacity: 0;
  transform: scale(0.8);
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.diary-card:hover .edit-btn {
  opacity: 1;
  transform: scale(1);
}

.edit-btn:hover {
  background: var(--accent-primary-hover);
  transform: scale(1.05);
}

.card-body {
  margin-bottom: 20px;
}

.diary-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.diary-preview {
  color: var(--text-muted);
  line-height: 1.6;
  font-size: 14px;
  font-weight: 500;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid var(--bg-tertiary);
}

.file-info {
  flex: 1;
}

.file-name {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
  background: var(--bg-secondary);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--bg-tertiary);
  display: inline-block;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.footer-right .diary-date {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 600;
  background: var(--bg-secondary);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--bg-tertiary);
}

.read-more {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 600;
  color: var(--accent-primary);
}

.arrow {
  font-size: 14px;
  transition: transform 0.2s ease;
}

.diary-card:hover .arrow {
  transform: translateX(4px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .list-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .header-actions {
    align-items: center;
  }

  .dropdown-menu {
    right: 0;
    left: auto;
  }

  .diary-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .diary-card {
    padding: 20px;
  }

  .page-title {
    font-size: 28px;
  }

  .empty-state {
    padding: 60px 20px;
  }
}

@media (max-width: 480px) {
  .header-actions {
    align-items: center;
  }

  .controls-section {
    justify-content: center;
  }
  
  .control-text {
    font-size: 13px;
  }
  
  .custom-date-picker {
    min-width: 280px;
    left: 50%;
    transform: translateX(-50%);
    right: auto;
  }
  
  .calendar-day {
    width: 32px;
    height: 32px;
    font-size: 13px;
  }
  
  .card-footer {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .footer-right {
    justify-content: space-between;
  }

  .read-more {
    justify-content: center;
  }
}
</style>