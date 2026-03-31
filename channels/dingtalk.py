#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""钉钉机器人通知渠道."""

import json
import time
import hmac
import hashlib
import base64
import urllib.parse
from typing import Dict, Any

from .base import BaseChannel


class DingtalkChannel(BaseChannel):
    """钉钉机器人通知渠道."""

    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        self.webhook = config.get("webhook", "")
        self.secret = config.get("secret", "")

    def send(self, data: Dict[str, Any]) -> bool:
        if not self.is_enabled() or not self.webhook:
            self.logger.warning("钉钉未启用或 webhook 为空")
            return False

        title = data.get("title", "Claude Code 通知")
        body = data.get("body", "")
        event = data.get("event", "notification")

        self.logger.info(f"[钉钉] 原始数据: title={title!r}, event={event!r}, body_len={len(body)}")

        emojis = {
            "permission": "🔐",
            "idle": "💤",
            "stop": "✅",
            "notification": "🔔",
        }
        emoji = emojis.get(event, "🔔")

        # 使用 text 格式避免 Markdown 渲染问题
        use_markdown = self.config.get("markdown", True)

        if use_markdown:
            # Markdown 格式
            markdown_text = f"#### {emoji} {title}\n\n{body}"
            payload = {
                "msgtype": "markdown",
                "markdown": {
                    "title": title,
                    "text": markdown_text,
                },
            }
        else:
            # 纯文本格式
            text_content = f"{emoji} {title}\n\n{body}"
            payload = {
                "msgtype": "text",
                "text": {
                    "content": text_content,
                },
            }

        self.logger.info(f"[钉钉] 发送 payload: {json.dumps(payload, ensure_ascii=False)!r}")

        url = self._sign_url()
        ok, text = self._http_post(url, payload)

        self.logger.info(f"[钉钉] HTTP 结果: ok={ok}, response={text!r}")

        if not ok:
            return False
        try:
            result = json.loads(text)
            self.logger.info(f"[钉钉] 解析结果: {result}")
            if result.get("errcode") == 0:
                return True
            self.logger.error(f"钉钉返回错误: {result}")
            return False
        except Exception as e:
            self.logger.error(f"解析钉钉响应失败: {e}, text={text!r}")
            return True

    def _sign_url(self) -> str:
        if not self.secret:
            return self.webhook
        timestamp = str(round(time.time() * 1000))
        secret_enc = self.secret.encode("utf-8")
        string_to_sign = f"{timestamp}\n{self.secret}"
        string_to_sign_enc = string_to_sign.encode("utf-8")
        hmac_code = hmac.new(secret_enc, string_to_sign_enc, digestmod=hashlib.sha256).digest()
        sign = urllib.parse.quote_plus(base64.b64encode(hmac_code))
        return f"{self.webhook}&timestamp={timestamp}&sign={sign}"
