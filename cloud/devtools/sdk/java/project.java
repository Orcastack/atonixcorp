package sdk;

public class Projects {

    private final Client client;

    public Projects(Client client) {
        this.client = client;
    }

    public String list() throws Exception {
        var res = client.request("GET", "/projects/list", null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }
}
