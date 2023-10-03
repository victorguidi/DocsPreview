from dotenv import load_dotenv
import os
from typing import Any

load_dotenv()

class PrivateLLM:
    def __init__(self) -> None:
        self.model = str( os.getenv("DEPLOYMENT_NAME") )
        self.api_key = str( os.getenv("OPENAI_API_KEY") )
        self.api_version = str( os.getenv("OPENAI_API_VERSION") )
        self.resource = str( os.getenv("RESOURCE") )

    def PromptForBulletPoint(self, prompt: str, type: str) -> Any:
        try:
            import requests
            import json

            chat = []
            chat.append(self.create_system_message(type))

            message_user = self.create_user_message(prompt)

            print(message_user)

            chat.append(message_user)

            url = 'https://{0}/openai/deployments/{1}/chat/completions?api-version={2}'.format(self.resource,
                                                                                               self.model, self.api_version)
            headers = {
                'Content-Type': 'application/json',
                'api-key': self.api_key
            }
            body_data = {
                "messages": chat
            }

            print(body_data)

            response = requests.post(url, headers=headers, json=body_data)
            text = response.text
            parsed_json = json.loads(text)

            history = parsed_json['choices'][0]['message']
            chat.append(history)

            return {"history": chat[-1]}

        except Exception as error:
            raise Exception(error)


    def create_user_message(self, message):
        return {"role": "user", "content": "{0}".format(message)}

    def create_assistant_message(self, message):
        return {"role": "assistant", "content": "{0}".format(message)}

    def create_system_message(self, type: str):
        return {"role": "assistant", "content": f"Você é uma IA especialidada em resumir os prompts que o usuário te passar e retornar em formato markdow, mas formatando os pontos principais de acordo com o que inferir no tipo: {type}"}


