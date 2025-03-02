# Explicação do Código do Despachante de Eventos

Este documento fornece uma explicação detalhada do código encontrado em `dispatcher.go`, que implementa um `EventDispatcher`. O `EventDispatcher` é responsável por gerenciar manipuladores de eventos em uma aplicação Go. Abaixo, vamos detalhar os principais componentes e métodos do `EventDispatcher`.

## Visão Geral do EventDispatcher

O `EventDispatcher` é uma struct que mantém um mapa de nomes de eventos para seus respectivos manipuladores. Cada evento pode ter múltiplos manipuladores associados, permitindo uma gestão flexível de eventos.

### Componentes Principais

1. **Definição da Struct**:
   ```go
   type EventDispatcher struct {
       handlers map[string][]EventHandlerInterface
   }
   ```
   - `handlers`: Um mapa onde a chave é uma string representando o nome do evento, e o valor é um slice de `EventHandlerInterface` que contém os manipuladores para esse evento.

2. **Construtor**:
   ```go
   func NewEventDispatcher() *EventDispatcher {
       return &EventDispatcher{
           handlers: make(map[string][]EventHandlerInterface),
       }
   }
   ```
   - Esta função inicializa um novo `EventDispatcher` com um mapa vazio de manipuladores.

## Métodos

### Register

```go
func (dispatcher *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {

	if _, ok := dispatcher.handlers[eventName]; ok {
		if slices.Contains(dispatcher.handlers[eventName], handler) {
			return ErrHandlerAlreadyExists
		}
	}

	dispatcher.handlers[eventName] = append(dispatcher.handlers[eventName], handler)
	return nil
}
```
- **Propósito**: Registra um novo manipulador para um evento específico.
- **Funcionalidade**:
  - Verifica se o evento já possui manipuladores registrados.
  - Se o manipulador já estiver registrado para aquele evento, retorna um erro (`ErrHandlerAlreadyExists`).
  - Caso contrário, adiciona o manipulador à lista de manipuladores para aquele evento.

### Clear

```go
func (dispatcher *EventDispatcher) Clear() {
	dispatcher.handlers = make(map[string][]EventHandlerInterface)
}
```
- **Propósito**: Limpa todos os manipuladores registrados.
- **Funcionalidade**: Reseta o mapa `handlers` para um novo mapa vazio.

### Has

```go
func (dispatcher *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	_, ok := dispatcher.handlers[eventName]

	if ok && slices.Contains(dispatcher.handlers[eventName], handler) {
		return true
	}

	return false
}
```
- **Propósito**: Verifica se um manipulador específico está registrado para um determinado evento.
- **Funcionalidade**: Retorna `true` se o manipulador existir para o evento; caso contrário, retorna `false`.

### Dispatch

```go
func (dispatcher *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers, ok := dispatcher.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}

	return nil
}
```
- **Propósito**: Despacha um evento para todos os manipuladores registrados.
- **Funcionalidade**:
  - Recupera os manipuladores para o evento com base no seu nome.
  - Usa um `sync.WaitGroup` para gerenciar a execução concorrente dos manipuladores.
  - Cada manipulador é executado em uma goroutine separada, permitindo o processamento assíncrono dos eventos.

### Remove

```go
func (dispatcher *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := dispatcher.handlers[eventName]; ok {
		for index, handlerOnDispatcher := range dispatcher.handlers[eventName] {
			if handlerOnDispatcher == handler {
				dispatcher.handlers[eventName] = append(dispatcher.handlers[eventName][:index], dispatcher.handlers[eventName][index+1:]...)
				return nil
			}
		}
	}

	return nil
}
```
- **Propósito**: Remove um manipulador específico de um evento.
- **Funcionalidade**:
  - Verifica se o evento possui manipuladores registrados.
  - Itera pelos manipuladores e remove o manipulador especificado se encontrado.

  
# Explicação do Código de Teste `dispatcher_test.go`

Este arquivo contém testes para a funcionalidade de um despachante de eventos (`EventDispatcher`). Vamos detalhar cada parte do código para entender melhor seu funcionamento.

## Estruturas Mock

### `EventMock`

```go
type EventMock struct {
	Name    string
	Payload any
}

func (eventMock *EventMock) GetName() string {
	return eventMock.Name
}

func (eventMock *EventMock) GetPayload() any {
	return eventMock.Payload
}

func (eventMock *EventMock) GetDateTime() time.Time {
	return time.Now()
}
```

`EventMock` é uma estrutura que simula um evento. Ela possui três métodos:
- `GetName()`: retorna o nome do evento.
- `GetPayload()`: retorna o payload do evento.
- `GetDateTime()`: retorna a data e hora atuais.

### `EventHandlerMock`

```go
type EventHandlerMock struct {
	ID string
}

func (eventHandlerMock *EventHandlerMock) Handle(event EventInterface, wg *sync.WaitGroup) {}
```

`EventHandlerMock` é uma estrutura que simula um manipulador de eventos. Ela possui um método:
- `Handle()`: método vazio que simula o manuseio de um evento.

## Suite de Testes

### `EventDispatcherSuiteTest`

```go
type EventDispatcherSuiteTest struct {
	suite.Suite
	event1          EventMock
	event2          EventMock
	handler1        EventHandlerMock
	handler2        EventHandlerMock
	handler3        EventHandlerMock
	eventDispatcher *EventDispatcher
}
```

`EventDispatcherSuiteTest` é uma estrutura que agrupa os testes. Ela possui:
- Eventos (`event1` e `event2`).
- Manipuladores de eventos (`handler1`, `handler2`, `handler3`).
- Uma instância do despachante de eventos (`eventDispatcher`).

### `SetupTest`

```go
func (suite *EventDispatcherSuiteTest) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler1 = EventHandlerMock{ID: "0001"}
	suite.handler2 = EventHandlerMock{ID: "0002"}
	suite.handler3 = EventHandlerMock{ID: "0003"}
	suite.event1 = EventMock{Name: "event mock 01", Payload: struct{ ID, name string }{ID: "01", name: "event mock 01"}}
	suite.event2 = EventMock{Name: "event mock 02", Payload: struct{ ID, name string }{ID: "02", name: "event mock 02"}}
}
```

`SetupTest` inicializa os objetos necessários para os testes.

### `TestSuite`

```go
func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherSuiteTest))
}
```

`TestSuite` executa a suite de testes.

## Testes

### `TestEventDispatcher_Register`

```go
func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Register() {
	// ...código de teste...
}
```

Testa o registro de manipuladores de eventos.

### `TestEventDispatcher_Register_SameHandler`

```go
func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Register_SameHandler() {
	// ...código de teste...
}
```

Testa o registro do mesmo manipulador para um evento.

### `TestEventDispatcher_Clear`

```go
func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Clear() {
	// ...código de teste...
}
```

Testa a limpeza de todos os manipuladores registrados.

### `TestEventDispatcher_Has`

```go
func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Has() {
	// ...código de teste...
}
```

Testa se um manipulador está registrado para um evento específico.

### `EventHandlerMocked`

```go
type EventHandlerMocked struct {
	mock.Mock
}

func (EventHandlerMocked *EventHandlerMocked) Handle(event EventInterface, wg *sync.WaitGroup) {
	EventHandlerMocked.Called(event)
	wg.Done()
}
```

`EventHandlerMocked` é uma estrutura que usa `mock.Mock` para simular o comportamento de um manipulador de eventos.

### `TestEventDispatcher_Dispatch`

```go
func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Dispatch() {
	// ...código de teste...
}
```

Testa o despacho de eventos para os manipuladores registrados.

### `TestEventDispatcher_Remove`

```go
func (suite *EventDispatcherSuiteTest) TestEventDispatcher_Remove() {
	// ...código de teste...
}
```

Testa a remoção de manipuladores de eventos.

# Explicação do Código `rabbitmq.go`

Este documento fornece uma explicação detalhada do código encontrado em `rabbitmq.go`, que implementa funções para interagir com o RabbitMQ. O RabbitMQ é um sistema de mensageria que permite que diferentes partes de uma aplicação se comuniquem de forma assíncrona. Abaixo, vamos detalhar os principais componentes e funções do código.

## Visão Geral

O código `rabbitmq.go` contém três funções principais:
1. `OpenChannel`: Abre uma conexão e um canal com o RabbitMQ.
2. `Consumer`: Consome mensagens de uma fila do RabbitMQ.
3. `Publish`: Publica mensagens em um exchange do RabbitMQ.

### OpenChannel

```go
func OpenChannel() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch, nil
}
```
- **Propósito**: Abre uma conexão e um canal com o RabbitMQ.
- **Funcionalidade**:
  - `amqp.Dial`: Estabelece uma conexão com o RabbitMQ usando a URL fornecida.
  - `conn.Channel()`: Cria um canal a partir da conexão estabelecida.
  - Retorna o canal aberto ou um erro, se houver.

### Consumer

```go
func Consumer(ch *amqp.Channel, out chan<- amqp.Delivery, queueName string) error {
	messages, err := ch.Consume(
		queueName,
		"go-consumer",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	for message := range messages {
		out <- message
	}
	return nil
}
```
- **Propósito**: Consome mensagens de uma fila do RabbitMQ e as envia para um canal Go.
- **Funcionalidade**:
  - `ch.Consume`: Inicia o consumo de mensagens da fila especificada (`queueName`).
  - As mensagens consumidas são enviadas para o canal `out`.
  - Retorna um erro, se houver.

### Publish

```go
func Publish(ch *amqp.Channel, body string, exchangeName string) error {
	err := ch.Publish(
		exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
```
- **Propósito**: Publica mensagens em um exchange do RabbitMQ.
- **Funcionalidade**:
  - `ch.Publish`: Publica uma mensagem no exchange especificado (`exchangeName`).
  - A mensagem é enviada com o corpo (`body`) e o tipo de conteúdo (`ContentType`) especificados.
  - Retorna um erro, se houver.


# Explicação do Código `main.go` do Consumidor de Mensagens RabbitMQ

Este documento fornece uma explicação detalhada do código encontrado em `main.go`, que implementa um consumidor de mensagens do RabbitMQ. Abaixo, vamos detalhar os principais componentes e funções do código.

## Visão Geral

O código `main.go` contém a função principal (`main`) que configura e executa um consumidor de mensagens do RabbitMQ.

### main

```go
func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	messages := make(chan amqp.Delivery)

	go rabbitmq.Consumer(ch, messages, "queue1")

	for message := range messages {
		fmt.Println(string(message.Body))
		message.Ack(false)
	}
}
```
- **Propósito**: Configura e executa um consumidor de mensagens do RabbitMQ.
- **Funcionalidade**:
  - `rabbitmq.OpenChannel()`: Abre uma conexão e um canal com o RabbitMQ.
  - `defer ch.Close()`: Garante que o canal será fechado ao final da execução.
  - `make(chan amqp.Delivery)`: Cria um canal Go para receber mensagens.
  - `go rabbitmq.Consumer(ch, messages, "queue1")`: Inicia uma goroutine que consome mensagens da fila `queue1` e as envia para o canal `messages`.
  - `for message := range messages`: Itera sobre as mensagens recebidas.
    - `fmt.Println(string(message.Body))`: Imprime o corpo da mensagem.
    - `message.Ack(false)`: Confirma o recebimento da mensagem.


# Explicação do Código `main.go` do Produtor de Mensagens RabbitMQ

Este documento fornece uma explicação detalhada do código encontrado em `main.go`, que implementa um produtor de mensagens do RabbitMQ. Abaixo, vamos detalhar os principais componentes e funções do código.

## Visão Geral

O código `main.go` contém a função principal (`main`) que configura e executa um produtor de mensagens do RabbitMQ.

### main

```go
func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	rabbitmq.Publish(ch, "Hello World", "amq.direct")
}
```
- **Propósito**: Configura e executa um produtor de mensagens do RabbitMQ.
- **Funcionalidade**:
  - `rabbitmq.OpenChannel()`: Abre uma conexão e um canal com o RabbitMQ.
  - `defer ch.Close()`: Garante que o canal será fechado ao final da execução.
  - `rabbitmq.Publish(ch, "Hello World", "amq.direct")`: Publica uma mensagem com o conteúdo "Hello World" no exchange `amq.direct`.
  