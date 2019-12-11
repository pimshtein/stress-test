package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
	"os"
	"sort"
	"stress-test/internal/amqp"
	"stress-test/internal/worker"
	"time"
)

var config worker.Config

type httpResult struct {
	index   int
	httpRe     *http.Response
	httpErr *int
	err     error
}

type rabbitResult struct {
	index   int
	err     error
}

var (
	payload []byte
	uris [] string
	producer *amqp.Producer
)

func init() {
	parser := flags.NewParser(nil, flags.Default)
	parser.ShortDescription = "Go stress tool"

	_, err := parser.AddGroup("Rabbit consumer options", "", &config.Rabbit)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = parser.AddGroup("Parallel requests options", "", &config.ParallelRequests)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	for i := 1; i <= config.ParallelRequests.RequestsNum; i++ {
		uris = append(uris, "interceptor.local")
	}

	payload = []byte(`
	{
      "header": {
		"id": "9a9e99a0-f30a-49e7-8c25-b93acb2fd939",
		"type": "/HbPro/Hotel/CommercialOffer/Event/CommercialOfferSendingCreated/Message",
		"time": "2019-12-10T10:37:32.676+00:00",
		"meta": {
		  "eventName": "CommercialOfferSendingCreated",
		  "locale": "ru"
		}
	  },
	  "body": {
		"recipient": "asdf@asdf.com",
		"city": "Москва",
		"checkIn": "2019-12-14",
		"checkOut": "2019-12-15",
		"attachments": [
		  {
			"name": "commercialOffer.pdf",
			"fileBase64": "aq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0SZU2qdImVdqkSptUaZMqbVKlTaq0S"
		  }
		]
	  }
	}
	`)

	producer = amqp.NewProducer(config.Rabbit.GetProducerSettings())
}

func main() {
	benchmark := func(uris []string, concurrency int) string {
		startTime := time.Now()
		results := produceParallelRequests(uris, concurrency)
		seconds := time.Since(startTime).Seconds()
		countError := 0
		for _, result := range results {
			if nil != result.err {
				countError++
				fmt.Println(result.err)
			}
		}
		template := "%d bounded parallel requests: %d/%d in %v, count of errors: %d"
		return fmt.Sprintf(template, concurrency, len(results), len(uris), seconds, countError)
	}

	fmt.Println(benchmark(uris, config.ParallelRequests.GoNum))
}

func produceParallelRequests(uris []string, concurrencyLimit int) []rabbitResult {
	semaphoreChan := make(chan struct{}, concurrencyLimit)

	resultsChan := make(chan *rabbitResult)

	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	for i, uri := range uris {

		go func(i int, url string) {

			// Buff canal (when it's full, goroutines are waiting)
			semaphoreChan <- struct{}{}

			err := producer.PublishMessage(payload, uri, "text/json")

			result := &rabbitResult{i, err}

			resultsChan <- result

			<-semaphoreChan

		}(i, uri)
	}

	var results []rabbitResult

	for {
		result := <-resultsChan
		results = append(results, *result)

		if len(results) == len(uris) {
			break
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	return results
}