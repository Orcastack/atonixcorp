import client from "./client.js";

export default {
  async snapshot() {
    return client.request("GET", "/telemetry");
  }
};
