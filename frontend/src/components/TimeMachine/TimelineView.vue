<template>
  <div class="timeline-view">
    <!-- 时间轴过滤器 -->
    <div class="timeline-filters">
      <div class="filter-group">
        <label class="filter-label">时间范围</label>
        <select v-model="selectedTimeRange" @change="applyFilters" class="filter-select">
          <option value="all">全部时间</option>
          <option value="week">最近一周</option>
          <option value="month">最近一个月</option>
          <option value="quarter">最近三个月</option>
          <option value="year">最近一年</option>
        </select>
      </div>
      
      <div class="filter-group">
        <label class="filter-label">情感筛选</label>
        <select v-model="selectedEmotion" @change="applyFilters" class="filter-select">
          <option value="all">全部情感</option>
          <option value="positive">积极情感</option>
          <option value="neutral">平静情感</option>
          <option value="negative">消极情感</option>
        </select>
      </div>
      
      <div class="filter-group">
        <label class="filter-label">排序方式</label>
        <select v-model="sortOrder" @change="applyFilters" class="filter-select">
          <option value="desc">最新在前</option>
          <option value="asc">最早在前</option>
        </select>
      </div>
      
      <button @click="resetFilters" class="reset-btn" title="重置筛选">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 12a9 9 0 013-6.364l1.5 1.5"/>
          <path d="M21 12a9 9 0 01-3 6.364l-1.5-1.5"/>
          <path d="M8 7.5l-1.5-1.5M16 16.5l1.5 1.5"/>
        </svg>
      </button>
    </div>

    <!-- 时间轴主体 -->
    <div class="timeline-container" ref="timelineContainer">
      <div v-if="filteredDiaries.length === 0" class="empty-timeline">
        <div class="empty-icon">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="3"/>
            <path d="M12 1v6m0 6v6"/>
            <path d="m21 12-6-6v12l6-6"/>
            <path d="m3 12 6-6v12l-6-6"/>
          </svg>
        </div>
        <h3>没有符合条件的日记</h3>
        <p>尝试调整筛选条件或创建新的日记</p>
      </div>
      
      <div v-else class="timeline-content">
        <!-- 时间轴线 -->
        <div class="timeline-line"></div>
        
        <!-- 时间轴项目 -->
        <div
          v-for="(group, index) in groupedDiaries"
          :key="`group-${index}`"
          class="timeline-group"
        >
          <!-- 日期标签 -->
          <div class="date-marker">
            <div class="date-circle">
              <span class="date-number">{{ group.day }}</span>
            </div>
            <div class="date-info">
              <span class="date-text">{{ group.dateText }}</span>
              <span class="diary-count">{{ group.diaries.length }} 篇</span>
            </div>
          </div>
          
          <!-- 该日期的日记列表 -->
          <div class="diaries-for-date">
            <div
              v-for="(diary, diaryIndex) in group.diaries"
              :key="diary.id"
              :class="[
                'timeline-item',
                { 'last-item': diaryIndex === group.diaries.length - 1 }
              ]"
              @click="selectDiary(diary)"
            >
              <!-- 时间点 -->
              <div class="timeline-point">
                <div :class="['timeline-dot', `emotion-${getEmotionLevel(diary)}`]"></div>
                <div class="timeline-time">{{ formatTime(diary.createdAt) }}</div>
              </div>
              
              <!-- 日记内容卡片 -->
              <div class="timeline-card">
                <div class="card-header">
                  <h3 class="card-title">{{ diary.title }}</h3>
                  <div class="card-badges">
                    <div :class="['emotion-badge', `emotion-${getEmotionLevel(diary)}`]">
                      {{ getEmotionText(diary) }}
                    </div>
                    <div v-if="diary.tags && diary.tags.length > 0" class="tags-badge">
                      {{ diary.tags.length }} 标签
                    </div>
                  </div>
                </div>
                
                <div class="card-content">
                  <p class="content-preview">{{ getPreview(diary.content) }}</p>
                </div>
                
                <div v-if="diary.tags && diary.tags.length > 0" class="card-tags">
                  <span
                    v-for="tag in diary.tags.slice(0, 4)"
                    :key="tag"
                    class="tag"
                  >
                    {{ tag }}
                  </span>
                  <span v-if="diary.tags.length > 4" class="more-tags">
                    +{{ diary.tags.length - 4 }}
                  </span>
                </div>
                
                <div class="card-footer">
                  <div class="card-stats">
                    <span class="word-count">{{ getWordCount(diary.content) }} 字</span>
                    <span class="read-time">约 {{ getReadTime(diary.content) }} 分钟</span>
                  </div>
                  
                  <button class="read-more-btn" @click.stop="selectDiary(diary)">
                    阅读全文
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="9,18 15,12 9,6"/>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 滚动到顶部按钮 -->
    <Transition name="scroll-btn">
      <button
        v-if="showScrollTop"
        @click="scrollToTop"
        class="scroll-top-btn"
        title="回到顶部"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="18,15 12,9 6,15"/>
        </svg>
      </button>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'

const props = defineProps({
  diaries: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['diary-select'])

// 状态管理
const selectedTimeRange = ref('all')
const selectedEmotion = ref('all')
const sortOrder = ref('desc')
const showScrollTop = ref(false)
const timelineContainer = ref(null)

// 过滤和排序后的日记
const filteredDiaries = computed(() => {
  let filtered = [...props.diaries]
  
  // 时间范围过滤
  if (selectedTimeRange.value !== 'all') {
    const now = new Date()
    const filterDate = new Date()
    
    switch (selectedTimeRange.value) {
      case 'week':
        filterDate.setDate(now.getDate() - 7)
        break
      case 'month':
        filterDate.setMonth(now.getMonth() - 1)
        break
      case 'quarter':
        filterDate.setMonth(now.getMonth() - 3)
        break
      case 'year':
        filterDate.setFullYear(now.getFullYear() - 1)
        break
    }
    
    filtered = filtered.filter(diary => 
      new Date(diary.createdAt) >= filterDate
    )
  }
  
  // 情感过滤
  if (selectedEmotion.value !== 'all') {
    filtered = filtered.filter(diary => 
      getEmotionLevel(diary) === selectedEmotion.value
    )
  }
  
  // 排序
  filtered.sort((a, b) => {
    const dateA = new Date(a.createdAt)
    const dateB = new Date(b.createdAt)
    return sortOrder.value === 'desc' ? dateB - dateA : dateA - dateB
  })
  
  return filtered
})

// 按日期分组的日记
const groupedDiaries = computed(() => {
  const groups = []
  const dateGroups = new Map()
  
  filteredDiaries.value.forEach(diary => {
    const date = new Date(diary.createdAt)
    const dateKey = date.toDateString()
    
    if (!dateGroups.has(dateKey)) {
      dateGroups.set(dateKey, {
        date,
        day: date.getDate(),
        dateText: formatDateText(date),
        diaries: []
      })
    }
    
    dateGroups.get(dateKey).diaries.push(diary)
  })
  
  // 按照sortOrder排序日期组
  const sortedGroups = Array.from(dateGroups.values()).sort((a, b) => {
    return sortOrder.value === 'desc' ? b.date - a.date : a.date - b.date
  })
  
  return sortedGroups
})

// 方法
const applyFilters = () => {
  // 触发重新计算
  nextTick(() => {
    scrollToTop()
  })
}

const resetFilters = () => {
  selectedTimeRange.value = 'all'
  selectedEmotion.value = 'all'
  sortOrder.value = 'desc'
  applyFilters()
}

const selectDiary = (diary) => {
  emit('diary-select', diary)
}

const scrollToTop = () => {
  if (timelineContainer.value) {
    timelineContainer.value.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
  }
}

const handleScroll = () => {
  if (timelineContainer.value) {
    showScrollTop.value = timelineContainer.value.scrollTop > 300
  }
}

const formatDateText = (date) => {
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(today.getDate() - 1)
  
  if (date.toDateString() === today.toDateString()) {
    return '今天'
  } else if (date.toDateString() === yesterday.toDateString()) {
    return '昨天'
  } else {
    return date.toLocaleDateString('zh-CN', {
      month: 'long',
      day: 'numeric',
      weekday: 'long'
    })
  }
}

const formatTime = (dateString) => {
  return new Date(dateString).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getEmotionLevel = (diary) => {
  if (diary.emotionScore) {
    if (diary.emotionScore >= 0.7) return 'positive'
    if (diary.emotionScore <= -0.7) return 'negative'
    return 'neutral'
  }
  
  // 基于内容长度和关键词的简单判断
  const content = diary.content.toLowerCase()
  const positiveWords = ['开心', '高兴', '快乐', '幸福', '满足', '成功', '好', '棒', '爱']
  const negativeWords = ['难过', '伤心', '痛苦', '失败', '困难', '压力', '烦恼', '焦虑', '失望']
  
  const positiveCount = positiveWords.filter(word => content.includes(word)).length
  const negativeCount = negativeWords.filter(word => content.includes(word)).length
  
  if (positiveCount > negativeCount) return 'positive'
  if (negativeCount > positiveCount) return 'negative'
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
  return content.length > 150 ? content.slice(0, 150) + '...' : content
}

const getWordCount = (content) => {
  return content.length
}

const getReadTime = (content) => {
  const wordsPerMinute = 300 // 中文阅读速度大约每分钟300字
  const minutes = Math.ceil(content.length / wordsPerMinute)
  return Math.max(1, minutes)
}

// 生命周期
onMounted(() => {
  if (timelineContainer.value) {
    timelineContainer.value.addEventListener('scroll', handleScroll)
  }
})

onUnmounted(() => {
  if (timelineContainer.value) {
    timelineContainer.value.removeEventListener('scroll', handleScroll)
  }
})
</script>

<style scoped>
.timeline-view {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 过滤器区域 */
.timeline-filters {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
  padding: 20px;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.filter-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.filter-select {
  padding: 8px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 120px;
}

.filter-select:hover {
  border-color: var(--accent-primary);
}

.filter-select:focus {
  outline: none;
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 3px var(--accent-primary-alpha);
}

.reset-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.2s ease;
  margin-top: 18px;
}

.reset-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  transform: rotate(180deg);
}

/* 时间轴容器 */
.timeline-container {
  flex: 1;
  overflow-y: auto;
  padding-right: 8px;
  position: relative;
}

.timeline-container::-webkit-scrollbar {
  width: 6px;
}

.timeline-container::-webkit-scrollbar-track {
  background: var(--bg-tertiary);
  border-radius: 3px;
}

.timeline-container::-webkit-scrollbar-thumb {
  background: var(--accent-primary);
  border-radius: 3px;
}

.timeline-container::-webkit-scrollbar-thumb:hover {
  background: var(--accent-secondary);
}

/* 空状态 */
.empty-timeline {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 40px;
  text-align: center;
}

.empty-icon {
  color: var(--text-muted);
  opacity: 0.6;
  margin-bottom: 20px;
}

.empty-timeline h3 {
  font-size: 20px;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.empty-timeline p {
  color: var(--text-muted);
  font-size: 14px;
}

/* 时间轴内容 */
.timeline-content {
  position: relative;
  padding-left: 80px;
  padding-bottom: 40px;
}

/* 时间轴主线 */
.timeline-line {
  position: absolute;
  left: 64px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: linear-gradient(to bottom, var(--accent-primary), var(--accent-primary-alpha));
  opacity: 0.6;
}

/* 时间轴组 */
.timeline-group {
  margin-bottom: 40px;
  position: relative;
}

/* 日期标记 */
.date-marker {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  position: relative;
  z-index: 2;
}

.date-circle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: var(--accent-primary);
  border-radius: 50%;
  color: white;
  font-weight: 700;
  position: absolute;
  left: -80px;
  box-shadow: 0 0 0 4px var(--bg-primary), var(--shadow-md);
}

.date-number {
  font-size: 16px;
}

.date-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-left: -16px;
}

.date-text {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.diary-count {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
}

/* 日记列表 */
.diaries-for-date {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-left: -16px;
}

/* 时间轴项目 */
.timeline-item {
  display: flex;
  gap: 20px;
  position: relative;
  margin-bottom: 16px;
}

.timeline-item:not(.last-item)::after {
  content: '';
  position: absolute;
  left: -64px;
  top: 40px;
  bottom: -16px;
  width: 2px;
  background: var(--bg-tertiary);
  opacity: 0.5;
}

/* 时间点 */
.timeline-point {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  position: relative;
  z-index: 2;
  margin-left: -80px;
  min-width: 80px;
}

.timeline-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 3px solid var(--bg-primary);
  position: relative;
  box-shadow: var(--shadow-sm);
}

.timeline-dot.emotion-positive {
  background: var(--success-primary);
}

.timeline-dot.emotion-negative {
  background: var(--error-primary);
}

.timeline-dot.emotion-neutral {
  background: var(--warning-primary);
}

.timeline-time {
  font-size: 11px;
  color: var(--text-muted);
  font-weight: 600;
  text-align: center;
  background: var(--bg-secondary);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--bg-tertiary);
  white-space: nowrap;
}

/* 时间轴卡片 */
.timeline-card {
  flex: 1;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-lg);
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: var(--shadow-sm);
}

.timeline-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
  border-color: var(--accent-primary);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
  gap: 16px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  flex: 1;
  line-height: 1.3;
}

.card-badges {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.emotion-badge {
  font-size: 11px;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  font-weight: 600;
  white-space: nowrap;
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

.tags-badge {
  font-size: 11px;
  padding: 4px 8px;
  background: var(--bg-tertiary);
  color: var(--text-muted);
  border-radius: var(--radius-sm);
  font-weight: 600;
  white-space: nowrap;
}

.card-content {
  margin-bottom: 16px;
}

.content-preview {
  font-size: 15px;
  line-height: 1.6;
  color: var(--text-secondary);
  margin: 0;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 16px;
}

.tag {
  font-size: 11px;
  padding: 4px 8px;
  background: var(--bg-tertiary);
  color: var(--text-muted);
  border-radius: var(--radius-sm);
  font-weight: 500;
}

.more-tags {
  font-size: 11px;
  padding: 4px 8px;
  background: var(--accent-primary-alpha);
  color: var(--accent-primary);
  border-radius: var(--radius-sm);
  font-weight: 600;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid var(--bg-tertiary);
}

.card-stats {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
}

.read-more-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: var(--accent-primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.read-more-btn:hover {
  background: var(--accent-secondary);
  transform: translateX(2px);
}

/* 滚动到顶部按钮 */
.scroll-top-btn {
  position: fixed;
  bottom: 32px;
  right: 32px;
  width: 48px;
  height: 48px;
  background: var(--accent-primary);
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  box-shadow: var(--shadow-lg);
  transition: all 0.3s ease;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
}

.scroll-top-btn:hover {
  background: var(--accent-secondary);
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

/* 滚动按钮动画 */
.scroll-btn-enter-active,
.scroll-btn-leave-active {
  transition: all 0.3s ease;
}

.scroll-btn-enter-from,
.scroll-btn-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.8);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .timeline-content {
    padding-left: 60px;
  }
  
  .date-circle {
    left: -60px;
    width: 40px;
    height: 40px;
  }
  
  .timeline-point {
    margin-left: -60px;
    min-width: 60px;
  }
  
  .timeline-line {
    left: 54px;
  }
}

@media (max-width: 768px) {
  .timeline-filters {
    padding: 16px;
    gap: 16px;
  }
  
  .filter-group {
    flex: 1;
    min-width: 100px;
  }
  
  .filter-select {
    min-width: auto;
  }
  
  .timeline-content {
    padding-left: 40px;
  }
  
  .date-circle {
    left: -40px;
    width: 32px;
    height: 32px;
  }
  
  .date-number {
    font-size: 14px;
  }
  
  .timeline-point {
    margin-left: -40px;
    min-width: 40px;
  }
  
  .timeline-line {
    left: 34px;
  }
  
  .timeline-item:not(.last-item)::after {
    left: -44px;
  }
  
  .timeline-card {
    padding: 20px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .card-badges {
    align-self: flex-end;
  }
}

@media (max-width: 480px) {
  .timeline-filters {
    padding: 12px;
    gap: 12px;
    flex-wrap: wrap;
  }
  
  .filter-group {
    flex: 1 1 auto;
    min-width: 80px;
  }
  
  .timeline-content {
    padding-left: 32px;
  }
  
  .date-circle {
    left: -32px;
    width: 28px;
    height: 28px;
  }
  
  .date-number {
    font-size: 12px;
  }
  
  .timeline-point {
    margin-left: -32px;
    min-width: 32px;
  }
  
  .timeline-line {
    left: 26px;
  }
  
  .timeline-item:not(.last-item)::after {
    left: -36px;
  }
  
  .timeline-card {
    padding: 16px;
  }
  
  .card-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .scroll-top-btn {
    width: 40px;
    height: 40px;
    bottom: 24px;
    right: 24px;
  }
}
</style>
