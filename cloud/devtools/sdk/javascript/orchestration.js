import client from "./client.js";

export default {
  async events() {
    return client.request("GET", "/orchestration/events");
  },

  async metrics() {
    return client.request("GET", "/orchestration/metrics");
  }
};
