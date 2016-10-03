package netstorage

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
	"fmt"
    "io"
    "io/ioutil"
    "math/rand"
	"net/http"
    "net/url"
    "os"
    "path"
    "strconv"
    "strings"
    "time"
)


type Netstorage struct {
    hostname    string
    keyname     string
    key         string
    ssl         string
}

func New(hostname, keyname, key string, ssl bool) *Netstorage {
    s := ""
    if ssl {
        s = "s"
    }
    return &Netstorage{hostname, keyname, key, s}
}

func (ns *Netstorage) _request(kwargs map[string]string) (*http.Response, error) {
    ns_path := kwargs["path"]
    if (!strings.HasPrefix(ns_path, "/")) {
        // Exception    
    }
    
    ns_path = strconv.Quote(ns_path)

    acs_action := fmt.Sprintf("version=1&action=%s", kwargs["action"])
    acs_auth_data := fmt.Sprintf("5, 0.0.0.0, 0.0.0.0, %d, %d, %s",
                                    time.Now().Unix(),
                                    rand.Intn(100000),
                                    ns.keyname)

    sign_string := fmt.Sprintf("%s\nx-akamai-acs-action:%s\n", ns_path, acs_action)
    mac := hmac.New(sha256.New, []byte(ns.key))
    mac.Write([]byte(acs_auth_data + sign_string))
    acs_auth_sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

    var data io.Reader = nil
    if kwargs["action"] == "upload" {
        bArr, err := ioutil.ReadFile(kwargs["soruce"])
        if err != nil {
            return nil, err    
        }

        data = bytes.NewReader(bArr)
    }

    request, err := http.NewRequest(kwargs["method"], 
        fmt.Sprintf("http%s://%s%s", ns.ssl, ns.hostname, ns_path), data)
    
    if err != nil {
		return nil, err
	}

    request.Header.Add("X-Akamai-ACS-Action", acs_action)
    request.Header.Add("X-Akamai-ACS-Auth-Data", acs_auth_data)
    request.Header.Add("X-Akamai-ACS-Auth-Sign", acs_auth_sign)
    request.Header.Add("Accept-Encoding", "identity")
    request.Header.Add("User-Agent", "NetStorageKit-Golang")

    client := &http.Client{}
    response, err := client.Do(request)
    
    if err != nil {
		return nil, err
	}
    
    defer response.Body.Close()

    if kwargs["action"] == "download" {
        body, err := ioutil.ReadAll(response.Body)
        if err != nil {
            return nil, err
        }
        err = ioutil.WriteFile(kwargs["destination"], body, 0666)
        if err != nil {
            return nil, err
        }
    }
    
    return response, nil
}

func (ns *Netstorage) dir(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "dir&format=xml",
        "method": "GET",
        "path": ns_path,
    })
}

func (ns *Netstorage) download(ns_source, local_destination string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "download",
        "method": "GET",
        "path": ns_source,
        "destination": local_destination,
    })
}

func (ns *Netstorage) du(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "du&format=xml",
        "method": "GET",
        "path": ns_path,
    })
}

func (ns *Netstorage) stat(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "stat&format=xml",
        "method": "GET",
        "path": ns_path,
    })
}

func (ns *Netstorage) mkdir(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "mkdir",
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) rmdir(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "rmdir",
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) mtime(ns_path string, mtime int64) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": fmt.Sprintf("mtime&format=xml&mtime=%d", mtime),
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) delete(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "delete",
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) quick_delete(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "quick-delete&quick-delete=imreallyreallysure",
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) rename(ns_target, ns_destination string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "rename&destination=" + url.QueryEscape(ns_destination),
        "method": "POST",
        "path": ns_target,
    })
}

func (ns *Netstorage) symlink(ns_target, ns_destination string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "symlink&target=" + url.QueryEscape(ns_target),
        "method": "POST",
        "path": ns_destination,
    })
}

func (ns *Netstorage) upload(local_source, ns_destination string) (*http.Response, error) {
    s, err := os.Stat(local_source)

    if err != nil {
        // do something
    }   

    if s.Mode().IsRegular() {    
        if strings.HasSuffix(ns_destination, "/") {
            ns_destination = ns_destination + path.Base(local_source)
        }
    } else {
        // do something
    }
    
    return ns._request(map[string]string{
        "action": "upload",
        "method": "PUT",
        "source": local_source,
        "path": ns_destination,
    })
}