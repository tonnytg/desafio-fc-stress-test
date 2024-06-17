package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Definindo os parâmetros CLI
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 1, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("A URL do serviço é obrigatória.")
		flag.Usage()
		return
	}

	// Variáveis para armazenar resultados
	var totalRequests, successRequests, totalTime int64
	statusCodes := make(map[int]int)
	var mutex sync.Mutex

	// Função para executar um request
	worker := func(wg *sync.WaitGroup, ch chan struct{}) {
		defer wg.Done()
		for range ch {
			start := time.Now()
			resp, err := http.Get(*url)
			elapsed := time.Since(start).Milliseconds()

			mutex.Lock()
			totalRequests++
			totalTime += elapsed

			statusCodes[resp.StatusCode]++

			if err == nil {

				if resp.StatusCode == http.StatusOK {
					successRequests++
				}
				resp.Body.Close()
			}
			mutex.Unlock()
		}
	}

	// Canal para distribuir as tarefas
	ch := make(chan struct{}, *requests)
	var wg sync.WaitGroup

	// Iniciando workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(&wg, ch)
	}

	// Enviando requests para o canal
	for i := 0; i < *requests; i++ {
		ch <- struct{}{}
	}
	close(ch)

	// Aguardando a finalização dos workers
	wg.Wait()

	// Gerando o relatório
	fmt.Println("Relatório de Teste de Carga:")
	fmt.Printf("Tempo total gasto: %d ms\n", totalTime)
	fmt.Printf("Total de requests realizados: %d\n", totalRequests)
	fmt.Printf("Total de requests com status 200: %d\n", successRequests)
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for code, count := range statusCodes {
		if code != http.StatusOK {
			fmt.Printf("\t\t\t\tStatus %d: %d\n", code, count)
		}
	}
}
