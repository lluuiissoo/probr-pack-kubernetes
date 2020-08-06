package features

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"citihub.com/probr/internal/coreengine"
)

//this is the "TEST HANDLER" impl  and will get called when probr is invoked from the CLI or API
//all we do here is set the godog args based on what has been supplied (e.g. output path)
//and call to the "feature" implementation (i.e the same impl when godog / go test is invoked)

//GodogTestHandler ...
func GodogTestHandler(gd *coreengine.GodogTest) (int, error) {

	f, err := getFeaturesPath(gd)
	if err != nil {
		return -1, err
	}

	o, err := GetOutputPath(&gd.TestDescriptor.Name)
	if err != nil {
		return -2, err
	}

	opts := godog.Options{
		Format: "cucumber",
		Output: colors.Colored(o),
		Paths:  []string{f},
	}

	status := godog.TestSuite{
		Name:                 gd.TestDescriptor.Name,
		TestSuiteInitializer: gd.TestSuiteInitializer,
		ScenarioInitializer:  gd.ScenarioInitializer,
		Options:              &opts,
	}.Run()

	return status, nil
}

func getFeaturesPath(gd *coreengine.GodogTest) (string, error) {
	r, err := GetRootDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine root directory - not able to perform tests")
	}

	if gd.FeaturePath != nil {
		//if we've been given a feature path, add to root and return:
		return filepath.Join(r, *gd.FeaturePath), nil
	}

	//otherwise derive it from the group and category data:
	var g = gd.TestDescriptor.Group.String()
	var c = gd.TestDescriptor.Category.String()

	return filepath.Join(r, "test", "features",
		strings.ReplaceAll(strings.ToLower(g), " ", ""),
		strings.ReplaceAll(strings.ToLower(c), " ", ""), "features"), nil

}