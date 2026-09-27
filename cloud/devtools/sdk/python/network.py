from .client import Client

class Network:
    def __init__(self, client: Client):
        self.client = client

    def create(self, name: str, cidr: str):
        return self.client.request("POST", "/network/create", {
            "name": name,
            "cidr": cidr
        })

    def delete(self, network_id: str):
        return self.client.request("DELETE", f"/network/delete/{network_id}")

    def list(self):
        return self.client.request("GET", "/network/list")
