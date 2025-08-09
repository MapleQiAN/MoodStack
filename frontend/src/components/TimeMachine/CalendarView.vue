<template>
  <div class="calendar-view">
    <!-- 日历头部 -->
    <div class="calendar-header">
      <div class="calendar-navigation">
        <button @click="previousMonth" class="nav-btn" title="上个月">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15,18 9,12 15,6"/>
          </svg>
        </button>
        
        <div class="current-month">
          <h2 class="month-title">{{ currentMonthYear }}</h2>
          <button @click="goToToday" class="today-btn" v-if="!isCurrentMonth">今天</button>
        </div>
        
        <button @click="nextMonth" class="nav-btn" title="下个月">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9,6 15,12 9,18"/>
          </svg>
        </button>
      </div>
      
      <!-- 年月快速选择 -->
      <div class="quick-select">
        <select v-model="currentYear" @change="updateCalendar" class="year-select">
          <option v-for="year in availableYears" :key="year" :value="year">
            {{ year }}年
          </option>
        </select>
        <select v-model="currentMonth" @change="updateCalendar" class="month-select">
          <option v-for="(month, index) in monthNames" :key="index" :value="index">
            {{ month }}
          </option>
        </select>
      </div>
    </div>

    <!-- 日历主体 -->
    <div class="calendar-grid">
      <!-- 星期头部 -->
      <div class="weekday-header">
        <div v-for="day in weekDays" :key="day" class="weekday">
          {{ day }}
        </div>
      </div>
      
      <!-- 日期网格 -->
      <div class="days-grid">
        <div
          v-for="(day, index) in calendarDays"
          :key="index"
          :class="[
            'day-cell',
            {
              'other-month': !day.isCurrentMonth,
              'today': day.isToday,
              'has-diary': day.hasDiary,
              'selected': day.isSelected,
              'weekend': day.isWeekend
            }
          ]"
          @click="selectDate(day)"
        >
          <span class="day-number">{{ day.date }}</span>
          
          <!-- 日记指示器 -->
          <div v-if="day.diaries.length > 0" class="diary-indicators">
            <div
              v-for="(diary, diaryIndex) in day.diaries.slice(0, 3)"
              :key="diary.id"
              :class="['diary-dot', `emotion-${getEmotionLevel(diary)}`]"
              :title="diary.title"
              @click.stop="selectDiary(diary)"
            ></div>
            <div v-if="day.diaries.length > 3" class="more-indicator">
              +{{ day.diaries.length - 3 }}
            </div>
          </div>

          <!-- 悬浮预览：前三篇标题与时间 -->
          <div v-if="day.diaries.length > 0" class="day-hover-preview">
            <div
              v-for="diary in day.diaries.slice(0, 3)"
              :key="diary.id"
              class="preview-item"
            >
              <span class="preview-time">{{ formatTime(diary.createdAt) }}</span>
              <span class="preview-title">{{ diary.title || '无标题' }}</span>
            </div>
            <div v-if="day.diaries.length > 3" class="preview-more">
              还有 {{ day.diaries.length - 3 }} 篇...
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 选中日期的日记列表 -->
    <div v-if="selectedDay && selectedDay.diaries.length > 0" class="selected-day-diaries">
      <div class="selected-day-header">
        <h3>{{ formatSelectedDate(selectedDay) }}</h3>
        <span class="diary-count">{{ selectedDay.diaries.length }} 篇日记</span>
      </div>
      
      <div class="diary-list">
        <div
          v-for="diary in selectedDay.diaries"
          :key="diary.id"
          class="diary-item"
          @click="$emit('diary-select', diary)"
        >
          <div class="diary-header">
            <h4 class="diary-title">{{ diary.title }}</h4>
            <div class="diary-meta">
              <span class="diary-time">{{ formatTime(diary.createdAt) }}</span>
              <div :class="['emotion-badge', `emotion-${getEmotionLevel(diary)}`]">
                {{ getEmotionText(diary) }}
              </div>
            </div>
          </div>
          
          <p class="diary-preview">{{ getPreview(diary.content) }}</p>
          
          <div v-if="diary.tags && diary.tags.length > 0" class="diary-tags">
            <span v-for="tag in diary.tags.slice(0, 3)" :key="tag" class="tag">
              {{ tag }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态提示 -->
    <div v-else-if="selectedDay && selectedDay.diaries.length === 0" class="empty-day">
      <div class="empty-day-icon">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14,2 14,8 20,8"/>
          <line x1="12" y1="18" x2="12" y2="12"/>
          <line x1="9" y1="15" x2="15" y2="15"/>
        </svg>
      </div>
      <p>这天还没有写日记</p>
      <button @click="createDiary(selectedDay)" class="create-diary-btn">
        为这一天写日记
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'

const props = defineProps({
  diaries: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['diary-select', 'date-select', 'create-diary'])

// 日期相关状态
const currentDate = ref(new Date())
const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth())
const selectedDay = ref(null)

// 常量
const weekDays = ['日', '一', '二', '三', '四', '五', '六']
const monthNames = [
  '1月', '2月', '3月', '4月', '5月', '6月',
  '7月', '8月', '9月', '10月', '11月', '12月'
]

// 计算属性
const currentMonthYear = computed(() => {
  return `${currentYear.value}年${monthNames[currentMonth.value]}`
})

const isCurrentMonth = computed(() => {
  const today = new Date()
  return currentYear.value === today.getFullYear() && 
         currentMonth.value === today.getMonth()
})

const availableYears = computed(() => {
  const years = new Set()
  const currentYear = new Date().getFullYear()
  
  // 添加当前年份和前后几年
  for (let i = currentYear - 5; i <= currentYear + 1; i++) {
    years.add(i)
  }
  
  // 添加有日记的年份
  props.diaries.forEach(diary => {
    const year = new Date(diary.createdAt).getFullYear()
    years.add(year)
  })
  
  return Array.from(years).sort((a, b) => b - a)
})

const calendarDays = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startDate = new Date(firstDay)
  
  // 调整到周日开始
  startDate.setDate(startDate.getDate() - startDate.getDay())
  
  const days = []
  const today = new Date()
  
  // 生成42天（6周）
  for (let i = 0; i < 42; i++) {
    const date = new Date(startDate)
    date.setDate(startDate.getDate() + i)
    
    const isCurrentMonth = date.getMonth() === month
    const isToday = date.toDateString() === today.toDateString()
    const isWeekend = date.getDay() === 0 || date.getDay() === 6
    
    // 获取该日期的日记
    const dayDiaries = props.diaries.filter(diary => {
      const diaryDate = new Date(diary.createdAt)
      return diaryDate.toDateString() === date.toDateString()
    })
    
    days.push({
      date: date.getDate(),
      fullDate: new Date(date),
      isCurrentMonth,
      isToday,
      isWeekend,
      hasDiary: dayDiaries.length > 0,
      diaries: dayDiaries,
      isSelected: selectedDay.value && 
                 selectedDay.value.fullDate.toDateString() === date.toDateString()
    })
  }
  
  return days
})

// 方法
const previousMonth = () => {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
  updateCalendar()
}

const nextMonth = () => {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
  updateCalendar()
}

const goToToday = () => {
  const today = new Date()
  currentYear.value = today.getFullYear()
  currentMonth.value = today.getMonth()
  updateCalendar()
}

const updateCalendar = () => {
  // 触发日历重新计算
  currentDate.value = new Date(currentYear.value, currentMonth.value, 1)
}

const selectDate = (day) => {
  selectedDay.value = day
  emit('date-select', day.fullDate)
}

const selectDiary = (diary) => {
  emit('diary-select', diary)
}

const createDiary = (day) => {
  emit('create-diary', { date: day.fullDate })
}

const formatSelectedDate = (day) => {
  return day.fullDate.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
}

const formatTime = (dateString) => {
  return new Date(dateString).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getEmotionLevel = (diary) => {
  // 简单的情感分析，可以根据实际的情感分析结果调整
  if (diary.emotionScore) {
    if (diary.emotionScore >= 0.7) return 'positive'
    if (diary.emotionScore <= -0.7) return 'negative'
    return 'neutral'
  }
  
  // 基于内容长度的简单判断
  if (diary.content.length > 500) return 'positive'
  if (diary.content.length < 100) return 'negative'
  return 'neutral'
}

const getEmotionText = (diary) => {
  const level = getEmotionLevel(diary)
  const emotions = {
    positive: '😊 积极',
    negative: '😔 消极',
    neutral: '😐 平静'
  }
  return emotions[level] || '😐 平静'
}

const getPreview = (content) => {
  return content.length > 120 ? content.slice(0, 120) + '...' : content
}

// 监听器
watch([currentYear, currentMonth], () => {
  selectedDay.value = null
})

onMounted(() => {
  // 如果有日记，选择最近的一天
  if (props.diaries.length > 0) {
    const latestDiary = props.diaries.reduce((latest, diary) => {
      return new Date(diary.createdAt) > new Date(latest.createdAt) ? diary : latest
    })
    
    const latestDate = new Date(latestDiary.createdAt)
    currentYear.value = latestDate.getFullYear()
    currentMonth.value = latestDate.getMonth()
    updateCalendar()
  }
})
</script>

<style scoped>
.calendar-view {
  max-width: 100%;
}

/* 日历头部 */
.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.calendar-navigation {
  display: flex;
  align-items: center;
  gap: 16px;
}

.nav-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: none;
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.2s ease;
}

.nav-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  transform: scale(1.05);
}

.current-month {
  display: flex;
  align-items: center;
  gap: 12px;
}

.month-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  min-width: 140px;
  text-align: center;
}

.today-btn {
  padding: 6px 12px;
  background: var(--accent-primary);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.today-btn:hover {
  background: var(--accent-secondary);
  transform: translateY(-1px);
}

.quick-select {
  display: flex;
  gap: 8px;
}

.year-select,
.month-select {
  padding: 8px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.year-select:hover,
.month-select:hover {
  border-color: var(--accent-primary);
}

.year-select:focus,
.month-select:focus {
  outline: none;
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 3px var(--accent-primary-alpha);
}

/* 日历网格 */
.calendar-grid {
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.weekday-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  background: var(--bg-tertiary);
}

.weekday {
  padding: 16px 8px;
  text-align: center;
  font-weight: 600;
  font-size: 14px;
  color: var(--text-secondary);
  border-right: 1px solid var(--bg-quaternary);
}

.weekday:last-child {
  border-right: none;
}

.days-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
}

.day-cell {
  min-height: 100px;
  padding: 8px;
  border-right: 1px solid var(--bg-tertiary);
  border-bottom: 1px solid var(--bg-tertiary);
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  background: var(--bg-primary);
}

.day-cell:nth-child(7n) {
  border-right: none;
}

.day-cell:hover {
  background: var(--surface-hover);
}

.day-cell.other-month {
  color: var(--text-disabled);
  background: var(--bg-secondary);
}

.day-cell.other-month .day-number {
  opacity: 0.5;
}

.day-cell.today {
  background: var(--accent-primary-alpha);
  border: 2px solid var(--accent-primary);
}

.day-cell.today .day-number {
  color: var(--accent-primary);
  font-weight: 700;
}

.day-cell.selected {
  background: var(--accent-secondary-alpha);
  border: 2px solid var(--accent-secondary);
}

.day-cell.weekend:not(.other-month) {
  background: var(--bg-glass);
}

.day-number {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.diary-indicators {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: auto;
}

.diary-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s ease;
}

.diary-dot:hover {
  transform: scale(1.3);
}

.diary-dot.emotion-positive {
  background: var(--success-primary);
}

.diary-dot.emotion-negative {
  background: var(--error-primary);
}

.diary-dot.emotion-neutral {
  background: var(--warning-primary);
}

.more-indicator {
  font-size: 10px;
  color: var(--text-muted);
  font-weight: 500;
  background: var(--bg-tertiary);
  padding: 2px 4px;
  border-radius: var(--radius-sm);
}

/* 悬浮预览 Tooltip */
.day-hover-preview {
  position: absolute;
  left: -4px;
  right: auto;
  top: 32px;
  min-width: 240px;
  max-width: 320px;
  background: var(--bg-primary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  padding: 12px;
  z-index: 100;
  opacity: 0;
  transform: translateY(-4px) scale(0.98);
  transition: opacity 0.15s ease, transform 0.15s ease;
  pointer-events: none;
}

.day-cell:hover .day-hover-preview {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.preview-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-primary);
  line-height: 1.5;
  margin-bottom: 4px;
}

.preview-item:last-child {
  margin-bottom: 0;
}

.preview-time {
  color: var(--text-muted);
  font-variant-numeric: tabular-nums;
  flex-shrink: 0;
}

.preview-title {
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  flex: 1;
}

.preview-more {
  margin-top: 6px;
  font-size: 11px;
  color: var(--text-muted);
}

/* 选中日期的日记列表 */
.selected-day-diaries {
  margin-top: 32px;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.selected-day-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--bg-tertiary);
}

.selected-day-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.diary-count {
  font-size: 14px;
  color: var(--text-muted);
  font-weight: 500;
}

.diary-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.diary-item {
  background: var(--bg-primary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.diary-item:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--accent-primary);
}

.diary-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  gap: 16px;
}

.diary-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  flex: 1;
  line-height: 1.3;
}

.diary-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.diary-time {
  font-size: 12px;
  color: var(--text-muted);
}

.emotion-badge {
  font-size: 11px;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  font-weight: 500;
}

.emotion-badge.emotion-positive {
  background: var(--success-alpha);
  color: var(--success-primary);
}

.emotion-badge.emotion-negative {
  background: var(--error-alpha);
  color: var(--error-primary);
}

.emotion-badge.emotion-neutral {
  background: var(--warning-alpha);
  color: var(--warning-primary);
}

.diary-preview {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin: 0 0 12px 0;
}

.diary-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tag {
  font-size: 11px;
  padding: 4px 8px;
  background: var(--bg-tertiary);
  color: var(--text-muted);
  border-radius: var(--radius-sm);
  font-weight: 500;
}

/* 空状态 */
.empty-day {
  margin-top: 32px;
  text-align: center;
  padding: 40px;
  background: var(--bg-glass);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-lg);
}

.empty-day-icon {
  color: var(--text-muted);
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-day p {
  font-size: 16px;
  color: var(--text-muted);
  margin-bottom: 20px;
}

.create-diary-btn {
  padding: 10px 20px;
  background: var(--accent-primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
}

.create-diary-btn:hover {
  background: var(--accent-secondary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .calendar-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .calendar-navigation {
    justify-content: center;
  }
  
  .quick-select {
    justify-content: center;
  }
  
  .day-cell {
    min-height: 80px;
    padding: 6px;
  }
  
  .month-title {
    font-size: 20px;
    min-width: 120px;
  }
  
  .selected-day-diaries {
    padding: 16px;
  }
  
  .diary-item {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .day-cell {
    min-height: 60px;
    padding: 4px;
  }
  
  .day-number {
    font-size: 14px;
  }
  
  .diary-dot {
    width: 6px;
    height: 6px;
  }
  
  .weekday {
    padding: 12px 4px;
    font-size: 12px;
  }
  
  .diary-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
