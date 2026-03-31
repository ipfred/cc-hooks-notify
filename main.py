#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Claude Code Notification Hook 入口.

ettings.json 配置示例:
{
  "hooks": {
    "Notification": [
      {
        "matcher": "idle_prompt|permission_prompt",
        "hooks": [
          { "type": "command", "command": "python /path/to/cc-hooks-notify/notify.py" }
        ]
      }
    ]
  }
}
"""

import sys
import os

# 将项目根目录加入路径，确保可直接运行
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))


from parser import parse_stdin
from notifier import notify


def main():
    # parsed = parse_stdin()
    # 默认事件类型为 notification，交给 notifier 内部再细化
    parsed = {'session_id': '216fcdf9-74e0-4ccc-b401-7291a25871e9', 'transcript_path': 'C:\\Users\\admin\\.claude\\projects\\E--my-work-github-pro-cc-hooks-notify\\216fcdf9-74e0-4ccc-b401-7291a25871e9.jsonl', 'cwd': 'E:\\my_work\\github_pro\\cc-hooks-notify', 'permission_mode': 'bypassPermissions', 'hook_event_name': 'Stop', 'stop_hook_active': False, 'last_assistant_message': '所有验证通过！钉钉 payload 包含正确的 UTF-8 字节。\n\n## 总结\n\nstdin 中文乱码问题已修复。改动如下：\n\n### `cc_hooks_notify/parser.py`\n- 新增 `_is_valid_utf8()` 检测有效 UTF-8 字节\n- 新增 `_decode_with_fallback()` 智能解码（优先 UTF-8，自动检测乱码）\n- `parse_stdin()` 现在使用 `sys.stdin.buffer.read()` 读取原始字节，避免 Python 用系统编码（GBK）自动解码\n\n### 关键修复点\n```python\n# 以前：Python 用 gbk 解码，中文变乱码\nraw = sys.stdin.read()  # GBK 解码 → 乱码\n\n# 现在：直接读字节，智能选择编码\nraw_bytes = sys.stdin.buffer.read()\nraw = _decode_with_fallback(raw_bytes)  # UTF-8 优先 → 正确\n```\n\n### 验证结果\n- ✅ `sys.stdin.buffer` 读取原始字节\n- ✅ 正确识别并解码 UTF-8 编码的中文\n- ✅ 日志文件写入正确 UTF-8 字节（无替换字符）\n- ✅ 钉钉 payload JSON 编码正确\n\n控制台显示乱码是 Windows GBK 终端限制，不影响实际数据。钉钉收到的消息中文应该正常显示。如果钉钉仍显示乱码，可能是钉钉客户端的问题，请检查钉钉消息的实际显示效果。'}
    notify("notification", parsed)
    # 必须以 0 退出，Claude Code 认为 hooks 执行成功
    sys.exit(0)


if __name__ == "__main__":
    main()
