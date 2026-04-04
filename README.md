# cc-hooks-notify

轻量级 Claude Code 通知插件，支持钉钉和飞书。任务完成、停止或需要权限时自动推送消息。

## 特性

- **Plugin 模式**：支持 Claude Code Plugin 标准，一条命令安装
- **可视化配置**：通过插件设置界面直接配置 webhook，无需编辑配置文件
- **多事件支持**：支持 Stop、Notification、TaskCompleted 等事件
- **渠道可扩展**：基于抽象基类，新增渠道只需继承 `BaseChannel`

## 安装

### 方式一：从 GitHub 安装（推荐）

```bash
# 1. 添加 marketplace
/plugin marketplace add https://github.com/ipfred/cc-hooks-notify

# 2. 安装插件
/plugin install cc-hooks-notify@cc-hooks-notify --scope user

# 3. 重新加载
/reload-plugins
```

### 方式二：本地开发安装

```bash
# 1. 克隆仓库
git clone https://github.com/ipfred/cc-hooks-notify.git
cd cc-hooks-notify

# 2. 安装依赖
pip install pyyaml requests

# 3. 添加本地 marketplace
/plugin marketplace add $(pwd)

# 4. 安装插件
/plugin install cc-hooks-notify@local-plugins --scope user

# 5. 重新加载
/reload-plugins
```

### 方式三：使用 --plugin-dir 测试

无需安装，直接加载：

```bash
claude --plugin-dir /path/to/cc-hooks-notify
```

## 配置

### 通过插件界面配置（推荐）

安装后运行 `/plugin`，找到 `cc-hooks-notify`，点击 ⚙️ 设置图标：

| 配置项 | 说明 |
|--------|------|
| 启用钉钉通知 | 是否启用钉钉消息推送 |
| 钉钉 Webhook URL | 钉钉机器人的 Webhook 地址 |
| 钉钉签名密钥 | 钉钉机器人的加签密钥（可选） |
| 启用飞书通知 | 是否启用飞书消息推送 |
| 飞书 Webhook URL | 飞书机器人的 Webhook 地址 |
| 飞书签名密钥 | 飞书机器人的加签密钥（可选） |

### 获取 Webhook 地址

**钉钉**：
1. 打开钉钉群 → 群设置 → 智能群助手 → 添加机器人
2. 选择「自定义」机器人
3. 安全设置选择「加签」，复制密钥
4. 复制 Webhook 地址

**飞书**：
1. 打开飞书群 → 设置 → 群机器人 → 添加机器人
2. 选择「自定义机器人」
3. 复制 Webhook 地址

## 通知事件

| 事件 | 触发时机 |
|------|----------|
| Stop | Claude 完成响应时 |
| Notification | 需要用户确认权限、空闲提示时 |
| TaskCompleted | 任务完成时 |

## 目录结构

```
cc-hooks-notify/
├── .claude-plugin/
│   ├── plugin.json        # Plugin 元数据和配置项
│   └── marketplace.json   # Marketplace 配置
├── hooks/
│   ├── hooks.json         # Hook 事件配置
│   └── scripts/
│       └── notify.sh      # 入口脚本
├── channels/              # 通知渠道
│   ├── base.py            # 抽象基类
│   ├── dingtalk.py        # 钉钉
│   ├── feishu.py          # 飞书
│   └── __init__.py        # 渠道注册
├── main.py                # 统一入口
├── notifier.py            # 核心通知逻辑
├── parser.py              # Claude Code JSON 解析
├── config.py              # 配置加载
└── config.yaml.example    # 配置模板（可选）
```

## 扩展渠道

新增渠道只需三步：

1. 在 `channels/` 下新建文件，继承 `BaseChannel`
2. 实现 `send(self, data)` 方法
3. 在 `channels/__init__.py` 的 `CHANNEL_REGISTRY` 中注册

示例：

```python
# channels/wecom.py
from .base import BaseChannel

class WecomChannel(BaseChannel):
    def send(self, data):
        # 实现企业微信发送逻辑
        return True
```

```python
# channels/__init__.py
from .wecom import WecomChannel

CHANNEL_REGISTRY = {
    "dingtalk": DingtalkChannel,
    "feishu": FeishuChannel,
    "wecom": WecomChannel,
}
```

## 故障排查

查看日志：

```bash
# 日志文件位置
cat logs/cc_hooks_notify.log

# 或查看最新日志
tail -f logs/cc_hooks_notify.log
```

常见问题：

1. **没有收到通知**
   - 检查 `/plugin` 中配置是否正确
   - 确认 webhook URL 是否有效
   - 查看日志是否有错误

2. **插件加载失败**
   - 确保已安装依赖：`pip install pyyaml requests`
   - 运行 `/reload-plugins` 重新加载

## License

MIT
