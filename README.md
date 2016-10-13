NetstorageAPI: Akamai Netstorage API for GO
================================================

NetstorageAPI is Akamai Netstorage (File/Object Store) API for Go.
  
  
Installation
------------

To install Netstorage API for Go:  

```bash
$ go get github.com/AstinCHOI/NetStorageKit-GoLang/akamai/netstorage
```
  
  
Example
-------

```GoLang
package main

import (
  "fmt"
  "github.com/AstinCHOI/NetStorageKit-GoLang/akamai/netstorage"
  "./secrets"
)

func main() {
  NS_HOSTNAME := "astin-nsu.akamaihd.net"
  NS_KEYNAME  := "astinastin"
  NS_KEY := secrets.KEY // Don't expose NS_KEY on public repository.
  NS_CPCODE := "360949"

  ns := netstorage.NewNetstorage(NS_HOSTNAME, NS_KEYNAME, NS_KEY, false)

  local_source := "hello.txt"
  ns_destination := fmt.Sprintf("/%s/hello.txt", NS_CPCODE) 

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

```GoLang
ns.Delete(NETSTORAGE_PATH)
ns.Dir(NETSTORAGE_PATH)
ns.Download(NETSTORAGE_SOURCE, LOCAL_DESTINATION)
ns.Du(NETSTORAGE_PATH)
ns.List(NETSTORAGE_PATH)
ns.Mkdir("#{NETSTORAGE_PATH}/#{DIRECTORY_NAME}")
ns.Mtime(NETSTORAGE_PATH, TIME) # ex) TIME: time.Now().Unix()
ns.Quick_delete(NETSTORAGE_DIR) # needs to be enabled on the CP Code
ns.Rename(NETSTORAGE_TARGET, NETSTORAGE_DESTINATION)
ns.Rmdir(NETSTORAGE_DIR) # remove empty direcoty
ns.Stat(NETSTORAGE_PATH)
ns.Symlink(NETSTORAGE_TARGET, NETSTORAGE_DESTINATION)
ns.Upload(LOCAL_SOURCE, NETSTORAGE_DESTINATION)

// INFO: can "upload" Only a single file, not directory.
```
  
  
Test
----
You can test all above methods with [Unit Test Script](https://github.com/AstinCHOI/NetStorageKit-Golang/blob/master/netstorage_test.go). It uses [Testify](https://github.com/stretchr/testify) for the test:


```bash
$ go get github.com/stretchr/testify
...
$ go test netstorage_test.go -v

=== RUN TestExampleTestSuite
=== RUN TestNetstorage
[TEST] dir /360949 done
[TEST] mkdir /360949/nst_1476344471 done
[TEST] upload nst_1476344471.txt to /360949/nst_1476344471/nst_1476344471.txt done
[TEST] du /360949/nst_1476344471 done
[TEST] mtime /360949/nst_1476344471/nst_1476344471.txt done
[TEST] stat /360949/nst_1476344471/nst_1476344471.txt done
[TEST] symlink /360949/nst_1476344471/nst_1476344471.txt to /360949/nst_1476344471/nst_1476344471.txt_lnk done
[TEST] rename /360949/nst_1476344471/nst_1476344471.txt to /360949/nst_1476344471/nst_1476344471.txt_rename done
[TEST] download /360949/nst_1476344471/nst_1476344471.txt done
[TEST] delete /360949/nst_1476344471/nst_1476344471.txt_rename done
[TEST] delete /360949/nst_1476344471/nst_1476344471.txt_lnk done
[TEST] rmdir /360949/nst_1476344471 done
--- PASS: TestNetstorage (11.00 seconds)
[TEARDOWN] remove nst_1476344471.txt from local done
[TEARDOWN] remove nst_1476344471.txt_rename from local done
--- PASS: TestExampleTestSuite (xx.xx seconds)
PASS
ok  	command-line-arguments	xx.xxxs
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