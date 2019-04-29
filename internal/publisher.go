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

const envDockerExperimental = "DOCKER_CLI_EXPERIMENTAL=enabled"

// Publisher Publish multi-arch image.
type Publisher struct {
	Builds           []*exec.Cmd
	Push             []*exec.Cmd
	ManifestAnnotate []*exec.Cmd
	ManifestCreate   *exec.Cmd
	ManifestPush     *exec.Cmd
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
		option, ok := buildOptions[target]
		if !ok {
			return Publisher{}, fmt.Errorf("unsupported platform: %s", target)
		}

		descriptor, err := findManifestDescriptor(option, manifest.Manifests)
		if err != nil {
			return Publisher{}, err
		}

		dockerfile := fmt.Sprintf("%s-%s-%s.Dockerfile", option.OS, option.GoARCH, option.GoARM)

		err = createDockerfile(dockerfile, option, descriptor, baseImageName)
		if err != nil {
			return Publisher{}, err
		}

		dBuild := exec.Command("docker", "build",
			"-t", fmt.Sprintf("%s:%s-%s", imageName, version, target),
			"-f", dockerfile,
			".")
		publisher.Builds = append(publisher.Builds, dBuild)

		dPush := exec.Command("docker", "push", fmt.Sprintf(`%s:%s-%s`, imageName, version, target))
		publisher.Push = append(publisher.Push, dPush)

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

		cmdMA := exec.Command("docker", ma...)
		cmdMA.Env = append(cmdMA.Env, envDockerExperimental)
		publisher.ManifestAnnotate = append(publisher.ManifestAnnotate, cmdMA)
	}

	mc := []string{
		"manifest", "create", "--amend", fmt.Sprintf("%s:%s", imageName, version),
	}
	for _, target := range targets {
		mc = append(mc, fmt.Sprintf("%s:%s-%s", imageName, version, target))
	}

	cmdMC := exec.Command("docker", mc...)
	cmdMC.Env = append(cmdMC.Env, envDockerExperimental)
	publisher.ManifestCreate = cmdMC

	cmdMP := exec.Command("docker", "manifest", "push", fmt.Sprintf("%s:%s", imageName, version))
	cmdMP.Env = append(cmdMP.Env, envDockerExperimental)
	publisher.ManifestPush = cmdMP

	return publisher, nil
}

func (b Publisher) execute(dryRun bool) error {
	for _, cmd := range b.Builds {
		if err := execCmd(cmd, dryRun); err != nil {
			return fmt.Errorf("failed to build: %v: %v", cmd, err)
		}
	}

	for _, cmd := range b.Push {
		if err := execCmd(cmd, dryRun); err != nil {
			return fmt.Errorf("failed to push: %v: %v", cmd, err)
		}
	}

	if err := execCmd(b.ManifestCreate, dryRun); err != nil {
		return fmt.Errorf("failed to create manifest: %v: %v", b.ManifestCreate, err)
	}

	for _, cmd := range b.ManifestAnnotate {
		if err := execCmd(cmd, dryRun); err != nil {
			return fmt.Errorf("failed to annotate manifest: %v: %v", cmd, err)
		}
	}

	if err := execCmd(b.ManifestPush, dryRun); err != nil {
		return fmt.Errorf("failed to push manifest: %v: %v", b.ManifestPush, err)
	}

	return nil
}

func createDockerfile(dockerfile string, option buildOption, descriptor manifestlist.ManifestDescriptor, baseImageName string) error {
	base := template.New("tmpl.Dockerfile")
	parse, err := base.ParseFiles("./internal/tmpl.Dockerfile")
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"GoOS":         option.OS,
		"GoARCH":       option.GoARCH,
		"GoARM":        option.GoARM,
		"RuntimeImage": fmt.Sprintf("%s@%s", baseImageName, descriptor.Digest),
	}

	file, err := os.Create(dockerfile)
	if err != nil {
		return err
	}

	return parse.Execute(file, data)
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

func execCmd(cmd *exec.Cmd, dryRun bool) error {
	if dryRun {
		fmt.Println(cmd.Path, strings.Join(cmd.Args, " "))
		return nil
	}

	output, err := cmd.CombinedOutput()

	if len(output) != 0 {
		log.Println(string(output))
	}

	if err != nil {
		return err
	}
	return nil
}
