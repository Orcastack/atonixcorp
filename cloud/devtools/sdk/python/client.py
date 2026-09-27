import requests
import json

class Client:
    def __init__(self, endpoint: str, token: str, timeout: int = 30):
        self.endpoint = endpoint
        self.token = token
        self.timeout = timeout

    def request(self, method: str, path: str, body=None):
        url = self.endpoint + path
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json",
        }

        data = json.dumps(body) if body else None

        resp = requests.request(
            method=method,
            url=url,
            headers=headers,
            data=data,
            timeout=self.timeout,
        )

        if resp.status_code >= 400:
            raise Exception(f"API Error {resp.status_code}: {resp.text}")

        return resp.json()
