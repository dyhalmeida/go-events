
# Projeto Go Events

Este repositório contém um exemplo de implementação de um sistema de eventos utilizando RabbitMQ em Go. O projeto inclui tanto um produtor quanto um consumidor de mensagens, além de um despachante de eventos (`EventDispatcher`) para gerenciar manipuladores de eventos.

## Estrutura de Pastas

```plaintext
├── cmd
│   ├── consumer
│   │   └── main.go
│   └── producer
│       └── main.go
├── pkg
│   ├── events
│   │   ├── dispatcher.go
│   │   └── dispatcher_test.go
│   │   └── interfaces.go
│   └── rabbitmq
│       └── rabbitmq.go
├── docs.md
└── README.md
```

## Descrição dos Componentes

- **cmd/producer/main.go**: Código do produtor de mensagens RabbitMQ.
- **cmd/consumer/main.go**: Código do consumidor de mensagens RabbitMQ.
- **pkg/events/dispatcher.go**: Implementação do despachante de eventos.
- **pkg/events/dispatcher_test.go**: Testes para o despachante de eventos.
- - **pkg/events/interfaces.go**: Interfaces para o despachante de eventos.
- **pkg/rabbitmq/rabbitmq.go**: Funções para interagir com o RabbitMQ.
- **docs.md**: Explicação detalhada do código desse projeto.

## Instruções para Rodar o Projeto

### Pré-requisitos

- [Go](https://golang.org/doc/install) instalado.
- [Docker](https://www.docker.com/) instalado

### Rodando o container do RabbitMQ

1. Suba o container do RabbitMQ:
   ```sh
   docker compose up
   ```

### Rodando o Produtor

1. Navegue até o diretório do producer:
   ```sh
   cd cmd/consumer
   ```

2. Execute o producer:
   ```sh
   go run main.go
   ```

### Rodando o Consumidor

1. Navegue até o diretório do consumer:
   ```sh
   cd cmd/consumer
   ```

2. Execute o consumer:
   ```sh
   go run main.go
   ```

## Documentação

Para uma explicação detalhada do código, consulte o arquivo [docs.md](./docs.md).
