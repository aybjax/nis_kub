package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	errs := make(chan error)
	// Connect to NATS
	fmt.Printf("NATS_URL=%s:%s\n", os.Getenv("NATS_URL"), os.Getenv("NATS_PORT"))
	nc, err := nats.Connect(fmt.Sprintf("%s:%s", os.Getenv("NATS_URL"), os.Getenv("NATS_PORT")))
	if err != nil {
		fmt.Println(err)
		// panic(err)
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		fmt.Println(err)
		// panic(err)
	}

	// Create a Stream
	js.AddStream(&nats.StreamConfig{
		Name:        "StudentDataModified",
		Description: "Student data updated, so need to update courses data. Publishing to external service",
		Subjects:    []string{"student.modified"},
	})
	defer js.DeleteStream("StudentDataModified")
	js.AddStream(&nats.StreamConfig{
		Name:        "StudentToModify",
		Description: "Student data updated. Need to transform data to modify it for sending. Internal queue",
		Subjects:    []string{"student.to_modify"},
	})
	defer js.DeleteStream("StudentToModify")
	js.AddStream(&nats.StreamConfig{
		Name:        "CourseModified",
		Description: "Course data updated, so need to update students data. Publishing for external service",
		Subjects:    []string{"course.modified"},
	})
	defer js.DeleteStream("CourseModified")
	js.AddStream(&nats.StreamConfig{
		Name:        "StudentModified",
		Description: "Student Data is modified, so need to update courses data. Consuming from external service",
		Subjects:    []string{"student.modified"},
	})
	defer js.DeleteStream("StudentModified")
	js.AddStream(&nats.StreamConfig{
		Name:        "CourseToModify",
		Description: "Course data updated. Need to transform data to modify it. Internal queue",
		Subjects:    []string{"course.to_modify"},
	})
	defer js.DeleteStream("CourseToModify")

	log.Printf("Initialized and waiting")
	log.Fatalf("exited: %s\n", <-errs)
}
