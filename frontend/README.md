# CherryAuctions - Frontend

## Overview

This is the frontend service for CherryAuctions, built as a single page application
using Rolldown Vite and Vue 3, along with commonly used libraries like Pinia, Vue
Router and Zod.

## Caveats: Unsupported Linux Distributions

This project uses **Playwright** for end-to-end testing, and Playwright is only
officially supported for Ubuntu 22, Ubuntu 24 and Debian as Linux distributions.
The project's owner is currently using an Arch Linux machine, so there have to
be some weird setups to let Playwright work.

Playwright provides this through a contained browser image, hosted on Microsoft.
At the time of writing this, the latest stable version is 1.57 (you may use Jammy
instead of Noble).

```bash
docker pull mcr.microsoft.com/playwright:v1.57.0-noble
```

Use the script inside `scripts` folder to run Playwright. This browser container
works by making a websocket server for our Playwright in NPM to connect to, and
run tests on browsers that way. But due to Docker isolation, you can't actually
just run it, stream it from NPM and expect it to work.

```bash
docker run --rm --init -it --workdir /home/pwuser --user pwuser --network host \
 mcr.microsoft.com/playwright:v1.57.0-noble \
 /bin/sh -c "cd /home/pwuser && npx -y playwright@1.57.0 run-server --port 3000 --host 0.0.0.0"
```

Simple explanations: This creates a random container with playwright 1.57, that
removes itself after use, adds it to your Host network and run the Playwright server.

Beware that this command will add the server into your host, and it would be able
to read host network, so use at your own risk.

To run Playwright tests,

```bash
PW_TEST_CONNECT_WS_ENDPOINT=ws://127.0.0.1:3000/ pnpm playwright test
```

Provide the environment variable `PW_TEST_CONNECT_WS_ENDPOINT` to the Playwright
server. It may differ here, if you hosted that Playwright Server on another
machine, or on another port.
