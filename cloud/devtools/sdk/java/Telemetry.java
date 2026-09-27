package sdk;

public class Telemetry {

    private final Client client;

    public Telemetry(Client client) {
        this.client = client;
    }

    public String snapshot() throws Exception {
        var res = client.request("GET", "/telemetry", null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }
}
