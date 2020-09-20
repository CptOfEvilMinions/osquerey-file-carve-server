# Test cases

## Test case 01: Upload file to disk
1. ``
1. ``
1. ``

## Test case 02: Upload file to Mongo
1. ``
1. ``
1. ``

## Test case 03: Retrieve file from disk
1. ``
1. ``
1. ``

## Test case 03: Retrieve file from Mongo
1. ``
1. ``
1. ``

# Tests for writing/retrieving files to Mongo

1. `curl http://ipv4.download.thinkbroadband.com/50MB.zip --output /tmp/50MB.zip`
1. `openssl dgst -sha256 /tmp/50MB.zip`
1. `ENV=debug ./osquery-file-carve-server &`
1. `osqueryi --flagfile conf/osquery/osquery.test.flags --json "SELECT * FROM carves WHERE path like '/tmp/50MB.zip' AND carve=1; SELECT * FROM carves WHERE path like '/tmp/%';"`
1. `curl -X GET -k https://kolide.hackinglab.local:8000/download -d '{"file_carve_guid": "<GUID>"}'`

## References
* []()
* []()
* []()
* []()
* []()
* []()
* []()