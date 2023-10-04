from flask import Flask, request, jsonify
from llm.Privatellm import PrivateLLM

app = Flask(__name__)

@app.route("/", methods=["POST"])
def call():
    try:
        body = request.get_json()
        print(body)
        l = PrivateLLM()
        resp = l.PromptForBulletPoint(body["prompt"], body["t"])
        return jsonify({"response": resp["content"]}), 201

    except Exception as error:
        return error, 500

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5555)
