#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# CREATE THE PRIVATE KEY FOR OUR CUSTOM CA
openssl genrsa -out certs/ca.key 2048

# GENERATE A CA CERT WITH THE PRIVATE KEY
openssl req -new -x509 -key certs/ca.key -out certs/ca.crt -config certs/lbp_config.txt

# CREATE THE PRIVATE KEY FOR OUR GRUMPY SERVER
openssl genrsa -out certs/lbp-key.pem 2048

# CREATE A CSR FROM THE CONFIGURATION FILE AND OUR PRIVATE KEY
openssl req -new -key certs/lbp-key.pem -subj "/CN=lbp.default.svc" -out certs/lbp.csr -config certs/lbp_config.txt

# CREATE THE CERT SIGNING THE CSR WITH THE CA CREATED BEFORE
openssl x509 -req -in certs/lbp.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/lbp-crt.pem

# INJECT CA IN THE WEBHOOK CONFIGURATION
export CA_BUNDLE=$(cat certs/ca.crt | base64 | tr -d '\n')
export LBP_KEY=$(cat certs/lbp-key.pem | base64 | tr -d '\n')
export LBP_CRT=$(cat certs/lbp-crt.pem | base64 | tr -d '\n')
cat ../deploy/ValidatingWebhookConfiguration.yaml.tmpl | envsubst > ../deploy/ValidatingWebhookConfiguration.yaml
cat ../deploy/secret.yaml.tmpl | envsubst > ../deploy/secret.yaml
