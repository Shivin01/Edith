#!/bin/bash

set -ex

# it will start jenkins and postgres server which is required for python backend
docker-compose up -d

# Sleeping for 60 seconds for postgres to start fully
sleep 60

# Run migration
python manage.py migrate

# Run backend server
python manage runserver
