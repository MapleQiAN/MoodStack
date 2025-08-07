<template>
  <div class="viewer-container">
    <div class="viewer-header">
      <div class="header-top">
        <button @click="$emit('back')" class="back-button">
        <span class="back-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15,18 9,12 15,6"/>
          </svg>
        </span>
        <span>返回列表</span>
      </button>
      
      <div class="diary-meta">
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
      </div>
      
      <div class="header-actions">
        <button 
          v-if="!isEditing"
          @click="startEditing" 
          class="edit-button"
          title="编辑日记"
        >
          <span class="edit-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
          </span>
          <span>编辑</span>
        </button>
        
        <div v-if="isEditing" class="edit-actions">
          <button @click="cancelEditing" class="cancel-button">
            <span>取消</span>
          </button>
          <button 
            @click="saveChanges" 
            class="save-button"
            :disabled="!hasChanges || isSaving"
          >
            <span v-if="!isSaving">保存</span>
            <span v-else>保存中...</span>
          </button>
        </div>
        </div>
      </div>
    </div>
    
    <!-- 查看模式 -->
    <div v-if="!isEditing" class="viewer-content">
      <article class="article-container">
        <header class="article-header">
          <h1 class="article-title">{{ diary.title }}</h1>
        </header>
        
        <div class="article-body">
          <div class="markdown-content" v-html="renderedContent"></div>
        </div>
      </article>
    </div>
    
    <!-- 编辑模式 -->
    <div v-if="isEditing" class="editor-content">
      <div class="editor-container">
        <div class="title-section">
          <input 
            v-model="editingTitle"
            class="title-input"
            placeholder="输入日记标题..."
            @input="markAsChanged"
          />
        </div>
        
        <div class="markdown-editor">
          <div class="editor-toolbar">
            <div class="toolbar-group">
              <button @click="insertBold" class="toolbar-btn" title="粗体 (Ctrl+B)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M6 4h8a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"/>
                  <path d="M6 12h9a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"/>
                </svg>
              </button>
              <button @click="insertItalic" class="toolbar-btn" title="斜体 (Ctrl+I)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polygon points="14,6 10,14 8,14 12,6"/>
                  <line x1="6" y1="4" x2="14" y2="4"/>
                  <line x1="4" y1="20" x2="12" y2="20"/>
                </svg>
              </button>
              <button @click="insertStrikethrough" class="toolbar-btn" title="删除线">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M16 4H9a3 3 0 0 0-2.83 4"/>
                  <path d="M14 12a4 4 0 0 1 0 8H6"/>
                  <line x1="4" y1="12" x2="20" y2="12"/>
                </svg>
              </button>
              <button @click="insertInlineCode" class="toolbar-btn" title="行内代码 (Ctrl+E)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="16,18 22,12 16,6"/>
                  <polyline points="8,6 2,12 8,18"/>
                </svg>
              </button>
            </div>
            
            <div class="toolbar-divider"></div>
            
            <div class="toolbar-group">
              <button @click="insertHeading(1)" class="toolbar-btn" title="一级标题">H1</button>
              <button @click="insertHeading(2)" class="toolbar-btn" title="二级标题">H2</button>
              <button @click="insertHeading(3)" class="toolbar-btn" title="三级标题">H3</button>
            </div>
            
            <div class="toolbar-divider"></div>
            
            <div class="toolbar-group">
              <button @click="insertUnorderedList" class="toolbar-btn" title="无序列表">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="8" y1="6" x2="21" y2="6"/>
                  <line x1="8" y1="12" x2="21" y2="12"/>
                  <line x1="8" y1="18" x2="21" y2="18"/>
                  <line x1="3" y1="6" x2="3" y2="6"/>
                  <line x1="3" y1="12" x2="3" y2="12"/>
                  <line x1="3" y1="18" x2="3" y2="18"/>
                </svg>
              </button>
              <button @click="insertOrderedList" class="toolbar-btn" title="有序列表">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="10" y1="6" x2="21" y2="6"/>
                  <line x1="10" y1="12" x2="21" y2="12"/>
                  <line x1="10" y1="18" x2="21" y2="18"/>
                  <text x="3" y="9" font-size="12" fill="currentColor">1</text>
                  <text x="3" y="15" font-size="12" fill="currentColor">2</text>
                  <text x="3" y="21" font-size="12" fill="currentColor">3</text>
                </svg>
              </button>
              <button @click="insertQuote" class="toolbar-btn" title="引用">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 21c3 0 7-1 7-8V5c0-1.25-.756-2.017-2-2H4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2 1 0 1 0 1 1v1c0 1-1 2-2 2s-1 .008-1 1.031V20c0 1 0 1 1 1z"/>
                  <path d="M15 21c3 0 7-1 7-8V5c0-1.25-.757-2.017-2-2h-4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2h.75c0 2.25.25 4-2.75 4v3c0 1 0 1 1 1z"/>
                </svg>
              </button>
              <button @click="insertLink" class="toolbar-btn" title="链接 (Ctrl+K)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.72"/>
                  <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.72-1.72"/>
                </svg>
              </button>
            </div>
            
            <div class="toolbar-divider"></div>
            
            <div class="toolbar-group">
              <button @click="togglePreview" class="toolbar-btn" :class="{ active: showPreview }" title="预览">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
              </button>
            </div>
          </div>
          
          <div class="editor-main" :class="{ 'preview-mode': showPreview }">
            <div class="editor-pane">
              <textarea
                ref="textareaRef"
                v-model="editingContent"
                class="markdown-textarea"
                placeholder="在这里写下你的故事..."
                @input="markAsChanged"
                @keydown="handleKeyDown"
              ></textarea>
            </div>
            
            <div v-if="showPreview" class="preview-pane">
              <div class="preview-content markdown-content" v-html="previewContent"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { marked } from 'marked'
import { UpdateDiary } from '../../wailsjs/go/main/App'

const props = defineProps({
  diary: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['back', 'diary-updated'])

// 编辑状态
const isEditing = ref(false)
const hasChanges = ref(false)
const isSaving = ref(false)
const showPreview = ref(false)

// 编辑内容
const editingTitle = ref('')
const editingContent = ref('')
const originalTitle = ref('')
const originalContent = ref('')

// 元素引用
const textareaRef = ref(null)

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

// 预览内容
const previewContent = computed(() => {
  if (!editingContent.value) return ''
  return marked.parse(editingContent.value)
})

// 编辑功能
const startEditing = () => {
  isEditing.value = true
  editingTitle.value = props.diary.title || ''
  editingContent.value = props.diary.content || ''
  originalTitle.value = props.diary.title || ''
  originalContent.value = props.diary.content || ''
  hasChanges.value = false
  
  nextTick(() => {
    if (textareaRef.value) {
      textareaRef.value.focus()
    }
  })
}

const cancelEditing = () => {
  if (hasChanges.value) {
    if (confirm('有未保存的修改，确定要取消编辑吗？')) {
      isEditing.value = false
      showPreview.value = false
      hasChanges.value = false
    }
  } else {
    isEditing.value = false
    showPreview.value = false
  }
}

const markAsChanged = () => {
  const titleChanged = editingTitle.value !== originalTitle.value
  const contentChanged = editingContent.value !== originalContent.value
  hasChanges.value = titleChanged || contentChanged
}

const saveChanges = async () => {
  if (!hasChanges.value || isSaving.value) return
  
  try {
    isSaving.value = true
    
    const updatedDiary = {
      ...props.diary,
      title: editingTitle.value.trim() || '无标题',
      content: editingContent.value,
      updatedAt: new Date()
    }
    
    await UpdateDiary(updatedDiary)
    
    // 更新原始值
    originalTitle.value = editingTitle.value
    originalContent.value = editingContent.value
    hasChanges.value = false
    isEditing.value = false
    showPreview.value = false
    
    // 通知父组件日记已更新
    emit('diary-updated', updatedDiary)
  } catch (error) {
    console.error('保存日记失败:', error)
    alert('保存失败，请稍后重试')
  } finally {
    isSaving.value = false
  }
}

// Markdown工具栏功能
const insertAtCursor = (before, after = '', placeholder = '') => {
  const textarea = textareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = editingContent.value.substring(start, end)
  const replacement = selectedText || placeholder
  
  const newText = before + replacement + after
  const newValue = editingContent.value.substring(0, start) + newText + editingContent.value.substring(end)
  
  editingContent.value = newValue
  markAsChanged()
  
  nextTick(() => {
    const newStart = start + before.length
    const newEnd = newStart + replacement.length
    textarea.setSelectionRange(newStart, newEnd)
    textarea.focus()
  })
}

const insertBold = () => insertAtCursor('**', '**', '粗体文本')
const insertItalic = () => insertAtCursor('*', '*', '斜体文本')
const insertStrikethrough = () => insertAtCursor('~~', '~~', '删除线文本')
const insertInlineCode = () => insertAtCursor('`', '`', '代码')

const insertHeading = (level) => {
  const hashes = '#'.repeat(level)
  insertAtCursor(`${hashes} `, '', '标题文本')
}

const insertUnorderedList = () => insertAtCursor('- ', '', '列表项')
const insertOrderedList = () => insertAtCursor('1. ', '', '列表项')
const insertQuote = () => insertAtCursor('> ', '', '引用文本')

const insertLink = () => {
  const url = prompt('请输入链接地址:', 'https://')
  if (url) {
    insertAtCursor('[', `](${url})`, '链接文本')
  }
}

const togglePreview = () => {
  showPreview.value = !showPreview.value
}

// 编辑历史管理（用于撤销功能）
const editHistory = ref([])
const historyIndex = ref(-1)
const maxHistorySize = 50

const saveToHistory = () => {
  // 移除当前位置之后的历史记录
  editHistory.value.splice(historyIndex.value + 1)
  
  // 添加新的历史记录
  editHistory.value.push({
    title: editingTitle.value,
    content: editingContent.value,
    timestamp: Date.now()
  })
  
  // 限制历史记录大小
  if (editHistory.value.length > maxHistorySize) {
    editHistory.value.shift()
  } else {
    historyIndex.value++
  }
}

const undo = () => {
  if (historyIndex.value > 0) {
    historyIndex.value--
    const state = editHistory.value[historyIndex.value]
    editingTitle.value = state.title
    editingContent.value = state.content
    markAsChanged()
  }
}

const redo = () => {
  if (historyIndex.value < editHistory.value.length - 1) {
    historyIndex.value++
    const state = editHistory.value[historyIndex.value]
    editingTitle.value = state.title
    editingContent.value = state.content
    markAsChanged()
  }
}

// 键盘快捷键
const handleKeyDown = (event) => {
  if (event.ctrlKey || event.metaKey) {
    switch (event.key) {
      case 'b':
        event.preventDefault()
        insertBold()
        break
      case 'i':
        event.preventDefault()
        insertItalic()
        break
      case 'k':
        event.preventDefault()
        insertLink()
        break
      case 'e':
        event.preventDefault()
        insertInlineCode()
        break
      case 's':
        event.preventDefault()
        if (hasChanges.value) {
          saveChanges()
        }
        break
      case 'z':
        event.preventDefault()
        if (event.shiftKey) {
          redo()
        } else {
          undo()
        }
        break
      case 'y':
        event.preventDefault()
        redo()
        break
    }
  }
  
  if (event.key === 'Tab') {
    event.preventDefault()
    insertAtCursor('  ', '', '')
  }
}

// 初始化编辑历史
const initializeHistory = () => {
  editHistory.value = [{
    title: editingTitle.value,
    content: editingContent.value,
    timestamp: Date.now()
  }]
  historyIndex.value = 0
}

// 监听内容变化并保存历史
let historyTimer = null
const scheduleHistorySave = () => {
  if (historyTimer) {
    clearTimeout(historyTimer)
  }
  historyTimer = setTimeout(() => {
    saveToHistory()
  }, 1000) // 1秒后保存到历史
}

// 监听编辑内容变化
watch([editingTitle, editingContent], () => {
  if (isEditing.value) {
    scheduleHistorySave()
  }
})

// 处理粘贴事件
const handlePaste = (event) => {
  const clipboardData = event.clipboardData || window.clipboardData
  const pastedData = clipboardData.getData('text')
  
  // 如果粘贴的是URL，自动转换为链接格式
  const urlRegex = /^https?:\/\/.+/i
  if (urlRegex.test(pastedData.trim())) {
    event.preventDefault()
    insertAtCursor(`[链接文本](${pastedData.trim()})`, '[', ']()')
  }
}

// 组件挂载时的处理
onMounted(() => {
  initializeHistory()
  
  // 绑定键盘事件
  document.addEventListener('keydown', handleKeyDown)
  
  // 绑定粘贴事件
  const textarea = document.querySelector('.content-editor')
  if (textarea) {
    textarea.addEventListener('paste', handlePaste)
  }
})

// 组件卸载时的清理
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
  
  const textarea = document.querySelector('.content-editor')
  if (textarea) {
    textarea.removeEventListener('paste', handlePaste)
  }
  
  if (historyTimer) {
    clearTimeout(historyTimer)
  }
})
</script>

<style scoped>
.viewer-container {
  max-width: 100%;
  box-sizing: border-box;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: visible;
}

.viewer-header {
  position: sticky;
  top: 0;
  z-index: 1000;
  padding: 16px 32px;
  background: var(--bg-header);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--bg-tertiary);
  box-shadow: var(--shadow-sm);
  width: 100%;
  box-sizing: border-box;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.diary-meta {
  display: flex;
  align-items: center;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-muted);
  font-size: 14px;
  font-weight: 500;
}

.meta-icon {
  display: flex;
  align-items: center;
  opacity: 0.7;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.2s ease;
}

.back-button:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.back-icon {
  font-size: 16px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.edit-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--accent-primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.2s ease;
}

.edit-button:hover {
  background: var(--accent-secondary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.edit-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cancel-button, .save-button {
  padding: 8px 16px;
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.2s ease;
}

.cancel-button {
  background: var(--bg-secondary);
  color: var(--text-secondary);
}

.cancel-button:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.save-button {
  background: var(--accent-primary);
  color: white;
}

.save-button:hover:not(:disabled) {
  background: var(--accent-secondary);
}

.save-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
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
  padding: 32px;
  flex: 1;
  overflow-y: auto;
}

.article-container {
  padding: 0;
}

.article-header {
  padding: 32px 0 24px 0;
  border-bottom: 1px solid var(--bg-tertiary);
  margin-bottom: 32px;
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
  padding: 0;
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
  border-bottom: 2px solid var(--bg-tertiary);
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
  background: var(--bg-secondary);
  padding: 16px 20px;
  margin: 20px 0;
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
  font-style: italic;
  color: var(--text-muted);
}

:deep(.markdown-content code) {
  background: var(--bg-secondary);
  color: var(--text-primary);
  padding: 3px 6px;
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: 14px;
  border: 1px solid var(--bg-tertiary);
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
  border: 1px solid var(--bg-tertiary);
}

:deep(.markdown-content th) {
  background: var(--bg-secondary);
  color: var(--text-primary);
  padding: 16px;
  text-align: left;
  font-weight: 600;
  border-bottom: 1px solid var(--bg-tertiary);
}

:deep(.markdown-content td) {
  padding: 14px 16px;
  border-bottom: 1px solid var(--bg-tertiary);
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
  border: 1px solid var(--bg-tertiary);
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

/* 编辑器样式 */
.editor-content {
  padding: 32px;
  flex: 1;
  overflow-y: auto;
}

.editor-container {
  padding: 0;
}

.title-section {
  padding: 24px 0;
  border-bottom: 1px solid var(--bg-tertiary);
  margin-bottom: 24px;
}

.title-input {
  width: 100%;
  padding: 12px 0;
  border: none;
  background: transparent;
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.3;
  letter-spacing: -0.025em;
  outline: none;
}

.title-input::placeholder {
  color: var(--text-muted);
}

.markdown-editor {
  display: flex;
  flex-direction: column;
}

.editor-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 0;
  border-bottom: 1px solid var(--bg-tertiary);
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.toolbar-group {
  display: flex;
  align-items: center;
  gap: 4px;
}

.toolbar-divider {
  width: 1px;
  height: 24px;
  background: var(--bg-tertiary);
  margin: 0 4px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 12px;
  font-weight: 600;
}

.toolbar-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.toolbar-btn.active {
  background: var(--accent-primary);
  color: white;
}

.editor-main {
  display: flex;
  min-height: 500px;
}



.editor-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.preview-pane {
  flex: 1;
  border-left: 1px solid var(--bg-tertiary);
  background: var(--nord6);
  overflow-y: auto;
}

.markdown-textarea {
  flex: 1;
  width: 100%;
  padding: 24px 0;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 14px;
  line-height: 1.6;
  resize: none;
  outline: none;
}

.markdown-textarea::placeholder {
  color: var(--text-muted);
}

.preview-content {
  padding: 24px 0;
  color: var(--text-secondary);
  line-height: 1.7;
  font-size: 14px;
  max-width: none;
}

/* 编辑器响应式设计 */
@media (max-width: 768px) {
  .editor-main.preview-mode {
    flex-direction: column;
  }
  
  .preview-pane {
    border-left: none;
    border-top: 1px solid var(--bg-tertiary);
    max-height: 300px;
  }
  
  .toolbar-group {
    gap: 2px;
  }
  
  .toolbar-btn {
    width: 28px;
    height: 28px;
  }
  
  .title-section {
    padding: 20px 16px;
  }
  
  .title-input {
    font-size: 24px;
  }
  
  .markdown-textarea {
    padding: 20px 16px;
  }
  
  .preview-content {
    padding: 20px 16px;
  }
}
</style>