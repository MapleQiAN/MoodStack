<template>
  <div class="time-machine-container">
    <!-- 头部区域 -->
    <div class="time-machine-header">
      <div class="header-title">
        <h1 class="page-title">时光机</h1>
        <p class="page-description">回忆曾经写过的日记，用不同的方式重新发现你的情感历程</p>
      </div>
      
      <!-- 视图切换器 -->
      <div class="view-switcher">
        <button 
          @click="currentView = 'calendar'" 
          :class="['view-btn', { active: currentView === 'calendar' }]"
          title="日历视图"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
            <line x1="16" y1="2" x2="16" y2="6"/>
            <line x1="8" y1="2" x2="8" y2="6"/>
            <line x1="3" y1="10" x2="21" y2="10"/>
          </svg>
          <span>日历</span>
        </button>
        
        <button 
          @click="currentView = 'timeline'" 
          :class="['view-btn', { active: currentView === 'timeline' }]"
          title="时间轴视图"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M12 1v6m0 6v6"/>
            <path d="m21 12-6-6v12l6-6"/>
            <path d="m3 12 6-6v12l-6-6"/>
          </svg>
          <span>时间轴</span>
        </button>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="stats-section" v-if="diaries.length > 0">
      <div class="stat-card">
        <div class="stat-icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14,2 14,8 20,8"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ diaries.length }}</span>
          <span class="stat-label">篇日记</span>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12,6 12,12 16,14"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ daysSinceFirst }}</span>
          <span class="stat-label">天旅程</span>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ activeDays }}</span>
          <span class="stat-label">活跃天数</span>
        </div>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="main-content">
      <Transition name="view-transition" mode="out-in">
        <!-- 日历视图 -->
        <CalendarView
          v-if="currentView === 'calendar'"
          key="calendar"
          :diaries="diaries"
          @diary-select="handleDiarySelect"
          @date-select="handleDateSelect"
        />
        
        <!-- 时间轴视图 -->
        <TimelineView
          v-else-if="currentView === 'timeline'"
          key="timeline"
          :diaries="diaries"
          @diary-select="handleDiarySelect"
        />
      </Transition>
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && diaries.length === 0" class="empty-state">
      <div class="empty-visual">
        <div class="empty-icon">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="3"/>
            <path d="M12 1v6m0 6v6"/>
            <path d="m21 12-6-6v12l6-6"/>
            <path d="m3 12 6-6v12l-6-6"/>
          </svg>
        </div>
        <div class="empty-decoration">
          <div class="decoration-circle"></div>
          <div class="decoration-circle"></div>
          <div class="decoration-circle"></div>
        </div>
      </div>
      <div class="empty-content">
        <h3>时光机还很空</h3>
        <p>开始写日记，让时光机记录你的情感历程</p>
        <button @click="$emit('create-diary')" class="empty-action">
          开始写日记
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p class="loading-text">正在加载时光记录...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import CalendarView from './TimeMachine/CalendarView.vue'
import TimelineView from './TimeMachine/TimelineView.vue'

const props = defineProps({
  diaries: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['diary-select', 'create-diary'])

const currentView = ref('calendar')
const loading = ref(false)

// 计算统计数据
const daysSinceFirst = computed(() => {
  if (props.diaries.length === 0) return 0
  
  const firstDiary = props.diaries.reduce((earliest, diary) => {
    return new Date(diary.createdAt) < new Date(earliest.createdAt) ? diary : earliest
  })
  
  const firstDate = new Date(firstDiary.createdAt)
  const today = new Date()
  const diffTime = Math.abs(today - firstDate)
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
})

const activeDays = computed(() => {
  if (props.diaries.length === 0) return 0
  
  const uniqueDates = new Set()
  props.diaries.forEach(diary => {
    const date = new Date(diary.createdAt).toDateString()
    uniqueDates.add(date)
  })
  
  return uniqueDates.size
})

// 处理日记选择
const handleDiarySelect = (diary) => {
  emit('diary-select', diary)
}

// 处理日期选择
const handleDateSelect = (date) => {
  // 可以在这里添加日期选择后的逻辑，比如过滤或高亮显示
  console.log('选择了日期:', date)
}
</script>

<style scoped>
.time-machine-container {
  padding: 32px;
  min-height: 100%;
  max-width: 1400px;
  margin: 0 auto;
}

/* 头部区域 */
.time-machine-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
  gap: 24px;
}

.header-title {
  flex: 1;
}

.page-title {
  font-family: var(--font-display);
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
  letter-spacing: -0.025em;
  background: linear-gradient(135deg, var(--accent-primary), var(--nord15));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-description {
  font-size: 16px;
  color: var(--text-muted);
  font-weight: 500;
  line-height: 1.5;
  max-width: 500px;
}

/* 视图切换器 */
.view-switcher {
  display: flex;
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-lg);
  padding: 4px;
  box-shadow: var(--shadow-sm);
}

.view-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  color: var(--text-secondary);
  transition: all 0.2s ease;
  position: relative;
}

.view-btn:hover {
  color: var(--text-primary);
  background: var(--bg-tertiary);
}

.view-btn.active {
  background: var(--accent-primary);
  color: white;
  box-shadow: var(--shadow-sm);
}

.view-btn svg {
  display: block;
  flex-shrink: 0;
}

/* 统计信息区域 */
.stats-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-lg);
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: var(--shadow-sm);
  transition: all 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--accent-primary);
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: var(--accent-primary-alpha);
  border-radius: var(--radius-md);
  color: var(--accent-primary);
  flex-shrink: 0;
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: var(--text-muted);
  font-weight: 500;
}

/* 主内容区域 */
.main-content {
  position: relative;
  min-height: 400px;
}

/* 视图切换动画 */
.view-transition-enter-active,
.view-transition-leave-active {
  transition: all 0.3s ease;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
}

.view-transition-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.view-transition-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

.view-transition-enter-to,
.view-transition-leave-from {
  opacity: 1;
  transform: translateX(0);
}

/* 空状态 */
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
  z-index: 1;
}

.decoration-circle {
  position: absolute;
  border: 2px solid var(--bg-tertiary);
  border-radius: 50%;
  opacity: 0.3;
}

.decoration-circle:nth-child(1) {
  width: 80px;
  height: 80px;
  top: -40px;
  left: -40px;
  animation: float 6s ease-in-out infinite;
}

.decoration-circle:nth-child(2) {
  width: 120px;
  height: 120px;
  top: -60px;
  left: -60px;
  animation: float 8s ease-in-out infinite reverse;
}

.decoration-circle:nth-child(3) {
  width: 160px;
  height: 160px;
  top: -80px;
  left: -80px;
  animation: float 10s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-10px) rotate(180deg);
  }
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

.empty-action {
  padding: 12px 24px;
  background: var(--accent-primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
}

.empty-action:hover {
  background: var(--accent-secondary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 40px;
  gap: 20px;
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

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .stats-section {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 16px;
  }
  
  .stat-card {
    padding: 20px;
  }
}

@media (max-width: 768px) {
  .time-machine-container {
    padding: 24px 20px;
  }
  
  .time-machine-header {
    flex-direction: column;
    align-items: stretch;
    gap: 20px;
  }
  
  .page-title {
    font-size: 28px;
  }
  
  .stats-section {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .stat-card {
    padding: 16px;
  }
  
  .view-switcher {
    align-self: center;
  }
}

@media (max-width: 480px) {
  .time-machine-container {
    padding: 20px 16px;
  }
  
  .view-btn {
    padding: 10px 16px;
    font-size: 13px;
  }
  
  .view-btn span {
    display: none;
  }
  
  .stat-card {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }
  
  .empty-state {
    padding: 60px 20px;
  }
}
</style>