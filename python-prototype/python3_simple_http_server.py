from flask import Flask, request, jsonify
import random
import string
import json
import base64

app = Flask(__name__)

print ('Starting prototype Python server')

FILE_CARVE_DIR = '/tmp/'
FILE_CARVE_MAP = {}

@app.route('/')
def index():
    return "Hello World!"

@app.route('/start_uploads', methods=['POST'])
def start_upload():
    """
    Input: Takes in an HTTP POST request from a client for a file upload. 
    Post body: `{"block_count":1,"block_size":30000,"carve_size":2048,"carve_id":"b1733cf9-94c1-4a98-aa4e-03acf9953c15","request_id":"","node_key":"sYHI8O/nlPFRa8d1Ifpq0q7UraLUGReD"}`
    Output: Return a session ID for the file upload
    """
    json_data = request.json
    print ("#" * 20 + "Start upload" + "#" *20)
    sid = ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(10))
    response  = {'session_id': sid}
    FILE_CARVE_MAP[sid] = {
            'block_count': int(json_data['block_count']),
            'block_size': int(json_data['block_size']),
            'blocks_received': {},
            'carve_size': int(json_data['carve_size']),
            'carve_guid': json_data['carve_id'],
        }
    return json.dumps(response)


@app.route('/upload_blocks', methods=['POST'])
def upload_blocks():
    """
    Input: HTTP POST body that includes a payload chunk encoded with base64 
    Output: Return status code 204 with no other content.
    """
    print ("#" * 20 + "upload_blocks" + "#" *20)

    json_data = request.json
    print (json_data)

    # First check if we have already received this block
    if json_data['block_id'] in FILE_CARVE_MAP[json_data['session_id']][
            'blocks_received']:
        return

    # Store block data to be reassembled later
    FILE_CARVE_MAP[json_data['session_id']]['blocks_received'][int(
        json_data['block_id'])] = json_data['data']

    # Are we expecting to receive more blocks?
    if len(FILE_CARVE_MAP[json_data['session_id']]['blocks_received']
            ) < FILE_CARVE_MAP[json_data['session_id']]['block_count']:
        return

    # If not, let's reassemble everything
    out_file_name = FILE_CARVE_DIR + FILE_CARVE_MAP[json_data['session_id']]['carve_guid']

    # Check the first four bytes for the zstd header. If not no
    # compression was used, it's an uncompressed .tar
    if (base64.standard_b64decode(FILE_CARVE_MAP[json_data['session_id']][
            'blocks_received'][0])[0:4] == b'\x28\xB5\x2F\xFD'):
        out_file_name += '.zst'
    else:
        out_file_name += '.tar'
    f = open(out_file_name, 'wb')
    for x in range(0,
                    FILE_CARVE_MAP[json_data['session_id']]['block_count']):
        f.write(
            base64.standard_b64decode(FILE_CARVE_MAP[json_data['session_id']]
                                        ['blocks_received'][x]))
    f.close()
    print("File successfully carved to: %s" % out_file_name)
    FILE_CARVE_MAP[json_data['session_id']] = {}

    return '', 204