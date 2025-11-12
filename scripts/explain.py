import json

with open("response.json", "r") as f:
    d = json.loads(f.read())

agent_name = d[0]["author"]
print(agent_name)

for i in d[0]["content"]["parts"]:
    function_name = i["text"]
    print(f"\t{function_name}")
