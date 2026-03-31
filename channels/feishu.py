#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""飞书机器人通知渠道."""

import json
import time
import hmac
import hashlib
import base64
from typing import Dict, Any

from .base import BaseChannel


class FeishuChannel(BaseChannel):
    """飞书机器人通知渠道."""

    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        self.webhook = config.get("webhook", "")
        self.secret = config.get("secret", "")

    def send(self, data: Dict[str, Any]) -> bool:
        if not self.is_enabled() or not self.webhook:
            self.logger.warning("飞书未启用或 webhook 为空")
            return False

        title = data.get("title", "Claude Code 通知")
        body = data.get("body", "")
        event = data.get("event", "notification")

        self.logger.info(f"[飞书] 原始数据: title={title!r}, event={event!r}, body_len={len(body)}")

        emojis = {
            "permission": "🔐",
            "idle": "💤",
            "stop": "✅",
            "notification": "🔔",
        }
        emoji = emojis.get(event, "🔔")

        colors = {
            "permission": "orange",
            "idle": "blue",
            "stop": "green",
            "notification": "blue",
        }
        color = colors.get(event, "blue")

        # 支持配置选择消息格式: card(卡片) 或 text(纯文本)
        use_card = self.config.get("card", True)

        if use_card:
            # 卡片消息格式
            payload = {
                "msg_type": "interactive",
                "card": {
                    "config": {"wide_screen_mode": True},
                    "header": {
                        "template": color,
                        "title": {
                            "content": f"{emoji} {title}",
                            "tag": "plain_text",
                        },
                    },
                    "elements": [
                        {
                            "tag": "div",
                            "text": {
                                "content": body,
                                "tag": "lark_md",
                            },
                        }
                    ],
                },
            }
        else:
            # 纯文本消息格式
            payload = {
                "msg_type": "text",
                "content": {
                    "text": f"{emoji} {title}\n\n{body}",
                },
            }

        self.logger.info(f"[飞书] 发送 payload: {json.dumps(payload, ensure_ascii=False)!r}")

        if self.secret:
            timestamp = str(int(time.time()))
            sign = self._sign(timestamp)
            payload["timestamp"] = timestamp
            payload["sign"] = sign

        ok, text = self._http_post(self.webhook, payload)

        self.logger.info(f"[飞书] HTTP 结果: ok={ok}, response={text!r}")

        if not ok:
            return False
        try:
            result = json.loads(text)
            self.logger.info(f"[飞书] 解析结果: {result}")
            if result.get("code") == 0:
                return True
            self.logger.error(f"飞书返回错误: {result}")
            return False
        except Exception as e:
            self.logger.error(f"解析飞书响应失败: {e}, text={text!r}")
            return True

    def _sign(self, timestamp: str) -> str:
        string_to_sign = f"{timestamp}\n{self.secret}"
        hmac_code = hmac.new(
            self.secret.encode("utf-8"),
            string_to_sign.encode("utf-8"),
            digestmod=hashlib.sha256,
        ).digest()
        return base64.b64encode(hmac_code).decode("utf-8")
