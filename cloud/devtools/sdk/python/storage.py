from .client import Client

class Storage:
    def __init__(self, client: Client):
        self.client = client

    def create(self, name: str, size_gb: int):
        return self.client.request("POST", "/storage/create", {
            "name": name,
            "size_gb": size_gb
        })

    def delete(self, volume_id: str):
        return self.client.request("DELETE", f"/storage/delete/{volume_id}")

    def list(self):
        return self.client.request("GET", "/storage/list")
