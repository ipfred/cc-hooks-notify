"""通知渠道注册与发现."""

from .base import BaseChannel
from .dingtalk import DingtalkChannel
from .feishu import FeishuChannel

CHANNEL_REGISTRY = {
    "dingtalk": DingtalkChannel,
    "feishu": FeishuChannel,
}


def get_channel(name: str, config: dict) -> BaseChannel | None:
    """根据名称实例化对应渠道."""
    klass = CHANNEL_REGISTRY.get(name)
    if klass is None:
        return None
    return klass(config)


def list_channels() -> list[str]:
    """返回已注册渠道名称列表."""
    return list(CHANNEL_REGISTRY.keys())
