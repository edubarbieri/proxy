# Proxy server

Servidor http proxy construído com go lang com as seguintes funcionalidades:
- Access logs dos requests processados;
- Estatisticas do servidor:
  - Total de request atendidos;
  - Tempo de resposta médio;
  - Total de requests por http status;
  - Total de requests por url path;
- Rate limite
  - Possibilita configurar o limite por: path de destino, ip de origem e valor de header http;
  - Utiliza redis para possibilitar que varias instancias do proxy compartilhem as informações.
- Rotas configuráveis:
  - Configurável por path 
  - Múltiplos backends
  - Possibilita configuração de middlewares por rota 
- Api Rest para administração das configurações e estatísticas

## Desenho
<p align="center">
<img src="proxy.png" alt="Diagram" title="Diagram" />
</p>

## Iniciar todos os serviços
```shell
docker-compose.exe up --build --remove-orphans
```
## API Rest
A documentação para api resto pode ser encontrada [aqui](swagger.yml).

### Update proxy config
```shell
curl --request PUT \
  --url http://localhost:8080/proxy-admin/config \
  --header 'Content-Type: application/json' \
  --data '{
  "routes": [
    {
      "pattern": "/backend1",
      "backends": [
        "http://backend1:3000"
      ]
    },
    {
      "pattern": "/backend2",
      "backends": [
        "http://backend2a:3001",
        "http://backend2b:3002",
				"http://backend2c:3003",
        "http://backend2d:3004"
      ]
    }
  ],
  "limits": [
    {
      "id": "1",
      "requestMin": 100000,
      "targetPath": "/backend2",
      "sourceIp": true,
      "headerValue": ""
    }    
  ]
}'
```
### Get current proxy config
```shell
curl --request GET \
  --url http://localhost:8080/proxy-admin/config
```
### Get proxy statistics
```shell
curl --request GET \
  --url http://localhost:8080/proxy-admin/stats
```
### Reset proxy statistics
```shell
curl --request DELETE \
  --url http://localhost:8080/proxy-admin/stats
```
### Get backend1
```shell
curl --request GET \
  --url http://localhost:8080/backend1
```
### Get backend2
```shell
curl --request GET \
  --url http://localhost:8080/backend2
```
## Test

Para testes básicos de desempenho utilizei a ferramenta [bombardier](https://github.com/codesenberg/bombardier), disparando para uma rota configurada com 4 replicas do [sample-backend](sample-backend), os resultados obtidos podem ser vistos a seguir.

### Comando executado
```shell
bombardier -c 10 -d 10m -l http://localhost:8080/backend2
```
<p align="center">
<img src="stresstest.png" alt="Diagram" title="Diagram" />
</p>



