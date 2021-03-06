package docker

import(
	"errors"
	"path/filepath"
	"io/ioutil"
	"fmt"
)

//------------------------------------------------------------------------------
/**
	@type DockerImage
	
	Represents a docker image.
	Currently, the only validation it performs is to check the file extension to be .dkf and that the file isn't empty.
*/
type DockerImage struct {
	path string
	name string
	content string
	isLoaded bool
}

//------------------------------------------------------------------------------
/**
	Use this constructor to create a new docker image!!!
*/
func NewDockerImage() DockerImage {
	return DockerImage{isLoaded : false}
}

//------------------------------------------------------------------------------
/**
	Opens a docker image.
	This will open the file handle, read the file into the image content field, then check it for validation.
*/
func (image* DockerImage) Open() error {

	// attempt to open file
	image.isLoaded = false
	content, err := ioutil.ReadFile(image.path)
	
	// return error if this didn't succeed
	if (err != nil) {
		return err
	}

	// now set string in object
	image.content = string(content)
	
	// validate directly
	err = image.Validate()
	
	if (err == nil) {
		fmt.Printf("Docker image '%s' loaded successfully!\n", image.path)
		image.isLoaded = true
		return nil
	} else {
		return err
	}

	return nil
}

//------------------------------------------------------------------------------
/**
	Validates if this image is a valid docker file.
	First checks that the file has the correct extension (.dkf).
	Then further validates file.
*/
func (image* DockerImage) Validate() error {
	
	// first validate extension
	compare := ".dkf"
	if (filepath.Ext(image.path) != compare) {
		msg := fmt.Sprintf("Docker image file '%s' does not have extension .dkf", image.path)
		return errors.New(msg)
	}
	
	// check that the file isn't empty
	if (len(image.content) == 0) {
		msg := fmt.Sprintf("Docker image file is '%s' is empty", image.path)
		return errors.New(msg) 
	}
	
	return nil
}
