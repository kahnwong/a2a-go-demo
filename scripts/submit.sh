#!/bin/bash

export HURL_AGENT_NAME="weather_time_agent"
HURL_SESSION_ID=$(openssl rand -base64 6)
export HURL_SESSION_ID
hurl hurl/01_init_session.hurl

hurl hurl/02_submit.hurl >response.json
echo ""
jq '.[-1].content.parts[0].text' response.json

python3 scripts/explain.py
