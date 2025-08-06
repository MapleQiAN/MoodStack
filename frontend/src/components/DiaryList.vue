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
  margin-bottom: 40px;
  padding: 32px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(30px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.1);
}

.header-content {
  flex: 1;
}

.list-title {
  font-size: 32px;
  font-weight: 800;
  color: rgba(255, 255, 255, 0.95);
  margin-bottom: 12px;
  letter-spacing: -0.5px;
}

.list-subtitle {
  color: rgba(255, 255, 255, 0.7);
  font-size: 17px;
  font-weight: 600;
  letter-spacing: 0.2px;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 24px;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 600;
  font-size: 15px;
  backdrop-filter: blur(20px);
}

.refresh-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-3px) scale(1.05);
  box-shadow: 0 10px 30px rgba(255, 255, 255, 0.2);
}

.refresh-icon {
  font-size: 16px;
}

.empty-state {
  text-align: center;
  color: rgba(255, 255, 255, 0.9);
  margin-top: 100px;
  padding: 80px 50px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(30px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 24px;
  opacity: 0.7;
}

.empty-state h3 {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 16px;
  color: rgba(255, 255, 255, 0.95);
  letter-spacing: -0.3px;
}

.empty-state p {
  font-size: 17px;
  color: rgba(255, 255, 255, 0.7);
  margin-bottom: 40px;
  font-weight: 500;
}

.empty-action-btn {
  padding: 16px 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 700;
  font-size: 16px;
  letter-spacing: 0.5px;
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.empty-action-btn:hover {
  background: linear-gradient(135deg, #5a67d8 0%, #6b46a1 100%);
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 15px 40px rgba(102, 126, 234, 0.4);
}

.diary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(420px, 1fr));
  gap: 32px;
}

.diary-card {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 32px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(30px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
}

.diary-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea, #764ba2, #667eea);
  opacity: 0;
  transition: all 0.4s ease;
  border-radius: 20px 20px 0 0;
}

.diary-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.15);
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
  font-size: 20px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.95);
  line-height: 1.4;
  margin-bottom: 12px;
  letter-spacing: -0.2px;
}

.diary-meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
}

.diary-filename, .diary-date {
  display: flex;
  align-items: center;
  gap: 4px;
}

.diary-type-badge {
  font-size: 24px;
  opacity: 0.9;
  padding: 12px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  min-width: 48px;
  text-align: center;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.diary-preview {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.7;
  margin-bottom: 24px;
  min-height: 54px;
  font-size: 15px;
  font-weight: 500;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.view-hint {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.arrow-icon {
  color: rgba(255, 255, 255, 0.8);
  font-size: 18px;
  font-weight: 700;
  transition: transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.diary-card:hover .arrow-icon {
  transform: translateX(8px) scale(1.2);
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