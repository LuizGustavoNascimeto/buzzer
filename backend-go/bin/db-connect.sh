#! /usr/bin/bash
if [ "$1" = "prod" ]; then
  echo "Running in production mode"
  URL=$AWS_RDS
else
  URL=$POSTGRES_URL
fi

psql $URL