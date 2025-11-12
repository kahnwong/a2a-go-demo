import json

with open("response.json", "r") as f:
    d = json.loads(f.read())

agent_name = d[1]["author"]
print(agent_name)

for i in d[1]["content"]["parts"]:
    function_name = i["functionResponse"]["name"]
    print(f"\t{function_name}")
