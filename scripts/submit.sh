#!/bin/bash

export HURL_AGENT_NAME="root_agent"
HURL_SESSION_ID=$(openssl rand -base64 6)
export HURL_SESSION_ID
hurl hurl/01_init_session.hurl

hurl hurl/02_submit.hurl >response.json
echo ""
cat response.json | jq | grep text
