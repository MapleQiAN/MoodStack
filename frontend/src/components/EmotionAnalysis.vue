<script setup>
import { ref, onMounted, computed, nextTick } from 'vue'
import { 
  GetUserEmotionTrends, 
  GetUserEmotionStatistics, 
  AnalyzeAllDiariesEmotion
} from '../../wailsjs/go/main/App'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  RadialLinearScale
} from 'chart.js'
import { Line, Doughnut, PolarArea } from 'vue-chartjs'

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  RadialLinearScale
)

const loading = ref(false)
const analyzing = ref(false)
const emotionTrends = ref([])
const emotionStats = ref(null)
const selectedTimeRange = ref(30) // days
const analysisMethod = ref('programmatic') // 'programmatic' or 'ai'
const ollamaURL = ref('http://localhost:11434')
const showSettings = ref(false)
const forceReAnalysis = ref(false)
const rootEl = ref(null)

// Custom notification system
const notification = ref({
  show: false,
  title: '',
  message: '',
  type: 'info', // 'success', 'error', 'warning', 'info'
  details: []
})

const resolveCssVar = (name, fallback) => {
  try {
    const target = rootEl.value || document.documentElement
    const val = getComputedStyle(target).getPropertyValue(name)
    return val && val.trim().length > 0 ? val.trim() : fallback
  } catch {
    return fallback
  }
}

// Emotion colors mapping - modern cohesive palette
const emotionColors = {
  joy: '#F4D35E',       // warm yellow
  sadness: '#4D7EA8',   // calm blue
  anger: '#E76F51',     // soft red
  fear: '#6C5CE7',      // indigo
  love: '#EF476F',      // pink red
  surprise: '#06D6A0',  // teal
  disgust: '#8D99AE'    // muted slate
}

// Emotion labels in Chinese
const emotionLabels = {
  joy: '快乐',
  sadness: '悲伤',
  anger: '愤怒',
  fear: '恐惧',
  love: '爱意',
  surprise: '惊讶',
  disgust: '厌恶'
}

const timeRangeOptions = [
  { value: 7, label: '最近7天' },
  { value: 30, label: '最近30天' },
  { value: 90, label: '最近90天' },
  { value: 365, label: '最近一年' },
  { value: 0, label: '全部时间' }
]

// Computed properties for charts
const emotionTrendChartData = computed(() => {
  if (!emotionTrends.value || emotionTrends.value.length === 0) {
    return {
      labels: [],
      datasets: []
    }
  }

  // Group data by date and calculate daily averages
  const dailyData = {}
  
  emotionTrends.value.forEach(analysis => {
    const date = new Date(analysis.createdAt).toLocaleDateString('zh-CN')
    if (!dailyData[date]) {
      dailyData[date] = {
        joy: [],
        sadness: [],
        anger: [],
        fear: [],
        love: [],
        surprise: [],
        disgust: []
      }
    }
    
    dailyData[date].joy.push(analysis.joy)
    dailyData[date].sadness.push(analysis.sadness)
    dailyData[date].anger.push(analysis.anger)
    dailyData[date].fear.push(analysis.fear)
    dailyData[date].love.push(analysis.love)
    dailyData[date].surprise.push(analysis.surprise)
    dailyData[date].disgust.push(analysis.disgust)
  })

  const dates = Object.keys(dailyData).sort()
  
  const datasets = Object.keys(emotionColors).map(emotion => ({
    label: emotionLabels[emotion],
    data: dates.map(date => {
      const values = dailyData[date][emotion]
      return values.length > 0 ? values.reduce((a, b) => a + b, 0) / values.length : 0
    }),
    borderColor: emotionColors[emotion],
    backgroundColor: emotionColors[emotion] + '20',
    tension: 0.4,
    fill: false
  }))

  return {
    labels: dates,
    datasets
  }
})

const emotionDistributionChartData = computed(() => {
  if (!emotionStats.value || !emotionStats.value.emotionDistribution) {
    return {
      labels: [],
      datasets: []
    }
  }

  const distribution = emotionStats.value.emotionDistribution
  const labels = Object.keys(distribution).map(emotion => emotionLabels[emotion] || emotion)
  const data = Object.values(distribution)
  const backgroundColor = Object.keys(distribution).map(emotion => emotionColors[emotion] || '#666')

  return {
    labels,
    datasets: [{
      data,
      backgroundColor,
      borderWidth: 2,
      borderColor: '#fff'
    }]
  }
})

const averageEmotionsChartData = computed(() => {
  if (!emotionStats.value || !emotionStats.value.averageEmotions) {
    return {
      labels: [],
      datasets: []
    }
  }

  const averages = emotionStats.value.averageEmotions
  const labels = Object.keys(averages).map(emotion => emotionLabels[emotion])
  const data = Object.values(averages)
  const backgroundColor = Object.keys(averages).map(emotion => emotionColors[emotion])

  return {
    labels,
    datasets: [{
      label: '平均情感强度',
      data,
      backgroundColor,
      borderWidth: 2,
      borderColor: '#fff'
    }]
  }
})

const sentimentDistributionChartData = computed(() => {
  if (!emotionStats.value || !emotionStats.value.sentimentDistribution) {
    return {
      labels: [],
      datasets: []
    }
  }

  const distribution = emotionStats.value.sentimentDistribution
  const labels = ['积极', '消极', '中性']
  const data = [
    distribution.positive || 0,
    distribution.negative || 0,
    distribution.neutral || 0
  ]
  const backgroundColor = ['#2DD4BF', '#F43F5E', '#94A3B8']

  return {
    labels,
    datasets: [{
      data,
      backgroundColor,
      borderWidth: 2,
      borderColor: '#fff'
    }]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top',
      labels: {
        usePointStyle: true,
        padding: 20,
        color: () => resolveCssVar('--text-primary', '#E5E7EB')
      }
    },
    tooltip: {
      backgroundColor: () => resolveCssVar('--bg-secondary', '#111827'),
      titleColor: () => resolveCssVar('--text-primary', '#E5E7EB'),
      bodyColor: () => resolveCssVar('--text-secondary', '#CBD5E1'),
      borderColor: () => resolveCssVar('--border-color', '#1F2937'),
      borderWidth: 1,
      callbacks: {
        labelColor(ctx){
          return { borderColor: ctx.dataset.borderColor, backgroundColor: ctx.dataset.borderColor };
        }
      }
    }
  },
  scales: {
    x: {
      grid: {
        color: () => resolveCssVar('--bg-tertiary', '#111827')
      },
      ticks: {
        color: () => resolveCssVar('--text-secondary', '#94A3B8')
      }
    },
    y: {
      grid: {
        color: () => resolveCssVar('--bg-tertiary', '#111827')
      },
      ticks: {
        color: () => resolveCssVar('--text-secondary', '#94A3B8')
      },
      beginAtZero: true,
      max: 1
    }
  }
}

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'right',
      labels: {
        usePointStyle: true,
        padding: 20,
        color: () => resolveCssVar('--text-primary', '#E5E7EB')
      }
    },
    tooltip: {
      backgroundColor: () => resolveCssVar('--bg-secondary', '#111827'),
      titleColor: () => resolveCssVar('--text-primary', '#E5E7EB'),
      bodyColor: () => resolveCssVar('--text-secondary', '#CBD5E1'),
      borderColor: () => resolveCssVar('--border-color', '#1F2937'),
      borderWidth: 1,
      callbacks: {
        labelColor(ctx){
          return { borderColor: ctx.dataset.borderColor || '#000', backgroundColor: ctx.dataset.backgroundColor?.[ctx.dataIndex] || '#000' };
        }
      }
    }
  }
}

const polarAreaOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'right',
      labels: {
        usePointStyle: true,
        padding: 20,
        color: () => resolveCssVar('--text-primary', '#E5E7EB')
      }
    },
    tooltip: {
      backgroundColor: () => resolveCssVar('--bg-secondary', '#111827'),
      titleColor: () => resolveCssVar('--text-primary', '#E5E7EB'),
      bodyColor: () => resolveCssVar('--text-secondary', '#CBD5E1'),
      borderColor: () => resolveCssVar('--border-color', '#1F2937'),
      borderWidth: 1,
      callbacks: {
        labelColor(ctx){
          return { borderColor: ctx.dataset.borderColor || '#000', backgroundColor: ctx.dataset.backgroundColor?.[ctx.dataIndex] || '#000' };
        }
      }
    }
  },
  scales: {
    r: {
      grid: {
        color: () => resolveCssVar('--bg-tertiary', '#111827')
      },
      ticks: {
        color: () => resolveCssVar('--text-secondary', '#94A3B8'),
        backdropColor: 'transparent'
      },
      beginAtZero: true,
      max: 1
    }
  }
}

const wordCloudData = computed(() => {
  if (!emotionStats.value || !emotionStats.value.topKeywords) {
    return []
  }
  
  // Convert keywords to word cloud format
  return emotionStats.value.topKeywords.map((keyword, index) => ({
    text: keyword,
    size: Math.max(10, 40 - index * 2) // Decreasing size
  }))
})

const loadEmotionData = async () => {
  try {
    loading.value = true
    
    const [trends, stats] = await Promise.all([
      GetUserEmotionTrends(selectedTimeRange.value),
      GetUserEmotionStatistics(selectedTimeRange.value)
    ])
    
    emotionTrends.value = trends || []
    emotionStats.value = stats || {}
    
    // Render word cloud after data is loaded
    await nextTick()
    renderWordCloud()
  } catch (error) {
    console.error('加载情感数据失败:', error)
  } finally {
    loading.value = false
  }
}

const analyzeAllDiaries = async () => {
  try {
    analyzing.value = true
    
    const useAI = analysisMethod.value === 'ai'
    const result = await AnalyzeAllDiariesEmotion(useAI, ollamaURL.value)
    
    console.log('分析结果:', result)
    
    // Show detailed result message
    if (!result) {
      showNotification('分析失败', '分析结果为空，请重试', 'error')
      return
    }
    
    if (result.total === 0) {
      showNotification('无日记', result.message || '没有找到日记进行分析', 'warning')
      return
    }
    
    // Prepare summary message
    let summary = `总计 ${result.total} 篇日记，成功 ${result.success} 篇`
    
    if (result.skipped > 0) {
      summary += `，跳过 ${result.skipped} 篇`
    }
    
    if (result.failures > 0) {
      summary += `，失败 ${result.failures} 篇`
    }
    
    // Prepare details
    let details = []
    
    if (result.results && result.results.length > 0) {
      details = result.results.map(item => {
        let detail = `${item.title}: ${item.status}`
        if (item.emotion) {
          detail += ` (${item.emotion})`
        }
        if (item.error) {
          detail += ` - ${item.error}`
        }
        return detail
      })
    }
    
    if (result.lastError) {
      details.push(`最后错误: ${result.lastError}`)
    }
    
    // Determine notification type
    const type = result.failures > 0 ? 'warning' : 'success'
    
    showNotification('分析完成', summary, type, details)
    
    // Reload data after analysis
    await loadEmotionData()
  } catch (error) {
    console.error('分析失败:', error)
    showNotification('分析失败', '情感分析过程中发生错误', 'error', [error.message])
  } finally {
    analyzing.value = false
  }
}

const onTimeRangeChange = () => {
  loadEmotionData()
}

// Notification functions
const showNotification = (title, message, type = 'info', details = []) => {
  notification.value = {
    show: true,
    title,
    message,
    type,
    details
  }
  
  // Auto close after 5 seconds (except for errors)
  if (type !== 'error') {
    setTimeout(() => {
      closeNotification()
    }, 5000)
  }
}

const closeNotification = () => {
  notification.value.show = false
}

const renderWordCloud = () => {
  if (!wordCloudData.value || wordCloudData.value.length === 0) {
    return
  }
  
  const canvas = document.getElementById('wordcloud-canvas')
  if (!canvas) return
  
  // Clear canvas
  const ctx = canvas.getContext('2d')
  ctx.clearRect(0, 0, canvas.width, canvas.height)
  
  // Simple word cloud implementation (fallback if wordcloud library not available)
  if (typeof WordCloud === 'undefined') {
    // Fallback: simple text rendering
    const words = wordCloudData.value.slice(0, 20)
    const centerX = canvas.width / 2
    const centerY = canvas.height / 2
    
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    
    words.forEach((word, index) => {
      const angle = (index / words.length) * 2 * Math.PI
      const radius = 30 + (index % 3) * 40
      const x = centerX + Math.cos(angle) * radius
      const y = centerY + Math.sin(angle) * radius
      
      ctx.font = `${word.size}px Arial`
      ctx.fillStyle = Object.values(emotionColors)[index % Object.keys(emotionColors).length]
      ctx.fillText(word.text, x, y)
    })
    return
  }
  
  // Use wordcloud library if available
  try {
    WordCloud(canvas, {
      list: wordCloudData.value.map(item => [item.text, item.size]),
      gridSize: 8,
      weightFactor: 1,
      fontFamily: 'Arial, sans-serif',
      color: function() {
        const colors = Object.values(emotionColors)
        return colors[Math.floor(Math.random() * colors.length)]
      },
      rotateRatio: 0.3,
      backgroundColor: 'transparent'
    })
  } catch (error) {
    console.warn('WordCloud rendering failed, using fallback', error)
    // Use fallback method above
  }
}

onMounted(() => {
  // 将根元素引用到当前容器，以便读取局部 CSS 变量
  rootEl.value = document.querySelector('.emotion-analysis')
  loadEmotionData()
})
</script>

<template>
  <div class="emotion-analysis">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <span class="title-icon" aria-hidden="true">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 21h18"/>
              <path d="M7 16V9"/>
              <path d="M12 16V5"/>
              <path d="M17 16v-4"/>
            </svg>
          </span>
          情感成长分析
        </h1>
        <p class="page-subtitle">通过AI和数据分析了解你的情感变化</p>
      </div>
      
      <div class="header-actions">
        <button 
          @click="showSettings = !showSettings"
          class="btn btn-settings"
          :class="{ 'active': showSettings }"
          title="分析设置"
        >
          <svg class="icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/>
            <circle cx="12" cy="12" r="3"/>
          </svg>
          设置
        </button>
        
        <button 
          @click="analyzeAllDiaries" 
          :disabled="analyzing"
          class="btn btn-analyze"
        >
          <svg v-if="!analyzing" class="icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 3v18h18"/>
            <path d="M18.7 8l-5.1 5.2-2.8-2.7L7 14.3"/>
            <path d="M9 12l2 2 4-4"/>
          </svg>
          <div v-else class="spinner"></div>
          {{ analyzing ? '分析中...' : '开始分析' }}
        </button>
      </div>
    </div>

    <!-- Settings Panel -->
    <div v-if="showSettings" class="settings-panel">
      <div class="settings-content">
        <h3>分析设置</h3>
        
        <div class="setting-group">
          <label>时间范围</label>
          <select v-model="selectedTimeRange" @change="onTimeRangeChange" class="select-input">
            <option v-for="option in timeRangeOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
        
        <div class="setting-group">
          <label>分析方法</label>
          <div class="radio-group">
            <label class="radio-option">
              <input type="radio" v-model="analysisMethod" value="programmatic" />
              <span>程序化分析（快速）</span>
            </label>
            <label class="radio-option">
              <input type="radio" v-model="analysisMethod" value="ai" />
              <span>AI分析（精确）</span>
            </label>
          </div>
        </div>
        
        <div v-if="analysisMethod === 'ai'" class="setting-group">
          <label>Ollama服务地址</label>
          <input 
            type="text" 
            v-model="ollamaURL" 
            class="text-input"
            placeholder="http://localhost:11434"
          />
          <small class="help-text">确保Ollama服务正在运行且已安装中文模型</small>
        </div>
        
        <div class="setting-group">
          <label class="radio-option">
            <input type="checkbox" v-model="forceReAnalysis" />
            <span>强制重新分析（重新分析已分析过的日记）</span>
          </label>
          <small class="help-text">勾选此项将重新分析所有日记，包括已经分析过的</small>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>正在加载情感数据...</p>
    </div>

    <div v-else-if="!emotionStats || emotionStats.totalEntries === 0" class="empty-state">
      <div class="empty-icon" aria-hidden="true">
        <svg class="icon xl" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 21h18"/>
          <path d="M5 17l4-4 3 3 6-6 2 2"/>
        </svg>
      </div>
      <h3>暂无情感数据</h3>
      <p>点击"开始分析"按钮来分析你的日记情感</p>
    </div>

    <div v-else class="analysis-content">
      <!-- Statistics Overview -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon" aria-hidden="true">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <rect x="5" y="3" width="14" height="18" rx="2"/>
              <path d="M8 8h8"/>
              <path d="M8 12h8"/>
              <path d="M8 16h5"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ emotionStats.totalEntries }}</div>
            <div class="stat-label">已分析日记</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon" aria-hidden="true">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="9"/>
              <path d="M8.5 10h.01"/>
              <path d="M15.5 10h.01"/>
              <path d="M8.5 14c1.3 1 2.7 1 4 0s2.7-1 4 0"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">
              {{ Math.round((emotionStats.sentimentDistribution?.positive || 0) / emotionStats.totalEntries * 100) }}%
            </div>
            <div class="stat-label">积极情感</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon" aria-hidden="true">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 21s-6-4.5-8-7.5C2 10 3.8 6.5 7 6.5c2 0 3 1.3 5 2 2-.7 3-2 5-2 3.2 0 5 3.5 3 7-2 3-8 7.5-8 7.5z"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">
              {{ (emotionStats.averageEmotions?.joy || 0).toFixed(2) }}
            </div>
            <div class="stat-label">平均快乐度</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon" aria-hidden="true">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M7 7h4l7 7-4 4-7-7V7z"/>
              <circle cx="9" cy="9" r="1.3"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ emotionStats.topKeywords?.length || 0 }}</div>
            <div class="stat-label">情感关键词</div>
          </div>
        </div>
      </div>

      <!-- Charts Section -->
      <div class="charts-section">
        <!-- Emotion Trends -->
        <div class="chart-card">
          <div class="chart-header">
            <h3>情感趋势变化</h3>
            <p>展示不同情感随时间的变化趋势</p>
          </div>
          <div class="chart-container">
            <Line 
              v-if="emotionTrendChartData.labels.length > 0"
              :data="emotionTrendChartData" 
              :options="chartOptions" 
            />
            <div v-else class="chart-empty">
              <p>暂无趋势数据</p>
            </div>
          </div>
        </div>

        <!-- Emotion Distribution -->
        <div class="chart-card">
          <div class="chart-header">
            <h3>主导情感分布</h3>
            <p>显示不同情感出现的频率</p>
          </div>
          <div class="chart-container">
            <Doughnut 
              v-if="emotionDistributionChartData.labels.length > 0"
              :data="emotionDistributionChartData" 
              :options="doughnutOptions" 
            />
            <div v-else class="chart-empty">
              <p>暂无分布数据</p>
            </div>
          </div>
        </div>

        <!-- Average Emotions -->
        <div class="chart-card">
          <div class="chart-header">
            <h3>平均情感强度</h3>
            <p>各种情感的平均强度水平</p>
          </div>
          <div class="chart-container">
            <PolarArea 
              v-if="averageEmotionsChartData.labels.length > 0"
              :data="averageEmotionsChartData" 
              :options="polarAreaOptions" 
            />
            <div v-else class="chart-empty">
              <p>暂无强度数据</p>
            </div>
          </div>
        </div>

        <!-- Sentiment Distribution -->
        <div class="chart-card">
          <div class="chart-header">
            <h3>情感倾向分布</h3>
            <p>积极、消极、中性情感的比例</p>
          </div>
          <div class="chart-container">
            <Doughnut 
              v-if="sentimentDistributionChartData.labels.length > 0"
              :data="sentimentDistributionChartData" 
              :options="doughnutOptions" 
            />
            <div v-else class="chart-empty">
              <p>暂无倾向数据</p>
            </div>
          </div>
        </div>

        <!-- Word Cloud -->
        <div class="chart-card word-cloud-card">
          <div class="chart-header">
            <h3>情感关键词云</h3>
            <p>从日记中提取的高频情感词汇</p>
          </div>
          <div class="chart-container">
            <canvas 
              id="wordcloud-canvas" 
              width="600" 
              height="400"
              class="wordcloud-canvas"
            ></canvas>
          </div>
        </div>
      </div>
    </div>

    <!-- Toast Notification -->
    <div v-if="notification.show" class="toast-notification" :class="notification.type">
      <div class="toast-icon">
        <svg v-if="notification.type === 'success'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 12l2 2 4-4"/>
          <circle cx="12" cy="12" r="9"/>
        </svg>
        <svg v-else-if="notification.type === 'error'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="9"/>
          <path d="M15 9l-6 6"/>
          <path d="M9 9l6 6"/>
        </svg>
        <svg v-else-if="notification.type === 'warning'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
          <path d="M12 9v4"/>
          <path d="M12 17h.01"/>
        </svg>
        <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="9"/>
          <path d="M12 16v-4"/>
          <path d="M12 8h.01"/>
        </svg>
      </div>
      
      <div class="toast-content">
        <div class="toast-title">{{ notification.title }}</div>
        <div class="toast-message">{{ notification.message }}</div>
      </div>
      
      <button @click="closeNotification" class="toast-close">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6L6 18"/>
          <path d="M6 6l12 12"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.emotion-analysis {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  /* Local theming overrides */
  --accent-primary: #6366F1; /* indigo-500 */
  --accent-hover: #4F46E5;   /* indigo-600 */
  --border-color: color-mix(in srgb, var(--accent-primary) 10%, var(--bg-tertiary));
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
  gap: 24px;
}

.header-content {
  flex: 1;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  display: inline-flex;
  align-items: center;
  line-height: 0;
}

.icon {
  width: 1.25em;
  height: 1.25em;
  display: inline-block;
  stroke: currentColor;
  fill: none;
}

.page-title .icon {
  width: 28px;
  height: 28px;
  color: var(--text-primary);
}

.empty-icon .icon.xl {
  width: 64px;
  height: 64px;
}

.page-subtitle {
  font-size: 16px;
  color: var(--text-muted);
  margin: 0;
  line-height: 1.5;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
  position: relative;
  overflow: hidden;
}

.btn .icon {
  transition: transform 0.2s ease;
}

.btn:hover .icon {
  transform: scale(1.1);
}

.btn-analyze {
  background: linear-gradient(135deg, var(--nord8), var(--nord9));
  color: white;
  border: 1px solid var(--nord8);
  box-shadow: 0 4px 12px rgba(136, 192, 208, 0.3);
  position: relative;
  overflow: hidden;
}

.btn-analyze::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.btn-analyze:hover:not(:disabled)::before {
  left: 100%;
}

.btn-analyze:hover:not(:disabled) {
  background: linear-gradient(135deg, var(--nord9), var(--nord10));
  border-color: var(--nord9);
  transform: translateY(-1px);
  box-shadow: 0 8px 20px rgba(136, 192, 208, 0.4);
}

.btn-analyze:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-settings {
  background: var(--bg-glass);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 8px rgba(46, 52, 64, 0.08);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.btn-settings:hover {
  background: var(--bg-glass-hover);
  color: var(--text-primary);
  border-color: var(--nord8);
  box-shadow: 0 4px 12px rgba(136, 192, 208, 0.15);
  transform: translateY(-1px);
}

.btn-settings.active {
  background: var(--accent-primary-alpha);
  color: var(--nord8);
  border-color: var(--nord8);
  box-shadow: 0 4px 12px rgba(136, 192, 208, 0.2);
  animation: settings-pulse 2s ease-in-out infinite;
}

@keyframes settings-pulse {
  0%, 100% {
    box-shadow: 0 4px 12px rgba(136, 192, 208, 0.2);
  }
  50% {
    box-shadow: 0 4px 16px rgba(136, 192, 208, 0.3);
  }
}

.btn:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--nord8) 24%, transparent);
}

.btn-analyze:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--nord8) 24%, transparent), 0 8px 20px rgba(136, 192, 208, 0.4);
}

.btn-settings:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--nord8) 24%, transparent), 0 4px 12px rgba(136, 192, 208, 0.15);
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.btn-analyze .spinner {
  border-color: rgba(255, 255, 255, 0.2);
  border-top-color: white;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.settings-panel {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px;
  margin-bottom: 24px;
}

.settings-content h3 {
  margin: 0 0 20px 0;
  color: var(--text-primary);
  font-size: 18px;
  font-weight: 600;
}

.setting-group {
  margin-bottom: 20px;
}

.setting-group label {
  display: block;
  margin-bottom: 8px;
  color: var(--text-secondary);
  font-weight: 500;
  font-size: 14px;
}

.select-input,
.text-input {
  width: 100%;
  max-width: 300px;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  transition: box-shadow .2s ease, border-color .2s ease;
}

.select-input:focus,
.text-input:focus {
  outline: none;
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--accent-primary) 24%, transparent);
}

.radio-option input {
  accent-color: var(--accent-primary);
}

.radio-group {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.radio-option {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-secondary);
}

.radio-option input {
  margin: 0;
}

.help-text {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-muted);
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  text-align: center;
}

.loading-state .spinner {
  width: 32px;
  height: 32px;
  border-width: 3px;
  margin-bottom: 16px;
}

.empty-state .empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.6;
  color: var(--text-secondary);
}

.empty-state h3 {
  margin: 0 0 8px 0;
  color: var(--text-primary);
  font-size: 20px;
}

.empty-state p {
  margin: 0;
  color: var(--text-muted);
  font-size: 14px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: transform 0.2s ease, box-shadow .2s ease, border-color .2s ease;
  box-shadow: 0 6px 16px rgba(0,0,0,0.06);
}

.stat-card:hover {
  transform: translateY(-2px);
  border-color: color-mix(in srgb, var(--accent-primary) 24%, var(--border-color));
  box-shadow: 0 10px 24px rgba(0,0,0,0.12);
}

.stat-icon {
  font-size: 32px;
  flex-shrink: 0;
  color: var(--text-secondary);
}

.stat-icon .icon {
  color: inherit;
}

.header-actions .btn .icon {
  color: inherit;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.charts-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: 24px;
}

.chart-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px;
  min-height: 400px;
  box-shadow: 0 8px 20px rgba(0,0,0,0.06);
}

.word-cloud-card {
  grid-column: 1 / -1;
}

.chart-header {
  margin-bottom: 20px;
}

.chart-header h3 {
  margin: 0 0 4px 0;
  color: var(--text-primary);
  font-size: 18px;
  font-weight: 600;
}

.chart-header p {
  margin: 0;
  color: var(--text-muted);
  font-size: 14px;
}

.chart-container {
  height: 300px;
  position: relative;
}

.word-cloud-card .chart-container {
  height: 400px;
}

.chart-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-muted);
  font-size: 14px;
}

.wordcloud-canvas {
  width: 100%;
  height: 100%;
  border-radius: var(--radius-md);
}

/* Responsive Design */
@media (max-width: 768px) {
  .emotion-analysis {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-start;
  }
  
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 16px;
  }
  
  .charts-section {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .chart-card {
    min-height: 350px;
    padding: 20px;
  }
  
  .radio-group {
    flex-direction: column;
    gap: 12px;
  }
}

@media (max-width: 480px) {
  .page-title {
    font-size: 24px;
  }
  
  .title-icon {
    font-size: 28px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .stat-card {
    padding: 16px;
  }
  
  .chart-card {
    min-height: 300px;
    padding: 16px;
  }
  
  .chart-container {
    height: 250px;
  }
}

/* Toast Notification Styles - Nord Theme */
.toast-notification {
  position: fixed;
  top: 24px;
  right: 24px;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  padding: 16px;
  max-width: 400px;
  min-width: 320px;
  z-index: 1000;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  animation: toast-slide-in 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.toast-notification.success {
  border-left: 4px solid var(--success-primary);
}

.toast-notification.error {
  border-left: 4px solid var(--error-primary);
}

.toast-notification.warning {
  border-left: 4px solid var(--warning-primary);
}

.toast-notification.info {
  border-left: 4px solid var(--accent-info);
}

.toast-icon {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  margin-top: 2px;
}

.toast-notification.success .toast-icon {
  color: var(--success-primary);
}

.toast-notification.error .toast-icon {
  color: var(--error-primary);
}

.toast-notification.warning .toast-icon {
  color: var(--warning-primary);
}

.toast-notification.info .toast-icon {
  color: var(--accent-info);
}

.toast-content {
  flex: 1;
  min-width: 0;
}

.toast-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  font-family: var(--font-display);
  margin-bottom: 4px;
  line-height: 1.4;
}

.toast-message {
  font-size: 13px;
  color: var(--text-secondary);
  font-family: var(--font-body);
  line-height: 1.5;
  word-wrap: break-word;
}

.toast-close {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  margin-top: -2px;
}

.toast-close:hover {
  background: var(--surface-hover);
  color: var(--text-secondary);
}

@keyframes toast-slide-in {
  from {
    opacity: 0;
    transform: translateX(100%);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* Mobile responsiveness for toast */
@media (max-width: 768px) {
  .toast-notification {
    top: 16px;
    right: 16px;
    left: 16px;
    max-width: none;
    min-width: auto;
  }
}

@media (max-width: 480px) {
  .toast-notification {
    top: 12px;
    right: 12px;
    left: 12px;
    padding: 14px;
  }
  
  .toast-title {
    font-size: 13px;
  }
  
  .toast-message {
    font-size: 12px;
  }
}
</style>

