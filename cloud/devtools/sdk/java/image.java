package sdk;

public class Images {

    private final Client client;

    public Images(Client client) {
        this.client = client;
    }

    public String list() throws Exception {
        var res = client.request("GET", "/images/list", null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }
}
