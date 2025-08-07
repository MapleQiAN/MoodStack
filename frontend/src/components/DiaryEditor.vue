<template>
  <div class="editor-container">
    <div class="editor-header">
      <div class="header-left">
        <button @click="$emit('back')" class="back-button">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M19 12H5"/>
            <path d="M12 19l-7-7 7-7"/>
          </svg>
          <span>返回</span>
        </button>
        <div class="editor-title">
          <h1 v-if="diary.id">编辑日记</h1>
          <h1 v-else>新建日记</h1>
          <p class="last-saved" v-if="lastSaved">
            {{ lastSaved }}
          </p>
        </div>
      </div>
      
      <div class="header-actions">
        <button 
          @click="showEncryptionSettings = !showEncryptionSettings" 
          class="encryption-button"
          :class="{ active: encryptionOptions.mode !== 'unified' }"
          title="加密设置"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
            <circle cx="12" cy="16" r="1"/>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
          </svg>
          <span>加密</span>
        </button>
        <button 
          @click="saveDiary" 
          :disabled="isSaving || !hasChanges"
          class="save-button"
        >
          <span v-if="!isSaving">保存</span>
          <span v-else>保存中...</span>
        </button>
      </div>
    </div>
    
    <!-- 加密设置面板 -->
    <div v-if="showEncryptionSettings" class="encryption-panel">
      <div class="encryption-options">
        <h3>加密设置</h3>
        <div class="encryption-modes">
          <label class="mode-option">
            <input 
              type="radio" 
              value="unified" 
              v-model="encryptionOptions.mode"
              @change="onEncryptionModeChange"
            />
            <div class="mode-info">
              <span class="mode-title">统一密码</span>
              <span class="mode-desc">使用您的主密码加密</span>
            </div>
          </label>
          
          <label class="mode-option">
            <input 
              type="radio" 
              value="individual" 
              v-model="encryptionOptions.mode"
              @change="onEncryptionModeChange"
            />
            <div class="mode-info">
              <span class="mode-title">单独密码</span>
              <span class="mode-desc">为此日记设置独立密码</span>
            </div>
          </label>
          
          <label class="mode-option">
            <input 
              type="radio" 
              value="biometric" 
              v-model="encryptionOptions.mode"
              @change="onEncryptionModeChange"
            />
            <div class="mode-info">
              <span class="mode-title">Windows Hello</span>
              <span class="mode-desc">使用生物识别解锁</span>
            </div>
          </label>
        </div>
        
        <div v-if="encryptionOptions.mode === 'individual'" class="individual-password">
          <label for="individual-password">单独密码:</label>
          <input 
            id="individual-password"
            type="password" 
            v-model="encryptionOptions.individualPassword"
            placeholder="请输入此日记的专用密码"
            class="password-input"
            @input="markAsChanged"
          />
        </div>
        
        <div class="encryption-actions">
          <button @click="showEncryptionSettings = false" class="btn-secondary">
            关闭
          </button>
        </div>
      </div>
    </div>

    <div class="editor-content">
      <div class="title-input-wrapper">
        <input 
          v-model="localDiary.title"
          type="text"
          placeholder="请输入标题..."
          class="title-input"
          @input="markAsChanged"
        />
      </div>
      
      <div class="content-editor">
        <div class="editor-toolbar">
          <!-- <div class="toolbar-info">
            <span class="toolbar-hint">提示: Ctrl+S 保存, Ctrl+B 粗体, Ctrl+I 斜体, Ctrl+K 链接</span>
          </div> -->
          <div class="toolbar-group">
            <button @click="insertMarkdown('**', '**')" class="toolbar-btn" title="粗体 (Ctrl+B)">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M6 4h8a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"/>
                <path d="M6 12h9a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"/>
              </svg>
            </button>
            <button @click="insertMarkdown('*', '*')" class="toolbar-btn" title="斜体 (Ctrl+I)">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="10,4 14,4 6,20 2,20"/>
              </svg>
            </button>
            <button @click="insertMarkdown('~~', '~~')" class="toolbar-btn" title="删除线">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M6 16L18 8"/>
                <path d="M6 8L18 16"/>
                <path d="M6 12h12"/>
              </svg>
            </button>
            <button @click="insertMarkdown('`', '`')" class="toolbar-btn" title="行内代码">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="16,18 22,12 16,6"/>
                <polyline points="8,6 2,12 8,18"/>
              </svg>
            </button>
          </div>
          
          <div class="toolbar-group">
            <button @click="insertMarkdown('# ', '')" class="toolbar-btn" title="一级标题">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M4 12h8m-8-6v12"/>
                <path d="M20 6v12"/>
                <circle cx="20" cy="12" r="1"/>
              </svg>
            </button>
            <button @click="insertMarkdown('## ', '')" class="toolbar-btn" title="二级标题">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M4 12h8m-8-6v12"/>
                <path d="M20 6v12"/>
                <circle cx="16" cy="12" r="1"/>
                <circle cx="20" cy="12" r="1"/>
              </svg>
            </button>
            <button @click="insertMarkdown('### ', '')" class="toolbar-btn" title="三级标题">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M4 12h8m-8-6v12"/>
                <path d="M20 6v12"/>
                <circle cx="14" cy="12" r="1"/>
                <circle cx="17" cy="12" r="1"/>
                <circle cx="20" cy="12" r="1"/>
              </svg>
            </button>
          </div>
          
          <div class="toolbar-group">
            <button @click="insertMarkdown('- ', '')" class="toolbar-btn" title="无序列表">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="8" y1="6" x2="21" y2="6"/>
                <line x1="8" y1="12" x2="21" y2="12"/>
                <line x1="8" y1="18" x2="21" y2="18"/>
                <line x1="3" y1="6" x2="3.01" y2="6"/>
                <line x1="3" y1="12" x2="3.01" y2="12"/>
                <line x1="3" y1="18" x2="3.01" y2="18"/>
              </svg>
            </button>
            <button @click="insertMarkdown('1. ', '')" class="toolbar-btn" title="有序列表">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="10" y1="6" x2="21" y2="6"/>
                <line x1="10" y1="12" x2="21" y2="12"/>
                <line x1="10" y1="18" x2="21" y2="18"/>
                <path d="M4 6h1v4"/>
                <path d="M4 10h2"/>
                <path d="M6 18H4c0-1 1-2 2-2s1 1 1 2"/>
              </svg>
            </button>
            <button @click="insertMarkdown('> ', '')" class="toolbar-btn" title="引用">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 21c3 0 7-1 7-8V5c0-1.25-.756-2.017-2-2H4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2 1 0 1 0 1 1v1c0 1-1 2-2 2s-1 .008-1 1.031V20c0 1 0 1 1 1z"/>
                <path d="M15 21c3 0 7-1 7-8V5c0-1.25-.757-2.017-2-2h-4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2h.75c0 2.25.25 4-2.75 4v3c0 1 0 1 1 1z"/>
              </svg>
            </button>
            <button @click="insertLink" class="toolbar-btn" title="链接">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.72"/>
                <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.72-1.72"/>
              </svg>
            </button>
          </div>
          
          <div class="toolbar-group">
            <button @click="showPreview = !showPreview" class="toolbar-btn" :class="{ active: showPreview }" title="预览">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
            </button>
          </div>
        </div>
        
        <div class="editor-main" :class="{ 'split-view': showPreview }">
          <div class="text-editor">
            <textarea 
              ref="contentTextarea"
              v-model="localDiary.content"
              placeholder="开始写下你的心情故事..."
              class="content-textarea"
              @input="markAsChanged"
              @scroll="syncScroll"
            ></textarea>
          </div>
          
          <div v-if="showPreview" class="preview-panel">
            <div ref="previewContent" class="preview-content" v-html="renderedContent"></div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-if="saveError" class="message error">
      <div class="message-icon">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="15" y1="9" x2="9" y2="15"/>
          <line x1="9" y1="9" x2="15" y2="15"/>
        </svg>
      </div>
      <div class="message-content">
        <h4>保存失败</h4>
        <p>{{ saveError }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { CreateDiary, UpdateDiary, CreateDiaryWithEncryption, UpdateDiaryWithEncryption } from '../../wailsjs/go/main/App'

const props = defineProps({
  diary: {
    type: Object,
    default: () => ({
      id: '',
      title: '',
      content: '',
      tags: []
    })
  }
})

const emit = defineEmits(['back', 'diary-saved'])

const localDiary = ref({ ...props.diary })
const hasChanges = ref(false)
const isSaving = ref(false)
const saveError = ref('')
const lastSaved = ref('')
const showPreview = ref(false)
const contentTextarea = ref(null)
const previewContent = ref(null)

// 加密相关
const showEncryptionSettings = ref(false)
const encryptionOptions = ref({
  mode: 'unified', // 'unified', 'individual', 'biometric'
  individualPassword: ''
})

// 监听props变化
watch(() => props.diary, (newDiary) => {
  localDiary.value = { ...newDiary }
  hasChanges.value = false
}, { deep: true })

onMounted(() => {
  // 如果是新日记，设置默认内容并聚焦到标题输入框
  if (!localDiary.value.id && !localDiary.value.content.trim()) {
    const defaultContent = `# 今天的心情

## 今日感想
写下你今天的想法和感受...

## 重要事件
- 
- 
- 

## 明天计划
> 明天我想要...

---
*记录于 ${new Date().toLocaleDateString('zh-CN')}*`
    
    localDiary.value.content = defaultContent
    hasChanges.value = true // 标记为已更改，允许保存
    
    nextTick(() => {
      const titleInput = document.querySelector('.title-input')
      if (titleInput) {
        titleInput.focus()
      }
    })
  }
})

const markAsChanged = () => {
  hasChanges.value = true
}

const onEncryptionModeChange = () => {
  markAsChanged()
  // 切换模式时清空单独密码
  if (encryptionOptions.value.mode !== 'individual') {
    encryptionOptions.value.individualPassword = ''
  }
}

const saveDiary = async () => {
  if (!localDiary.value.title.trim() && !localDiary.value.content.trim()) {
    saveError.value = '请至少填写标题或内容'
    return
  }
  
  // 验证单独密码模式下的密码
  if (encryptionOptions.value.mode === 'individual' && !encryptionOptions.value.individualPassword) {
    saveError.value = '请输入单独密码'
    return
  }
  
  try {
    isSaving.value = true
    saveError.value = ''
    
    if (localDiary.value.id) {
      // 更新现有日记
      if (encryptionOptions.value.mode !== 'unified') {
        await UpdateDiaryWithEncryption(localDiary.value, encryptionOptions.value)
      } else {
        await UpdateDiary(localDiary.value)
      }
    } else {
      // 创建新日记
      if (encryptionOptions.value.mode !== 'unified') {
        const result = await CreateDiaryWithEncryption(localDiary.value.title, localDiary.value.content, encryptionOptions.value)
        localDiary.value.id = result.id
      } else {
        const result = await CreateDiary(localDiary.value.title, localDiary.value.content)
        localDiary.value.id = result.id
      }
    }
    
    hasChanges.value = false
    lastSaved.value = `最后保存于 ${new Date().toLocaleTimeString('zh-CN')}`
    
    // 保存成功后清空单独密码（避免在界面上保留敏感信息）
    if (encryptionOptions.value.mode === 'individual') {
      encryptionOptions.value.individualPassword = ''
    }
    
    // 短暂延迟后触发保存成功事件
    setTimeout(() => {
      emit('diary-saved', localDiary.value)
    }, 500)
    
  } catch (error) {
    console.error('Save error:', error)
    saveError.value = error.message || '保存失败，请重试'
  } finally {
    isSaving.value = false
  }
}

const insertMarkdown = (before, after = '') => {
  const textarea = contentTextarea.value
  if (!textarea) return
  
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = localDiary.value.content.substring(start, end)
  
  // 对于行级元素，确保在新行开始
  const isLineElement = before.includes('#') || before.includes('-') || before.includes('1.') || before.includes('>')
  let prefix = ''
  
  if (isLineElement && start > 0 && localDiary.value.content[start - 1] !== '\n') {
    prefix = '\n'
  }
  
  const newText = prefix + before + selectedText + after
  localDiary.value.content = 
    localDiary.value.content.substring(0, start) + 
    newText + 
    localDiary.value.content.substring(end)
  
  markAsChanged()
  
  // 重新设置光标位置
  nextTick(() => {
    textarea.focus()
    const newPosition = start + prefix.length + before.length + selectedText.length + after.length
    textarea.setSelectionRange(newPosition, newPosition)
  })
}

const insertLink = () => {
  const textarea = contentTextarea.value
  if (!textarea) return
  
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = localDiary.value.content.substring(start, end)
  
  const linkText = selectedText || '链接文本'
  const linkUrl = 'https://example.com'
  
  const markdown = `[${linkText}](${linkUrl})`
  
  localDiary.value.content = 
    localDiary.value.content.substring(0, start) + 
    markdown + 
    localDiary.value.content.substring(end)
  
  markAsChanged()
  
  // 选中URL部分便于编辑
  nextTick(() => {
    textarea.focus()
    const urlStart = start + linkText.length + 3 // [text](
    const urlEnd = urlStart + linkUrl.length
    textarea.setSelectionRange(urlStart, urlEnd)
  })
}

const syncScroll = () => {
  if (!showPreview.value || !previewContent.value) return
  
  const textarea = contentTextarea.value
  const preview = previewContent.value
  
  const scrollPercentage = textarea.scrollTop / (textarea.scrollHeight - textarea.clientHeight)
  preview.scrollTop = scrollPercentage * (preview.scrollHeight - preview.clientHeight)
}

// 改进的Markdown渲染器
const renderedContent = computed(() => {
  if (!localDiary.value.content.trim()) {
    return '<p class="empty-preview">在左侧输入内容以查看预览...</p>'
  }
  
  let html = localDiary.value.content
  
  // 处理代码块（需要在其他替换之前处理）
  html = html.replace(/```([\s\S]*?)```/g, '<pre><code>$1</code></pre>')
  
  // 处理标题
  html = html.replace(/^### (.*$)/gim, '<h3>$1</h3>')
  html = html.replace(/^## (.*$)/gim, '<h2>$1</h2>')
  html = html.replace(/^# (.*$)/gim, '<h1>$1</h1>')
  
  // 处理文本样式
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/(?<!\*)\*([^*]+)\*(?!\*)/g, '<em>$1</em>')
  html = html.replace(/~~(.*?)~~/g, '<del>$1</del>')
  html = html.replace(/`([^`]+)`/g, '<code>$1</code>')
  
  // 处理链接
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>')
  
  // 处理引用
  html = html.replace(/^> (.*$)/gim, '<blockquote>$1</blockquote>')
  
  // 处理列表
  html = html.replace(/^(\d+)\. (.*$)/gim, '<li class="ordered">$2</li>')
  html = html.replace(/^- (.*$)/gim, '<li class="unordered">$2</li>')
  
  // 包装列表项
  html = html.replace(/(<li class="ordered">.*?<\/li>)/gs, '<ol>$1</ol>')
  html = html.replace(/(<li class="unordered">.*?<\/li>)/gs, '<ul>$1</ul>')
  
  // 清理类名
  html = html.replace(/class="(ordered|unordered)"/g, '')
  
  // 处理段落和换行
  html = html.replace(/\n\n/g, '</p><p>')
  html = html.replace(/\n/g, '<br>')
  html = '<p>' + html + '</p>'
  
  // 清理空段落和修复嵌套问题
  html = html.replace(/<p><\/p>/g, '')
  html = html.replace(/<p>(<h[1-6]>.*?<\/h[1-6]>)<\/p>/g, '$1')
  html = html.replace(/<p>(<blockquote>.*?<\/blockquote>)<\/p>/g, '$1')
  html = html.replace(/<p>(<[uo]l>.*?<\/[uo]l>)<\/p>/g, '$1')
  html = html.replace(/<p>(<pre><code>.*?<\/code><\/pre>)<\/p>/g, '$1')
  
  return html
})

// 键盘快捷键
const handleKeydown = (event) => {
  if (event.ctrlKey || event.metaKey) {
    switch (event.key.toLowerCase()) {
      case 's':
        event.preventDefault()
        if (hasChanges.value) {
          saveDiary()
        }
        break
      case 'b':
        event.preventDefault()
        insertMarkdown('**', '**')
        break
      case 'i':
        event.preventDefault()
        insertMarkdown('*', '*')
        break
      case 'k':
        event.preventDefault()
        insertLink()
        break
      case 'e':
        event.preventDefault()
        insertMarkdown('`', '`')
        break
    }
  }
  
  // Tab键处理
  if (event.key === 'Tab') {
    event.preventDefault()
    const textarea = contentTextarea.value
    if (!textarea) return
    
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    
    // 插入制表符或空格
    const tabChar = '  ' // 使用2个空格代替制表符
    localDiary.value.content = 
      localDiary.value.content.substring(0, start) + 
      tabChar + 
      localDiary.value.content.substring(end)
    
    markAsChanged()
    
    nextTick(() => {
      textarea.focus()
      textarea.setSelectionRange(start + tabChar.length, start + tabChar.length)
    })
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

// 清理事件监听器
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.editor-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px;
  border-bottom: 1px solid var(--bg-tertiary);
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: transparent;
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 500;
}

.back-button:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.editor-title h1 {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.last-saved {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.encryption-button {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: transparent;
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 500;
  font-size: 14px;
}

.encryption-button:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.encryption-button.active {
  background: var(--accent-primary);
  color: white;
  border-color: var(--accent-primary);
}

/* 加密面板样式 */
.encryption-panel {
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--bg-tertiary);
  padding: 24px 32px;
}

.encryption-options h3 {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 20px;
}

.encryption-modes {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 20px;
}

.mode-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all 0.2s ease;
}

.mode-option:hover {
  background: var(--bg-tertiary);
  border-color: var(--accent-primary);
}

.mode-option input[type="radio"] {
  width: 18px;
  height: 18px;
  accent-color: var(--accent-primary);
}

.mode-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mode-title {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.mode-desc {
  font-size: 12px;
  color: var(--text-muted);
}

.individual-password {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 20px;
}

.individual-password label {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.password-input {
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  font-size: 14px;
  transition: border-color 0.2s ease;
}

.password-input:focus {
  outline: none;
  border-color: var(--accent-primary);
}

.encryption-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-secondary {
  padding: 8px 16px;
  background: transparent;
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 500;
}

.btn-secondary:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.save-button {
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

.save-button:hover:not(:disabled) {
  background: var(--accent-secondary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.save-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.editor-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.title-input-wrapper {
  padding: 24px 32px 0;
}

.title-input {
  width: 100%;
  border: none;
  background: transparent;
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 600;
  color: var(--text-primary);
  outline: none;
  padding: 16px 0;
  border-bottom: 2px solid transparent;
  transition: border-color 0.2s ease;
}

.title-input:focus {
  border-bottom-color: var(--accent-primary);
}

.title-input::placeholder {
  color: var(--text-muted);
  opacity: 0.6;
}

.content-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0 32px 32px;
}

.editor-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid var(--bg-tertiary);
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 16px;
}

.toolbar-info {
  display: flex;
  align-items: center;
  order: -1;
  width: 100%;
  margin-bottom: 8px;
}

.toolbar-hint {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
  background: var(--bg-secondary);
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--bg-tertiary);
}

.toolbar-group {
  display: flex;
  gap: 8px;
}

.toolbar-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid var(--bg-tertiary);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.toolbar-btn:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.toolbar-btn.active {
  background: var(--accent-primary);
  color: white;
  border-color: var(--accent-primary);
}

.editor-main {
  flex: 1;
  display: flex;
  gap: 16px;
  overflow: hidden;
}

.editor-main.split-view .text-editor {
  flex: 1;
}

.text-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.content-textarea {
  flex: 1;
  border: none;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-md);
  padding: 20px;
  font-family: var(--font-mono);
  font-size: 15px;
  line-height: 1.6;
  color: var(--text-primary);
  outline: none;
  resize: none;
  border: 1px solid var(--bg-tertiary);
  transition: border-color 0.2s ease;
}

.content-textarea:focus {
  border-color: var(--accent-primary);
}

.content-textarea::placeholder {
  color: var(--text-muted);
  opacity: 0.6;
}

.preview-panel {
  flex: 1;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-md);
  border: 1px solid var(--bg-tertiary);
  overflow: hidden;
}

.preview-content {
  height: 100%;
  overflow-y: auto;
  padding: 20px;
  font-family: var(--font-body);
  font-size: 15px;
  line-height: 1.6;
  color: var(--text-primary);
}

.preview-content h1 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-primary);
}

.preview-content h2 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--text-primary);
}

.preview-content h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.preview-content strong {
  font-weight: 600;
  color: var(--text-primary);
}

.preview-content em {
  font-style: italic;
  color: var(--text-secondary);
}

.preview-content code {
  background: var(--bg-tertiary);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: 13px;
}

.preview-content ul,
.preview-content ol {
  margin: 12px 0;
  padding-left: 20px;
}

.preview-content li {
  margin-bottom: 4px;
}

.preview-content ol {
  list-style-type: decimal;
}

.preview-content ul {
  list-style-type: disc;
}

.preview-content blockquote {
  margin: 16px 0;
  padding: 12px 16px;
  border-left: 4px solid var(--accent-primary);
  background: var(--bg-secondary);
  color: var(--text-secondary);
  font-style: italic;
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
}

.preview-content pre {
  margin: 16px 0;
  padding: 16px;
  background: var(--nord0);
  border-radius: var(--radius-md);
  overflow-x: auto;
  border: 1px solid var(--bg-tertiary);
}

.preview-content pre code {
  background: none;
  padding: 0;
  font-size: 14px;
  color: var(--nord6);
  line-height: 1.5;
}

.preview-content a {
  color: var(--accent-primary);
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-color 0.2s ease;
}

.preview-content a:hover {
  border-bottom-color: var(--accent-primary);
}

.preview-content del {
  color: var(--text-muted);
  opacity: 0.7;
}

.preview-content p {
  margin-bottom: 12px;
  line-height: 1.7;
}

.preview-content h1,
.preview-content h2,
.preview-content h3 {
  margin-top: 24px;
  margin-bottom: 12px;
}

.preview-content h1:first-child,
.preview-content h2:first-child,
.preview-content h3:first-child {
  margin-top: 0;
}

.empty-preview {
  color: var(--text-muted);
  font-style: italic;
  text-align: center;
  margin-top: 40px;
}

.message {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin: 16px 32px;
  padding: 16px;
  border-radius: var(--radius-md);
  border: 1px solid;
}

.message.error {
  background: rgba(191, 97, 106, 0.1);
  border-color: var(--accent-error);
  color: var(--accent-error);
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

/* 响应式设计 */
@media (max-width: 768px) {
  .editor-header {
    padding: 16px 20px;
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .header-left {
    gap: 16px;
  }
  
  .editor-title h1 {
    font-size: 20px;
  }
  
  .title-input {
    font-size: 24px;
  }
  
  .content-editor {
    padding: 0 20px 20px;
  }
  
  .title-input-wrapper {
    padding: 16px 20px 0;
  }
  
  .editor-main.split-view {
    flex-direction: column;
  }
  
  .preview-panel {
    min-height: 200px;
  }
  
  .message {
    margin: 16px 20px;
  }
}

@media (max-width: 480px) {
  .toolbar-group {
    gap: 4px;
  }
  
  .toolbar-btn {
    width: 32px;
    height: 32px;
  }
  
  .content-textarea {
    font-size: 14px;
  }
  
  .preview-content {
    font-size: 14px;
  }
}
</style>