"""
Feeder Simulation Engine API
提供電力系統模擬服務（powerflow, reliability）
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from src.powerflow_stub import PowerflowStub
from src.reliability_stub import ReliabilityStub
from src.penetration_stub import PenetrationStub

app = FastAPI(
    title="Feeder Simulation Engine",
    description="Simulation & Analysis API for feeder topologies",
    version="0.1.0"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # 開發環境，生產環境應限制
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 初始化 stub 服務
powerflow_stub = PowerflowStub()
reliability_stub = ReliabilityStub()
penetration_stub = PenetrationStub()


@app.get("/health")
async def health_check():
    """健康檢查"""
    return {"status": "ok"}


@app.post("/simulate/powerflow")
async def simulate_powerflow(request: dict):
    """
    執行潮流分析（stub 版本）
    
    接收拓樸 JSON，回傳節點電壓與線路載流率
    """
    topology = request.get("topology", {})
    nodes = topology.get("nodes", [])
    lines = topology.get("lines", [])
    
    result = powerflow_stub.run_powerflow(nodes, lines)
    return result


@app.post("/simulate/reliability")
async def simulate_reliability(request: dict):
    """
    執行可靠度分析（stub 版本）
    
    接收拓樸 JSON + 假設參數，回傳 SAIDI/SAIFI
    """
    topology = request.get("topology", {})
    parameters = request.get("parameters", {})
    
    result = reliability_stub.run_reliability_analysis(topology, parameters)
    return result


@app.post("/simulate/penetration")
async def simulate_penetration(request: dict):
    """
    執行滲透測試（stub 版本）
    
    接收拓樸 JSON + 攻擊場景列表，回傳攻擊結果
    """
    topology = request.get("topology", {})
    attack_scenarios = request.get("attack_scenarios", [])
    target_nodes = request.get("target_nodes", None)
    
    result = penetration_stub.run_penetration_test(
        topology=topology,
        attack_scenarios=attack_scenarios,
        target_nodes=target_nodes
    )
    return result


@app.get("/simulate/penetration/scenarios")
async def get_penetration_scenarios():
    """
    取得所有可用的攻擊場景列表
    """
    scenarios = penetration_stub.get_available_scenarios()
    return {"scenarios": scenarios}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8081)

