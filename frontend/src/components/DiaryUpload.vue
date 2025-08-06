<template>
  <div class="upload-container">
    <div class="upload-header">
      <h1 class="page-title">上传新日记</h1>
      <p class="page-description">支持 TXT、Markdown、Word 和 PDF 文件，最大 10MB</p>
    </div>
    
    <div class="upload-card">
      <div 
        class="upload-zone"
        :class="{ 'drag-over': isDragOver, 'uploading': isUploading }"
        @drop="handleDrop"
        @dragover="handleDragOver"
        @dragleave="handleDragLeave"
        @click="selectFile"
      >
        <div v-if="!isUploading" class="upload-content">
          <div class="upload-visual">
            <div class="upload-icon">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10,9 9,9 8,9"/>
              </svg>
            </div>
            <div class="upload-circle"></div>
          </div>
          <div class="upload-text">
            <h3>选择文件或拖拽到此处</h3>
            <p>支持 .txt、.md、.docx、.pdf 格式</p>
          </div>
        </div>
        
        <div v-else class="uploading-content">
          <div class="upload-spinner"></div>
          <p>正在处理文件...</p>
        </div>
      </div>
      
      <input 
        ref="fileInput"
        type="file"
        accept=".txt,.md,.docx,.pdf"
        @change="handleFileSelect"
        style="display: none"
      />
      
      <div v-if="selectedFile" class="file-preview">
        <div class="file-info">
          <div class="file-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14,2 14,8 20,8"/>
            </svg>
          </div>
          <div class="file-details">
            <span class="file-name">{{ selectedFile.name }}</span>
            <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
          </div>
        </div>
        <button @click="uploadFile" :disabled="isUploading" class="upload-button">
          <span v-if="!isUploading">开始上传</span>
          <span v-else>上传中...</span>
        </button>
      </div>
      
      <div v-if="uploadError" class="message error">
        <div class="message-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
        </div>
        <div class="message-content">
          <h4>上传失败</h4>
          <p>{{ uploadError }}</p>
        </div>
      </div>
      
      <div v-if="uploadSuccess" class="message success">
        <div class="message-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="20,6 9,17 4,12"/>
          </svg>
        </div>
        <div class="message-content">
          <h4>上传成功！</h4>
          <p>文件已成功处理，即将跳转到日记列表</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { UploadDiary } from '../../wailsjs/go/main/App'

const emit = defineEmits(['upload-success'])

const fileInput = ref(null)
const selectedFile = ref(null)
const isDragOver = ref(false)
const isUploading = ref(false)
const uploadError = ref('')
const uploadSuccess = ref(false)

const selectFile = () => {
  fileInput.value.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    selectedFile.value = file
    uploadError.value = ''
    uploadSuccess.value = false
  }
}

const handleDrop = (event) => {
  event.preventDefault()
  isDragOver.value = false
  
  const files = event.dataTransfer.files
  if (files.length > 0) {
    selectedFile.value = files[0]
    uploadError.value = ''
    uploadSuccess.value = false
  }
}

const handleDragOver = (event) => {
  event.preventDefault()
  isDragOver.value = true
}

const handleDragLeave = (event) => {
  event.preventDefault()
  isDragOver.value = false
}

// 移除emoji图标函数，改用SVG图标

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
    return
  }
  
  // Check file type
  const allowedTypes = ['.txt', '.md', '.docx', '.pdf']
  const fileExt = '.' + selectedFile.value.name.split('.').pop().toLowerCase()
  if (!allowedTypes.includes(fileExt)) {
    uploadError.value = '不支持的文件类型，请选择 .txt、.md、.docx 或 .pdf 文件'
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
    
    uploadSuccess.value = true
    
    // Reset form
    setTimeout(() => {
      selectedFile.value = null
      fileInput.value.value = ''
      uploadSuccess.value = false
      emit('upload-success')
    }, 1500)
    
  } catch (error) {
    console.error('Upload error:', error)
    uploadError.value = error.message || '上传失败，请重试'
  } finally {
    isUploading.value = false
  }
}
</script>

<style scoped>
.upload-container {
  max-width: 100%;
  padding: 32px;
  box-sizing: border-box;
}

.upload-header {
  margin-bottom: 32px;
}

.page-title {
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 12px;
  letter-spacing: -0.025em;
}

.page-description {
  font-size: 16px;
  color: var(--text-muted);
  font-weight: 500;
}

.upload-card {
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-xl);
  padding: 32px;
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--nord4);
}

.upload-zone {
  border: 2px dashed var(--nord3);
  border-radius: var(--radius-lg);
  padding: 48px 32px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 24px;
  position: relative;
  overflow: hidden;
}

.upload-zone:hover {
  border-color: var(--accent-primary);
  background: rgba(136, 192, 208, 0.1);
}

.upload-zone.drag-over {
  border-color: var(--accent-primary);
  background: rgba(136, 192, 208, 0.15);
  transform: scale(1.01);
}

.upload-zone.uploading {
  border-color: var(--nord3);
  cursor: not-allowed;
  background: var(--nord5);
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

.upload-visual {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-icon {
  color: var(--accent-primary);
  z-index: 2;
  position: relative;
}

.upload-circle {
  position: absolute;
  width: 80px;
  height: 80px;
  border: 2px solid var(--nord4);
  border-radius: 50%;
  z-index: 1;
}

.upload-text h3 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.upload-text p {
  font-size: 14px;
  color: var(--text-muted);
  font-weight: 500;
}

.uploading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: var(--text-muted);
}

.upload-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--nord4);
  border-top: 3px solid var(--accent-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.file-preview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: var(--nord5);
  border-radius: var(--radius-md);
  margin-bottom: 24px;
  border: 1px solid var(--nord4);
}

.file-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.file-icon {
  color: var(--accent-primary);
}

.file-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 15px;
}

.file-size {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 500;
}

.upload-button {
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

.message {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border-radius: var(--radius-md);
  border: 1px solid;
}

.message.error {
  background: rgba(191, 97, 106, 0.1);
  border-color: var(--accent-error);
  color: var(--accent-error);
}

.message.success {
  background: rgba(163, 190, 140, 0.1);
  border-color: var(--accent-success);
  color: var(--accent-success);
}

.message-icon {
  font-size: 20px;
  margin-top: 2px;
}

.message-content h4 {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 4px;
}

.message-content p {
  font-size: 13px;
  font-weight: 500;
  opacity: 0.8;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .upload-card {
    padding: 24px;
  }
  
  .upload-zone {
    padding: 32px 20px;
  }
  
  .page-title {
    font-size: 28px;
  }
  
  .file-preview {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .upload-button {
    width: 100%;
  }
}
</style>