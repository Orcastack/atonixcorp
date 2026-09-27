from .client import Client

class Projects:
    def __init__(self, client: Client):
        self.client = client

    def list(self):
        return self.client.request("GET", "/projects/list")
