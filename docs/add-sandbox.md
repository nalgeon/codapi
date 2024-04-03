# Adding a sandbox

A _sandbox_ is an isolated execution environment for running code snippets. A sandbox is typically implemented as one or more Docker containers. A sandbox supports at least one _command_ (usually the `run` one), but it can support more (like `test` or any other).

Codapi comes with a single `sh` sandbox preinstalled, but you can easily add others. Let's see some examples.

## Python

First, let's create a Docker image capable of running Python with some third-party packages:

```sh
cd /opt/codapi
mkdir images/python
touch images/python/Dockerfile
touch images/python/requirements.txt
```

Fill the `Dockerfile`:

```Dockerfile
FROM python:3.11-alpine

RUN adduser --home /sandbox --disabled-password sandbox

COPY requirements.txt /tmp
RUN pip install --no-cache-dir -r /tmp/requirements.txt && rm -f /tmp/requirements.txt

USER sandbox
WORKDIR /sandbox

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
```

And the `requirements.txt`:

```
numpy==1.26.2
pandas==2.1.3
```

Build the image:

```sh
docker build --file images/python/Dockerfile --tag codapi/python:latest images/python/
```

And register the image as a Codapi _box_ in `configs/boxes.json`:

```js
{
    // ...
    "python": {
        "image": "codapi/python"
    }
}
```

Finally, let's configure what happens when the client executes the `run` command in the `python` sandbox. To do this, we create `configs/commands/python.json`:

```js
{
    "run": {
        "engine": "docker",
        "entry": "main.py",
        "steps": [
            {
                "box": "python",
                "command": ["python", "main.py"]
            }
        ]
    }
}
```

This is essentially what it says:

> When the client executes the `run` command in the `python` sandbox, save their code to the `main.py` file, then run it in the `python` box (Docker container) using the `python main.py` shell command.

What if we want to add another command (say, `test`) to the same sandbox? Let's edit `configs/commands/python.json` again:

```js
{
    "run": {
        // ...
    },
    "test": {
        "engine": "docker",
        "entry": "test_main.py",
        "steps": [
            {
                "box": "python",
                "command": ["python", "-m", "unittest"],
                "noutput": 8192
            }
        ]
    }
}
```

Besides configuring a different shell command, here we increased the maximum output size to 8Kb, as tests tend to be quite chatty (you can see the default value in `configs/config.json`).

To apply the changed configuration, restart Codapi (as root):

```sh
systemctl restart codapi.service
```

And try running some Python code:

```sh
curl -H "content-type: application/json" -d '{ "sandbox": "python", "command": "run", "files": {"": "print(42)" }}' http://localhost:1313/v1/exec
```

Which produces the following output:

```json
{
    "id": "python_run_7683de5a",
    "ok": true,
    "duration": 252,
    "stdout": "42\n",
    "stderr": ""
}
```

## Go

First, let's create a Docker image capable of running Go:

```sh
cd /opt/codapi
mkdir images/go
touch images/go/Dockerfile
```

Fill the `Dockerfile`:

```Dockerfile
FROM golang:1.22.1-alpine3.19

RUN adduser --home /sandbox --disabled-password sandbox
USER sandbox
WORKDIR /sandbox
```

Build the image:

```sh
docker build --file images/go/Dockerfile --tag codapi/go:latest images/go/
```

And register the image as a Codapi _box_ in `configs/boxes.json`:

```json
{
    // ...
    "go": {
        "image": "codapi/go",
        "runtime": "runc",
        "cpu": 2,
        "memory": 512,
        "network": "none",
        "writable": true,
        "volume": "%s:/sandbox:rw",
        "cap_drop": ["all"],
        "ulimit": ["nofile=96"],
        "nproc": 64
    }
}
```

Finally, let's configure what happens when the client executes the `run` command in the `go` sandbox. To do this, we create `configs/commands/go.json`:

```json
{
    "run": {
        "engine": "docker",
        "entry": "main.go",
        "steps": [
            {
                "box": "go",
                "command": ["go", "run", "main.go"]
            }
        ]
    }
}
```

This is essentially what it says:

> When the client executes the `run` command in the `go` sandbox, save their code to the `main.go` file, then run it in the `go` box (Docker container) using the `go run main.go` shell command.

To apply the changed configuration, restart Codapi (as root):

```sh
systemctl restart codapi.service
```

And try running some go code:

```sh
curl -H "content-type: application/json" -d '{"sandbox":"go","version":"","command":"run","files":{"":"package main\nimport (\n    \"fmt\"\n)\n\nfunc main() {\n    fmt.Println(\"hello\")\n}"}}' http://localhost:1313/v1/exec
```

Which produces the following output:

```json
{
    "id": "go_run_f9592410",
    "ok": true,
    "duration": 10839,
    "stdout": "hello",
    "stderr": ""
}
```