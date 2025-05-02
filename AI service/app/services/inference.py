from fastapi import APIRouter
import asyncio
from agent.agent import run_comparator

router = APIRouter()  

def run_agent(original: str , query: str):
    """Run the agent with the provided query."""
    # Run tests when executed directly
    asyncio.run(run_comparator(original, query))
