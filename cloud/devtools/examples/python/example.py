from sdk.python.client import Client
from sdk.python.compute import Compute
from sdk.python.identity import Identity
from sdk.python.telemetry import Telemetry

client = Client(
    endpoint="https://api.atcloud.world/v1",
    token="YOUR_TOKEN"
)

compute = Compute(client)
identity = Identity(client)
telemetry = Telemetry(client)

print("Compute:", compute.list())
print("Identity:", identity.info())
print("Telemetry:", telemetry.snapshot())
