#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""解析 Claude Code 传入的 JSON 数据."""

import json
import sys
from typing import Dict, Any


def parse_stdin() -> Dict[str, Any]:
    """从标准输入读取并解析 JSON（Claude Code 使用 UTF-8 编码）."""
    try:
        # Windows 下 sys.stdin 默认用系统编码（gbk），需要直接读字节再用 UTF-8 解码
        if hasattr(sys.stdin, 'buffer'):
            raw = sys.stdin.buffer.read().decode('utf-8', errors='replace').strip()
        else:
            raw = sys.stdin.read().strip()
    except Exception as e:
        return {"parse_error": str(e)}

    if not raw:
        return {}
    try:
        return json.loads(raw)
    except json.JSONDecodeError as e:
        return {"raw_input": raw, "parse_error": str(e)}


def get_project_name(cwd: str) -> str:
    """提取项目名（取路径最后一部分）."""
    if not cwd:
        return "unknown"
    cwd_parts = cwd.replace("\\", "/").rstrip("/").split("/")
    cwd_short = "/".join(cwd_parts[-2:]) if len(cwd_parts) >= 2 else cwd_parts[-1]
    return cwd_short


def get_message_preview(msg: str, limit: int = 200) -> str:
    """截取消息预览."""
    if not msg:
        return ""
    return msg if len(msg) <= limit else msg[: limit - 3] + "..."


def sanitize_for_dingtalk(text: str) -> str:
    """清理钉钉 Markdown 中的特殊字符，避免渲染异常."""
    if not text:
        return ""
    for old, new in {"`": "'", "**": "*", "~~": "", "[": "【", "]": "】"}.items():
        text = text.replace(old, new)
    return text
