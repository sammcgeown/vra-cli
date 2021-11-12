/*
Package helpers Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package helpers

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hokaccha/go-prettyjson"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/vra-sdk-go/pkg/models"
)

// IsInputFromPipe - pipeline detection
func IsInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

// // readInputFromPipe
// func readInputFromPipe(r io.Reader, w io.Writer) error {
// 	pipeScanner := bufio.NewScanner(bufio.NewReader(r))
// 	for pipeScanner.Scan() {
// 		_, e := fmt.Fprintln(w, pipeScanner.Text())
// 		if e != nil {
// 			return e
// 		}
// 	}
// 	return nil
// }

// PrettyPrint prints interfaces
func PrettyPrint(v interface{}) (err error) {
	s, err := prettyjson.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Println(string(s))

	// b, err := json.MarshalIndent(v, "", "  ")
	// if err == nil {
	// 	fmt.Println(string(b))
	// }
	return nil
}

// PrintTable prints an array of objects with table headers
func PrintTable(objects []interface{}, headers []string) {
	// Print result table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, object := range objects {
		t := reflect.TypeOf(object)
		fmt.Println(t)
		// var o t
		// mapstructure.Decode(object, &t)
		// var values []string
		// for _, header := range headers {
		// 	append(values, object.(*t).header)
		// }
		// 	table.Append([]string{c.ID, c.Name, c.Project})
	}
	table.Render()
}

// GetFilePaths - Get all files of type in a directory
func GetFilePaths(filePath string, filetype string) []string {
	var files []string
	// Read importPath
	stat, err := os.Stat(filePath)
	if err == nil && stat.IsDir() {
		log.Debugln("GetFilePaths -", filePath, "is a directory, looking for filetype", filetype)
		fileList, err := ioutil.ReadDir(filePath)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range fileList {
			if strings.Contains(f.Name(), filetype) {
				absPath, _ := filepath.Abs(filepath.Join(filePath, f.Name()))
				files = append(files, absPath)
			}
		}
	} else {
		log.Debugln("importPath is a file")
		absPath, _ := filepath.Abs(filePath)
		files = append(files, absPath)
	}
	return files
}

// RemoveDuplicateStrings - remote duplicate strings from a slice
func RemoveDuplicateStrings(elements []string) []string {
	encountered := map[string]bool{}
	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}
	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

// ZipFiles compresses one or many files into a single zip archive file.
func ZipFiles(filename string, files []string, basedir string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file, basedir); err != nil {
			return err
		}
	}
	return nil
}

// AddFileToZip - adds a file to a zip archive
func AddFileToZip(zipWriter *zip.Writer, filename string, basedir string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	fmt.Println(basedir)
	fmt.Println(filename)
	fmt.Println(strings.ReplaceAll(filename, basedir, ""))
	header.Name = strings.ReplaceAll(filename, basedir, "")

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// func GetInputsFromSchema(schema *cmd.CloudAssemblyCloudTemplateInputSchema) cmd.DeploymentInput {
// 	var inputs cmd.DeploymentInput
// 	for name := range schema.Properties {
// 		log.Infoln(name)
// 		// c := CloudAssemblyCloudTemplateInputProperty{}
// 		// mapstructure.Decode(value, &c)
// 		inputs.Inputs[name] = "test"
// 	}

// 	return inputs
// }

// AskForConfirmation - Credit - https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		log.Warnf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// promptUserForInputs
// func GetCatalogItemInputs(SchemaProperties map[string]cmd.CatalogItemSchemaProperties) map[string]string {
// 	inputs := make(map[string]string)
// 	for name, schema := range SchemaProperties {
// 		fmt.Printf(name + "[" + schema.Type + "]: ")
// 		var response string
// 		fmt.Scanln(&response)
// 		inputs[name] = response
// 	}
// 	return inputs
// }

// StringToTags - split a string into tags
func StringToTags(tags string) []*models.Tag {
	var tagsArray []*models.Tag
	if tags == "" {
		return tagsArray
	}
	for _, tag := range strings.Split(tags, ",") {
		tagKey := strings.Split(tag, ":")[0]
		tagValue := strings.Split(tag, ":")[1]
		tagsArray = append(tagsArray, &models.Tag{
			Key:   &tagKey,
			Value: &tagValue,
		})
	}
	return tagsArray
}

// StringArrayContains - check if a string array contains a string
func StringArrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// CreateUserArray - Create an array of users from emails
func CreateUserArray(emails []string) []*models.User {
	if emails[0] == "" {
		return nil
	}
	users := make([]*models.User, 0, len(emails))
	for i := range emails {
		user := models.User{
			Email: &emails[i],
		}
		users = append(users, &user)
	}
	return users
}

// UserArrayToString - Create an comma separated list of users from an array of users
func UserArrayToString(users []*models.User) string {
	var userList []string
	for user := range users {
		userList = append(userList, *users[user].Email)
	}
	return strings.Join(userList, ",")

}
