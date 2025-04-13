# Deploying to production

If you want to run a public Codapi instance, you'll need a server for that.

The server requirements mainly depend on the sandboxes you plan to install (e.g. just Python and SQLite or something heavier) and the amount of code runs you expect from your users (e.g. 1000 runs per day).

Personally, I prefer [DigitalOcean](https://www.digitalocean.com/) for hosting. For lighter sandboxes, a $6 (1CPU/1GB RAM) or $12 (1CPU/2GB RAM) [Basic](https://www.digitalocean.com/pricing/droplets#basic-droplets) droplet would do the trick. If you need something more powerful, a $42 [CPU-optimized droplet](https://www.digitalocean.com/pricing/droplets#cpu-optimized) (2CPU/4GB RAM) will definitely suffice.

I recommend using Debian as the operating system.

Follow these steps to set up a production server:

1. [Install Codapi](install.md) on the server.

2. Install [Nginx](https://www.digitalocean.com/community/tutorials/how-to-install-nginx-on-debian-11) and an [HTTPS certificate](https://certbot.eff.org/) (if you want Codapi to be accessible via HTTPS).

3. Connect [Nginx to Codapi](nginx.md).

You can also use Caddy or any other proxy you prefer instead of Nginx.

That's it!
