# API

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

`sandbox` is the name of the pre-configured sandbox, and `command` is the name of a command supported by that sandbox. See [Adding a sandbox](add-sandbox.md) for details on how to add a new sandbox.

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
-
