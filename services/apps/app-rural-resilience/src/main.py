"""
Rural Resilience Engine
預測性可靠度分析與升級建議
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from src.ingestion import FaultLogIngestion
from src.risk_model import RiskModel

app = FastAPI(
    title="Rural Resilience Engine",
    description="Predictive resilience analysis for rural feeders",
    version="0.1.0"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 初始化組件
ingestion = FaultLogIngestion()
risk_model = RiskModel(ingestion)

# 導入 API routes
from src.api import router, set_components
set_components(ingestion, risk_model)
app.include_router(router, prefix="/api/v1", tags=["resilience"])


@app.get("/health")
async def health_check():
    """健康檢查"""
    return {"status": "ok"}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8084)

