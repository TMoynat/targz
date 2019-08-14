package targz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func check_err(err error, msg string) {
	if err != nil {
		fmt.Println("Error : ", msg, " failed...")
	}
	os.Exit(-1)
}

func Ungzipator(gzipstream io.Reader) {
	fmt.Println("UntarGz in progress...")
	uncompressed, err := gzip.NewReader(gzipstream)
	check_err(err, "uncompression")
	tarstream := tar.NewReader(uncompressed)
	for true {
		header, err := tarstream.Next()
		if err == io.EOF {
			break
		}
		check_err(err, "iteration on files")
		switch header.Typeflag {
		case tar.TypeDir:
			err := os.Mkdir(header.Name, 0755)
			check_err(err, "folder creation")
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			check_err(err, "file creation")
			defer outFile.Close()
			bytes, err := io.Copy(outFile, tarstream)
			check_err(err, "copy")
			fmt.Println(bytes, " bytes written in : ", header.Name)
		default:
			fmt.Println("Error : unknown type: %s in %s", header.Typeflag, header.Name)
		}
	}
}