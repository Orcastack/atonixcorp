import client from "./client.js";

export default {
  async createContainer(name) {
    return client.request("POST", "/swift/container/create", { name });
  },

  async deleteContainer(name) {
    return client.request("DELETE", `/swift/container/delete/${name}`);
  },

  async uploadObject(container, name, data) {
    return client.request("POST", "/swift/object/upload", {
      container,
      name,
      data
    });
  },

  async deleteObject(container, name) {
    return client.request("DELETE", `/swift/object/delete/${container}/${name}`);
  }
};
