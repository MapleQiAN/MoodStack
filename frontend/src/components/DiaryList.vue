<template>
  <div class="list-container">
    <div class="list-header">
      <div class="header-content">
        <h2 class="list-title">📚 我的日记</h2>
        <p class="list-subtitle">共 {{ diaries.length }} 篇日记</p>
      </div>
      <button @click="$emit('refresh')" class="refresh-btn">
        <span class="refresh-icon">🔄</span>
        <span class="refresh-text">刷新</span>
      </button>
    </div>
    
    <div v-if="diaries.length === 0" class="empty-state">
      <div class="empty-icon">📝</div>
      <h3>还没有日记</h3>
      <p>开始上传你的第一篇日记吧！</p>
      <button @click="$emit('refresh')" class="empty-action-btn">
        刷新列表
      </button>
    </div>
    
    <div v-else class="diary-grid">
      <div 
        v-for="diary in sortedDiaries" 
        :key="diary.id"
        class="diary-card"
        @click="$emit('view-diary', diary)"
      >
        <div class="card-header">
          <div class="diary-info">
            <h3 class="diary-title">{{ diary.title }}</h3>
            <div class="diary-meta">
              <span class="diary-filename">📄 {{ diary.fileName }}</span>
              <span class="diary-date">📅 {{ formatDate(diary.createdAt) }}</span>
            </div>
          </div>
          <div class="diary-type-badge">
            {{ getFileTypeIcon(diary.fileType) }}
          </div>
        </div>
        
        <div class="diary-preview">
          {{ getPreview(diary.content) }}
        </div>
        
        <div class="card-footer">
          <span class="view-hint">点击查看完整内容</span>
          <span class="arrow-icon">→</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  diaries: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['view-diary', 'refresh'])

const sortedDiaries = computed(() => {
  return [...props.diaries].sort((a, b) => 
    new Date(b.createdAt) - new Date(a.createdAt)
  )
})

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getFileTypeIcon = (fileType) => {
  const icons = {
    '.txt': '📄',
    '.md': '📝',
    '.docx': '📋',
    '.pdf': '📕'
  }
  return icons[fileType] || '📄'
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
  
  return preview.length > 120 ? preview.substring(0, 120) + '...' : preview
}
</script>

<style scoped>
.list-container {
  width: 100%;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
  padding: 24px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.header-content {
  flex: 1;
}

.list-title {
  font-size: 28px;
  font-weight: 700;
  color: #2d3748;
  margin-bottom: 8px;
}

.list-subtitle {
  color: #718096;
  font-size: 16px;
  font-weight: 500;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
  font-size: 14px;
}

.refresh-btn:hover {
  background: rgba(102, 126, 234, 0.15);
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.2);
}

.refresh-icon {
  font-size: 16px;
}

.empty-state {
  text-align: center;
  color: #4a5568;
  margin-top: 80px;
  padding: 60px 40px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 24px;
  opacity: 0.7;
}

.empty-state h3 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #2d3748;
}

.empty-state p {
  font-size: 16px;
  color: #718096;
  margin-bottom: 32px;
}

.empty-action-btn {
  padding: 12px 24px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
}

.empty-action-btn:hover {
  background: rgba(102, 126, 234, 0.15);
  transform: translateY(-2px);
}

.diary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 24px;
}

.diary-card {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
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
  background: linear-gradient(90deg, #667eea, #764ba2);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.diary-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
  border-color: rgba(102, 126, 234, 0.2);
}

.diary-card:hover::before {
  opacity: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.diary-info {
  flex: 1;
  margin-right: 16px;
}

.diary-title {
  font-size: 18px;
  font-weight: 600;
  color: #2d3748;
  line-height: 1.4;
  margin-bottom: 8px;
}

.diary-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;
  color: #718096;
}

.diary-filename, .diary-date {
  display: flex;
  align-items: center;
  gap: 4px;
}

.diary-type-badge {
  font-size: 20px;
  opacity: 0.8;
  padding: 8px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 8px;
  min-width: 40px;
  text-align: center;
}

.diary-preview {
  color: #4a5568;
  line-height: 1.6;
  margin-bottom: 16px;
  min-height: 48px;
  font-size: 14px;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.08);
}

.view-hint {
  color: #667eea;
  font-size: 13px;
  font-weight: 500;
}

.arrow-icon {
  color: #667eea;
  font-size: 16px;
  font-weight: 600;
  transition: transform 0.3s ease;
}

.diary-card:hover .arrow-icon {
  transform: translateX(4px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .list-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .refresh-btn {
    align-self: flex-end;
  }
  
  .diary-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .diary-card {
    padding: 20px;
  }
  
  .list-title {
    font-size: 24px;
  }
}
</style>