import json

from .llm_logger import logger

from .llm_controller import AsyncOpenAIController


async def run_openai_async(
    prompt: str, 
    json_schema: dict = None, 
    temperature: float = 0.7, 
    model: str = 'gpt-4.1',
    system_prompt: str = 'You are a helpful assistant.',
    **kwargs
) -> dict:
    """
    Runs a single OpenAI API call.

    Args:
        model (str): The model to use.
        prompt (str): The prompt to process.
        response_format (dict, optional): Response format dictionary. Defaults to None.
        temperature (float, optional): Temperature setting for the model. Defaults to 0.7.
        **kwargs: Additional arguments to pass to the controller.

    Returns:
        result: The result of the API call.
    """
    async with AsyncOpenAIController(model=model) as controller:
        result = await controller.get_completion(
            system_prompt=system_prompt,
            prompt=prompt,
            response_format={"type": "json_schema", "json_schema": json_schema} if json_schema else None,
            temperature=temperature,
            **kwargs
        )
    
    if json_schema:
        try:
            result = json.loads(result)
        except Exception as e:
            logger.error(f"Failed to parse JSON response: {e}")
            logger.debug(f"Raw response: {result}")
            return result
            
    return result
