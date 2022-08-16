test:
	curl -X POST 127.0.0.1:8080 -H 'Content-Type: application/json' -d '{"record":{"value":""}}'
	curl -X POST 127.0.0.1:8080 -H 'Content-Type: application/json' -d '{"record":{"value":""}}'
	curl -X GET 127.0.0.1:8080 -d '{"offset":0}'
	curl -X GET 127.0.0.1:8080 -d '{"offset":1}'
	curl -X GET 127.0.0.1:8080 -d '{"offset":1}'