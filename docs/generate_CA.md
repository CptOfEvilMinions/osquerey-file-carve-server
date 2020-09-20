# Generate Root CA with OpenSSL

## Generate root CA
1. `openssl genrsa -aes256 -out conf/tls/root_ca/snakeoil_root_CA.key 2048`
  1. Generate root CA key
  1. Enter password
1. `openssl req -x509 -new -nodes -key conf/tls/root_ca/snakeoil_root_CA.key -sha256 -days 3650 -out conf/tls/root_ca/snakeoil_root_CA.crt`
  1. Generate root CA cert
  1. Enter details about org

## Generate leaf cert for Kolide
1. `openssl genrsa -out conf/tls/kolide/snakeoil_kolide_hackinglab_local.key 2048`
  1. Generate private key for Kolide server
1. `openssl req -new -key conf/tls/kolide/snakeoil_kolide_hackinglab_local.key -out conf/tls/kolide/snakeoil_kolide_hackinglab_local.csr`
  1. Generate signing request
  1. Enter details about Kolide
1. `openssl x509 -req -in conf/tls/kolide/snakeoil_kolide_hackinglab_local.csr -CA conf/tls/root_ca/snakeoil_root_CA.crt -CAkey conf/tls/root_ca/snakeoil_root_CA.key -CAcreateserial -out conf/tls/kolide/snakeoil_kolide_hackinglab_local.crt -days 1095 -sha256`
  1. Generate public cert for Kolide
1. `cat conf/tls/root_ca/snakeoil_root_CA.crt >> conf/tls/kolide/snakeoil_kolide_hackinglab_local.crt`
  1. Add root CA to chain

## Generate device cert
1. `openssl genrsa -out conf/tls/device/snakeoil_device01.key 2048`
  1. Generate private key for device
1. `openssl req -new -key conf/tls/device/snakeoil_device01.key -out conf/tls/device/snakeoil_device01.csr`
  1. Generate signing request
  1. Enter details about devive
1. `openssl x509 -req -in conf/tls/device/snakeoil_device01.csr -CA conf/tls/root_ca/snakeoil_root_CA.crt -CAkey conf/tls/root_ca/snakeoil_root_CA.key -CAcreateserial -out conf/tls/device/snakeoil_device01.crt -days 365 -sha256`


## Generate PKSC12
1. `openssl pkcs12 -export -out conf/tls/device/snakeoil_device01.p12 -inkey conf/tls/device/snakeoil_device01.key -in conf/tls/device/snakeoil_device01.crt -certfile conf/tls/root_ca/snakeoil_root_CA.crt`
  1. Generate P12
  1. Enter password

## Generate cert bundle for Osquery
1. `cat conf/tls/kolide/snakeoil_kolide_hackinglab_local.crt > conf/tls/osquery/snakeoil_osquery_server_cert_chain.crt`
1. `cat conf/tls/root_ca/snakeoil_root_CA.crt >> conf/tls/osquery/snakeoil_osquery_server_cert_chain.crt` 

## References
* [How to Create Your Own SSL Certificate Authority for Local HTTPS Development](https://deliciousbrains.com/ssl-certificate-authority-for-local-https-development/)
* [Creating the Intermediate CA](https://roll.urown.net/ca/ca_intermed_setup.html)