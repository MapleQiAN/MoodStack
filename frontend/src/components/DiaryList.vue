<template>
  <div class="list-container">
    <div class="list-header">
      <div class="header-content">
        <h1 class="page-title">我的日记</h1>
        <p class="page-description">共 {{ diaries.length }} 篇日记记录</p>
      </div>
              <button @click="selectFile" class="upload-button" :disabled="isUploading">
        <span class="upload-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14,2 14,8 20,8"/>
            <line x1="12" y1="11" x2="12" y2="17"/>
            <polyline points="9,14 12,11 15,14"/>
          </svg>
        </span>
        <span v-if="!isUploading">新增日记</span>
        <span v-else>上传中...</span>
      </button>
      
      <input 
        ref="fileInput"
        type="file"
        accept=".txt,.md,.docx,.pdf"
        @change="handleFileSelect"
        style="display: none"
      />
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
        <button @click="selectFile" class="empty-action">
          上传第一篇日记
        </button>
      </div>
    </div>
    
    <div v-else class="diary-grid">
      <div 
        v-for="diary in sortedDiaries" 
        :key="diary.id"
        class="diary-card"
        @click="$emit('view-diary', diary)"
      >
        <div class="card-header">
          <div class="diary-badge">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14,2 14,8 20,8"/>
            </svg>
          </div>
          <div class="diary-date">
            {{ formatDate(diary.createdAt) }}
          </div>
        </div>
        
        <div class="card-body">
          <h3 class="diary-title">{{ diary.title }}</h3>
          <p class="diary-preview">{{ getPreview(diary.content) }}</p>
        </div>
        
        <div class="card-footer">
          <div class="file-info">
            <span class="file-name">{{ diary.fileName }}</span>
          </div>
          <div class="read-more">
            <span>阅读</span>
            <span class="arrow">→</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { UploadDiary } from '../../wailsjs/go/main/App'

const props = defineProps({
  diaries: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['view-diary', 'refresh', 'upload-success'])

const fileInput = ref(null)
const selectedFile = ref(null)
const isUploading = ref(false)
const uploadError = ref('')

const sortedDiaries = computed(() => {
  return [...props.diaries].sort((a, b) => 
    new Date(b.createdAt) - new Date(a.createdAt)
  )
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

// 移除emoji图标函数，改用SVG图标

const selectFile = () => {
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

.upload-button {
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

.upload-button:hover:not(:disabled) {
  background: var(--accent-secondary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.upload-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.upload-icon {
  font-size: 16px;
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
  border: 2px solid var(--nord4);
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
  border: 1px solid var(--nord4);
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
}

.diary-badge {
  width: 40px;
  height: 40px;
  background: var(--nord5);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-primary);
  border: 1px solid var(--nord4);
}

.diary-date {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 600;
  background: var(--nord5);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--nord4);
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
  border-top: 1px solid var(--nord4);
}

.file-info {
  flex: 1;
}

.file-name {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
  background: var(--nord5);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--nord4);
  display: inline-block;
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
  
  .upload-button {
    align-self: flex-end;
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
  .card-footer {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .read-more {
    justify-content: center;
  }
}
</style>