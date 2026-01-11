# Docker Dev Mode (Hot Reload)
dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# Docker Production Mode
prod:
	docker compose up --build

# Stop Containers
down:
	docker compose down

frontend:
	cd frontend && npm run dev

backend:
	cd backend && air

.PHONY: frontend backend dev prod down
