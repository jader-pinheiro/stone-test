package utils

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"stone-test/internal/infra/data"
	"stone-test/internal/infra/entity"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/schollz/progressbar/v3"
)

const (
	folderPath    = "./file"
	batchSize     = 13000
	maxWorkers    = 7
	insertWorkers = 2
)

func ProcessFileContent(ctx context.Context, conn *pgxpool.Pool) (string, error) {
	startBenchmark := time.Now()

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return "", fmt.Errorf("erro ao ler pasta: %w", err)
	}

	var wg sync.WaitGroup
	stocksCh := make(chan entity.Stocks, 5000)
	sem := make(chan struct{}, maxWorkers)

	var batchPool = sync.Pool{
		New: func() interface{} {
			return make([]entity.Stocks, 0, batchSize)
		},
	}

	insertWg := sync.WaitGroup{}
	insertWg.Add(insertWorkers)
	for i := 0; i < insertWorkers; i++ {
		go func() {
			defer insertWg.Done()

			batch := batchPool.Get().([]entity.Stocks)[:0]

			for stock := range stocksCh {
				batch = append(batch, stock)
				if len(batch) >= batchSize {
					if err := data.InsertBatchCopy(ctx, conn, batch); err != nil {
						panic(fmt.Errorf("erro na inserção batch: %w", err))
					}
					batch = batch[:0]
				}
			}

			if len(batch) > 0 {
				if err := data.InsertBatchCopy(ctx, conn, batch); err != nil {
					panic(fmt.Errorf("erro na inserção batch final: %w", err))
				}
			}

			batchPool.Put(batch[:0])
		}()
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(file fs.DirEntry) {
			defer wg.Done()
			defer func() { <-sem }()

			path := filepath.Join(folderPath, file.Name())

			f, err := os.Open(path)
			if err != nil {
				fmt.Printf("Erro ao abrir arquivo %s: %v\n", path, err)
				return
			}
			defer f.Close()

			// linhas do arquivo para a progress bar
			lineCount := 0
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				lineCount++
			}
			if err := scanner.Err(); err != nil {
				fmt.Printf("Erro lendo arquivo %s: %v\n", path, err)
				return
			}

			f.Seek(0, 0)
			reader := bufio.NewReader(f)

			bar := progressbar.Default(int64(lineCount-1), fmt.Sprintf("📄 %s", file.Name()))
			lineNum := 0
			parseErrors := 0

			for {
				lineBytes, _, err := reader.ReadLine()
				if err != nil {
					if err.Error() != "EOF" {
						fmt.Printf("Erro lendo arquivo %s: %v\n", path, err)
					}
					break
				}

				lineNum++
				if lineNum == 1 {
					continue
				}

				line := string(lineBytes)
				stock, err := ParseLine(line)
				if err != nil {
					fmt.Printf("Erro parseando linha %d em %s: %v\n", lineNum, file.Name(), err)
					continue
				}

				stocksCh <- stock
				bar.Add(1)
			}
			if parseErrors == 0 {
				fmt.Printf("✅ Arquivo %s processado com sucesso. Removendo...\n", file.Name())
				if err := os.Remove(path); err != nil {
					fmt.Printf("❌ Erro ao remover arquivo %s: %v\n", path, err)
				}
			} else {
				fmt.Printf("⚠️ Arquivo %s processado com %d erros de parsing. Não será removido.\n", file.Name(), parseErrors)
			}
		}(file)
	}

	wg.Wait()
	close(stocksCh)
	insertWg.Wait()

	totalTime := time.Since(startBenchmark)
	fmt.Printf("⏱️ Tempo total de inserção: %s\n", totalTime)

	return "Processamento concluído com sucesso!", nil
}

func ParseLine(line string) (entity.Stocks, error) {
	fields := strings.Split(line, ";")
	if len(fields) < 11 {
		return entity.Stocks{}, fmt.Errorf("linha com menos campos que o esperado")
	}

	brLoc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return entity.Stocks{}, fmt.Errorf("erro ao carregar timezone do Brasil: %v", err)
	}

	businessDateRaw := fields[8]
	businessDate, err := time.ParseInLocation("2006-01-02", businessDateRaw, brLoc)
	if err != nil {
		return entity.Stocks{}, fmt.Errorf("erro ao parsear DataNegocio: %v", err)
	}

	// InstrumentCode (CodigoInstrumento), campo 1 (index 1)
	instrumentCode := fields[1]

	// BusinessPrice (PrecoNegocio), campo 3 (index 3)
	businessPrice, err := ParseBrazilianFloat(fields[3])
	if err != nil {
		return entity.Stocks{}, fmt.Errorf("erro ao parsear PrecoNegocio: %v", err)
	}

	negQty, err := strconv.ParseInt(fields[4], 10, 64)
	if err != nil {
		return entity.Stocks{}, fmt.Errorf("erro ao parsear QuantidadeNegociada como inteiro: %v", err)
	}

	closingStr := fields[5]
	if len(closingStr) != 9 {
		return entity.Stocks{}, fmt.Errorf("HoraFechamento inválida: %s", closingStr)
	}

	hh, err := strconv.Atoi(closingStr[0:2])
	if err != nil {
		return entity.Stocks{}, fmt.Errorf("hora inválida: %v", err)
	}
	mm, err := strconv.Atoi(closingStr[2:4])
	if err != nil {
		return entity.Stocks{}, fmt.Errorf("minuto inválido: %v", err)
	}
	ss, err := strconv.Atoi(closingStr[4:6])
	if err != nil {
		return entity.Stocks{}, fmt.Errorf("segundo inválido: %v", err)
	}
	ms, err := strconv.Atoi(closingStr[6:9])
	if err != nil {
		return entity.Stocks{}, fmt.Errorf("milissegundo inválido: %v", err)
	}

	closingTime := time.Date(
		businessDate.Year(), businessDate.Month(), businessDate.Day(),
		hh, mm, ss, ms*1_000_000, time.UTC,
	)

	return entity.Stocks{
		BusinessDate:       businessDate,
		InstrumentCode:     instrumentCode,
		BusinessPrice:      businessPrice,
		NegotiatedQuantity: negQty,
		ClosingTime:        closingTime,
	}, nil
}

func ParseBrazilianFloat(s string) (float64, error) {
	normalized := strings.ReplaceAll(s, ".", "")
	normalized = strings.ReplaceAll(normalized, ",", ".")
	return strconv.ParseFloat(normalized, 64)
}
