package files_deleter

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	loggermaker "github.com/xfiendx4life/gb_gobest/hw2/logger_maker"
	"go.uber.org/zap"
)

var (
	logLevel = zap.LevelFlag("loglevel", zap.InfoLevel, "set logging level")
)

func init() {
	logger = *loggermaker.InitLogger(logLevel, os.Getenv("FILELOG"))
}

func createNFiles(path string, n int) error {
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("%s/file%d.txt", path, i)
		file, err := os.Create(name)
		if err != nil {
			return fmt.Errorf("can't create file %s", err)
		}
		defer file.Close()
		_, err = file.Write([]byte(name))
		if err != nil {
			return fmt.Errorf("can't write to file %s", err)
		}
	}
	return nil
}

func createNCatalogs(path string, n int, createNFiles func(string, int) error, numFiles int) error {
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("%s/test%d", path, i)
		err := os.Mkdir(name, 0777)
		if err != nil {
			return err
		}
		errFiles := createNFiles(name, numFiles)
		if errFiles != nil {
			return errFiles
		}

	}
	return nil
}

func CreateDupl(pathToOrigin, pathToDest string) error {
	origin, err := os.Open(pathToOrigin)
	if err != nil {
		return err
	}
	defer origin.Close()
	dest, err := os.Create(pathToDest)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, origin)
	if err != nil {
		return err
	}
	return nil
}

func TestTwoCopiesInDifferentFolders(t *testing.T) {
	os.Mkdir("./test", 0777)
	createNCatalogs("./test", 3, createNFiles, 4)
	CreateDupl("./test/test1/file0.txt", "./test/test1/file0copy")
	CreateDupl("./test/test1/file2.txt", "./test/test2/file2_1copy")
	defer os.RemoveAll("./test")
	n, err := Delete("./test", false, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 2, n)
	}
}

func TestTwoCopiesInsameFolder(t *testing.T) {
	os.Mkdir("./test", 0777)
	createNCatalogs("./test", 3, createNFiles, 4)
	CreateDupl("./test/test1/file0.txt", "./test/test1/file0copy")
	CreateDupl("./test/test1/file2.txt", "./test/test1/file2_1copy")
	defer os.RemoveAll("./test")
	n, err := Delete("./test", false, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 2, n)
	}
}

func TestNoCopies(t *testing.T) {
	os.Mkdir("./test", 0777)
	createNCatalogs("./test", 3, createNFiles, 4)
	defer os.RemoveAll("./test")
	n, err := Delete("./test", false, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, n)
	}
}

func TestThreeCopiesInInternalFolders(t *testing.T) {
	os.Mkdir("./test", 0777)
	createNCatalogs("./test", 3, createNFiles, 4)
	createNCatalogs("./test/test1", 2, createNFiles, 2)
	CreateDupl("./test/test1/file0.txt", "./test/test1/file0copy")
	CreateDupl("./test/test1/test1/file0.txt", "./test/test2/file2_1copy")
	CreateDupl("./test/test1/test1/file1.txt", "./test/test1/test0/file3_1copy")
	defer os.RemoveAll("./test")
	n, err := Delete("./test", false, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 3, n)
	}
}

func TestEmptyFolders(t *testing.T) {
	os.Mkdir("./test", 0777)
	createNCatalogs("./test", 3, createNFiles, 0)
	defer os.RemoveAll("./test")
	n, err := Delete("./test", false, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, n)
	}
}

// _____________________________----------------------

func TestDeleteTwoCopiesInDifferentFolders(t *testing.T) {
	os.Mkdir("./test", 0777)
	createNCatalogs("./test", 3, createNFiles, 4)
	CreateDupl("./test/test1/file0.txt", "./test/test1/file0copy")
	CreateDupl("./test/test1/file2.txt", "./test/test2/file2_1copy")
	defer os.RemoveAll("./test")
	n, err := Delete("./test", true, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 2, n)
	}
	n, err = Delete("./test", false, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, n)
	}
}

func TestDeleteTwoCopiesInsameFolder(t *testing.T) {
	os.Mkdir("./test", 0777)
	createNCatalogs("./test", 3, createNFiles, 4)
	CreateDupl("./test/test1/file0.txt", "./test/test1/file0copy")
	CreateDupl("./test/test1/file2.txt", "./test/test1/file2_1copy")
	defer os.RemoveAll("./test")
	n, err := Delete("./test", true, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 2, n)
	}
	n, err = Delete("./test", false, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, n)
	}
}

func TestDeleteThreeCopiesInInternalFolders(t *testing.T) {
	os.Mkdir("./test", 0777)
	createNCatalogs("./test", 3, createNFiles, 4)
	createNCatalogs("./test/test1", 2, createNFiles, 2)
	CreateDupl("./test/test1/file0.txt", "./test/test1/file0copy")
	CreateDupl("./test/test1/test1/file0.txt", "./test/test2/file2_1copy")
	CreateDupl("./test/test1/test1/file1.txt", "./test/test1/test0/file3_1copy")
	defer os.RemoveAll("./test")
	n, err := Delete("./test", true, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 3, n)
	}
	n, err = Delete("./test", false, &logger)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, n)
	}
}

func ExampleDelete() {
	n, err := Delete("./", false, &logger)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Number of duplicates %d\n", n)
	// Output: Found duplicates .//2.txt, .//1.txt
	// Number of duplicates 1
}

func TestWrongDirectory(t *testing.T) {
	_, err := Delete("./wrongDir", false, &logger)
	assert.NotNil(t, err)
}
