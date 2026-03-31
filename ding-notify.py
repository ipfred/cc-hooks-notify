import sys
import json
import urllib.request
import urllib.error
import logging
from datetime import datetime

# 1. 配置 logging 模块
LOG_FILE = '/e/my_work/github_pro/cc-hooks-notify/claude-hook-debug.log'

logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S',
    filename=LOG_FILE,
    filemode='a',
    encoding='utf-8'
)
logger = logging.getLogger(__name__)

def main():
    logger.info("=" * 50)

    # 2. 从标准输入读取 Claude Code 传来的 JSON 数据
    input_data = sys.stdin.read().strip()
    logger.debug(f"收到 Claude Code 原始输入: {input_data}")

    # 3. 解析 JSON 数据
    try:
        parsed_data = json.loads(input_data) if input_data else {}
    except Exception as e:
        logger.error(f"JSON 解析失败: {e}")
        parsed_data = {}

    event_name = parsed_data.get("hook_event_name", "未知事件")
    cwd = parsed_data.get("cwd", "")
    last_msg = parsed_data.get("last_assistant_message", "")

    # 提取 cwd 最后两层路径
    if cwd:
        cwd_parts = cwd.rstrip("/").split("/")
        cwd_short = "/".join(cwd_parts[-2:]) if len(cwd_parts) >= 2 else cwd_parts[-1]
    else:
        cwd_short = "未知项目"

    # 提取 last_assistant_message 最后20个字
    msg_short = last_msg[-20:] if len(last_msg) > 20 else last_msg

    # 4. 构造通知文本 (精简版)
    content = f"Claude⚡ {cwd_short}\n…{msg_short}"

    # 钉钉机器人 Webhook URL (请替换为你自己的 access_token)
    webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=a94dd07767e608afd83f43fe8e987314409668d9b1f887b7d47fe2e59c91411f"

    # 5. 构造发送给钉钉的 JSON Payload
    payload_dict = {
        "msgtype": "text",
        "text": {
            "content": content
        }
    }
    # 将字典转为 JSON 字符串并编码为 bytes (urllib 要求)
    payload_bytes = json.dumps(payload_dict).encode('utf-8')
    logger.debug(f"准备发给钉钉的请求体: {json.dumps(payload_dict, ensure_ascii=False)}")

    # 6. 发送 HTTP POST 请求 (纯原生，替代 curl)
    req = urllib.request.Request(
        webhook_url,
        data=payload_bytes,
        headers={'Content-Type': 'application/json'}
    )

    try:
        with urllib.request.urlopen(req) as response:
            response_data = response.read().decode('utf-8')
            logger.info(f"钉钉服务器返回结果: {response_data}")
    except urllib.error.HTTPError as e:
        # 捕获 HTTP 错误 (如 400, 404, 500)
        error_msg = e.read().decode('utf-8')
        logger.error(f"HTTP 请求失败: {e.code} - {error_msg}")
    except Exception as e:
        # 捕获其他异常 (如网络不通)
        logger.error(f"请求异常: {e}")

if __name__ == "__main__":
    main()
    # 必须以0退出，让 Claude Code 认为钩子执行成功
    sys.exit(0)