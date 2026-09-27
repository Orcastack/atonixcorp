package sdk;

public class Orchestration {

    private final Client client;

    public Orchestration(Client client) {
        this.client = client;
    }

    public String events() throws Exception {
        var res = client.request("GET", "/orchestration/events", null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String metrics() throws Exception {
        var res = client.request("GET", "/orchestration/metrics", null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }
}
