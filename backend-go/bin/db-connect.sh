#! /usr/bin/bash
if [ "$1" = "prod" ]; then
  echo "Running in production mode"
  URL=$PROD_POSTGRES_URL
else
  URL=$POSTGRES_URL
fi

psql $URL