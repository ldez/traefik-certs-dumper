package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/distribution/manifest/manifestlist"
)

// Publisher Publish multi-arch image.
type Publisher struct {
	Builds           [][]string
	Push             [][]string
	ManifestAnnotate [][]string
	ManifestCreate   []string
	ManifestPush     []string
}

func newPublisher(imageName, version, baseImageName string, targets []string) (Publisher, error) {
	manifest, err := getManifest(baseImageName)
	if err != nil {
		return Publisher{}, err
	}

	buildOptions, err := getBuildOptions("./internal/build-options.json")
	if err != nil {
		return Publisher{}, err
	}

	publisher := Publisher{}

	for _, target := range targets {
		option := buildOptions[target]

		descriptor, err := findManifestDescriptor(option, manifest.Manifests)
		if err != nil {
			log.Fatal(err)
		}

		dockerfile := fmt.Sprintf("%s-%s-%s.Dockerfile", option.OS, option.GoARCH, option.GoARM)

		publisher.Builds = append(publisher.Builds, []string{
			"build",
			"-t", fmt.Sprintf("%s:%s-%s", imageName, version, target),
			"-f", dockerfile,
			".",
		})

		err = createDockerfile(dockerfile, option, descriptor, baseImageName)
		if err != nil {
			log.Fatal(err)
		}

		publisher.Push = append(publisher.Push, []string{"push", fmt.Sprintf(`%s:%s-%s`, imageName, version, target)})

		ma := []string{
			"manifest", "annotate",
			fmt.Sprintf("%s:%s", imageName, version),
			fmt.Sprintf("%s:%s-%s", imageName, version, target),
			fmt.Sprintf("--os=%s", option.OS),
			fmt.Sprintf("--arch=%s", option.GoARCH),
		}
		if option.Variant != "" {
			ma = append(ma, fmt.Sprintf("--variant=%s", option.Variant))
		}
		publisher.ManifestAnnotate = append(publisher.ManifestAnnotate, ma)
	}

	publisher.ManifestCreate = []string{
		"manifest", "create", "--amend", fmt.Sprintf("%s:%s", imageName, version),
	}

	for _, target := range targets {
		publisher.ManifestCreate = append(publisher.ManifestCreate, fmt.Sprintf("%s:%s-%s", imageName, version, target))
	}

	publisher.ManifestPush = []string{
		"manifest", "push", fmt.Sprintf("%s:%s", imageName, version),
	}

	return publisher, nil
}

func (b Publisher) execute(dryRun bool) error {
	for _, args := range b.Builds {
		if err := execDocker(args, dryRun); err != nil {
			return fmt.Errorf("failed to build: %v: %v", args, err)
		}
	}

	for _, args := range b.Push {
		if err := execDocker(args, dryRun); err != nil {
			return fmt.Errorf("failed to push: %v: %v", args, err)
		}
	}

	if err := execDocker(b.ManifestCreate, dryRun); err != nil {
		return fmt.Errorf("failed to create manifest: %v: %v", b.ManifestCreate, err)
	}

	for _, args := range b.ManifestAnnotate {
		if err := execDocker(args, dryRun); err != nil {
			return fmt.Errorf("failed to annotate manifest: %v: %v", args, err)
		}
	}

	if err := execDocker(b.ManifestPush, dryRun); err != nil {
		return fmt.Errorf("failed to push manifest: %v: %v", b.ManifestPush, err)
	}

	return nil
}

func getBuildOptions(source string) (map[string]buildOption, error) {
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}

	buildOptions := make(map[string]buildOption)

	err = json.NewDecoder(file).Decode(&buildOptions)
	if err != nil {
		return nil, err
	}

	return buildOptions, nil
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

		err = ioutil.WriteFile(manifestPath, output, 0666)
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

func execDocker(args []string, dryRun bool) error {
	if dryRun {
		fmt.Println("docker", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("docker", args...)
	cmd.Env = append(cmd.Env, "DOCKER_CLI_EXPERIMENTAL=enabled")

	output, err := cmd.CombinedOutput()

	if len(output) != 0 {
		log.Println(string(output))
	}

	if err != nil {
		return err
	}
	return nil
}
