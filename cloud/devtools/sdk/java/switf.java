package sdk;

public class Swift {

    private final Client client;

    public Swift(Client client) {
        this.client = client;
    }

    public String createContainer(String name) throws Exception {
        String body = String.format("{\"name\":\"%s\"}", name);

        var res = client.request("POST", "/swift/container/create", body);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String deleteContainer(String name) throws Exception {
        var res = client.request("DELETE", "/swift/container/delete/" + name, null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String uploadObject(String container, String name, String data) throws Exception {
        String body = String.format(
                "{\"container\":\"%s\",\"name\":\"%s\",\"data\":\"%s\"}",
                container, name, data
        );

        var res = client.request("POST", "/swift/object/upload", body);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }

    public String deleteObject(String container, String name) throws Exception {
        var res = client.request("DELETE", "/swift/object/delete/" + container + "/" + name, null);

        if (res.statusCode() >= 400)
            throw new RuntimeException(res.body());

        return res.body();
    }
}
