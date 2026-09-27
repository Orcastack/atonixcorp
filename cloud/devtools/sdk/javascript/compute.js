import client from "./client.js";

export default {
  async create(name, plan) {
    return client.request("POST", "/compute/create", { name, plan });
  },

  async delete(id) {
    return client.request("DELETE", `/compute/delete/${id}`);
  },

  async list() {
    return client.request("GET", "/compute/list");
  }
};
