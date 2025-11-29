#!/bin/sh

# This script is meant to be run after setting up the Playwright Server
# with `playwrightd.sh`.

PW_TEST_CONNECT_WS_ENDPOINT=ws://127.0.0.1:3000/ pnpm playwright test
