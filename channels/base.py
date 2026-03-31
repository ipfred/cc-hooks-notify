#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""通知渠道抽象基类，提供通用 HTTP 发送能力."""

import abc
import json
import logging
import urllib.request
import urllib.error
from typing import Dict, Any

# 可选依赖：requests
try:
    import requests

    HAS_REQUESTS = True
except ImportError:
    HAS_REQUESTS = False


class BaseChannel(abc.ABC):
    """通知渠道基础类."""

    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.logger = logging.getLogger(f"cc_hooks_notify.{self.__class__.__name__}")

    @abc.abstractmethod
    def send(self, data: Dict[str, Any]) -> bool:
        """发送通知，子类必须实现."""
        pass

    def is_enabled(self) -> bool:
        return self.config.get("enabled", False)

    def _http_post(self, url: str, payload: Dict[str, Any], headers: Dict[str, str] | None = None, timeout: int = 10) -> tuple[bool, str]:
        """发送 HTTP POST 请求，优先使用 requests，回退到 urllib.

        返回 (是否成功, 响应体文本)
        """
        _headers = {"Content-Type": "application/json"}
        if headers:
            _headers.update(headers)

        if HAS_REQUESTS:
            try:
                resp = requests.post(url, headers=_headers, json=payload, timeout=timeout)
                text = resp.text
                if resp.status_code == 200:
                    return True, text
                self.logger.error(f"[{self.__class__.__name__}] HTTP {resp.status_code}: {text}")
                return False, text
            except Exception as e:
                self.logger.error(f"[{self.__class__.__name__}] 请求异常: {e}")
                return False, str(e)
        else:
            try:
                req = urllib.request.Request(
                    url,
                    data=json.dumps(payload).encode("utf-8"),
                    headers=_headers,
                    method="POST",
                )
                with urllib.request.urlopen(req, timeout=timeout) as resp:
                    text = resp.read().decode("utf-8")
                    return True, text
            except urllib.error.HTTPError as e:
                text = e.read().decode("utf-8")
                self.logger.error(f"[{self.__class__.__name__}] HTTP {e.code}: {text}")
                return False, text
            except Exception as e:
                self.logger.error(f"[{self.__class__.__name__}] 请求异常: {e}")
                return False, str(e)
