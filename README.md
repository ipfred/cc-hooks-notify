# cc-notify

轻量级命令行 hook 通知工具，支持 Codex、Claude Code、Droid/Factory，将任务停止、权限请求等事件推送到钉钉或飞书。

## 安装

```bash
npm install -g @fredzhang/cc-notify
```

安装时 npm 会从 GitHub Releases 下载当前平台的 Go 二进制。安装后确认命令可用：

```bash
cc-notify version
```

## 配置通知

交互式生成配置：

```bash
cc-notify config init
```

配置默认保存到：

```text
~/.cc-notify/config.json
```

查看脱敏配置：

```bash
cc-notify config show
```

发送测试通知：

```bash
cc-notify test all
cc-notify test dingtalk
cc-notify test feishu
```

配置结构：

```json
{
  "log_level": "INFO",
  "channels": {
    "dingtalk": {
      "webhook": "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN",
      "secret": "",
      "markdown": true
    },
    "feishu": {
      "webhook": "https://open.feishu.cn/open-apis/bot/v2/hook/YOUR_HOOK",
      "secret": "",
      "card": true
    }
  },
  "events": {
    "notification": {
      "channels": ["dingtalk", "feishu"]
    },
    "stop": {
      "channels": ["dingtalk", "feishu"]
    }
  }
}
```

也可以通过环境变量覆盖 webhook：

```bash
export CC_NOTIFY_DINGTALK_WEBHOOK="https://oapi.dingtalk.com/robot/send?access_token=..."
export CC_NOTIFY_FEISHU_WEBHOOK="https://open.feishu.cn/open-apis/bot/v2/hook/..."
```

## Codex 插件安装

安装 Codex hooks 插件：

```bash
cc-notify install codex
```

这个命令会：

- 在 `~/.cc-notify/codex-marketplace/` 生成本地 Codex plugin。
- 添加本地 marketplace：`cc-notify-local`。
- 安装并启用 `cc-notify@cc-notify-local`。

然后新开一个 Codex 会话，运行：

```text
/hooks
```

检查并 trust 新 hook。Codex 插件只包含 hooks 配置，实际执行命令是：

```bash
cc-notify hook codex
```

不会写入 `~/.codex/hooks.json`。

## Claude Code / Droid 手动 hook 命令

Claude Code hook command：

```bash
cc-notify hook claude
```

Droid/Factory hook command：

```bash
cc-notify hook droid
```

示例：

```json
{
  "hooks": {
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "timeout": 10,
            "command": "cc-notify hook claude"
          }
        ]
      }
    ],
    "Notification": [
      {
        "matcher": "permission_prompt|idle_prompt",
        "hooks": [
          {
            "type": "command",
            "timeout": 10,
            "command": "cc-notify hook claude"
          }
        ]
      }
    ]
  }
}
```

## 本地开发

```bash
go test ./...
go build -o dist/cc-notify ./cmd/cc-notify
node npm/bin/cc-notify.js version
npm pack --dry-run
```

跳过 npm postinstall 下载：

```bash
CC_NOTIFY_SKIP_DOWNLOAD=1 npm install
```

## 通知渠道

### 钉钉

1. 打开钉钉群，进入群设置。
2. 添加自定义机器人。
3. 安全设置可选择加签，复制 webhook 和 secret。

### 飞书

1. 打开飞书群，进入群机器人设置。
2. 添加自定义机器人。
3. 复制 webhook；如开启签名校验，同时复制 secret。

## License

MIT
