# Learn-rabbitMQ

docker run -d --name rabbitMQ -p 5672:5672 -p 15672:15672 rabbitmq:4.0-management

<!-- Create vhost -->

docker exec rabbitMQ rabbitmqctl add_vhost customer

<!-- Declare exchange -->

docker exec rabbitMQ rabbitmqadmin declare exchange --vhost=customer name=customer_event type=topic -u guest -p guest durable=true

<!-- Set permission in new vhost -->

docker exec rabbitMQ rabbitmqctl set_permissions -p customer guest "._" "._" ".\*"
