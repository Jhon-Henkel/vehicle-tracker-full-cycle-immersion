# vehicle-tracker-full-cycle-immersion
Repositório destinado a uma aplicação com micro-serviços  simulando uma plataforma de rastreamento veicular

# Tecnologias
- TypeScript/JavaScript
- NestJS
- Prisma ORM
- MongoDB
- Google Maps API
- React
- React Server Components
- NextJS
- Route Handler
- Material UI
- Socket.io
- Kafka
- Docker
- Bulljs no NestJS
- Prometheus

# Temporário
Adicionar a linha abaixo no seu hosts (apenas linux e mac):

```
127.0.0.1 host.docker.internal
```

### Nest
- Iniciando o nest: npm run start:dev
- criando tablelas mongo: npx prisma generate
- acessando o banco via web: npx prisma studio

### Next
- Paginas:
    - new-route
    - driver
- Iniciando o next: npm run dev

### Kafka
- acesso pelo host http://localhost:9021/clusters

### GO
- Criar tableas no mysql do GO:
    ```
    CREATE TABLE routes (id VARCHAR(36) PRIMARY KEY, name VARCHAR(255) NOT NULL, distance FLOAT NOT NULL, status VARCHAR(255) NOT NULL, freight_price FLOAT, started_at DATETIME, finished_at DATETIME);
    ```
- Iniciando o go: go run cmd/freight/main.go
