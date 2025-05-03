from fastapi import APIRouter
from app.api.v1.endpoints import inference

router = APIRouter()
router.include_router(inference.router, tags=["Agent"])
