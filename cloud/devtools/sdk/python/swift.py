from .client import Client

class Swift:
    def __init__(self, client: Client):
        self.client = client

    def create_container(self, name: str):
        return self.client.request("POST", "/swift/container/create", {
            "name": name
        })

    def delete_container(self, name: str):
        return self.client.request("DELETE", f"/swift/container/delete/{name}")

    def upload_object(self, container: str, name: str, data: str):
        return self.client.request("POST", "/swift/object/upload", {
            "container": container,
            "name": name,
            "data": data
        })

    def delete_object(self, container: str, name: str):
        return self.client.request("DELETE", f"/swift/object/delete/{container}/{name}")
