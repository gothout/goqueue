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

	// Produzindo trabalhos para filas diferentes.
	emails := engine.GetOrCreateQueue("emails")
	sms := engine.GetOrCreateQueue("sms")

	emails.AddWork(goqueue.NewWork("emails", "enviar email de boas-vindas"))
	emails.AddWork(goqueue.NewWork("emails", "enviar fatura mensal"))
	sms.AddWork(goqueue.NewWork("sms", "enviar SMS de verificação"))

	// Consumindo os trabalhos com um contexto que expira após 5 segundos.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		work, key, ok := engine.NextWorkFromQueues(ctx, []string{"emails", "sms"})
		if !ok {
			fmt.Println("Nenhum trabalho encontrado ou contexto cancelado")
			return
		}

		fmt.Printf("Processando da fila %s: %s\n", key, work.GetInboundPayload())
		work.Start()

		// Simula um processamento curto e marca como concluído.
		time.Sleep(250 * time.Millisecond)
		work.Finish()

		fmt.Printf("Finalizado: %s (estado: %s)\n", work.UUID, work.GetState())

		if emails.CountPendingWorks()+sms.CountPendingWorks() == 0 {
			fmt.Println("Todos os trabalhos foram processados")
			return
		}
	}
}
