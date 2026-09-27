import client from "./client.js";

export default {
  async create(name, size_gb) {
    return client.request("POST", "/storage/create", { name, size_gb });
  },

  async delete(id) {
    return client.request("DELETE", `/storage/delete/${id}`);
  },

  async list() {
    return client.request("GET", "/storage/list");
  }
};
