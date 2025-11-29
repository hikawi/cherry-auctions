#!/bin/sh

# This is for unsupported Linux distributions. For Ubuntu or Debian, just install
# Playwright directly.

docker run --rm --init -it --workdir /home/pwuser --user pwuser --network host \
	mcr.microsoft.com/playwright:v1.57.0-noble \
	/bin/sh -c "cd /home/pwuser && npx -y playwright@1.57.0 run-server --port 3000 --host 0.0.0.0"
