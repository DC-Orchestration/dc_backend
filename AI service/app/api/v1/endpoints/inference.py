from fastapi import APIRouter
from pydantic import BaseModel
from app.services.agent.agent import run_comparator

router = APIRouter()

@router.post("/compare/")
def use_agent(request: str, original: str):
    result = run_comparator(request, original)
    return {"response": result}

@router.get("/health")
def health_check():
    return {"status": "ok"}

