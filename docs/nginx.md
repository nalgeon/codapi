# Connecting Nginx to Codapi

Nginx serves as the entry point for all HTTP(S) requests, while Codapi handles the requests and returns the results to Nginx, which in turn returns them to the end user.

Follow these steps to configure Nginx to work with Codapi:

1. Add a config file to Nginx (e.g. `/etc/nginx/sites-available/domain.org`):

```
# response cache
proxy_cache_path /var/cache/nginx keys_zone=codapi:10m levels=1:2 max_size=512m inactive=60m use_temp_path=off;

# allowed origins (domains)
map $http_origin $auth_origin {
    default 0;
    "http://localhost:3000" 1;
    "http://127.0.0.1:3000" 1;
    "https://domain.org" 1;
}

server {
    listen 80;
    listen 443 ssl;
    server_name domain.org;

    # limit request rate
    limit_req_zone $binary_remote_addr zone=sandbox:10m rate=10r/m;
    limit_req_status 429;

    # limit request size
    client_header_buffer_size 1k;
    client_body_buffer_size 16k;
    client_max_body_size 16k;
    large_client_header_buffers 2 1k;

    # persistent connections to upstream
    proxy_http_version 1.1;
    proxy_set_header Connection "";

    # pass client's ip to upstream
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    # https certificate
    ssl_certificate /etc/letsencrypt/live/domain.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/domain.org/privkey.pem;

    # response cache
    proxy_cache codapi;
    proxy_cache_methods GET HEAD POST;
    proxy_cache_key "$proxy_host$request_uri:$request_body";
    proxy_cache_valid 429 10s;
    proxy_cache_valid 404 500 502 503 1m;
    proxy_cache_valid any 60m;
    proxy_cache_use_stale updating error timeout http_500 http_502 http_503;
    proxy_cache_lock on;
    proxy_ignore_headers cache-control;
    add_header x-cache-status $upstream_cache_status;

    location / {
        if ($auth_origin = "0") {
            return 403;
        }
        limit_req zone=sandbox burst=20 nodelay;
        proxy_pass http://localhost:1313;
    }
}
```

Replace `domain.org` with your actual domain and make sure the `ssl_certificate` and `ssl_certificate_key` paths are correct (you can see them with certbot).

Set the allowed domains in the "allowed origins" section.

2. Activate the configuration:

```sh
sudo ln -s /etc/nginx/sites-available/domain.org /etc/nginx/sites-enabled/
```

3. Make sure there are no errors:

```sh
sudo nginx -t
```

4. Restart Nginx:

```sh
sudo systemctl restart nginx
```

5. Verify that Nginx+Codapi is working:

```sh
curl -H "content-type: application/json" -d '{ "sandbox": "sh", "command": "run", "files": {"": "echo hello" }}' https://domain.org/v1/exec
```
