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
    "strings"
    "time"
)


type Netstorage struct {
    hostname    string
    keyname     string
    key         string
    ssl         string
}

func NewNetstorage(hostname, keyname, key string, ssl bool) *Netstorage {
    s := ""
    if ssl {
        s = "s"
    }
    return &Netstorage{hostname, keyname, key, s}
}

func (ns *Netstorage) _request(kwargs map[string]string) (*http.Response, error) {
    ns_path := kwargs["path"]
    if u, err := url.Parse(ns_path); strings.HasPrefix(ns_path, "/") && err == nil {
        ns_path = u.RequestURI()
    } else {
        // Exception
    }

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
        bArr, err := ioutil.ReadFile(kwargs["source"])
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
        local_destination := kwargs["destination"]

        if strings.HasSuffix(kwargs["path"], "/") {
            // error
        } else if local_destination == "" {
            local_destination = path.Base(kwargs["path"]) 
        } else if s, err := os.Stat(local_destination); err == nil && s.IsDir() {
            local_destination = path.Join(local_destination, path.Base(kwargs["path"]))
        }

        out, err := os.Create(local_destination)
        if err != nil {
            return nil, err
        }
        defer out.Close()

        if _, err := io.Copy(out, response.Body); err != nil {
            return nil, err
        }
    }
    
    return response, nil
}

func (ns *Netstorage) Dir(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "dir&format=xml",
        "method": "GET",
        "path": ns_path,
    })
}

func (ns *Netstorage) Download(path ...string) (*http.Response, error) {
    ns_source := path[0]
    local_destination := ""
    if len(path) >= 2 {
        local_destination = path[1]
    }
    return ns._request(map[string]string{
        "action": "download",
        "method": "GET",
        "path": ns_source,
        "destination": local_destination,
    })
}

func (ns *Netstorage) Du(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "du&format=xml",
        "method": "GET",
        "path": ns_path,
    })
}

func (ns *Netstorage) Stat(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "stat&format=xml",
        "method": "GET",
        "path": ns_path,
    })
}

func (ns *Netstorage) Mkdir(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "mkdir",
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) Rmdir(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "rmdir",
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) Mtime(ns_path string, mtime int64) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": fmt.Sprintf("mtime&format=xml&mtime=%d", mtime),
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) Delete(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "delete",
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) Quick_delete(ns_path string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "quick-delete&quick-delete=imreallyreallysure",
        "method": "POST",
        "path": ns_path,
    })
}

func (ns *Netstorage) Rename(ns_target, ns_destination string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "rename&destination=" + url.QueryEscape(ns_destination),
        "method": "POST",
        "path": ns_target,
    })
}

func (ns *Netstorage) Symlink(ns_target, ns_destination string) (*http.Response, error) {
    return ns._request(map[string]string{
        "action": "symlink&target=" + url.QueryEscape(ns_target),
        "method": "POST",
        "path": ns_destination,
    })
}

func (ns *Netstorage) Upload(local_source, ns_destination string) (*http.Response, error) {
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