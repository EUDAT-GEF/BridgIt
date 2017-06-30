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
- Try to send a simple request by using CURL: `curl -X POST -d @FILE_NAME http://localhost:8080/jobs -v`

Configuration File
------------------
~~~~
{
  "StaticContentFolder": "./data/",
  "StaticContentURLPrefix": "/static",
  "StorageURL": HTTP://IP_ADDRESS_AND_PORT_OF_BRIDGIT,
  "GEFAddress": HTTPS://IP_ADDRESS_AND_PORT_OF_GEF,
  "PortNumber": BRIDGIT_PORT,
  "TimeOut": 1000
}
~~~~

`StaticContentFolder` - contains documents saved locally to be served to a GEF service

`StaticContentURLPrefix` - URL prefix for the local storage in BridgIT

`StorageURL` - URL of the local storage in BridgIT

`GEFAddress` - URL of the GEF instance

`PortNumber` - port number BridgIt is running on

`TimeOut` - system time out time