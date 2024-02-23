# Golang Microservices

This mini project showcases implementation of Microservices architecture, with multiple services written in Go lang

Services implemented are.
1. Broker service : a point of communication for the external services/frontend
2. Authentication service: Validating user credentials stored in the PostgresSQL db
3. Logger service: A logging service to be used by other internal service, that logs events in MongoDB db
4. Listener service: A event driven implementation using RabbitMQ. 
5. Frontend service: Hosts a web page used to test the services , request and response.
6. Mailer service: To send an email to a Mailhog mail server.

Communication between the services are done via a HTTP API call using JSON for most of the services.
RPC and gRPC implementation is done for logger service. 

I used docker to run the entire system. each folder has the individual docker file for it 

### How to run

1. Clone the repo
2. Navigate to the project folder
   In the terminal:
   ```sh
   $ docker compose up
   ```
   For the authenticate service to work. we will first need to run the users.sql file present in others folder. Use a SQL editor of your choice and run the sql script.

   Then re-run the above command.
4. For the frontend web page:
   Navigate to the front-end folder.
   ```sh
   $ go run ./cmd/web/main.go
   ```
   Open http://localhost in your web browser.
   Each service can be tested by clicking on the respective button.

<img width="1025" alt="image" src="https://github.com/rayancrasta/Golang_Microservices/assets/43010629/3d3fe052-3055-4de3-a812-6632405a5b29">



This is how the output will look like in the browser, when all services are working correctly.


This mini project was developed as a part of learning from [Trevor Sawlers Udemy Course](https://www.udemy.com/course/working-with-microservices-in-go/) on Go lang microservices. 


