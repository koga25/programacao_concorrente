Para essa atividade, é necessário que o RabbitMQ esteja instalado na máquina e rodando no localhost, na porta padrão (5672).

Para maior facilidade, é possível utilizar o container de docker do RabbitMQ:
```bash
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.10-management
```