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

    环境变量格式: CLAUDE_PLUGIN_OPTION_<KEY>
    例如: CLAUDE_PLUGIN_OPTION_DINGTALK_WEBHOOK
    """
    config = {"channels": {}}

    # 前缀 (Claude Code plugin 标准格式)
    prefix = "CLAUDE_PLUGIN_OPTION_"

    # 钉钉配置
    dingtalk_enabled = os.environ.get(f"{prefix}DINGTALK_ENABLED", "").lower()
    if dingtalk_enabled in ("true", "1", "yes"):
        config["channels"]["dingtalk"] = {
            "enabled": True,
            "webhook": os.environ.get(f"{prefix}DINGTALK_WEBHOOK", ""),
            "secret": os.environ.get(f"{prefix}DINGTALK_SECRET", ""),
        }

    # 飞书配置
    feishu_enabled = os.environ.get(f"{prefix}FEISHU_ENABLED", "").lower()
    if feishu_enabled in ("true", "1", "yes"):
        config["channels"]["feishu"] = {
            "enabled": True,
            "webhook": os.environ.get(f"{prefix}FEISHU_WEBHOOK", ""),
            "secret": os.environ.get(f"{prefix}FEISHU_SECRET", ""),
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
                os.path.join(os.getcwd(), "config.yaml"),
                os.path.join(pkg_dir, "config.json"),
                os.path.join(pkg_dir, "config.yaml"),
                os.path.expanduser("~/.cc-hooks-notify/config.json"),
                os.path.expanduser("~/.cc-hooks-notify/config.yaml"),
            ]
        )
    else:
        candidates = [path]

    yaml_candidates_without_parser = []
    for p in candidates:
        real_path = os.path.abspath(p)
        if os.path.exists(real_path):
            lower_path = real_path.lower()
            if lower_path.endswith(".json"):
                with open(real_path, "r", encoding="utf-8") as f:
                    return json.load(f) or {}

            if lower_path.endswith((".yaml", ".yml")):
                if not HAS_YAML:
                    yaml_candidates_without_parser.append(real_path)
                    continue
                with open(real_path, "r", encoding="utf-8") as f:
                    return yaml.safe_load(f) or {}

            # 未知扩展名：优先按 JSON 解析，失败后尝试 YAML
            try:
                with open(real_path, "r", encoding="utf-8") as f:
                    return json.load(f) or {}
            except Exception:
                if not HAS_YAML:
                    raise RuntimeError(
                        f"无法解析配置文件: {real_path}。请使用 JSON 格式，或安装 PyYAML: pip install pyyaml"
                    )
                with open(real_path, "r", encoding="utf-8") as f:
                    return yaml.safe_load(f) or {}

    if yaml_candidates_without_parser:
        raise RuntimeError(
            "检测到 YAML 配置文件但未安装 PyYAML。"
            "请改用 config.json（推荐，无第三方依赖）或安装 PyYAML: pip install pyyaml"
        )

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
