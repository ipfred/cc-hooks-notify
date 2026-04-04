#!/bin/bash
# Claude Code Notification Hook Script
# 此脚本由 Claude Code plugin 系统调用

set -e

# 获取插件根目录
PLUGIN_ROOT="${CLAUDE_PLUGIN_ROOT}"

# 设置 Python 路径
export PYTHONPATH="${PLUGIN_ROOT}:${PYTHONPATH:-}"

# 插件配置通过环境变量传递，格式为 CC_HOOKS_NOTIFY_<KEY>
# 例如: CC_HOOKS_NOTIFY_DINGTALK_WEBHOOK, CC_HOOKS_NOTIFY_FEISHU_WEBHOOK

# 调用 Python 主程序
python3 "${PLUGIN_ROOT}/main.py"
