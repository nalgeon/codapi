# Updating Codapi

Follow these steps to update Codapi to a specific version.

1. Backup the current version:

```sh
cd /opt/codapi
mv codapi codapi.bak
```

2. Download and install the new version (replace `0.x.x` with an actual version number):

```sh
cd /opt/codapi
export version="0.x.x"
curl -L -o codapi.tar.gz "https://github.com/nalgeon/codapi/releases/download/${version}/codapi_${version}_linux_amd64.tar.gz"
tar xvzf codapi.tar.gz
rm -f codapi.tar.gz
```

3. Restart Codapi:

```sh
sudo systemctl restart codapi.service
```

4. Verify that Codapi is working:

```sh
curl -H "content-type: application/json" -d '{ "sandbox": "ash", "command": "run", "files": {"": "echo hello" }}' http://localhost:1313/v1/exec
```
