# stress-test

How to use:  
* cd cmd/stress-test
* go build
* For example publish messages to amqp:  
```./stress-test --amqp-exchange-type=direct --amqp-exchange-name=exchange.name --amqp-url=amqp://login:password@127.0.0.1:5672/%2F --go-num=100 --requests-num=100000```

