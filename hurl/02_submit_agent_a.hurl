POST http://localhost:8080/api/run
Content-Type: application/json
{
    "appName": "{{AGENT_NAME}}",
    "userId": "u_123",
    "sessionId": "{{SESSION_ID}}",
    "newMessage": {
        "role": "user",
        "parts": [{
            "text": "Hey whats the current weather and time in new york today"
        }]
    }
}
