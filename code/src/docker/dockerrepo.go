package docker

import( 	
	"io/ioutil"
)

//------------------------------------------------------------------------------
/**
	@type DockerRepo
	
	A docker repo contains a list of all currently available docker files contained within a folder.
	This type loads all docker images in an entire folder. This type also allows for retrieving a docker image using the image name.
*/
type DockerRepo struct {
	images [] DockerImage
}

//------------------------------------------------------------------------------
/**
	Loads a folder containing docker files into local repository.
	Docker files can then be searched and retrieved.
*/
func (repo* DockerRepo) Load(path string) error {
	
	// open directory and list all files in it
	list, err := ioutil.ReadDir(path)
	
	// if this fails, return an error
	if (err != nil) {
		return err
	}
	
	// allocate array which will contain our files
	size := len(list)
	repo.images = make([] DockerImage, size, size)
	
	// loop through files and create docker images for them
	for index, element := range list {
	
		// create image
		image := NewDockerImage()
		
		// set local path to image
		image.path = path + "/" + element.Name()
		
		// save the name of the image file
		image.name = element.Name()
		
		// attempt to open it
		err := image.Open()
		
		if (err == nil) {
			// if we succeeded in opening (and validating) the file, add it to the repo
			repo.images[index] = image
		} else {
			// otherwise return error
			return err
		}
	}
	
	// return no error
	return nil
}

//------------------------------------------------------------------------------
/**
	Searches for a docker image with the given name, returns image upon success.
*/
func (repo* DockerRepo) Find(name string) *DockerImage {

	// search for image
	for _, element := range repo.images {
	
		// if the image is found, return it
		if (element.name == name) {
			return &element
		}
	}
	
	// if no image is found, return nil
	return nil
}

//------------------------------------------------------------------------------
/**
	Searches for a docker image, and returns true if it exists.
*/
func (repo* DockerRepo) Exists(name string) bool {

	// search for image
	for _, element := range repo.images {
	
		// if the image is found, return it
		if (element.name == name) {
			return true
		}
	}
	
	return false
}