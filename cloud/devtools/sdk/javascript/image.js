import client from "./client.js";

export default {
  async list() {
    return client.request("GET", "/images/list");
  }
};
