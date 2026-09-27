from .client import Client

class Compute:
    def __init__(self, client: Client):
        self.client = client

    def create(self, name: str, plan: str):
        return self.client.request("POST", "/compute/create", {
            "name": name,
            "plan": plan
        })

    def delete(self, compute_id: str):
        return self.client.request("DELETE", f"/compute/delete/{compute_id}")

    def list(self):
        return self.client.request("GET", "/compute/list")
