package sdk;

public class Network {

    private final Client client;

    public Network(Client client) {
        this.client = client;
    }

    public String create(String name, String cidr) throws Exception {
        String body = String.format("{\"name\":\"%s\",\"cidr\":\"%s\"}", name, cidr);

        var res = client.request("POST", "/network/create", body);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String delete(String id) throws Exception {
        var res = client.request("DELETE", "/network/delete/" + id, null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String list() throws Exception {
        var res = client.request("GET", "/network/list", null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }
}
