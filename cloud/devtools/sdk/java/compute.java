package sdk;

import com.fasterxml.jackson.databind.ObjectMapper;

public class Compute {

    private final Client client;
    private final ObjectMapper mapper = new ObjectMapper();

    public Compute(Client client) {
        this.client = client;
    }

    public String create(String name, String plan) throws Exception {
        String body = String.format("{\"name\":\"%s\",\"plan\":\"%s\"}", name, plan);

        var res = client.request("POST", "/compute/create", body);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String delete(String id) throws Exception {
        var res = client.request("DELETE", "/compute/delete/" + id, null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String list() throws Exception {
        var res = client.request("GET", "/compute/list", null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }
}
