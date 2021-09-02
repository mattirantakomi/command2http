#!/usr/bin/env bash

set -euo pipefail

docker build . -t "eu.gcr.io/$1/command2http-demo"
docker push "eu.gcr.io/$1/command2http-demo"

gcloud run deploy command2http-demo --image "eu.gcr.io/$1/command2http-demo" --project "$1" --region europe-north1 --allow-unauthenticated --cpu=1 --memory=128 --max-instances=3 --timeout=5 --concurrency=1