from .client import Client

class Orchestration:
    def __init__(self, client: Client):
        self.client = client

    def events(self):
        return self.client.request("GET", "/orchestration/events")

    def metrics(self):
        return self.client.request("GET", "/orchestration/metrics")
