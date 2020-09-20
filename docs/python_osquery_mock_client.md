# Python Osquery mock client

This README assumes you have a working stack. This Python script mocks an Osquery client and uploads files to stress test the osquery-file-carve server. 

## Download test files
1. Open a terminal
1. `cd /tmp`
1. `wget http://ipv4.download.thinkbroadband.com/10MB.zip`
1. `wget http://ipv4.download.thinkbroadband.com/100MB.zip`

## Test cases
The test cases below are meant to stres test the osquery-file-carve server.

### 10MB,zip single file upload
1. `osquery-file-carve-server/tests`
1. `python3 tests/osquery_mock_client.py --file /tmp/10MB.zip --threads=1 --base_url=https://<Kolide IP addr>:<port>`

### 10MB.zip multiple file upload
1. `osquery-file-carve-server/tests`
1. `python3 tests/osquery_mock_client.py --file /tmp/10MB.zip --threads=10 --base_url=https://<Kolide IP addr>:<port>`

### 100MB.zip multiple file upload
1. `osquery-file-carve-server/tests`
1. `python3 tests/osquery_mock_client.py --file /tmp/10MB.zip --threads=10 --base_url=https://<Kolide IP addr>:<port>`


## References
* [StackOverFlow - Suppress InsecureRequestWarning: Unverified HTTPS request is being made in Python2.6](https://stackoverflow.com/questions/27981545/suppress-insecurerequestwarning-unverified-https-request-is-being-made-in-pytho/33716188)
* [StackOverFlow - How to print colored text in Python?](https://stackoverflow.com/questions/287871/how-to-print-colored-text-in-python)
* [StackOverFlow - How do I change a UUID to a string?](https://stackoverflow.com/questions/37049289/how-do-i-change-a-uuid-to-a-string/37049434)
* [threading – Manage concurrent threads](https://pymotw.com/2/threading/)
* [UUID objects according to RFC 4122](https://docs.python.org/3/library/uuid.html)
* [The Python Modulo Operator - What Does the % Symbol Mean in Python? (Solved)](https://www.freecodecamp.org/news/the-python-modulo-operator-what-does-the-symbol-mean-in-python-solved/)
* [Getting file size in Python? [duplicate]](https://stackoverflow.com/questions/6591931/getting-file-size-in-python)
* [Chunking in Python---How to set the "chunk size" of read lines from file read with Python open()?](https://www.reddit.com/r/learnpython/comments/52vpjz/chunking_in_pythonhow_to_set_the_chunk_size_of/)
* []()
* []()