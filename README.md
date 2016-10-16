NetstorageKit-Golang: Akamai Netstorage API for Go
================================================

NetstorageKit-Golang is Akamai Netstorage (File/Object Store) API for Go.  
  
  
Installation
------------

To install Netstorage API for Go:  

```bash
$ go get github.com/astinchoi/netstoragekit-golang
```
  
  
Example
-------

```go
package main

import (
  "fmt"
  "github.com/astinchoi/netstoragekit-golang"
  "./secrets"
)

func main() {
  nsHostname := "astin-nsu.akamaihd.net"
  ns_Keyname  := "astinastin"
  ns_Key := secrets.KEY // Don't expose ns_Key on public repository.
  ns_Cpcode := "360949"

  ns := netstorage.NewNetstorage(nsHostname, ns_Keyname, ns_Key, false)

  local_source := "hello.txt"
  ns_destination := fmt.Sprintf("/%s/hello.txt", ns_Cpcode) 

  res, body, err := ns.Upload(local_source, ns_destination)
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
ns.List(netstoragePath)
ns.Mkdir(netstoragePath + newDirectory)
ns.Mtime(netstoragePath, mTime) # ex) mTime: time.Now().Unix()
ns.Quick_delete(netstorageDir) # needs to be enabled on the CP Code
ns.Rename(netstorageTarget, netstorageDestination)
ns.Rmdir(netstorageDir) # remove empty direcoty
ns.Stat(netstoragePath)
ns.Symlink(netstorageTarget, netstorageDestination)
ns.Upload(LOCAL_SOURCE, netstorageDestination)

// INFO: can "Upload" Only a single file, not directory.
```
  
  
Test
----
You can test all above methods with [unittest script](https://github.com/AstinCHOI/NetStorageKit-Golang/blob/master/netstorage_test.go):


```bash
$ go test
[TEST] Dir /360949 done
[TEST] Mkdir /360949/nst_1476598764 done
[TEST] Upload nst_1476598764.txt to /360949/nst_1476598764/nst_1476598764.txt done
[TEST] Du /360949/nst_1476598764 done
[TEST] Mtime /360949/nst_1476598764/nst_1476598764.txt done
[TEST] Stat /360949/nst_1476598764/nst_1476598764.txt done
[TEST] Symlink /360949/nst_1476598764/nst_1476598764.txt to /360949/nst_1476598764/nst_1476598764.txt_lnk done
[TEST] Rename /360949/nst_1476598764/nst_1476598764.txt to /360949/nst_1476598764/nst_1476598764.txt_rename done
[TEST] Download /360949/nst_1476598764/nst_1476598764.txt done
[TEST] delete /360949/nst_1476598764/nst_1476598764.txt_rename done
[TEST] delete /360949/nst_1476598764/nst_1476598764.txt_lnk done
[TEST] rmdir /360949/nst_1476598764 done
PASS
[TEARDOWN] remove nst_1476598764.txt from local done
[TEARDOWN] remove nst_1476598764.txt_rename from local done
ok  	github.com/astinchoi/netstoragekit-golang	x.xxxs
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