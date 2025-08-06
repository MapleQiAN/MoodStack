<template>
  <div class="viewer-container">
    <div class="viewer-header">
      <div class="header-content">
        <button @click="$emit('back')" class="back-btn">
          <span class="back-icon">←</span>
          <span class="back-text">返回列表</span>
        </button>
        <div class="diary-info">
          <h1 class="diary-title">{{ diary.title }}</h1>
          <div class="diary-meta">
            <span class="meta-item">
              <span class="meta-icon">📄</span>
              {{ diary.fileName }}
            </span>
            <span class="meta-item">
              <span class="meta-icon">📅</span>
              {{ formatDate(diary.createdAt) }}
            </span>
            <span class="meta-item">
              <span class="meta-icon">{{ getFileTypeIcon(diary.fileType) }}</span>
              {{ diary.fileType.toUpperCase() }}
            </span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="viewer-content">
      <div class="content-wrapper">
        <div class="markdown-content" v-html="renderedContent"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { marked } from 'marked'

const props = defineProps({
  diary: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['back'])

// Configure marked options for better rendering
marked.setOptions({
  breaks: true,
  gfm: true,
  headerIds: false,
  mangle: false
})

const renderedContent = computed(() => {
  if (!props.diary?.content) return ''
  return marked.parse(props.diary.content)
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
</script>

<style scoped>
.viewer-container {
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
}

.viewer-header {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.header-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.back-btn {
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
  align-self: flex-start;
}

.back-btn:hover {
  background: rgba(102, 126, 234, 0.15);
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.2);
}

.back-icon {
  font-size: 16px;
  font-weight: 600;
}

.diary-info {
  flex: 1;
}

.diary-title {
  font-size: 32px;
  font-weight: 700;
  color: #2d3748;
  margin-bottom: 16px;
  line-height: 1.2;
}

.diary-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  color: #718096;
  font-size: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.meta-icon {
  font-size: 16px;
}

.viewer-content {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.content-wrapper {
  padding: 32px;
  min-height: 500px;
}

:deep(.markdown-content) {
  color: #2d3748;
  line-height: 1.8;
  font-size: 16px;
  max-width: 100%;
}

:deep(.markdown-content h1) {
  font-size: 28px;
  font-weight: 700;
  color: #2d3748;
  margin: 32px 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 2px solid rgba(102, 126, 234, 0.2);
}

:deep(.markdown-content h2) {
  font-size: 24px;
  font-weight: 600;
  color: #4a5568;
  margin: 28px 0 12px 0;
}

:deep(.markdown-content h3) {
  font-size: 20px;
  font-weight: 600;
  color: #4a5568;
  margin: 24px 0 10px 0;
}

:deep(.markdown-content p) {
  margin-bottom: 16px;
  text-align: justify;
}

:deep(.markdown-content ul, .markdown-content ol) {
  margin: 16px 0;
  padding-left: 24px;
}

:deep(.markdown-content li) {
  margin-bottom: 8px;
}

:deep(.markdown-content blockquote) {
  border-left: 4px solid #667eea;
  background: rgba(102, 126, 234, 0.05);
  padding: 16px 20px;
  margin: 24px 0;
  border-radius: 0 12px 12px 0;
  font-style: italic;
  color: #4a5568;
}

:deep(.markdown-content code) {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  padding: 4px 8px;
  border-radius: 6px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

:deep(.markdown-content pre) {
  background: #2d3748;
  color: #f7fafc;
  padding: 20px;
  border-radius: 12px;
  overflow-x: auto;
  margin: 24px 0;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

:deep(.markdown-content pre code) {
  background: none;
  color: inherit;
  padding: 0;
  border-radius: 0;
  border: none;
}

:deep(.markdown-content table) {
  width: 100%;
  border-collapse: collapse;
  margin: 24px 0;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

:deep(.markdown-content th) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 16px;
  text-align: left;
  font-weight: 600;
}

:deep(.markdown-content td) {
  padding: 16px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

:deep(.markdown-content tr:nth-child(even)) {
  background: rgba(102, 126, 234, 0.05);
}

:deep(.markdown-content hr) {
  border: none;
  height: 2px;
  background: linear-gradient(to right, transparent, rgba(102, 126, 234, 0.3), transparent);
  margin: 32px 0;
}

:deep(.markdown-content a) {
  color: #667eea;
  text-decoration: none;
  border-bottom: 1px solid rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

:deep(.markdown-content a:hover) {
  color: #5a67d8;
  border-bottom-color: #5a67d8;
}

:deep(.markdown-content img) {
  max-width: 100%;
  height: auto;
  border-radius: 12px;
  margin: 20px 0;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .viewer-container {
    max-width: 100%;
  }
  
  .viewer-header, .content-wrapper {
    padding: 20px;
  }
  
  .diary-title {
    font-size: 24px;
  }
  
  .diary-meta {
    flex-direction: column;
    gap: 8px;
  }
  
  .meta-item {
    align-self: flex-start;
  }
  
  :deep(.markdown-content) {
    font-size: 15px;
  }
  
  :deep(.markdown-content h1) {
    font-size: 24px;
  }
  
  :deep(.markdown-content h2) {
    font-size: 20px;
  }
  
  :deep(.markdown-content h3) {
    font-size: 18px;
  }
}
</style>