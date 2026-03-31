#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""配置管理：加载并合并用户配置."""

import os
from typing import Dict, Any

try:
    import yaml

    HAS_YAML = True
except ImportError:
    HAS_YAML = False


def load_config(path: str | None = None) -> Dict[str, Any]:
    """加载配置文件.

    若未指定 path，则按以下顺序查找：
    1. 当前工作目录下的 config.yaml
    2. 脚本所在目录（项目根目录）下的 config.yaml
    3. 用户目录 ~/.cc-hooks-notify/config.yaml
    """
    if path is None:
        pkg_dir = os.path.dirname(os.path.abspath(__file__))
        candidates = [
            os.path.join(os.getcwd(), "config.yaml"),
            os.path.join(pkg_dir, "..", "config.yaml"),
            os.path.expanduser("~/.cc-hooks-notify/config.yaml"),
        ]
    else:
        candidates = [path]

    for p in candidates:
        real_path = os.path.abspath(p)
        if os.path.exists(real_path):
            if not HAS_YAML:
                raise RuntimeError("请安装 PyYAML: pip install pyyaml")
            with open(real_path, "r", encoding="utf-8") as f:
                return yaml.safe_load(f) or {}

    return {}
