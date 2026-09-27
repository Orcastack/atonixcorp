import client from "./client.js";

export default {
  async create(name, cidr) {
    return client.request("POST", "/network/create", { name, cidr });
  },

  async delete(id) {
    return client.request("DELETE", `/network/delete/${id}`);
  },

  async list() {
    return client.request("GET", "/network/list");
  }
};
