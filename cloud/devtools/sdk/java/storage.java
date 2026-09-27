package sdk;

public class Storage {

    private final Client client;

    public Storage(Client client) {
        this.client = client;
    }

    public String create(String name, int sizeGB) throws Exception {
        String body = String.format("{\"name\":\"%s\",\"size_gb\":%d}", name, sizeGB);

        var res = client.request("POST", "/storage/create", body);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String delete(String id) throws Exception {
        var res = client.request("DELETE", "/storage/delete/" + id, null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String list() throws Exception {
        var res = client.request("GET", "/storage/list", null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }
}
