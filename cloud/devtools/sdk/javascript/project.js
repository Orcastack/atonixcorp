import client from "./client.js";

export default {
  async list() {
    return client.request("GET", "/projects/list");
  }
};
