```
██╗░░██╗░█████╗░██╗░░██╗███████╗██╗██████╗░░█████╗░
██║░██╔╝██╔══██╗██║░██╔╝██╔════╝██║██╔══██╗██╔══██╗
█████═╝░███████║█████═╝░█████╗░░██║██████╦╝██║░░██║
██╔═██╗░██╔══██║██╔═██╗░██╔══╝░░██║██╔══██╗██║░░██║
██║░╚██╗██║░░██║██║░╚██╗███████╗██║██████╦╝╚█████╔╝
╚═╝░░╚═╝╚═╝░░╚═╝╚═╝░░╚═╝╚══════╝╚═╝╚═════╝░░╚════╝░
```

kakeibo is a minimalist expense tracker designed to help you manage your finances simply and effectively. The name comes from the Japanese art of saving money through mindful budgeting.

<img width="899" height="766" alt="Screenshot 2026-01-25 at 12 46 22 PM" src="https://github.com/user-attachments/assets/8f45449e-eff4-4650-b6c7-b2977e5bf721" />

## Installation

via docker-compose

```yaml
services:
  kakeibo:
    container_name: kakeibo
    image: ghcr.io/manosriram/kakeibo:latest
    # volumes:
      # - ./data:/app/data
    ports:
      - "5464:8080"
    env_file:
      - .env
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 30s
    restart: always
    networks:
      - kakeibo_network

  qdrant:
    image: qdrant/qdrant:latest
    container_name: qdrant_server
    restart: always
    ports:
      - "6333:6333" # HTTP API
      - "6334:6334" # gRPC API
    volumes:
      - ./qdrant_storage:/qdrant/storage
    networks:
      - kakeibo_network

networks:
  kakeibo_network:
    driver: bridge
```

## Usage

kakiebo supports query via bots:
  1. Telegram

## Steps to use telegram bot

1. Set TELEGRAM_BOT_ID in .env

2. From the bot, you'll have 2 commands:
    ```
    /track: Add an expense
    /summary: Summarize the current month's expenses
    ```


