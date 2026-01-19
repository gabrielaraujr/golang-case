.PHONY: configure logs test clean

configure:
	docker compose up -d --build

logs:
	docker compose logs -f

test:
	cd account && go test ./...

clean:
	docker compose down -v
	rm -rf account/bin risk-analysis/bin
	go clean -cache

create-proposal:
	curl -X POST http://localhost:8001/proposals \
		-H "Content-Type: application/json" \
		-d '{"full_name":"Test User","cpf":"12345678902","salary":5000.00,"email":"test@email.com","phone":"11999999999","birthdate":"02-06-2016","address":{"street":"Rua Teste 123","city":"Sao Paulo","state":"SP","zip_code":"01234567"}}'

check-queue:
	docker exec localstack awslocal sqs receive-message --queue-url http://localhost:4566/000000000000/proposals --max-number-of-messages 10

check-results:
	docker exec localstack awslocal sqs receive-message --queue-url http://localhost:4566/000000000000/risk-results --max-number-of-messages 10
