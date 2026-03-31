#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""核心通知逻辑：根据事件类型和配置，分发到各渠道."""

import os
import logging
from logging.handlers import TimedRotatingFileHandler
from typing import Dict, Any

from channels import get_channel, list_channels
from config import load_config
from parser import get_project_name, get_message_preview

logger = logging.getLogger("cc_hooks_notify")

# 日志目录: 项目根目录/logs
PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
LOG_DIR = os.path.join(PROJECT_ROOT, "logs")


def setup_logging(level: str = "INFO") -> None:
    """初始化日志：输出到控制台并按天轮转到 logs 目录."""
    log_level = getattr(logging, level.upper(), logging.INFO)
    logger.setLevel(log_level)

    # 避免重复添加 handler
    if logger.handlers:
        return

    os.makedirs(LOG_DIR, exist_ok=True)
    formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")

    # 按天轮转的日志文件，保留最近 7 天
    log_file = os.path.join(LOG_DIR, "cc-hooks-notify.log")
    file_handler = TimedRotatingFileHandler(
        log_file, when="midnight", interval=1, backupCount=7, encoding="utf-8"
    )
    file_handler.setFormatter(formatter)
    logger.addHandler(file_handler)

    # 控制台输出
    console_handler = logging.StreamHandler()
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)


def build_notification_data(event_type: str, parsed: Dict[str, Any]) -> Dict[str, Any]:
    """根据 Claude Code 传入的数据构造统一通知数据."""
    hook_event = parsed.get("hook_event_name", "")
    cwd = parsed.get("cwd", "")
    project: str = get_project_name(cwd)

    notification_type = parsed.get("notification_type", "")
    message = parsed.get("message", "")
    last_msg = parsed.get("last_assistant_message", "")

    # 优先使用 message，其次 last_assistant_message
    body_source = message or last_msg or "none msg"
    preview = get_message_preview(body_source, 100)

    titles = {
        "permission": f"Claude 需要权限确认 | {project}",
        "idle": f"Claude 正在等待 | {project}",
        "stop": f"Claude stop | {project}",
        "notification": f"Claude 通知 | {project}",
    }

    # 根据 notification_type 映射事件
    if "permission" in notification_type.lower() or notification_type.lower() == "permission_prompt":
        evt = "permission"
    elif "idle" in notification_type.lower() or notification_type.lower() == "idle_prompt":
        evt = "idle"
    elif hook_event.lower() == "stop":
        evt = "stop"
    else:
        evt = "notification"

    # 使用简单格式，避免钉钉 Markdown 渲染问题
    body_lines = []
    if preview:
        body_lines.append(f"{preview}")

    return {
        "title": titles.get(evt, titles["notification"]),
        "body": "\n".join(body_lines),
        "event": evt
    }


def notify(event_type: str, parsed: Dict[str, Any], config: Dict[str, Any] | None = None) -> None:
    """发送通知的主入口."""
    if config is None:
        config = load_config()

    log_level = config.get("log_level", "INFO")
    setup_logging(log_level)
    logger.info(f"input= {parsed}")
    data = build_notification_data(event_type, parsed)

    # 获取要发送的渠道列表
    channels_cfg = config.get("channels", {})
    default_enabled = [name for name, cfg in channels_cfg.items() if cfg.get("enabled", False)]
    event_cfg = config.get("events", {}).get(event_type) or {}
    target_channels = event_cfg.get("channels", default_enabled)

    if not target_channels:
        logger.warning("没有启用任何通知渠道，请在 config.yaml 中配置")
        return

    for name in target_channels:
        cfg = channels_cfg.get(name, {})
        channel = get_channel(name, cfg)
        if channel is None:
            logger.warning(f"未知的通知渠道: {name}")
            continue
        try:
            ok = channel.send(data)
            logger.info(f"渠道 {name}: {'成功' if ok else '失败'}")
        except Exception as e:
            logger.error(f"渠道 {name} 发送异常: {e}")
