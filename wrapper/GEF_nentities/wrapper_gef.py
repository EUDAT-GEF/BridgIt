# Python script to wrap a WebLicht Web Service
# C. Zinn, October 2017
#
# curl -H 'content-type: text/tcf+xml' -d @input.xml -X POST http://localhost:8080/service-lbjner/annotate/stream
 

import requests
import sys
import time
import io
import os, errno
import glob
import datetime

from subprocess import Popen

print("creating data directory")

try:
    os.makedirs('data')
except OSError as e:
    if e.errno != errno.EEXIST:
        raise
    
# start the converter service
print("starting the illinos NER service")
#child =
Popen(["sh", "/service-illinois-ner/bin/service-illinois-ner", "server", "/service-illinois-ner/service.yml"])
#streamdata = child.communicate()[0]
#rc = child.returncode
#print "return code: ", rc


print("datetime before sleeping 120 sec: ", datetime.datetime.now())
print os.environ['JAVA_OPTS']
print os.environ['_JAVA_OPTIONS']
time.sleep(120)

print("datetime after sleeping 120 sec: ", datetime.datetime.now())

# the url for the stanford tokenizer and its header
#systemURL = 'http://localhost:8080/service-illinois-ner/annotate/stream'
systemURL = 'http://localhost:8080/service-lbjner/annotate/stream'
#systemURL = 'https://weblicht.sfs.uni-tuebingen.de/rws/service-nentities/annotate/stream'
#headers = {'content-type': 'text/tcf+xml'}

print("getting the data from this input volume", datetime.datetime.now())
textFiles = glob.glob("/root/input/*.*")
for fi in range(len(textFiles)):
    print('input: ' + textFiles[fi], datetime.datetime.now())
    with open(textFiles[fi], 'rb') as f: 
        res = requests.post(systemURL,
                            data=f    # do s/f/data directly in case you don't want to store it first
        )

    # write result to file
    print("writing the result", datetime.datetime.now())
    if (res.status_code == 200):
        with io.open('/root/output/result.tcf', 'w', encoding='utf8') as f:                
            f.write(res.text)
            print("illinois NER status:", res.status_code)
            print("result written:", res.text)
    else:
        print("illinois NER error status:", res.status_code)

    
print("ending the script", datetime.datetime.now())

