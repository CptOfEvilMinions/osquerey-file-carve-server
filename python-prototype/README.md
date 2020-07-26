# Python prototype server

## Osqueryi
1. `echo 'hello' > /tmp/test.txt`
1. `osqueryi --flagfile conf/osquery/osquery.flags`
1. `select * from carves where path like '/tmp/test.txt' AND carve=1;`
  1. Upload file
1. `select * from carves WHERE path like '/tmp/test.txt';`
  1. Status of upload

## Tested Osquery versions
* `osquery version 4.3.0` 

## References
* [Flask return 204 No Content response](https://www.erol.si/2018/03/flask-return-204-no-content-response/)
* [Get the data received in a Flask request](https://stackoverflow.com/questions/10434599/get-the-data-received-in-a-flask-request)
* [Apk add python3 py3-pip](https://medium.com/@ssorcnafets/apk-add-python3-py3-pip-c3f91cd3d1e1)
* [Configure Flask dev server to be visible across the network](https://stackoverflow.com/questions/7023052/configure-flask-dev-server-to-be-visible-across-the-network/51164848)
* [Return JSON response from Flask view](https://stackoverflow.com/questions/13081532/return-json-response-from-flask-view)
* [Steps for Starting a New Flask Project using Python3](https://www.patricksoftwareblog.com/steps-for-starting-a-new-flask-project-using-python3/)
* [How To Make a Web Application Using Flask in Python 3](https://www.digitalocean.com/community/tutorials/how-to-make-a-web-application-using-flask-in-python-3)
* [Github Gist - Minimal JSON HTTP server in python](https://gist.github.com/nitaku/10d0662536f37a087e1b)
* [Github - Facebook - osquery/tools/tests/test_http_server.py ](https://github.com/osquery/osquery/blob/master/tools/tests/test_http_server.py)
* [Osquery docs - Remote authentication](https://osquery.readthedocs.io/en/stable/deployment/remote/)
* [Osquery docs - CLI flags](https://osquery.readthedocs.io/en/stable/installation/cli-flags/)
* [Understanding Nginx Server and Location Block Selection Algorithms](https://www.digitalocean.com/community/tutorials/understanding-nginx-server-and-location-block-selection-algorithms)
* []()