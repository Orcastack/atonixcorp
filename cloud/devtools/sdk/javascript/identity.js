import client from "./client.js";

export default {
  async info() {
    return client.request("GET", "/identity/info");
  }
};
