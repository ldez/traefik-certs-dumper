package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/docker/distribution/manifest/manifestlist"
)

type buildOption struct {
	OS      string
	GoARM   string
	GoARCH  string
	Variant string
}

type actions struct {
	Builds           [][]string
	Push             [][]string
	ManifestAnnotate [][]string
	ManifestCreate   []string
	ManifestPush     []string
}

func main() {
	log.SetFlags(log.Lshortfile)

	imageName := flag.String("image-name", "ldez/traefik-certs-dumper", "")
	version := flag.String("version", "", "")
	baseImageName := flag.String("base-image-name", "alpine:3.9", "")
	dryRun := flag.Bool("dry-run", true, "")

	flag.Parse()

	require("image-name", imageName)
	require("version", version)
	require("base-image-name", baseImageName)

	// FIXME
	// _, travisTag := os.LookupEnv("TRAVIS_TAG")
	// if !travisTag {
	// 	log.Println("Skipping deploy")
	// 	os.Exit(0)
	// }

	targets := []string{"arm.v6", "arm.v7", "arm.v8", "amd64", "386"}

	actions, err := buildActions(*imageName, *version, *baseImageName, targets)
	if err != nil {
		log.Fatal(err)
	}

	err = execute(actions, *dryRun)
	if err != nil {
		log.Fatal(err)
	}
}

func require(fieldName string, field *string) {
	if field == nil || *field == "" {
		log.Fatalf("%s is required", fieldName)
	}
}

func buildActions(imageName, version, baseImageName string, targets []string) (actions, error) {
	manifest, err := getManifest(baseImageName)
	if err != nil {
		return actions{}, err
	}

	buildOptions := map[string]buildOption{
		"arm.v5": {OS: "linux", GoARM: "5", GoARCH: "arm", Variant: "v5"},
		"arm.v6": {OS: "linux", GoARM: "6", GoARCH: "arm", Variant: "v6"},
		"arm.v7": {OS: "linux", GoARM: "7", GoARCH: "arm", Variant: "v7"},
		"arm.v8": {OS: "linux", GoARCH: "arm64", Variant: "v8"},
		"amd64":  {OS: "linux", GoARCH: "amd64"},
		"386":    {OS: "linux", GoARCH: "386"},
	}

	actions := actions{}

	for _, target := range targets {
		buildOption := buildOptions[target]

		descriptor, err := findManifestDescriptor(buildOption, manifest.Manifests)
		if err != nil {
			log.Fatal(err)
		}

		dockerfile := fmt.Sprintf("%s-%s-%s.Dockerfile", buildOption.OS, buildOption.GoARCH, buildOption.GoARM)

		actions.Builds = append(actions.Builds, []string{
			"build",
			"-t", fmt.Sprintf("%s:%s-%s", imageName, version, target),
			"-f", dockerfile,
			".",
		})

		err = createDockerfile(dockerfile, buildOption, descriptor, baseImageName)
		if err != nil {
			log.Fatal(err)
		}

		actions.Push = append(actions.Push, []string{"push", fmt.Sprintf(`%s:%s-%s`, imageName, version, target)})

		ma := []string{
			"manifest", "annotate",
			fmt.Sprintf(`"%s:%s"`, imageName, version),
			fmt.Sprintf(`"%s:%s-%s"`, imageName, version, target),
			fmt.Sprintf(`--os="%s"`, buildOption.OS),
			fmt.Sprintf(`--arch="%s"`, buildOption.GoARCH),
		}
		if buildOption.Variant != "" {
			ma = append(ma, fmt.Sprintf(`--variant="%s"`, buildOption.Variant))
		}
		actions.ManifestAnnotate = append(actions.ManifestAnnotate, ma)
	}

	actions.ManifestCreate = []string{
		"manifest", "create", "--amend",
		fmt.Sprintf("%s:%s", imageName, version),
	}

	for _, target := range targets {
		actions.ManifestCreate = append(actions.ManifestCreate, fmt.Sprintf(`"%s:%s-%s"`, imageName, version, target))
	}

	actions.ManifestPush = []string{
		"manifest", "push", fmt.Sprintf("%s:%s", imageName, version),
	}

	return actions, nil
}

func createDockerfile(dockerfile string, buildOption buildOption, descriptor manifestlist.ManifestDescriptor, baseImageName string) error {
	base := template.New("tmpl.Dockerfile")
	parse, err := base.ParseFiles("./internal/tmpl.Dockerfile")
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"GoOS":         buildOption.OS,
		"GoARCH":       buildOption.GoARCH,
		"GoARM":        buildOption.GoARM,
		"RuntimeImage": fmt.Sprintf("%s@%s", baseImageName, descriptor.Digest),
	}

	file, err := os.Create(dockerfile)
	if err != nil {
		return err
	}

	return parse.Execute(file, data)
}

func getManifest(baseImageName string) (*manifestlist.ManifestList, error) {
	manifestPath := "./manifest.json"

	if _, errExist := os.Stat(manifestPath); os.IsNotExist(errExist) {
		cmd := exec.Command("docker", "manifest", "inspect", baseImageName)
		cmd.Env = append(cmd.Env, "DOCKER_CLI_EXPERIMENTAL=enabled")

		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(manifestPath, output, 777)
		if err != nil {
			return nil, err
		}
	} else if errExist != nil {
		return nil, errExist
	}

	bytes, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	manifest := &manifestlist.ManifestList{}

	err = json.Unmarshal(bytes, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func findManifestDescriptor(criterion buildOption, descriptors []manifestlist.ManifestDescriptor) (manifestlist.ManifestDescriptor, error) {
	for _, descriptor := range descriptors {
		if descriptor.Platform.OS == criterion.OS &&
			descriptor.Platform.Architecture == criterion.GoARCH &&
			descriptor.Platform.Variant == criterion.Variant {
			return descriptor, nil
		}
	}
	return manifestlist.ManifestDescriptor{}, fmt.Errorf("not supported: %v", criterion)
}

func execute(actions actions, dryRun bool) error {
	for _, args := range actions.Builds {
		if err := execDocker(args, dryRun); err != nil {
			return err
		}
	}

	return nil

	// for _, args := range actions.Push {
	// 	if err := execDocker(args, dryRun); err != nil {
	// 		return err
	// 	}
	// }
	//
	// if err := execDocker(actions.ManifestCreate, dryRun); err != nil {
	// 	return err
	// }
	//
	// for _, args := range actions.ManifestAnnotate {
	// 	if err := execDocker(args, dryRun); err != nil {
	// 		return err
	// 	}
	// }
	//
	// return execDocker(actions.ManifestPush, dryRun)
}

func execDocker(args []string, dryRun bool) error {
	if dryRun {
		fmt.Println("docker", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("docker", args...)
	cmd.Env = append(cmd.Env, "DOCKER_CLI_EXPERIMENTAL=enabled")

	output, err := cmd.CombinedOutput()

	log.Println(string(output))

	if err != nil {
		return err
	}
	return nil
}
