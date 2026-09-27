package sdk;

import java.io.IOException;
import java.net.URI;
import java.net.http.*;
import java.time.Duration;

public class Client {

    private final String endpoint;
    private final String token;
    private final HttpClient http;

    public Client(String endpoint, String token) {
        this.endpoint = endpoint;
        this.token = token;

        this.http = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(10))
                .version(HttpClient.Version.HTTP_2)
                .build();
    }

    public HttpResponse<String> request(String method, String path, String body)
            throws IOException, InterruptedException {

        HttpRequest.Builder builder = HttpRequest.newBuilder()
                .uri(URI.create(endpoint + path))
                .header("Authorization", "Bearer " + token)
                .header("Content-Type", "application/json");

        switch (method) {
            case "GET" -> builder.GET();
            case "DELETE" -> builder.DELETE();
            case "POST" -> builder.POST(HttpRequest.BodyPublishers.ofString(body));
            default -> throw new IllegalArgumentException("Unsupported method: " + method);
        }

        return http.send(builder.build(), HttpResponse.BodyHandlers.ofString());
    }
}
