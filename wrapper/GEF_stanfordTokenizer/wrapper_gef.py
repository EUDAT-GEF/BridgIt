# Python script to wrap a WebLicht Web Service
# C. Zinn, October 2017
#
# curl -H 'content-type: text/tcf+xml' -d @input.xml -X POST http://localhost:8080/service-nentities/annotate/stream
 

import requests
import sys
import time
import io
import os, errno
import glob

from subprocess import Popen

print("creating data directory")

try:
    os.makedirs('data')
except OSError as e:
    if e.errno != errno.EEXIST:
        raise
    
# start the converter service
print("starting the stanford service locally")
Popen(["sh", "/service-stanford-tokenizer/bin/service-stanford-tokenizer", "server", "/service-stanford-tokenizer/service.yaml"])

print("Sleeping a little bit")
time.sleep(10)


# the url for the stanford tokenizer and its header
systemURL = 'http://localhost:8080/service-stanford-tokenizer/annotate/stream'
#systemURL = 'https://weblicht.sfs.uni-tuebingen.de/rws/service-stanford-tokenizer/annotate/stream'
headers = {'content-type': 'text/tcf+xml'}

print("getting the data")
textFiles = glob.glob("/root/input/*.*")
for fi in range(len(textFiles)):
    with open(textFiles[fi], 'rb') as f: 
        res = requests.post(systemURL,
                            data=f,    # do s/f/data directly in case you don't want to store it first
                            headers=headers)

    # write result to file
    print("writing the result")
    if (res.status_code == 200):
        with io.open('/root/output/result.tcf', 'w', encoding='utf8') as f:        
            f.write(res.text)
            print("stanford tokenizer status:", res.status_code)
            print("result written:", res.text)
    else:
        print("stanford tokenizer error status:", res.status_code)
        
print("ending the script")
