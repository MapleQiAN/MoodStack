<template>
  <div class="viewer-container">
    <div class="viewer-header">
      <button @click="$emit('back')" class="back-button">
        <span class="back-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15,18 9,12 15,6"/>
          </svg>
        </span>
        <span>返回列表</span>
      </button>
      
      <div class="diary-meta">
        <div class="meta-group">
          <div class="meta-item">
            <span class="meta-icon">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
              </svg>
            </span>
            <span class="meta-text">{{ diary.fileName }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-icon">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                <line x1="16" y1="2" x2="16" y2="6"/>
                <line x1="8" y1="2" x2="8" y2="6"/>
                <line x1="3" y1="10" x2="21" y2="10"/>
              </svg>
            </span>
            <span class="meta-text">{{ formatDate(diary.createdAt) }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-icon">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
              </svg>
            </span>
            <span class="meta-text">{{ diary.fileType.toUpperCase() }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="viewer-content">
      <article class="article-container">
        <header class="article-header">
          <h1 class="article-title">{{ diary.title }}</h1>
        </header>
        
        <div class="article-body">
          <div class="markdown-content" v-html="renderedContent"></div>
        </div>
      </article>
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

// 移除emoji图标函数，改用SVG图标
</script>

<style scoped>
.viewer-container {
  max-width: 100%;
  padding: 32px;
  box-sizing: border-box;
}

.viewer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding: 24px;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-xl);
  border: 1px solid var(--nord4);
  box-shadow: var(--shadow-sm);
}

.back-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--nord5);
  color: var(--text-secondary);
  border: 1px solid var(--nord4);
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.2s ease;
}

.back-button:hover {
  background: var(--nord4);
  color: var(--text-primary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.back-icon {
  font-size: 16px;
  font-weight: 600;
}

.diary-meta {
  flex: 1;
  display: flex;
  justify-content: flex-end;
}

.meta-group {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--nord5);
  border: 1px solid var(--nord4);
  border-radius: var(--radius-md);
  font-size: 13px;
}

.meta-icon {
  color: var(--accent-primary);
}

.meta-text {
  color: var(--text-muted);
  font-weight: 500;
}

.viewer-content {
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-xl);
  border: 1px solid var(--nord4);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.article-container {
  padding: 0;
}

.article-header {
  padding: 32px 32px 24px 32px;
  border-bottom: 1px solid var(--nord4);
  background: var(--nord5);
}

.article-title {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.3;
  letter-spacing: -0.025em;
  margin: 0;
}

.article-body {
  padding: 32px;
}

/* Markdown content styling */
:deep(.markdown-content) {
  color: var(--text-secondary);
  line-height: 1.7;
  font-size: 16px;
  font-weight: 400;
  max-width: none;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

:deep(.markdown-content h1) {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 32px 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--nord4);
  letter-spacing: -0.025em;
  font-family: var(--font-display);
}

:deep(.markdown-content h2) {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 24px 0 12px 0;
  letter-spacing: -0.02em;
  font-family: var(--font-display);
}

:deep(.markdown-content h3) {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 20px 0 10px 0;
  font-family: var(--font-display);
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
  margin-bottom: 6px;
}

:deep(.markdown-content blockquote) {
  border-left: 4px solid var(--accent-primary);
  background: var(--nord5);
  padding: 16px 20px;
  margin: 20px 0;
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
  font-style: italic;
  color: var(--text-muted);
}

:deep(.markdown-content code) {
  background: var(--nord5);
  color: var(--text-primary);
  padding: 3px 6px;
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: 14px;
  border: 1px solid var(--nord4);
  font-weight: 500;
}

:deep(.markdown-content pre) {
  background: var(--nord0);
  color: var(--nord6);
  padding: 20px;
  border-radius: var(--radius-md);
  overflow-x: auto;
  margin: 20px 0;
  border: 1px solid var(--nord3);
  box-shadow: var(--shadow-sm);
}

:deep(.markdown-content pre code) {
  background: none;
  color: inherit;
  padding: 0;
  border-radius: 0;
  border: none;
  font-weight: 400;
}

:deep(.markdown-content table) {
  width: 100%;
  border-collapse: collapse;
  margin: 24px 0;
  box-shadow: var(--shadow-sm);
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1px solid var(--nord4);
}

:deep(.markdown-content th) {
  background: var(--nord5);
  color: var(--text-primary);
  padding: 16px;
  text-align: left;
  font-weight: 600;
  border-bottom: 1px solid var(--nord4);
}

:deep(.markdown-content td) {
  padding: 14px 16px;
  border-bottom: 1px solid var(--nord4);
  color: var(--text-secondary);
}

:deep(.markdown-content tr:nth-child(even)) {
  background: var(--nord6);
}

:deep(.markdown-content tr:last-child td) {
  border-bottom: none;
}

:deep(.markdown-content hr) {
  border: none;
  height: 1px;
  background: linear-gradient(to right, transparent, var(--nord3), transparent);
  margin: 32px 0;
}

:deep(.markdown-content a) {
  color: var(--accent-primary);
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: all 0.2s ease;
  font-weight: 500;
}

:deep(.markdown-content a:hover) {
  border-bottom-color: var(--accent-primary);
}

:deep(.markdown-content img) {
  max-width: 100%;
  height: auto;
  border-radius: var(--radius-md);
  margin: 20px 0;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--nord4);
}

:deep(.markdown-content strong) {
  font-weight: 600;
  color: var(--text-primary);
}

:deep(.markdown-content em) {
  font-style: italic;
  color: var(--text-muted);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .viewer-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
    padding: 20px;
  }
  
  .diary-meta {
    justify-content: flex-start;
  }
  
  .meta-group {
    flex-direction: column;
    gap: 8px;
  }
  
  .meta-item {
    align-self: flex-start;
  }
  
  .article-header {
    padding: 24px 20px 20px 20px;
  }
  
  .article-title {
    font-size: 24px;
  }
  
  .article-body {
    padding: 24px 20px;
  }
  
  :deep(.markdown-content) {
    font-size: 15px;
  }
  
  :deep(.markdown-content h1) {
    font-size: 22px;
    margin: 24px 0 12px 0;
  }
  
  :deep(.markdown-content h2) {
    font-size: 18px;
  }
  
  :deep(.markdown-content h3) {
    font-size: 16px;
  }
}

@media (max-width: 480px) {
  .viewer-header {
    padding: 16px;
  }
  
  .article-header {
    padding: 20px 16px 16px 16px;
  }
  
  .article-title {
    font-size: 22px;
    line-height: 1.4;
  }
  
  .article-body {
    padding: 20px 16px;
  }
  
  :deep(.markdown-content) {
    font-size: 14px;
  }
  
  :deep(.markdown-content h1) {
    font-size: 20px;
    margin: 20px 0 10px 0;
  }
  
  :deep(.markdown-content h2) {
    font-size: 17px;
  }
  
  :deep(.markdown-content h3) {
    font-size: 15px;
  }
  
  :deep(.markdown-content pre) {
    padding: 16px;
  }
  
  :deep(.markdown-content table) {
    font-size: 13px;
  }
  
  :deep(.markdown-content th),
  :deep(.markdown-content td) {
    padding: 12px;
  }
}
</style>