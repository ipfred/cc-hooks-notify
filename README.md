# cc-hooks-notify

轻量级 Claude Code 通知 Hooks，支持钉钉和飞书。

## 特性

- **零侵入配置**：`.claude/settings.json` 配置固定后不再修改
- **项目级灵活控制**：通过 `config.yaml` 切换通知渠道、样式和开关
- **渠道可扩展**：基于抽象基类，新增渠道只需继承 `BaseChannel`
- **标准库优先**：HTTP 发送优先使用 `requests`，无依赖时自动回退到 `urllib`

## 目录结构

```
.
├── cc_hooks_notify/
│   ├── channels/          # 通知渠道
│   │   ├── base.py        # 抽象基类
│   │   ├── dingtalk.py    # 钉钉
│   │   ├── feishu.py      # 飞书
│   │   └── __init__.py    # 渠道注册
│   ├── main.py            # 统一入口（Notification / Stop 等 hook）
│   ├── complete.py        # Stop Hook 入口
│   ├── notifier.py        # 核心通知逻辑
│   ├── parser.py          # Claude Code JSON 解析
│   ├── config.py          # 配置加载
│   └── config.yaml.example # 用户配置模板（复制为 config.yaml 使用）
├── config.yaml.example    # 配置模板（不会被跟踪）
├── test/
│   └── ding-notify.py     # 钉钉调试脚本
├── requirements.txt
├── pyproject.toml
└── logs/                  # 日志目录（自动创建）
```

## 快速开始

### 1. 安装依赖

```bash
pip install -r requirements.txt
```

### 2. 配置 Claude Code settings.json

将以下配置写入项目根目录的 `.claude/settings.json`（或全局 `~/.claude/settings.json`）：

> 因为windows系统 claude 使用的 gitbash 路径写成样例中的那样
> 

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

### 3. 生成配置文件

根据模板 `config.yaml.example` 生成配置文件：

```bash
cp config.yaml.example config.yaml
```

> **注意**：`config.yaml.example` 是配置模板，已提交到版本控制。`config.yaml` 是实际使用的配置文件，包含敏感信息（如 webhook token），已被 `.gitignore` 忽略，不会提交到仓库。

编辑 `config.yaml`：

```yaml
# 日志配置
log_level: INFO                     # DEBUG | INFO | WARN | ERROR
# log_dir: "/path/to/custom/logs"  # 可选，默认项目根目录下的 logs/

# 通知渠道
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

> `log_dir` 取消注释后替换为自定义路径即可，日志会按天轮转、保留最近 7 天。

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
