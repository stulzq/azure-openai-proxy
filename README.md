# azure-openai-proxy
Azure OpenAI Service Proxy. Convert Azure OpenAI API(aoai) to the official OpenAI API(oai) standard.

## Get Start

### Retrieve key and endpoint

To successfully make a call against Azure OpenAI, you'll need the following:

| Name                  | Desc                                                         | Example                       |
| --------------------- | ------------------------------------------------------------ | ----------------------------- |
| AZURE_OPENAI_ENDPOINT | This value can be found in the **Keys & Endpoint** section when examining your resource from the Azure portal. Alternatively, you can find the value in **Azure OpenAI Studio** > **Playground** > **Code View**. An example endpoint is: `https://docs-test-001.openai.azure.com/`. | https://test.openai.azure.com |
| AZURE_OPENAI_API_VER  | [See here](https://learn.microsoft.com/en-us/azure/cognitive-services/openai/quickstart?tabs=command-line&pivots=rest-api) or Azure OpenAI Studio | 2023-03-15-preview            |
| AZURE_OPENAI_DEPLOY   | This value will correspond to the custom name you chose for your deployment when you deployed a model. This value can be found under **Resource Management** > **Deployments** in the Azure portal or alternatively under **Management** > **Deployments** in Azure OpenAI Studio. | gpt-35-turbo                  |

![Screenshot of the overview UI for an OpenAI Resource in the Azure portal with the endpoint & access keys location circled in red.](docs/assets/images/endpoint.png)

API Key: This value can be found in the **Keys & Endpoint** section when examining your resource from the Azure portal. You can use either `KEY1` or `KEY2`. 

> Use api key when call http api like Offical OpenAI.

### Use Docker

````shell
docker run -d --name=azure-openai-proxy stulzq/azure-openai-proxy
````

### Use ChatGPT-Web

ChatGPT Web: https://github.com/Chanzhaoyu/chatgpt-web

![image-20230324163357406](docs/assets/images/image-20230324163357406.png)

Envs:

- `OPENAI_API_KEY` Auzre OpenAI API Key
- `OPENAI_API_BASE_URL` Proxy addr

docker-compose.yml:

````yaml
version: '3'

services:
  chatgpt-web:
    image: chenzhaoyu94/chatgpt-web
    ports:
      - 3002:3002
    environment:
      OPENAI_API_KEY: <Auzre OpenAI API Key>
      OPENAI_API_BASE_URL: http://azure-openai:8080
      AUTH_SECRET_KEY: ""
      MAX_REQUEST_PER_HOUR: 1000
      TIMEOUT_MS: 60000
    depends_on:
      - azure-openai
    links:
      - azure-openai
    networks:
      - chatgpt-ns

  azure-openai:
    image: stulzq/azure-openai-proxy
    ports:
      - 8080:8080
    environment:
      AZURE_OPENAI_ENDPOINT: <Auzre OpenAI API Endpoint>
      AZURE_OPENAI_DEPLOY: <Auzre OpenAI API Deployment>
      AZURE_OPENAI_API_VER: 2023-03-15-preview
    networks:
      - chatgpt-ns

networks:
  chatgpt-ns:
    driver: bridge
````

Run:

````shell
docker compose up -d
````

## Proxy Api

| Api                  | Status |
| -------------------- | ------ |
| /v1/chat/completions | Ok     |
