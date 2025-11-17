# GoQueue

Pequena biblioteca em Go para criação de filas baseadas em *channels*, com observabilidade de trabalhos e controle manual de execução.

## Conceitos
- **Fila**: representada por um `QueueManager`, expõe um canal (`WorkQueue`) onde os trabalhos ficam disponíveis para consumo.
- **Work**: item processável; possui estados (`pending`, `running`, `done`, `failed`) e carrega *payloads* de entrada e saída.
- **Expiração**: cada fila pode ter um tempo de expiração configurado; você pode inspecionar quanto falta ou se já expirou.
- **Observação**: o `Engine` emite novos `QueueManager` via `WatchQueues` e cada fila pode ser inspecionada por contadores de estado.

## Instalação
```bash
go get github.com/gothout/goqueue
```

## Uso rápido
```go
package main

import (
"context"
"fmt"
"time"

"github.com/gothout/goqueue"
)

func main() {
        engine := goqueue.NewEngine()

        // cria ou recupera a fila e define expiração de 1 minuto
        queue := engine.GetOrCreateQueue("protocolo-123")
        queue.SetExpiration(time.Minute)

        // enfileira manualmente um trabalho
        work := goqueue.NewWork(queue.GetKey(), "payload inicial")
        _ = queue.AddWork(work)

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        // consumo manual: bloqueia até chegar um trabalho ou o contexto expirar
        next, ok := queue.NextWork(ctx)
        if !ok {
                panic("fila vazia ou cancelada")
        }

        // controle de execução manual
        next.Start()
        next.SetOutboundPayload("resultado")
        next.Finish()

        fmt.Printf("Trabalho %s terminou com estado %s\n", next.UUID, next.GetState())
}
```

## Padrão de trabalho manual
O fluxo típico é:
1. Enfileirar um `*Work` com `AddWork`.
2. Ler do canal `WorkQueue` (ou usar `NextWork` com `context.Context`).
3. Chamar `Start()` quando a execução iniciar, seguido de `Finish()` ou `Fail(reason)`.
4. Usar os contadores da fila para monitorar o progresso (`CountPendingWorks`, `CountRunningWorks`, etc.).

## Observando filas
```go
engine := goqueue.NewEngine()
updates := engine.WatchQueues()

// em outra goroutine, observe novas filas criadas
 go func() {
         for q := range updates {
                 fmt.Printf("fila criada: %s\n", q.GetKey())
         }
 }()

// cria fila e dispara evento
engine.GetOrCreateQueue("notificacoes")
```

## Expiração
```go
queue := engine.GetOrCreateQueue("uploads")
queue.SetExpiration(30 * time.Second)

for {
        if queue.HasExpired() {
                fmt.Println("fila expirada; considere descartá-la")
                break
        }
        fmt.Printf("faltam %s para expirar\n", queue.GetTimeToExpire())
        time.Sleep(5 * time.Second)
}
```

## Execução em processos separados
Como `WorkQueue` é um `chan *Work`, você pode compartilhar a fila entre processos/goroutines, desde que o transporte seja feito pelo seu mecanismo preferido (gRPC, HTTP, etc.). O consumidor apenas lê do canal e controla manualmente o estado do trabalho.

```go
// produtor
queue.AddWork(goqueue.NewWork("payout", "payload"))

// consumidor
ctx := context.Background()
if work, ok := engine.NextWork(ctx, "payout"); ok {
        work.Start()
        // ... executa ...
        work.Finish()
}
```

## Interfaces expostas
- `EngineManager` provê: `GetOrCreateQueue`, `GetQueue`, `DeleteQueue`, `ListQueues`, `WatchQueues`, `NextWork`.
- `QueueManagerInterface` compõe: expiração, operações (`AddWork`), inspeção.
- `WorkReader`, `WorkWriter`, `WorkController` permitem leitura e controle dos trabalhos.

## Notas
- `QueueManager.Close()` encerra o canal `WorkQueue` e impede novas inserções.
- `AddWork` respeita a capacidade configurada e reenvia em goroutine se o canal estiver momentaneamente cheio.
