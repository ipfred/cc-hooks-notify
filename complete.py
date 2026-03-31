#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Claude Code Stop Hook 入口.
在 Claude 停止/会话结束时发送完成通知.

settings.json 配置示例:
{
  "hooks": {
    "Stop": [
      {
        "hooks": [
          { "type": "command", "command": "python /path/to/cc-hooks-notify/complete.py" }
        ]
      }
    ]
  }
}
"""

import sys
import os

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from cc_hooks_notify.parser import parse_stdin
from cc_hooks_notify.notifier import notify


def main():
    parsed = parse_stdin()
    notify("stop", parsed)
    sys.exit(0)


if __name__ == "__main__":
    main()
