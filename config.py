#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""配置管理：支持插件环境变量和 config.yaml 两种配置方式."""

import os
import json
from typing import Dict, Any

try:
    import yaml

    HAS_YAML = True
except ImportError:
    HAS_YAML = False


def load_config_from_env() -> Dict[str, Any]:
    """从环境变量加载插件配置.

    环境变量格式: CLAUDE_PLUGIN_OPTION_<KEY> 或 CODEX_PLUGIN_OPTION_<KEY>
    例如: CLAUDE_PLUGIN_OPTION_DINGTALK_WEBHOOK

    规则: Webhook 地址非空即视为启用该渠道。
    """
    config = {"channels": {}}

    prefixes = (
        "CODEX_PLUGIN_OPTION_",
        "CLAUDE_PLUGIN_OPTION_",
        "CC_HOOKS_NOTIFY_",
    )

    def get_option(key: str) -> str:
        for prefix in prefixes:
            value = os.environ.get(f"{prefix}{key}", "").strip()
            if value:
                return value
        return ""

    # 钉钉配置: Webhook 非空即启用
    dingtalk_webhook = get_option("DINGTALK_WEBHOOK")
    if dingtalk_webhook:
        config["channels"]["dingtalk"] = {
            "enabled": True,
            "webhook": dingtalk_webhook,
            "secret": get_option("DINGTALK_SECRET"),
        }

    # 飞书配置: Webhook 非空即启用
    feishu_webhook = get_option("FEISHU_WEBHOOK")
    if feishu_webhook:
        config["channels"]["feishu"] = {
            "enabled": True,
            "webhook": feishu_webhook,
            "secret": get_option("FEISHU_SECRET"),
        }

    return config


def load_config_from_file(path: str | None = None) -> Dict[str, Any]:
    """从配置文件加载配置，支持 config.json / config.yaml.

    若未指定 path，则按以下顺序查找：
    1. 环境变量 CC_HOOKS_NOTIFY_CONFIG_PATH 指定的路径
    2. 当前工作目录下的 config.json / config.yaml
    3. 脚本所在目录（项目根目录）下的 config.json / config.yaml
    4. 用户目录 ~/.cc-hooks-notify/config.json / config.yaml
    """
    if path is None:
        pkg_dir = os.path.dirname(os.path.abspath(__file__))
        candidates = []

        # 优先使用环境变量指定的配置路径
        env_config = os.environ.get("CC_HOOKS_NOTIFY_CONFIG_PATH")
        if env_config:
            candidates.append(env_config)

        candidates.extend(
            [
                os.path.join(os.getcwd(), "config.json"),
                os.path.join(pkg_dir, "config.json"),
                os.path.expanduser("~/.cc-hooks-notify/config.json"),
            ]
        )
    else:
        candidates = [path]

    for p in candidates:
        real_path = os.path.abspath(p)
        if os.path.exists(real_path):
            lower_path = real_path.lower()
            if lower_path.endswith(".json"):
                with open(real_path, "r", encoding="utf-8") as f:
                    return json.load(f) or {}
    return {}


def merge_configs(env_config: Dict[str, Any], file_config: Dict[str, Any]) -> Dict[str, Any]:
    """合并环境变量配置和文件配置.

    环境变量配置优先级更高。
    """
    result = file_config.copy()

    # 合并 channels
    if "channels" not in result:
        result["channels"] = {}

    for channel, settings in env_config.get("channels", {}).items():
        if settings.get("enabled"):
            result["channels"][channel] = settings

    return result


def load_config(path: str | None = None) -> Dict[str, Any]:
    """加载配置，合并环境变量和文件配置.

    优先级：环境变量 > config.json / config.yaml
    """
    env_config = load_config_from_env()
    try:
        file_config = load_config_from_file(path)
    except RuntimeError as e:
        # 插件环境变量已提供渠道配置时，允许无 PyYAML 继续运行
        # （此时文件配置仅作为可选补充）
        if env_config.get("channels") and ("YAML" in str(e) or "PyYAML" in str(e)):
            file_config = {}
        else:
            raise

    # 如果环境变量有配置，优先使用
    if env_config.get("channels"):
        return merge_configs(env_config, file_config)

    return file_config
