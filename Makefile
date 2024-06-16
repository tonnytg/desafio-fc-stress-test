all:
	# Compilar a imagem Docker
	docker build -t stress-test .

run:
	# Executar a aplicação Docker
	docker run stress-test --url=http://google.com --requests=1000 --concurrency=10