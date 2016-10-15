//
package netstorage


import (
    "fmt"
    "io/ioutil"
    "os"
    "testing"
    "time"
    
    "./spike/secrets"
    // "github.com/AstinCHOI/netstoragekit-golang"
    "github.com/stretchr/testify/suite"
)

var NS_HOSTNAME string = "astin-nsu.akamaihd.net"
var NS_KEYNAME  string = "astinastin"
var NS_KEY string = secrets.KEY
var NS_CPCODE string = "360949"


type NetstorageTestSuite struct {
    suite.Suite
    ns *Netstorage
    tempNsDir string
    tempFile string
    tempNsFile string
}

func check(err error, exit bool) {
    if (err != nil) {
        if exit {
            panic(err)
        } else {
            fmt.Println(err)
        }
    }
}

func (suite *NetstorageTestSuite) SetupSuite() {
    suite.ns = NewNetstorage(NS_HOSTNAME, NS_KEYNAME, NS_KEY, false)
    suite.tempNsDir = fmt.Sprintf("/%s/nst_%d", NS_CPCODE, time.Now().Unix())
    suite.tempFile = fmt.Sprintf("nst_%d.txt", time.Now().Unix())
    suite.tempNsFile = fmt.Sprintf("%s/%s", suite.tempNsDir, suite.tempFile)
}

func (suite *NetstorageTestSuite) TearDownSuite() {
    // delete temp files for local
    if _, err := os.Stat(suite.tempFile); err == nil {
        err = os.Remove(suite.tempFile)
        check(err, false)
        fmt.Printf("[TEARDOWN] remove %s from local done\n", suite.tempFile)    
    }

    if _, err := os.Stat(suite.tempFile + "_rename"); err == nil {
        err = os.Remove(suite.tempFile + "_rename")
        check(err, false)
        fmt.Printf("[TEARDOWN] remove %s from local done\n", suite.tempFile + "_rename")    
    }

    // delete temp files for netstorage
    if res, _, err := suite.ns.Delete(suite.tempNsFile); res.StatusCode == 200 && err == nil {
        fmt.Printf("[TEARDOWN] delete %s done\n", suite.tempNsFile)
    }

    if res, _, err := suite.ns.Delete(suite.tempNsFile + "_lnk"); res.StatusCode == 200 && err == nil {
        fmt.Printf("[TEARDOWN] delete %s done\n", suite.tempNsFile + "_lnk")
    }

    if res, _, err := suite.ns.Delete(suite.tempNsFile + "_rename"); res.StatusCode == 200 && err == nil {
        fmt.Printf("[TEARDOWN] delete %s done\n", suite.tempNsFile + "_rename")
    }

    if res, _, err := suite.ns.Rmdir(suite.tempNsDir); res.StatusCode == 200 && err == nil {
        fmt.Printf("[TEARDOWN] rmdir %s done\n", suite.tempNsDir)
    }
}

func (suite *NetstorageTestSuite) TestNetstorage() {
    // Dir
    res, _, err := suite.ns.Dir("/" + NS_CPCODE)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[dir] StatusCode should be 200 OK")
    fmt.Printf("[TEST] dir /%s done\n", NS_CPCODE)
    
    
    // Mkdir
    res, _, err = suite.ns.Mkdir(suite.tempNsDir)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[mkdir] StatusCode should be 200 OK")
    fmt.Printf("[TEST] mkdir %s done\n", suite.tempNsDir)

    // Upload
    err = ioutil.WriteFile(suite.tempFile, []byte("Hello, Netstorage API World!"), 0666)
    check(err, true)
    res, _, err = suite.ns.Upload(suite.tempFile, suite.tempNsFile)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[upload] StatusCode should be 200 OK")
    fmt.Printf("[TEST] upload %s to %s done\n", suite.tempFile, suite.tempNsFile)

    // Du
    res, _, err = suite.ns.Du(suite.tempNsDir)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[du] StatusCode should be 200 OK")
    fmt.Printf("[TEST] du %s done\n", suite.tempNsDir)

    // Mtime
    currentTime := time.Now().Unix()
    res, _, err = suite.ns.Mtime(suite.tempNsFile, currentTime)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[mtime] StatusCode should be 200 OK")
    fmt.Printf("[TEST] mtime %s done\n", suite.tempNsFile)

    // Stat
    res, _, err = suite.ns.Stat(suite.tempNsFile)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[stat] StatusCode should be 200 OK")
    fmt.Printf("[TEST] stat %s done\n", suite.tempNsFile)

    // Symlink
    res, _, err = suite.ns.Symlink(suite.tempNsFile, suite.tempNsFile + "_lnk")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[symlink] StatusCode should be 200 OK")
    fmt.Printf("[TEST] symlink %s to %s done\n", suite.tempNsFile, suite.tempNsFile + "_lnk")

    // Rename
    res, _, err = suite.ns.Rename(suite.tempNsFile, suite.tempNsFile + "_rename")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[rename] StatusCode should be 200 OK")
    fmt.Printf("[TEST] rename %s to %s done\n", suite.tempNsFile, suite.tempNsFile + "_rename")

    // Download
    res, _, err = suite.ns.Download(suite.tempNsFile + "_rename")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[download] StatusCode should be 200 OK")
    fmt.Printf("[TEST] download %s done\n", suite.tempNsFile)

    // Delete
    res, _, err = suite.ns.Delete(suite.tempNsFile + "_rename")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[delete] StatusCode should be 200 OK")
    fmt.Printf("[TEST] delete %s done\n", suite.tempNsFile + "_rename")
    res, _, err = suite.ns.Delete(suite.tempNsFile + "_lnk")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[delete] StatusCode should be 200 OK")
    fmt.Printf("[TEST] delete %s done\n", suite.tempNsFile + "_lnk")

    // Rmdir
    res, _, err = suite.ns.Rmdir(suite.tempNsDir)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[rmdir] StatusCode should be 200 OK")
    fmt.Printf("[TEST] rmdir %s done\n", suite.tempNsDir)
}


func TestExampleTestSuite(t *testing.T) {
    suite.Run(t, new(NetstorageTestSuite))
}