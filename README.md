# Kakeibo

Kakeibo is a minimalist expense tracker designed to help you manage your finances simply and effectively. The name comes from the Japanese art of saving money through mindful budgeting.

## Installation

docker-compose
```yml
services:
  app:
    container_name: kakeibo
    image: ghcr.io/manosriram/kakeibo:latest
    volumes:
      - ./data:/app/data
    ports:
      - "5464:8080"
    env_file:
      - .env
```
