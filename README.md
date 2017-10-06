Bridgit
=======

Liaison between WebLicht and the GEF-enabled services

The GEF
--------
The GEF (Generic Execution Framework) is available at: https://github.com/EUDAT-GEF/GEF

Weblicht
--------
https://weblicht.sfs.uni-tuebingen.de

How does it work?
-----------------
- Copy the repository
- Run `go build`
- Modify the `config.json` from Bridgit the way you need
- Run `./bridgit`
- Start the GEF
- Start on the user profile page create a new API token and use it for all requests
- Try to send a simple request by using CURL: `curl -X POST http://localhost:8080/jobs\?service\=SERVICE_NAME\&token\=B2ACCESS_TOKEN\&input\=INPUT_FILE_PATH -v -o out3.txt`

Configuration File
------------------
~~~~
{
  "StaticContentFolder": "./data/",
  "StaticContentURLPrefix": "/static",
  "StorageURL": HTTP://IP_ADDRESS_OF_BRIDGIT,
  "GEFAddress": HTTPS://IP_ADDRESS_AND_PORT_OF_GEF,
  "StoragePortNumber": BRIDGIT_PORT,
  "TimeOut": 1000,
  "Apps": {
      "nltk": "91207aa2-5fb4-440d-bce0-79fb6027d928",
      "stanford": "ce5203bb-f161-49af-83e5-4486e408fce5"
  }
}
~~~~

`StaticContentFolder` - contains documents saved locally to be served to a GEF service

`StaticContentURLPrefix` - URL prefix for the local storage in BridgIT

`StorageURL` - URL of the local storage in BridgIT

`StoragePortNumber` - port number where BridgIt is serving static content

`GEFAddress` - URL of the GEF instance

`TimeOut` - system time out time

`Apps` - used to map human readable service names with their ID codes