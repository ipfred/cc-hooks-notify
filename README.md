# cc-hooks-notify

轻量级 Claude Code 通知 Hooks，支持钉钉和飞书。

## 特性

- **零侵入配置**：`.claude/settings.json` 配置固定后不再修改
- **项目级灵活控制**：通过 `config.yaml` 切换通知渠道、样式和开关
- **渠道可扩展**：基于抽象基类，新增渠道只需继承 `BaseChannel`
- **标准库优先**：HTTP 发送优先使用 `requests`，无依赖时自动回退到 `urllib`

## 目录结构

```
cc-hooks-notify/
├── cc_hooks_notify/
│   ├── channels/          # 通知渠道
│   │   ├── base.py        # 抽象基类
│   │   ├── dingtalk.py    # 钉钉
│   │   ├── feishu.py      # 飞书
│   │   └── __init__.py    # 渠道注册
│   ├── config.py          # 配置加载
│   ├── notifier.py        # 核心通知逻辑
│   └── parser.py          # Claude Code JSON 解析
├── notify.py              # Notification Hook 入口
├── complete.py            # Stop Hook 入口
├── config.yaml            # 用户配置文件
└── requirements.txt
```

## 快速开始

### 1. 安装依赖

```bash
cd cc-hooks-notify
pip install -r requirements.txt
```

### 2. 配置 Claude Code settings.json

将以下配置写入项目根目录的 `.claude/settings.json`（或全局 `~/.claude/settings.json`）：

```json
"hooks": {
    "TaskCompleted": [
      {
        "hooks": [
          { "type": "command", "timeout": 10, "async": true,
            "command": "python /e/my_work/github_pro/cc-hooks-notify/cc_hooks_notify/main.py"}
        ]
      }
    ],
    "Stop": [
      {
        "hooks": [
          { "type": "command", "timeout": 10, "async": true,
            "command": "python /e/my_work/github_pro/cc-hooks-notify/cc_hooks_notify/main.py"}
        ]
      }
    ],
    "Notification": [
      {
        "hooks": [
          { "type": "command", "timeout": 10, "async": true,
            "command": "python /e/my_work/github_pro/cc-hooks-notify/cc_hooks_notify/main.py"}
        ]
      }
    ]
  },
```

> 请把 `python /absolute/path/to/cc-hooks-notify/notify.py` 替换为实际绝对路径。

### 3. 修改项目配置

编辑 `cc-hooks-notify/config.yaml`：

```yaml
channels:
  dingtalk:
    enabled: true
    webhook: "https://oapi.dingtalk.com/robot/send?access_token=xxx"
    secret: ""

  feishu:
    enabled: true
    webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"
    secret: ""
```

重启 Claude Code 即可生效。

## 扩展渠道

新增渠道只需三步：

1. 在 `cc_hooks_notify/channels/` 下新建文件，继承 `BaseChannel`
2. 实现 `send(self, data)` 方法
3. 在 `cc_hooks_notify/channels/__init__.py` 的 `CHANNEL_REGISTRY` 中注册

示例：

```python
# cc_hooks_notify/channels/wecom.py
from .base import BaseChannel

class WecomChannel(BaseChannel):
    def send(self, data):
        # 实现企业微信发送逻辑
        return True
```

```python
# cc_hooks_notify/channels/__init__.py
from .wecom import WecomChannel

CHANNEL_REGISTRY = {
    "dingtalk": DingtalkChannel,
    "feishu": FeishuChannel,
    "wecom": WecomChannel,
}
```

## License

MIT


后续工作计划
1. 开发一个plugin的形式安装
2. 更灵活的配置使用方式
3. 支持更新的方式