from .client import Client

class Telemetry:
    def __init__(self, client: Client):
        self.client = client

    def snapshot(self):
        return self.client.request("GET", "/telemetry")
