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
- Run `./bridgit`
- Start the GEF
- Modify the `config.json` from Bridgit the way you need
- Try to send a simple request by using CURL: `curl -X POST -d @FILE_NAME http://localhost:8080/jobs -v`