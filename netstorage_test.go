package netstorage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
	//"./spike/secrets"
)

var nsHostname = "astin-nsu.akamaihd.net"
var nsKeyname = "astinapi"

// var nsKey = secrets.KEY // Don't expose nsKey on your public repository
var nsKey = os.Getenv("NS_KEY")
var nsCpcode = "360949"

var ns = NewNetstorage(nsHostname, nsKeyname, nsKey, false)
var tempNsDir = fmt.Sprintf("/%s/nst_%d", nsCpcode, time.Now().Unix())
var tempFile = fmt.Sprintf("nst_%d.txt", time.Now().Unix())
var tempNsFile = fmt.Sprintf("%s/%s", tempNsDir, tempFile)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func teardown() {
	// delete temp files for local
	if _, err := os.Stat(tempFile); err == nil {
		err = os.Remove(tempFile)
		check(err)
		fmt.Printf("[TEARDOWN] remove %s from local done\n", tempFile)
	}

	if _, err := os.Stat(tempFile + "_rename"); err == nil {
		err = os.Remove(tempFile + "_rename")
		check(err)
		fmt.Printf("[TEARDOWN] remove %s from local done\n", tempFile+"_rename")
	}

	// delete temp files for netstorage
	if res, _, err := ns.Delete(tempNsFile); res.StatusCode == 200 && err == nil {
		fmt.Printf("[TEARDOWN] delete %s done\n", tempNsFile)
	}

	if res, _, err := ns.Delete(tempNsFile + "_lnk"); res.StatusCode == 200 && err == nil {
		fmt.Printf("[TEARDOWN] delete %s done\n", tempNsFile+"_lnk")
	}

	if res, _, err := ns.Delete(tempNsFile + "_rename"); res.StatusCode == 200 && err == nil {
		fmt.Printf("[TEARDOWN] delete %s done\n", tempNsFile+"_rename")
	}

	if res, _, err := ns.Rmdir(tempNsDir); res.StatusCode == 200 && err == nil {
		fmt.Printf("[TEARDOWN] rmdir %s done\n", tempNsDir)
	}
}

func assertEqual(t *testing.T, got, expected interface{}, funcName, success, fail string, err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	if got != expected {
		t.Error(
			fmt.Sprintf("\n"),
			fmt.Sprintf("Function: %s\n", funcName),
			fmt.Sprintf("Expected: %v\n", expected),
			fmt.Sprintf("Got: %v\n", got),
			fmt.Sprintf("Message: %s\n", fail),
		)
		return true
	}
	fmt.Printf(success)
	return false
}

func TestNetstorage(t *testing.T) {
	fmt.Println("### Netstorage Test ###")
	res, body, err := ns.Dir("/" + nsCpcode)
	wrong := assertEqual(t, res.StatusCode, 200,
		"Dir",
		fmt.Sprintf("[TEST] Dir /%s done\n", nsCpcode),
		body,
		err,
	)
	if wrong {
		return
	}

	res, body, err = ns.Mkdir(tempNsDir)
	wrong = assertEqual(t, res.StatusCode, 200,
		"Mkdir",
		fmt.Sprintf("[TEST] Mkdir %s done\n", tempNsDir),
		body,
		err,
	)
	if wrong {
		return
	}

	testString := "Hello, Netstorage API World!"
	err = ioutil.WriteFile(tempFile, []byte(testString), 0666)
	check(err)
	res, body, err = ns.Upload(tempFile, tempNsFile)
	wrong = assertEqual(t, res.StatusCode, 200,
		"Upload",
		fmt.Sprintf("[TEST] Upload %s to %s done\n", tempFile, tempNsFile),
		body,
		err,
	)
	if wrong {
		return
	}

	res, body, err = ns.UploadContent(bytes.NewBufferString(testString), tempNsFile)
	wrong = assertEqual(t, res.StatusCode, 200,
		"Upload",
		fmt.Sprintf("[TEST] Upload content '%s' to %s done\n", testString, tempNsFile),
		body,
		err,
	)
	if wrong {
		return
	}

	res, body, err = ns.Du(tempNsDir)
	wrong = assertEqual(t, res.StatusCode, 200,
		"Du",
		fmt.Sprintf("[TEST] Du %s done\n", tempNsDir),
		body,
		err,
	)
	if wrong {
		return
	}

	currentTime := time.Now().Unix()
	res, body, err = ns.Mtime(tempNsFile, currentTime)
	wrong = assertEqual(t, res.StatusCode, 200,
		"Mtime",
		fmt.Sprintf("[TEST] Mtime %s done\n", tempNsFile),
		body,
		err,
	)
	if wrong {
		return
	}

	res, body, err = ns.Stat(tempNsFile)
	wrong = assertEqual(t, res.StatusCode, 200,
		"Stat",
		fmt.Sprintf("[TEST] Stat %s done\n", tempNsFile),
		body,
		err,
	)
	if wrong {
		return
	}

	res, body, err = ns.Symlink(tempNsFile, tempNsFile+"_lnk")
	wrong = assertEqual(t, res.StatusCode, 200,
		"Symlink",
		fmt.Sprintf("[TEST] Symlink %s to %s done\n", tempNsFile, tempNsFile+"_lnk"),
		body,
		err,
	)
	if wrong {
		return
	}

	res, body, err = ns.Rename(tempNsFile, tempNsFile+"_rename")
	wrong = assertEqual(t, res.StatusCode, 200,
		"Rename",
		fmt.Sprintf("[TEST] Rename %s to %s done\n", tempNsFile, tempNsFile+"_rename"),
		body,
		err,
	)
	if wrong {
		return
	}

	res, body, err = ns.Download(tempNsFile + "_rename")
	data, err := ioutil.ReadFile(tempFile + "_rename")
	check(err)
	wrong = assertEqual(t, string(data), testString,
		"Download",
		fmt.Sprintf("[TEST] Download %s done\n", tempNsFile),
		"Download Fail",
		err,
	)
	if wrong {
		return
	}

	res, body, err = ns.Delete(tempNsFile + "_rename")
	wrong = assertEqual(t, res.StatusCode, 200,
		"Delete",
		fmt.Sprintf("[TEST] delete %s done\n", tempNsFile+"_rename"),
		body,
		err,
	)
	if wrong {
		return
	}
	res, body, err = ns.Delete(tempNsFile + "_lnk")
	wrong = assertEqual(t, res.StatusCode, 200,
		"Delete",
		fmt.Sprintf("[TEST] delete %s done\n", tempNsFile+"_lnk"),
		body,
		err,
	)
	if wrong {
		return
	}

	res, body, err = ns.Rmdir(tempNsDir)
	wrong = assertEqual(t, res.StatusCode, 200,
		"Rmdir",
		fmt.Sprintf("[TEST] rmdir %s done\n", tempNsDir),
		body,
		err,
	)
	if wrong {
		return
	}
	fmt.Println("")

}

func TestNetstorageError(t *testing.T) {
	fmt.Println("### Error Test ###")
	_, body, err := ns.Dir("invalid ns path")
	wrong := assertEqual(t, reflect.TypeOf(err).String(), "*errors.errorString",
		"Dir",
		fmt.Sprintf("[TEST] Dir: netstorage invalid path test done\n"),
		body,
		nil,
	)
	if wrong {
		return
	}

	_, body, err = ns.Upload("invalid local path", tempNsFile)
	wrong = assertEqual(t, reflect.TypeOf(err).String(), "*os.PathError",
		"Upload",
		fmt.Sprintf("[TEST] Upload: local invalid path test done\n"),
		body,
		nil,
	)
	if wrong {
		return
	}

	_, body, err = ns.Download("/123456/directory/", tempFile)
	wrong = assertEqual(t, reflect.TypeOf(err).String(), "*errors.errorString",
		"Download",
		fmt.Sprintf("[TEST] Download: netstorage directory path test done\n"),
		body,
		nil,
	)
	if wrong {
		return
	}

	fmt.Println("")
}

func TestMain(m *testing.M) {
	retCode := m.Run()

	teardown()

	os.Exit(retCode)
}
