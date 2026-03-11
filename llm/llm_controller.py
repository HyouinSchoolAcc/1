import asyncio
import time
from tqdm.asyncio import tqdm_asyncio
from openai import (
    AsyncOpenAI,
    OpenAIError,
)
import os
from typing import Any

from .llm_logger import logger


class AsyncOpenAIController:
    """
    An asynchronous controller for OpenAI's ChatCompletion API that processes
    multiple prompts in parallel. Implements async context management for proper client closure.
    """
    def __init__(self, model: str = "gpt-4.1-nano", api_key: str | None = None):
        
        self.model = model
        
        if api_key is None:
            # api_key = os.getenv('OPENAI_API_KEY')
            api_key = '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'
        if api_key is None:
            raise ValueError("OpenAI API key not found. Set OPENAI_API_KEY environment variable or pass it during initialization.")

        try:
            # Initialize client here, but don't connect yet
            self._api_key = api_key
            self.client = None # Will be initialized in __aenter__
        except ImportError:
            raise ImportError("OpenAI package not found. Install it with: pip install openai")
        except Exception as e:
            # Catch potential issues during basic setup if any
            raise RuntimeError(f"Failed during initial setup: {e}")

    async def __aenter__(self):
        """Initialize the async client when entering the context."""
        try:
            self.client = AsyncOpenAI(api_key=self._api_key)
        except Exception as e:
            raise RuntimeError(f"Failed to initialize AsyncOpenAI client: {e}")
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Close the async client when exiting the context."""
        if self.client:
            logger.info("Closing OpenAI client...")
            try:
                await self.client.close()
                logger.info("OpenAI client closed.")
            except Exception as e:
                logger.error(f"Error closing OpenAI client: {e}")
        # Return False to propagate exceptions, True to suppress them (usually False is desired)
        return False

    async def _get_single_completion_with_timing(
        self,
        system_prompt: str,
        prompt: str,
        response_format: dict,
        temperature: float,
        max_tokens: int
    ) -> dict[str, Any]:

        if not self.client:
             raise RuntimeError("Client not initialized. Use 'async with AsyncOpenAIController(...)'")

        start_time = time.perf_counter()
        result = {}
        try:
            if max_tokens:
                response = await self.client.chat.completions.create(
                    model=self.model,
                    messages=[
                        {"role": "system", "content": system_prompt},
                        {"role": "user", "content": prompt}
                    ],
                    response_format=response_format,
                    temperature=temperature,
                    max_tokens=max_tokens
                )
            else:
                response = await self.client.chat.completions.create(
                    model=self.model,
                    messages=[
                        {"role": "system", "content": system_prompt},
                        {"role": "user", "content": prompt}
                    ],
                    response_format=response_format,
                )
            content = response.choices[0].message.content
            result = {"content": content, "response": response}
        except OpenAIError as e:
            # Handle API errors specifically
            logger.error(f"OpenAI API Error for prompt '{prompt[:50]}...': {e}")
            result = {"error": f"OpenAI API Error: {e}"}
        except Exception as e:
            # Handle other potential errors (network issues, etc.)
            logger.error(f"Unexpected error for prompt '{prompt[:50]}...': {e}")
            result = {"error": f"An unexpected error occurred: {e}"}
        finally:
            end_time = time.perf_counter()
            result["duration"] = end_time - start_time
            return result

    
    async def get_completion(
        self,
        system_prompt: str,
        prompt: str,
        response_format: dict | None = None,
        temperature: float = 0.7,
        max_tokens: int = None,
    ) -> str:
        """
        处理单个 prompt，返回内容的字符串
        """
        if not self.client:
            raise RuntimeError("Client not initialized. Use 'async with AsyncOpenAIController(...)'")
        # 如果用户不传 response_format，默认为空 dict（即返回原始 text）
        rf = response_format
        result = await self._get_single_completion_with_timing(
            system_prompt=system_prompt,
            prompt=prompt,
            response_format=rf,
            temperature=temperature,
            max_tokens=max_tokens
        )
        return result.get("content", None) if result else None
