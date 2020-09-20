import requests
import argparse
import json
import os
import uuid
import threading
import base64
from colorama import Style
from colorama import Fore
import urllib3
import math
import tarfile
urllib3.disable_warnings()

class FileUpload:
  def __init__(self, file_path, block_size, base_url):
    self.file_path = file_path                        # Specify the file path for the file to upload
    self.block_size = block_size
    self.base_url = base_url
    self.file_guid = self.generate_file_guid()        # Generate File GUID
    self.tar_file_path = self.generate_tar()          # Generate TAR
    self.block_count = self.generate_block_count()    # Generate of blocks for TAR file upload
    self.session_id = self.get_session_id()           # Request Session ID
    self.upload_status = True

  def generate_block_count(self):
    """
    Based on the file size divide by block size
    """
    print ("Block count", math.ceil((os.path.getsize(self.tar_file_path) / self.block_size)))
    return math.ceil((os.path.getsize(self.tar_file_path) / self.block_size))

  def generate_file_guid(self):
    file_guid = str( uuid.uuid4() ) 
    print (f"File GUID: {file_guid}")
    return file_guid

  def generate_tar(self):
    """
    Takes in a file path and TARs up the file in the tmp directory
    """
    tar_file_path = f"/tmp/{self.file_guid}.tar"
    with tarfile.open(tar_file_path, "w") as tar_handle:
      tar_handle.add(os.path.join(self.file_path))
    return tar_file_path


  def get_session_id(self):
    """
    {"block_count":1,"block_size":10000000,"carve_size":2048,"carve_id":"b028316a-8365-407a-8e00-6d16138d5826","request_id":"","node_key":"yLtsQVZ0Qt0SVaaqSjardnkrc2TmapnE"}
    """
    data = {
      "block_count": self.block_count,
      "block_size": self.block_size,
      "carve_size":2048,
      "carve_id": self.file_guid,
      "request_id":"",
      "node_key": "5"
    }

    url = f"{self.base_url}/start_uploads"
    r = requests.get(url=url, data=json.dumps(data), timeout=10, verify=False)
    
    print (f"""Session ID: {r.json()["session_id"]}""")
    return r.json()["session_id"]

  

def post_block(fileUpload, block_data, block_id):
  """
  Iterate over the file in chunks and send each chunk
  {"block_id":0,"session_id":"d64c2d46-1984-4e63-8acf-9a108127baf3","request_id":"","data":"aGVsbG......."}
  """

  data = {
    "block_id": block_id,
    "session_id": fileUpload.session_id,
    "request_id": "",
    "data": block_data
  }


  url = f"{fileUpload.base_url}/upload_blocks"
  r = requests.post(url=url, data=json.dumps(data), verify=False)

  if r.status_code != 200:
    print ( f"{Fore.RED}{block_id}{r.status_code}{r.text}{Style.RESET_ALL}")
    fileUpload.upload_status = False
  else:
    print ( f"{Fore.GREEN}{block_id} - {r.status_code} - {r.text}{Style.RESET_ALL}")

  

def chunk_file(fileUpload):
  print ("Uploading file")

  # Read the file in chunks
  with open(fileUpload.tar_file_path, 'rb') as f:
    
    # get session ID
    i = 0
    while True:
      read_data = f.read(fileUpload.block_size)
      if not read_data:
        break # done
      base64_encoded_data = base64.b64encode(read_data)
      base64_message = base64_encoded_data.decode('utf-8')
      post_block(fileUpload, base64_message, i)
      i = i +1

  # Delete tar before we exit thread
  os.remove(fileUpload.tar_file_path)
  return

if __name__ == "__main__":
  my_parser = argparse.ArgumentParser()
  my_parser.add_argument('--file', type=str, required=True, help='Test file to send', )
  my_parser.add_argument('--block_size', type=int, default=1000000, required=True, help='Size of each block')
  my_parser.add_argument('--base_url', type=int, required=True, help='Size of each block')
  my_parser.add_argument('--threads', type=int, default=1, help='Size of each block')
  args = my_parser.parse_args()

  fileUploads = list()
  # Start threads
  threads = [] 
  for i in range(args.threads):
    fileUpload = FileUpload(args.file, args.block_size, args.base_url)
    fileUploads.append(fileUpload)
    t = threading.Thread(target=chunk_file, args=(fileUpload,))
    threads.append(t)
    t.start()

  for x in threads:
    x.join()

  sucess = 0
  failure = 0
  for fileUpload in fileUploads:
    if fileUpload.upload_status == True:
      sucess = sucess + 1
    else: 
      failure = failure +1

  print ( f"{Fore.RED}Failure: {failure}{Style.RESET_ALL}")
  print ( f"{Fore.GREEN}Sucess: {sucess}{Style.RESET_ALL}")
  