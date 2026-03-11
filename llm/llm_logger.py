
import logging

# 创建 logger
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

# 创建 handler（输出到控制台）
console_handler = logging.StreamHandler()
console_handler.setLevel(logging.INFO)

# 设置日志输出格式
formatter = logging.Formatter('%(asctime)s [%(levelname)s] %(name)s: %(message)s')
console_handler.setFormatter(formatter)

# 添加 handler 到 logger
if not logger.hasHandlers():
    logger.addHandler(console_handler)