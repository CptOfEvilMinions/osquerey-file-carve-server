# Osquery

This README contains the necessary information for how to configure Osquery

## Config tempaltes
* [osquery.flags](conf/osquery/osquery.flags)
* [osquery.test.flags](conf/osquery/osquery.test.flags)

## Config values
* `--tls_hostname=` - Needs to be set to the URL to reach Kolide. The same URL to reach Kolide will be the same for the osquery-file-carve server. NGINX will handle directing HTTP requests.
* `--read_max=` - Needs to be set to the maximum value in bytes for which Osquery should read files. By default we set this to `1GB` so Osquery cannot read or upload files larger than 1GB.
* `--carver_block_size`- Needs to be set to the maximum size of each data block that will be sent to the osquery-file-carve server. Keep in mind this value can NOT be bigger than the `client_max_body_size` defined in the NGINX config.
