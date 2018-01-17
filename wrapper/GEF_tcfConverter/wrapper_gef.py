# Python script to wrap a WebLicht Web Service
# C. Zinn, October 2017
#
# Variant of wrapper.py where the input file resides in the given GEF input volume,
# and the output resides in the given output volume

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
print("starting the converter service locally")
Popen(["sh", "/service-converter/bin/service-converter", "server", "/service-converter/service.yml"])

print("Sleeping a little bit (10)")
time.sleep(10)


# the url for "To TCF converter" and its parameters
systemURL = 'http://localhost:8080/service-converter/convert/qp'
#systemURL = 'https://weblicht.sfs.uni-tuebingen.de/rws/service-converter/convert/qp'
systemParams = { "informat" : "plaintext",
                 "outformat" : "tcf04",
                 "language" : "de"}

# the url that the wrapper uses to fetch the data
print("getting the data")
textFiles = glob.glob("/root/input/*.txt")
for fi in range(len(textFiles)):
    with open(textFiles[fi], 'rb') as f: 
        res = requests.post(systemURL,
                            data=f,    # do s/f/data directly in case you don't want to store it first
                            params=systemParams)

        # write result to file
        print("writing the result")
        if (res.status_code == 200):
            with io.open('/root/output/result.tcf', 'w', encoding='utf8') as f:
                f.write(res.text)
                print("tcf status:", res.status_code)
                print("result written:", res.text)
        else:
            print("tcf error status:", res.status_code)
            

print("ending the script")
