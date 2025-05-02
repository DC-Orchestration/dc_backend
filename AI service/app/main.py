from fastapi import FastAPI
from fastapi_redis import FastAPIRedis
from app.api.v1.api_router import router as api_router
import asyncio
import redis.asyncio as redis
import json
import logging


logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()
app.include_router(api_router, prefix="/api/v1")


REDIS_URL = "redis://localhost:6379"

async def redis_listener():
    """Subscribe to Redis channel and process messages"""
    r = redis.from_url(REDIS_URL)
    pubsub = r.pubsub()
    await pubsub.subscribe("AI")
    
    logger.info("Subscribed to 'AI' channel. Waiting for messages...")
    
    async for message in pubsub.listen():
        try:
            if message["type"] == "message":
                data = json.loads(message["data"])
                logger.info(f"Received message: {data}")
                
                # Process your message here
                await process_ai_message(data)
                
        except Exception as e:
            logger.error(f"Error processing message: {e}")

async def process_ai_message(data: dict):
    """Process incoming AI messages"""
    logger.info(f"Processing AI message: {data}")
    if data.get("type") == "text_comparison":
        from app.services.agent.agent import TextComparator
        comparator = TextComparator()
        result = await comparator.compare_texts(
            data["original_text"],
            data["modified_text"]
        )
        logger.info(f"Comparison result: {result}")

@app.on_event("startup")
async def startup_event():
    """Start Redis listener when app starts"""
    asyncio.create_task(redis_listener())

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)