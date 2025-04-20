# Interactive code examples

_for documentation, education and fun_ 🎉

Codapi is a platform for embedding interactive code snippets directly into your product documentation, online course or blog post. It's also useful for experimenting with new languages, databases, or tools in a sandbox.

```
┌───────────────────────────────┐
│ def greet(name):              │
│   print(f"Hello, {name}!")    │
│                               │
│ greet("World")                │
└───────────────────────────────┘
  Run ►  Edit  ✓ Done
┌───────────────────────────────┐
│ Hello, World!                 │
└───────────────────────────────┘
```

Codapi manages sandboxes (isolated execution environments) and provides an API to execute code in these sandboxes. It also provides a JavaScript widget [codapi-js](https://github.com/nalgeon/codapi-js) for easier integration.

For an introduction to Codapi, see this post: [Interactive code examples for fun and profit](https://antonz.org/code-examples/).

## Installation

To run Codapi locally, follow these steps:

1. Install Docker (or Podman/OrbStack) for your operating system.
2. Install the [latest](https://github.com/nalgeon/codapi/releases/latest) Codapi release (change the `linux_amd64` part according to your OS):

```sh
mkdir ~/codapi && cd ~/codapi
curl -L -o codapi.tar.gz "https://github.com/nalgeon/codapi/releases/download/v0.11.0/codapi_0.11.0_linux_amd64.tar.gz"
tar xvzf codapi.tar.gz
rm -f codapi.tar.gz
```

3. Build the sample `ash` sandbox image:

```sh
docker build --file sandboxes/ash/Dockerfile --tag codapi/ash:latest sandboxes/ash
```

4. Start the server:

```sh
./codapi
```

## Usage

See [Adding a sandbox](docs/add-sandbox.md) to add a sandbox from the [registry](https://github.com/nalgeon/sandboxes) or create a custom one.

See [API](docs/api.md) to run sandboxed code using the HTTP API.

See [codapi-js](https://github.com/nalgeon/codapi-js) to embed the JavaScript widget into a web page.

## Production

Running in production is a bit more involved. See these guides:

-   [Installing Codapi](docs/install.md)
-   [Updating Codapi](docs/update.md)
-   [Deploying to production](docs/production.md)

## Contributing

Contributions are welcome. For anything other than bugfixes, please first open an issue to discuss what you want to change.

Be sure to add or update tests as appropriate.

## Support

Codapi is mostly a [one-man](https://antonz.org/) project, not backed by a VC fund or anything.

If you find Codapi useful, please star it on GitHub and spread the word among your peers. It really helps to move the project forward.

★ [Subscribe](https://antonz.org/subscribe/) to stay on top of new features.
