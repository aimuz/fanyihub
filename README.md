# FanyiHub (翻译中心)

FanyiHub 是一个基于大语言模型（LLM）的智能翻译工具，使用 Go 语言和 Wails 框架开发，提供跨平台的桌面翻译体验。

## 功能特点

- 🌐 **多语言支持**：自动检测多种语言，包括中文、英语、日语、韩语等 
- 🤖 **LLM 翻译**：基于大语言模型的高质量翻译，支持多种 LLM 提供商
- ⌨️ **全局快捷键**：通过自定义快捷键快速唤起翻译窗口
- ⚙️ **灵活配置**：可配置多个 LLM 提供商，包括 OpenAI 和兼容 OpenAI API 的其他服务
- 🔄 **智能语言检测**：自动检测输入文本的语言，并选择合适的目标语言
- 🖥️ **跨平台**：支持 macOS、Windows 和 Linux 系统

## 安装

### 预构建版本

从 [Releases](https://github.com/aimuz/fanyihub/releases) 页面下载适合您系统的预构建版本。

### 从源码构建

#### 前提条件

- Go 1.21+
- Wails CLI

#### 安装步骤

1. 安装 Wails CLI:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

2. 克隆仓库:

```bash
git clone https://github.com/aimuz/fanyihub.git
cd fanyihub
```

3. 构建应用:

```bash
wails build
```

## 使用方法

1. 启动应用后，您可以通过全局快捷键（默认为 `Cmd+Shift+Space` 或 `Ctrl+Shift+Space`）唤起翻译窗口
2. 在输入框中输入需要翻译的文本
3. 应用会自动检测文本语言并翻译到目标语言
4. 您可以在设置中配置 LLM 提供商和其他选项

## 配置 LLM 提供商

FanyiHub 支持多种 LLM 提供商，包括：

1. **OpenAI API**：使用 OpenAI 的 GPT 模型
2. **兼容 OpenAI API 的服务**：如 Azure OpenAI、Claude、本地部署的大模型等

在应用设置中，您可以添加、编辑和删除 LLM 提供商配置。

## 许可证

[MIT 许可证](LICENSE)

## 贡献

欢迎贡献代码、报告问题或提出功能建议！请查看 [贡献指南](CONTRIBUTING.md) 了解更多信息。

## 联系方式

如有问题或建议，请通过 GitHub Issues 联系我们。 