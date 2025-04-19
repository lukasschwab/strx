# strx

![*Asterix und das Atomkraftwerk*](https://beyondnuclearinternational.org/wp-content/uploads/2018/12/asterix-slider.png)

Database-free URL shortener for Tailscale. It supports three operations:

1. **Creating a new alias:** Allows users to create a short alias for a URL. If no alias is provided, a random one is generated.
2. **Resolving an alias:** Redirects users to the original URL when they access the alias.
3. **Displaying all aliases:** Lists all stored aliases and their corresponding URLs.

## Install

```console
$ go install github.com/lukasschwab/strx@latest
```

## Run

```console
$ strx --port 3000
```

## API

1. `GET /` Display All Aliases

    Returns a JSON object containing all stored aliases and their corresponding URLs.

    Example Response:

    ```json
    {
        "example": "https://example.com",
        "example2": "https://example.com"
    }
    ```

2. `POST /create` Create a New Alias

    Creates a new alias for a given URL. If no alias is provided, a random one is generated.

    Request Body:

    ```json
    {
        "url": "https://example.com",
        "alias": "example" // Optional
    }
    ```

    Example Request:

    ```bash
    curl -X POST http://localhost:3000/create \
        -H "Content-Type: application/json" \
        -d '{"url":"https://example.com","alias":"example"}'
    ```

    Example Response:

    ```json
    {
        "alias": "example",
        "url": "https://example.com"
    }
    ```

3. `GET /:alias` Resolve an Alias

    Redirects to the URL associated with the given alias.

    Example: Visiting `http://localhost:{strx port}/example` in your browser will redirect to `https://example.com`.
