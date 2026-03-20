.PHONY: up down restart logs clean certs

certs:
	@if [ ! -f nginx/certs/selfsigned.crt ]; then \
		mkdir -p nginx/certs && \
		openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
			-keyout nginx/certs/selfsigned.key \
			-out nginx/certs/selfsigned.crt \
			-subj "/CN=*.api.co.th" \
			-addext "subjectAltName=DNS:*.api.co.th,DNS:*.127.0.0.1.nip.io" \
			2>/dev/null && \
		echo "SSL certs generated"; \
	fi

up: certs
	docker compose -f docker/docker-compose.yml up --build -d
	@echo ""
	@echo "===== API Ready ====="
	@echo "Hospital A: https://hospital-a.127.0.0.1.nip.io"
	@echo "Hospital B: https://hospital-b.127.0.0.1.nip.io"
	@echo ""
	@echo "Login:  POST /login  {\"username\":\"admin\",\"password\":\"password123\"}"
	@echo "====================="

down:
	docker compose -f docker/docker-compose.yml down

restart: down up

logs:
	docker compose -f docker/docker-compose.yml logs -f

clean:
	docker compose -f docker/docker-compose.yml down -v
	rm -rf nginx/certs
