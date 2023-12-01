# Embeddable code playgrounds

_for education, documentation, and fun_ 🎉

Codapi is a platform for embedding interactive code snippets directly into your product documentation, online course, or blog post.

```
  python
┌───────────────────────────────┐
│ msg = "Hello, World!"         │
│ print(msg)                    │
│                               │
│                               │
│ run ►                         │
└───────────────────────────────┘
  ✓ Done
┌───────────────────────────────┐
│ Hello, World!                 │
└───────────────────────────────┘
```

Codapi manages sandboxes (isolated execution environments) and provides an API to execute code in these sandboxes. It also provides a JavaScript widget [codapi-js](https://github.com/nalgeon/codapi-js) for easier integration.

Highlights:

-   Custom sandboxes for any programming language, database, or software.
-   Available as a cloud service and as a self-hosted version.
-   Open source. Uses the permissive Apache-2.0 license.
-   Lightweight and easy to integrate.

Learn more at [codapi.org](https://codapi.org/)

## Installation

See [Installing Codapi](docs/install.md) for details.

## Usage (API)

Call `/v1/exec` to run the code in a sandbox:

```http
POST https://api.codapi.org/v1/exec
content-type: application/json

{
    "sandbox": "python",
    "command": "run",
    "files": {
        "": "print('hello world')"
    }
}
```

`sandbox` is the name of the pre-configured sandbox, and `command` is the name of a command supported by that sandbox. See [Adding a sandbox](docs/add-sandbox.md) for details on how to add a new sandbox.

`files` is a map, where the key is a filename and the value is its contents. When executing a single file, it should either be named as the `command` expects, or be an empty string (as in the example above).

Response:

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": "python_run_9b7b1afd",
  "ok": true,
  "duration": 314,
  "stdout": "hello world\n",
  "stderr": ""
}
```

-   `id` is the unique execution identifier.
-   `ok` is `true` if the code executed without errors, or `false` otherwise.
-   `duration` is the execution time in milliseconds.
-   `stdout` is what the code printed to the standard output.
-   `stderr` is what the code printed to the standard error, or a compiler/os error (if any).

## Usage (JavaScript)

See [codapi-js](https://github.com/nalgeon/codapi-js) to embed the JavaScript widget into a web page. The widget uses exactly the same API as described above.

## Contributing

Pull requests are welcome. For new features and improvements, please first open an issue to discuss what you would like to change.

Be sure to add or update tests as appropriate.

## License

Copyright 2023 [Anton Zhiyanov](https://antonz.org/).

The software is available under the Apache-2.0 license.

## Stay tuned

★ [Subscribe](https://antonz.org/subscribe/) to stay on top of new features.
