<template>
  <div class="upload-container">
    <div class="upload-header">
      <h2 class="upload-title">📝 上传日记文件</h2>
      <p class="upload-description">支持 .txt、.md、.docx、.pdf 格式文件，最大 10MB</p>
    </div>
    
    <div class="upload-card">
      <div 
        class="upload-area"
        :class="{ 'drag-over': isDragOver, 'uploading': isUploading }"
        @drop="handleDrop"
        @dragover="handleDragOver"
        @dragleave="handleDragLeave"
        @click="selectFile"
      >
        <div v-if="!isUploading" class="upload-content">
          <div class="upload-icon">📁</div>
          <p class="upload-text">
            <strong>点击选择文件</strong> 或将文件拖拽到此处
          </p>
          <p class="upload-hint">支持多种格式，自动识别内容</p>
        </div>
        
        <div v-else class="uploading-content">
          <div class="upload-spinner"></div>
          <p>正在上传和转换文件...</p>
        </div>
      </div>
      
      <input 
        ref="fileInput"
        type="file"
        accept=".txt,.md,.docx,.pdf"
        @change="handleFileSelect"
        style="display: none"
      />
      
      <div v-if="selectedFile" class="file-info">
        <div class="file-details">
          <div class="file-icon">📄</div>
          <div class="file-info-content">
            <span class="file-name">{{ selectedFile.name }}</span>
            <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
          </div>
        </div>
        <button @click="uploadFile" :disabled="isUploading" class="upload-btn">
          <span v-if="!isUploading">上传文件</span>
          <span v-else>上传中...</span>
        </button>
      </div>
      
      <div v-if="uploadError" class="error-message">
        <div class="error-icon">❌</div>
        <div class="error-content">
          <h4>上传失败</h4>
          <p>{{ uploadError }}</p>
        </div>
      </div>
      
      <div v-if="uploadSuccess" class="success-message">
        <div class="success-icon">✅</div>
        <div class="success-content">
          <h4>上传成功！</h4>
          <p>文件上传成功！正在跳转到日记列表...</p>
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
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
}

.upload-header {
  text-align: center;
  margin-bottom: 32px;
  padding: 24px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.upload-title {
  font-size: 28px;
  font-weight: 700;
  color: #2d3748;
  margin-bottom: 8px;
}

.upload-description {
  color: #718096;
  font-size: 16px;
  font-weight: 500;
}

.upload-card {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20px;
  padding: 32px;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
}

.upload-area {
  border: 3px dashed rgba(102, 126, 234, 0.3);
  border-radius: 16px;
  padding: 48px 32px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 24px;
  background: rgba(255, 255, 255, 0.5);
}

.upload-area:hover {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.05);
  transform: translateY(-2px);
}

.upload-area.drag-over {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.1);
  transform: scale(1.02);
}

.upload-area.uploading {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.05);
  cursor: not-allowed;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.upload-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.8;
}

.upload-text {
  color: #4a5568;
  font-size: 18px;
  margin-bottom: 8px;
  font-weight: 500;
}

.upload-hint {
  color: #a0aec0;
  font-size: 14px;
}

.uploading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #667eea;
}

.upload-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(102, 126, 234, 0.2);
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

.file-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 12px;
  margin-bottom: 24px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.file-details {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-icon {
  font-size: 24px;
  opacity: 0.8;
}

.file-info-content {
  display: flex;
  flex-direction: column;
}

.file-name {
  font-weight: 600;
  color: #2d3748;
  font-size: 16px;
}

.file-size {
  font-size: 13px;
  color: #718096;
}

.upload-btn {
  padding: 12px 24px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 14px;
}

.upload-btn:hover:not(:disabled) {
  background: #5a67d8;
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.upload-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message, .success-message {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 16px;
}

.error-message {
  background: rgba(245, 101, 101, 0.1);
  border: 1px solid rgba(245, 101, 101, 0.3);
}

.success-message {
  background: rgba(72, 187, 120, 0.1);
  border: 1px solid rgba(72, 187, 120, 0.3);
}

.error-icon, .success-icon {
  font-size: 20px;
  margin-top: 2px;
}

.error-content h4, .success-content h4 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
}

.error-content h4 {
  color: #e53e3e;
}

.success-content h4 {
  color: #38a169;
}

.error-content p, .success-content p {
  font-size: 14px;
  color: #718096;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .upload-container {
    max-width: 100%;
  }
  
  .upload-card {
    padding: 24px;
  }
  
  .upload-area {
    padding: 32px 20px;
  }
  
  .upload-title {
    font-size: 24px;
  }
  
  .file-info {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .upload-btn {
    width: 100%;
  }
}
</style>