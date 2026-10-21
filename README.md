# Stress Test CLI

CLI em Go para testes de carga em serviços web.

## Parâmetros

| Flag            | Descrição                              |
|-----------------|----------------------------------------|
| `--url`         | URL do serviço a ser testado           |
| `--requests`    | Número total de requisições (exato)    |
| `--concurrency` | Número de chamadas simultâneas         |
| `--timeout`     | (opcional) Timeout por request, def. 30s |

## Build da imagem

```sh
docker build -t stresstest .
```

## Execução

```sh
docker run stresstest --url=http://google.com --requests=1000 --concurrency=10
```

## Execução local (sem Docker)

```sh
go run ./cmd/stresstest --url=http://google.com --requests=100 --concurrency=10
```

## Testes

```sh
go test ./...
```

## Relatório

Ao final são exibidos: tempo total, total de requests, quantidade de status 200,
distribuição dos demais status HTTP e a quantidade de falhas de conexão/timeout.
