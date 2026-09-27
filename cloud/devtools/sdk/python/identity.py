from .client import Client

class Identity:
    def __init__(self, client: Client):
        self.client = client

    def info(self):
        return self.client.request("GET", "/identity/info")
