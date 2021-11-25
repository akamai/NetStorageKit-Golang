[![Build Status](https://travis-ci.org/akamai/NetStorageKit-Golang.svg?branch=master)](https://travis-ci.org/akamai/NetStorageKit-Golang)
[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/NetStorageKit-Golang)](https://goreportcard.com/report/github.com/akamai/NetStorageKit-Golang)
[![GoDoc](https://godoc.org/github.com/akamai/NetStorageKit-Golang?status.svg)](https://godoc.org/github.com/akamai/NetStorageKit-Golang)
[![License](http://img.shields.io/:license-apache-blue.svg)](https://github.com/akamai/akamai/NetStorageKit-Golang/blob/master/LICENSE)

NetstorageKit-Golang: Akamai Netstorage API for Go
==================================================

NetstorageKit-Golang is Akamai Netstorage (File/Object Store) API for Go 1.4+. 
  
Important
------------

Akamai does not maintain or regulate this package. While it can be incorporated to assist you in API use, Akamai Technical Support will not offer assistance and Akamai cannot be held liable if issues arise from its use. 
  
Installation
------------

To install Netstorage API for Go:  

```bash
$ go get github.com/akamai/netstoragekit-golang
```
  
  
Example
-------

```go
package main

import (
  "fmt"
  "github.com/akamai/netstoragekit-golang"
  "./secrets" // in the .gitignore file
)

func main() {
  nsHostname := "astin-nsu.akamaihd.net"
  nsKeyname  := "astinapi"
  nsKey := secrets.KEY // Don't expose nsKey on public repository.
  nsCpcode := "360949"

  ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, false)

  localSource := "hello.txt"
  nsDestination := fmt.Sprintf("/%s/hello.txt", nsCpcode) // or "/%s/" is same. 

  res, body, err := ns.Upload(localSource, nsDestination)
  if err != nil {
      // Do something
  }

  if res.StatusCode == 200 {
      fmt.Printf(body)
  }
}
```
  
  
Methods
-------

```go
ns.Delete(netstoragePath)
ns.Dir(netstoragePath)
ns.Download(netstorageSource, localDestintation)
ns.Du(netstoragePath)
ns.Mkdir(netstoragePath + newDirectory)
ns.Mtime(netstoragePath, mTime) // ex) mTime: time.Now().Unix()
ns.QuickDelete(netstorageDir) // needs to the privilege on the CP Code
ns.Rename(netstorageTarget, netstorageDestination)
ns.Rmdir(netstorageDir) // remove empty direcoty
ns.Stat(netstoragePath)
ns.Symlink(netstorageTarget, netstorageDestination)
ns.Upload(localSource, netstorageDestination) // upload single file
ns.UploadAndIndexZip(localSource, netstorageDestination) // upload and index zip archive

// INFO: can "Upload" Only a single file, not directory.
```
  
  
Test
----
You can test all above methods with the [unittest script](https://github.com/akamai/NetStorageKit-Golang/blob/master/netstorage_test.go) (NOTE: You should input nsHostname, nsKeyname, nsKey and nsCpcode in the script):


```bash
$ go test
### Netstorage Test ###
[TEST] Dir /360949 done
[TEST] Mkdir /360949/nst_1477474457 done
[TEST] Upload nst_1477474457.txt to /360949/nst_1477474457/nst_1477474457.txt done
[TEST] Du /360949/nst_1477474457 done
[TEST] Mtime /360949/nst_1477474457/nst_1477474457.txt done
[TEST] Stat /360949/nst_1477474457/nst_1477474457.txt done
[TEST] Symlink /360949/nst_1477474457/nst_1477474457.txt to /360949/nst_1477474457/nst_1477474457.txt_lnk done
[TEST] Rename /360949/nst_1477474457/nst_1477474457.txt to /360949/nst_1477474457/nst_1477474457.txt_rename done
[TEST] Download /360949/nst_1477474457/nst_1477474457.txt done
[TEST] delete /360949/nst_1477474457/nst_1477474457.txt_rename done
[TEST] delete /360949/nst_1477474457/nst_1477474457.txt_lnk done
[TEST] rmdir /360949/nst_1477474457 done

### Error Test ###
[TEST] Dir: netstorage invalid path test done
[TEST] Upload: local invalid path test done
[TEST] Download: netstorage directory path test done

PASS
[TEARDOWN] remove nst_1477474457.txt from local done
[TEARDOWN] remove nst_1477474457.txt_rename from local done
ok  	github.com/akamai/netstoragekit-golang	x.xxxs
```
  
  
Author
------

Astin Choi (achoi@akamai.com)  
  
  
License
-------

Copyright 2016 Akamai Technologies, Inc.  All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
