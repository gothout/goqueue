package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gothout/goqueue"
)

// Este exemplo mostra um fluxo simples de produção e consumo utilizando a Engine central.
func main() {
	engine := goqueue.NewEngine()

	// Produzindo trabalhos para uma fila chamada "emails".
	queue := engine.GetOrCreateQueue("emails")
	queue.AddWork(goqueue.NewWork("emails", "enviar email de boas-vindas"))
	queue.AddWork(goqueue.NewWork("emails", "enviar fatura mensal"))

	// Consumindo os trabalhos com um contexto que expira após 5 segundos.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		work, ok := engine.NextWork(ctx, "emails")
		if !ok {
			fmt.Println("Nenhum trabalho encontrado ou contexto cancelado")
			return
		}

		fmt.Printf("Processando: %s\n", work.GetInboundPayload())
		work.Start()

		// Simula um processamento curto e marca como concluído.
		time.Sleep(250 * time.Millisecond)
		work.Finish()

		fmt.Printf("Finalizado: %s (estado: %s)\n", work.UUID, work.GetState())

		if queue.CountPendingWorks() == 0 {
			fmt.Println("Todos os trabalhos foram processados")
			return
		}
	}
}
