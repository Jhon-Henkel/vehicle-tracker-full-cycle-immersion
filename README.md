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