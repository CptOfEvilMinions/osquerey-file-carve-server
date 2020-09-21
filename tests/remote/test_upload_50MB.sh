#!/bin/bash

# If file does not exist download test file
if [ ! -f /tmp/orig_50MB.zip ]; then
  echo "[*] - ${date} - Downloading 50MB.zip test file"
  wget http://ipv4.download.thinkbroadband.com/50MB.zip -O /tmp/orig_50MB.zip
fi

if [ -f /tmp/orig_50MB.zip ] && [ ! -f /tmp/50MB.zip ]; then
  echo "[*] - ${date} - Make copy of test file"
  cp /tmp/orig_50MB.zip /tmp/50MB.zip 
fi

# Issue Osqueryi file carve
echo '[+] - Start Osquery file carve'
osquery_file_carve_output=`osqueryi --flagfile conf/osquery/osquery.flags --json "SELECT * FROM carves WHERE path like '/tmp/50MB.zip' AND carve=1; SELECT * FROM carves WHERE path like '/tmp/%';"`
wait
echo $osquery_file_carve_output
osquery_carve_id=`echo ${osquery_file_carve_output} | jq '.[].carve_guid' | tr -d '"'`
echo ${osquery_carve_id}
echo '[+] - ${date} - Osquery file carve finished'

# Get Vault token
token_accessor=`vault token lookup --format=json | jq -r '.data.accessor'`
vault_token=`vault token lookup --format=json | jq -r '.data.id'`

# Download test file
curl -k https://kolide.hackinglab.local:8443/file_request \
-A "osquery/4.3.0" \
--cert conf/tls/device/snakeoil_device01.crt \
--key conf/tls/device/snakeoil_device01.key \
--cacert conf/tls/root_ca/snakeoil_root_CA.crt \
-d "{ \"file_carve_guid\": \"${osquery_carve_id}\", \"token_accessor\": \"${token_accessor}\", \"token\": \"${vault_token}\"}" \
--output "/tmp/${osquery_carve_id}.tar"


# Check file carve exists
if [ -f "/tmp/${osquery_carve_id}.tar" ]; then
  echo "[+] - ${date} - File Carve upload exists - /tmp/${osquery_carve_id}.tar"
else
  echo "[-] - ${date} - File Carve upload does NOT exists"
  exit 1
fi
