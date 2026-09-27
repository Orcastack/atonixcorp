package sdk;

public class Identity {

    private final Client client;

    public Identity(Client client) {
        this.client = client;
    }

    public String info() throws Exception {
        var res = client.request("GET", "/identity/info", null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }
}
