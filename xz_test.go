// Copyright 2014-2021 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xz_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ulikunitz/xz"
)

func TestPanic(t *testing.T) {
	data := []byte([]uint8{253, 55, 122, 88, 90, 0, 0, 0, 255, 18, 217, 65, 0, 189, 191, 239, 189, 191, 239, 48})
	t.Logf("%q", string(data))
	t.Logf("0x%02x", data)
	r, err := xz.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Logf("xz.NewReader error %s", err)
		return
	}
	_, err = ioutil.ReadAll(r)
	if err != nil {
		t.Logf("ioutil.ReadAll(r) error %s", err)
		return
	}
}

func FuzzXZ(f *testing.F) {
	/*
		Simple compression -> decompression fuzz test
		Also checks if files are compatible with the
		xz binary
	*/

	if _, err := exec.Command("xz", "-h").CombinedOutput(); err != nil {
		f.Log("xz binary not found. Skipping.")
		f.Skip()
	}

	f.Add([]byte{})
	f.Add([]byte("Hello world"))
	f.Add([]byte("Pure golang package for reading and writing xz-compressed files"))
	rm, err := os.ReadFile("README.md")
	if err != nil {
		f.Fatal(err.Error())
	}
	f.Add(rm)

	f.Fuzz(func(t *testing.T, text []byte) {
		var buf bytes.Buffer
		w, err := xz.NewWriter(&buf)
		if err != nil {
			t.Fatal(err.Error())
		}
		if _, err := w.Write(text); err != nil {
			t.Fatal(err.Error())
		}
		// close w
		if err := w.Close(); err != nil {
			t.Fatal(err.Error())
		}

		// read
		r, err := xz.NewReader(&buf)
		if err != nil {
			t.Fatal(err.Error())
		}
		read, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err.Error())
		}
		if !bytes.Equal(text, read) {
			t.Fatal("For this seed, output did not match.")
		}

		/*
			Write text to file and run xz on it.
		*/
		file1, err := os.CreateTemp("", "")
		if err != nil {
			t.Fatal(err.Error())
		}
		if _, err := file1.Write(text); err != nil {
			t.Fatal(err.Error())
		}
		file1.Close()
		defer os.Remove(file1.Name())
		defer os.Remove(file1.Name() + ".xz")
		if out, err := exec.Command("xz", "-z", file1.Name()).CombinedOutput(); err != nil {
			t.Logf("%s", out)
			t.Fatal(err.Error())
		}
		// read this compressed file and then attempt to decompress it
		compressedBytes, err := os.ReadFile(file1.Name() + ".xz")
		if err != nil {
			t.Fatal(err.Error())
		}
		compressedBytesReader, err := xz.NewReader(bytes.NewBuffer(compressedBytes))
		if err != nil {
			t.Fatal(err.Error())
		}
		compressedBytesReaderOut, err := io.ReadAll(compressedBytesReader)
		if err != nil {
			t.Fatal(err.Error())
		}
		if !bytes.Equal(compressedBytesReaderOut, text) {
			t.Fatal("For file compressed by xz, decompressed by github.com/ulikunitz/xz, result does not match")
		}

		/*
			Test file compressed by github.com/ulikunitz/xz, decompressed by xz
		*/

		file2, err := os.CreateTemp("", "*.xz")
		if err != nil {
			t.Fatal(err.Error())
		}
		defer os.Remove(file2.Name())
		defer os.Remove(strings.TrimSuffix(file2.Name(), ".xz"))

		var buf2 bytes.Buffer
		w2, err := xz.NewWriter(&buf2)
		if err != nil {
			t.Fatal(err.Error())
		}
		if _, err := w2.Write(text); err != nil {
			t.Fatal(err.Error())
		}
		// close w
		if err := w2.Close(); err != nil {
			t.Fatal(err.Error())
		}
		if _, err := io.Copy(file2, &buf2); err != nil {
			t.Fatal(err.Error())
		}
		file2.Close()

		// decompress using xz
		if out, err := exec.Command("xz", "-d", file2.Name()).CombinedOutput(); err != nil {
			t.Logf("%s", out)
			t.Fatal(err.Error())
		}

		decompressedBytes, err := os.ReadFile(strings.TrimSuffix(file2.Name(), ".xz"))
		if err != nil {
			t.Fatal(err.Error())
		}
		if !bytes.Equal(decompressedBytes, text) {
			t.Fatal("For file compressed by github.com/ulikunitz/xz, decompressed by xz, result does not match")
		}
	})
}
