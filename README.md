Kakeibo is a minimalist expense tracker designed to help you manage your finances simply and effectively. The name comes from the Japanese art of saving money through mindful budgeting.

<img width="984" height="767" alt="Screenshot 2026-01-20 at 2 26 06 AM" src="https://github.com/user-attachments/assets/e1f733aa-ce5a-4be3-a0d2-bd4cbb1c2e98" />

## Installation

via docker-compose

```bash
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
