#!/bin/bash

sudo docker-compose down
sudo docker rmi project_authentication-service project_logger-service project_broker-service project_mailer-service
sudo docker-compose up

