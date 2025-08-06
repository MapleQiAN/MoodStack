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
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 32px;
  margin-bottom: 32px;
  backdrop-filter: blur(30px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.back-btn {
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
  align-self: flex-start;
  backdrop-filter: blur(20px);
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-3px) scale(1.05);
  box-shadow: 0 10px 30px rgba(255, 255, 255, 0.2);
}

.back-icon {
  font-size: 16px;
  font-weight: 600;
}

.diary-info {
  flex: 1;
}

.diary-title {
  font-size: 36px;
  font-weight: 800;
  color: rgba(255, 255, 255, 0.95);
  margin-bottom: 20px;
  line-height: 1.2;
  letter-spacing: -0.5px;
}

.diary-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 15px;
  font-weight: 500;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(20px);
}

.meta-icon {
  font-size: 16px;
}

.viewer-content {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(30px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 12px 50px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.content-wrapper {
  padding: 40px;
  min-height: 600px;
}

:deep(.markdown-content) {
  color: rgba(255, 255, 255, 0.9);
  line-height: 1.8;
  font-size: 17px;
  max-width: 100%;
  font-weight: 500;
}

:deep(.markdown-content h1) {
  font-size: 32px;
  font-weight: 800;
  color: rgba(255, 255, 255, 0.95);
  margin: 40px 0 20px 0;
  padding-bottom: 16px;
  border-bottom: 2px solid rgba(255, 255, 255, 0.3);
  letter-spacing: -0.3px;
}

:deep(.markdown-content h2) {
  font-size: 26px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  margin: 32px 0 16px 0;
  letter-spacing: -0.2px;
}

:deep(.markdown-content h3) {
  font-size: 22px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  margin: 28px 0 12px 0;
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
  border-left: 4px solid rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.1);
  padding: 20px 24px;
  margin: 28px 0;
  border-radius: 0 16px 16px 0;
  font-style: italic;
  color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
}

:deep(.markdown-content code) {
  background: rgba(255, 255, 255, 0.15);
  color: rgba(255, 255, 255, 0.95);
  padding: 6px 10px;
  border-radius: 8px;
  font-family: 'JetBrains Mono', 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 15px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  font-weight: 600;
}

:deep(.markdown-content pre) {
  background: rgba(0, 0, 0, 0.3);
  color: rgba(255, 255, 255, 0.95);
  padding: 24px;
  border-radius: 16px;
  overflow-x: auto;
  margin: 32px 0;
  border: 1px solid rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(20px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
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
  margin: 32px 0;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(20px);
}

:deep(.markdown-content th) {
  background: rgba(255, 255, 255, 0.15);
  color: rgba(255, 255, 255, 0.95);
  padding: 20px;
  text-align: left;
  font-weight: 700;
  backdrop-filter: blur(20px);
}

:deep(.markdown-content td) {
  padding: 18px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.9);
}

:deep(.markdown-content tr:nth-child(even)) {
  background: rgba(255, 255, 255, 0.05);
}

:deep(.markdown-content hr) {
  border: none;
  height: 2px;
  background: linear-gradient(to right, transparent, rgba(255, 255, 255, 0.4), transparent);
  margin: 40px 0;
  border-radius: 1px;
}

:deep(.markdown-content a) {
  color: rgba(255, 255, 255, 0.9);
  text-decoration: none;
  border-bottom: 1px solid rgba(255, 255, 255, 0.4);
  transition: all 0.3s ease;
  font-weight: 600;
}

:deep(.markdown-content a:hover) {
  color: rgba(255, 255, 255, 1);
  border-bottom-color: rgba(255, 255, 255, 0.8);
  text-shadow: 0 0 8px rgba(255, 255, 255, 0.3);
}

:deep(.markdown-content img) {
  max-width: 100%;
  height: auto;
  border-radius: 16px;
  margin: 24px 0;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
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