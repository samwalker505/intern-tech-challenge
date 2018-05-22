package main

import (
	"context"
	"fmt"
	"sort"
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

type VersionSlice []*semver.Version

func (versions VersionSlice) Len() int {
	return len(versions)
}

func (versions VersionSlice) Less(i, j int) bool {
	return (*versions[i]).Compare(*versions[j]) > 0
}

func (versions VersionSlice) Swap(i, j int) {
	versions[i], versions[j] = versions[j], versions[i]
}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	var currentVer *semver.Version
	sort.Sort(VersionSlice(releases))
	// fmt.Println(releases)
	for _, v := range releases {
		if v.Compare(*minVersion) >= 0 {
			if currentVer == nil {
				versionSlice = append(versionSlice, v)
				currentVer = v
			} else if currentVer.Major != v.Major || currentVer.Minor != v.Minor {
				versionSlice = append(versionSlice, v)
				currentVer = v
			}
		}
	}
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, "kubernetes", "kubernetes", opt)
	if err != nil {
		panic(err) // is this really a good way?
	}
	minVersion := semver.New("1.8.0")
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	versionSlice := LatestVersions(allReleases, minVersion)

	fmt.Printf("latest versions of kubernetes/kubernetes: %s", versionSlice)
}
