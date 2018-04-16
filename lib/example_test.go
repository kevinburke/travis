package travis_test

import (
	"context"
	"fmt"
	"log"

	travis "github.com/kevinburke/travis/lib"
)

func Example() {
	token, err := travis.GetToken("kevinburke")
	if err != nil {
		log.Fatal(err)
	}
	c := travis.NewClient(token)
	build, err := c.Builds.Get(context.TODO(), 366686564, "build.jobs", "job.config")
	if err != nil {
		log.Fatal(err)
	}
	stats, err := c.BuildSummary(context.TODO(), build)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stats)
}

func ExampleListResponse() {
	token, err := travis.GetToken("kevinburke")
	if err != nil {
		log.Fatal(err)
	}
	client := travis.NewClient(token)
	req, err := client.NewRequest("GET", "/repo/rails%2Frails/builds?branch.name=master", nil)
	if err != nil {
		log.Fatal(err)
	}
	req = req.WithContext(context.TODO())
	builds := make([]*travis.Build, 0)
	resp := &travis.ListResponse{
		Data: &builds,
	}
	if err := client.Do(req, resp); err != nil {
		log.Fatal(err)
	}
	for i := range builds {
		fmt.Println(builds[i].ID)
	}
}
