# 📔 MoodStack - 基于 Go + Wails 的本地日记管理工具

> 你的心事，只让你和本地大模型知道 💬🧠
>  安全地写日记，私密地管理，智能地理解。

------

## 🧠 项目简介

**MoodStack** 是一款跨平台桌面应用，使用 Go + Wails 构建，支持用户导入各类格式的日记文档（.txt, .docx, .pdf, .md 等），将其统一转化为 Markdown 格式展示，同时集成本地 `Ollama` 大模型，实现以下智能功能：

- 🔍 **全文关键字搜索**
- 📈 **情绪分析与变化趋势可视化**
- 🧾 **日记内容自动摘要（本地LLM，无需联网）**
- 🧑‍💻 **程序员专属技术成长追踪**
- 🔒 **绝对私密，数据全程本地运行，不上传、不联网**

------

## 💡 核心功能

| 功能                | 描述                                                 |
| ------------------- | ---------------------------------------------------- |
| 🗂️ 日记上传          | 支持导入 `.md`、`.txt`、`.docx`、`.pdf` 等格式       |
| 📝 Markdown 预览     | 将上传的日记统一转换为 Markdown 格式渲染展示         |
| 🔍 关键字搜索        | 可搜索全文内容，按关键字筛选日记                     |
| 💬 本地摘要          | 使用 Ollama 本地大模型对单篇日记生成简要总结         |
| 📊 情绪分析          | 从日记中提取情绪倾向，展示折线图/词云                |
| 🧑‍💻 技术成长追踪     | 自动识别日记中的技术关键词，构建知识成长曲线与雷达图 |
| 📁 项目开发分析      | 结合本地 Git 仓库，分析代码提交记录与开发趋势        |
| 😤 Bug情绪识别       | 检测抱怨语句与 Bug 关联，生成 "崩溃点" 热力图        |
| 🧾 周报/成长报告生成 | 每周自动总结技术关键词/情绪关键词与建议              |
| 🔐 完全离线          | 所有处理均在本地完成，隐私完全保障                   |

------

## 🛠️ 技术栈

| 技术               | 用途                                        |
| ------------------ | ------------------------------------------- |
| Go (Wails)         | 后端逻辑与桌面端应用封装                    |
| Vue3 / React       | 前端 UI（可选 Wails + Svelte/Vue/React）    |
| Tesseract / Pandoc | 文本格式统一转换                            |
| SQLite / LiteFS    | 本地轻量化数据库存储                        |
| Ollama             | 集成本地 LLM 模型（如 Mistral, LLaMA, etc） |
| Tailwind CSS       | 快速构建美观 UI 界面                        |
| Chart.js / Echarts | 情绪变化可视化分析                          |
| Git CLI / go-git   | 项目代码提交数据分析                        |

------

## 📁 目录结构

```bash
moodstack/
├── app/
│   ├── main.go           # Wails 主入口
│   ├── convert.go        # 文档转 Markdown 模块
│   ├── llm.go            # 调用 Ollama 模块
│   ├── git.go            # 本地项目分析模块
│   └── db.go             # 本地数据库逻辑
├── frontend/             # 前端（Vue/React/Svelte）
│   ├── components/
│   ├── views/
│   └── App.vue
├── static/               # 默认图标/资源
├── docs/                 # 项目文档与使用指南
├── go.mod
└── README.md             # 当前文档
```

------

## 🚀 快速开始

### 1️⃣ 安装依赖

确保你已安装：

- [Go](https://golang.org/) >= 1.21
- [Wails CLI](https://wails.io/)
- [Node.js](https://nodejs.org/) >= 18
- [Ollama](https://ollama.com/) 本地模型运行环境

```bash
# 安装 wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2️⃣ 克隆项目

```bash
git clone https://github.com/yourname/moodstack.git
cd moodstack
```

### 3️⃣ 启动 Ollama 并拉取模型

```bash
ollama run mistral
# 或者 ollama pull llama3
```

### 4️⃣ 启动应用

```bash
wails dev
```

------

## 🧪 示例截图（构想）

> 可通过以下方式展示情绪分析与成长轨迹：

- 折线图：记录每日日记情绪变化
- 词云图：展示出现频率最高的情绪/技术词汇
- 技术成长图谱：雷达图 or 关键词热度趋势图
- 崩溃热力图：展示 bug 触发点

------

## 📦 打包发布

```bash
wails build
```

打包后即可得到 Windows/macOS/Linux 下的 `.app` 或 `.exe` 文件！

------

## ✨ TODO / 待开发功能

- 支持多账户/密码保护
- 支持语音转录写日记（自动语音转文字）
- 自动技术关键词提取与学习轨迹图谱
- 情绪热力图生成 + Bug词追踪
- 暗黑模式/自定义主题切换
- 自动生成个人周报/成长总结
- VSCode 插件 or CLI 工具日记接入

------

## ❤️ 项目愿景

让日记不只是记录，更是 **自我观察与理解的镜子**。
 本地模型赋能，情绪洞察加持，成长趋势追踪，
 让每一行文字、每一次崩溃、每一个 commit，
 都成为看见自己的方式。

> ✨ 你的秘密，只属于你自己，成长却值得被记录。🧠📓

------

## 📜 License

MIT License - 属于所有需要安心表达的你 🕊️