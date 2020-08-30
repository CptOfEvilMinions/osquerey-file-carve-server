#!/bin/bash


# Check if test file exists
# Download test file
test_file_url="http://ipv4.download.thinkbroadband.com/50MB.zip"
if [ ! -f /tmp/test_file_50MB.zip ]; then
  echo '[-] - Test file does NOT exist - Start download'
  curl ${test_file_url} --output /tmp/test_file_50MB.zip
  chmod 444 /tmp/test_file_50MB.zip
fi

# Make copy of data for upload
if [ ! -f /tmp/50MB.zip ]; then
  cp /tmp/test_file_50MB.zip /tmp/50MB.zip
  echo '[*] - Creating copy of test file'
  chmod 644 /tmp/50MB.zip
fi


# Generate test file SHA256 hash
test_file_sha256=`openssl dgst -sha256 /tmp/test_file_50MB.zip | awk '{print $2}'`

# Start server and background the proc
ENV=debug ./osquery-file-carve-server --config conf/osquery-file-carve-dev.yml &
server_proc_id=`ps aux | grep osquery-file-carve-server | grep -v 'grep' | awk '{print $2}'`

echo "[*] - osquery-file-carve-server process ID: ${server_proc_id}"

# Issue Osqueryi file carve
echo '[+] - Start Osquery file carve'
osquery_file_carve_output=`osqueryi --flagfile conf/osquery/osquery.test.flags --json "SELECT * FROM carves WHERE path like '/tmp/50MB.zip' AND carve=1; SELECT * FROM carves WHERE path like '/tmp/%';"`
osquery_carve_id=`echo ${osquery_file_carve_output} | jq '.[].carve_guid' | tr -d '"'`
rm /tmp/50MB.zip
echo '[+] - Osquery file carve finished'

# Check file carve exists
if [ -f "/tmp/${osquery_carve_id}.tar" ]; then
  echo '[+] - File Carve upload exists'
else
  echo '[-] - File Carve upload does NOT exists'
  exit 1
fi

# Untar file carve and check hash
cd /tmp && tar -xvf /tmp/${osquery_carve_id}.tar
file_carve_sha256=`openssl dgst -sha256 /tmp/50MB.zip | awk '{print $2}'`

if [ "$file_carve_sha256" = "$test_file_sha256" ]; then
  echo '[+] - SHA256 file hashes match - Orginal file uploaded'
else
  echo '[-] - SHA256 file hashes DO NOT match - NOT the orginal file uploaded'
  exit 1
fi

# Clean up
rm /tmp/50MB.zip
echo "[+] - Clean up - Deleted /tmp/50MB.zip"

rm /tmp/${osquery_carve_id}.tar
echo "[+] - Clean up - Deleted /tmp/${osquery_carve_id}.tar"

kill ${server_proc_id}
echo "[+] - Clean up - Stopped osquery-file-carve-server"