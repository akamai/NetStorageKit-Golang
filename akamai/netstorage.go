package netstorage

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
	"fmt"
    "io/ioutil"
    "math/rand"
	"net/http"
    "net/url"
    "strings"
    "time"
)


type Netstorage struct {
    hostname    string
    keyname     string
    key         string
    ssl         bool
}

func NewNetstorage(hostname, keyname, key string, ssl bool) *Netstorage {
    return &Netstorage{hostname, keyname, key, ssl}
}

func (ns *Netstorage) _request(kwargs map[string]string) (*http.Response, error) {
    path := kwargs["path"]
    if (!strings.HasPrefix(path, "/")) {
        // Exception    
    }
    

    return nil, nil
}

