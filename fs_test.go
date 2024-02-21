package jsondb_test

import (
	"strconv"
	"testing"

	"github.com/liamawhite/jsondb"
)

type testObject struct {
    Id string `json:"id"`
	Value int `json:"int"`
}

func (t testObject) ID() string {
    return t.Id
}

func TestFsClient_Write(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a new fsClient
	client, err := jsondb.NewFS[testObject](tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Write a new object
	err = client.Write(testObject{Id: "1", Value: 1})
	if err != nil {
		t.Fatal(err)
	}

	// Check if the file was created
    obj, err := client.Read("1")
    if err != nil {
        t.Fatal(err)
    }
    if obj.Value != 1 {
        t.Fatalf("expected 1, got %d", obj.Value)
    }

    // Ensure files are overwritten
    err = client.Write(testObject{Id: "1", Value: 2})
    if err != nil {
        t.Fatal(err)
    }
    obj, err = client.Read("1")
    if err != nil {
        t.Fatal(err)
    }
    if obj.Value != 2 {
        t.Fatalf("expected 2, got %d", obj.Value)
    }
}

func TestFsClient_Read(t *testing.T) {
    // Create a temporary directory
    tmpDir := t.TempDir()

    // Create a new fsClient
    client, err := jsondb.NewFS[testObject](tmpDir)
    if err != nil {
        t.Fatal(err)
    }

    // Write a new object
    err = client.Write(testObject{Id: "1", Value: 1})
    if err != nil {
        t.Fatal(err)
    }

    // Read the object
    obj, err := client.Read("1")
    if err != nil {
        t.Fatal(err)
    }
    if obj.Value != 1 {
        t.Fatalf("expected 1, got %d", obj.Value)
    }

    // Ensure an error is returned when the object doesn't exist
    _, err = client.Read("2")
    if _, ok := err.(jsondb.NotFoundError); !ok {
        t.Fatalf("expected NotFoundError, got %T", err)
    }
}

func TestFsClient_List(t *testing.T) {
    // Create a temporary directory
    tmpDir := t.TempDir()

    // Create a new fsClient
    client, err := jsondb.NewFS[testObject](tmpDir)
    if err != nil {
        t.Fatal(err)
    }

    // Write a few objects
    for i := 0; i < 10; i++ {
        err = client.Write(testObject{Id: strconv.Itoa(i), Value: i})
        if err != nil {
            t.Fatal(err)
        }
    }

    // List the objects
    objs, err := client.List()
    if err != nil {
        t.Fatal(err)
    }

    // Ensure the list is correct
    if len(objs) != 10 {
        t.Fatalf("expected 3, got %d", len(objs))
    }
}

func TestFsClient_Delete(t *testing.T) {
    // Create a temporary directory
    tmpDir := t.TempDir()

    // Create a new fsClient
    client, err := jsondb.NewFS[testObject](tmpDir)
    if err != nil {
        t.Fatal(err)
    }

    // Write a new object
    err = client.Write(testObject{Id: "1", Value: 1})
    if err != nil {
        t.Fatal(err)
    }

    // Delete the object
    err = client.Delete("1")
    if err != nil {
        t.Fatal(err)
    }

    // Ensure the object was deleted
    _, err = client.Read("1")
    if _, ok := err.(jsondb.NotFoundError); !ok {
        t.Fatalf("expected NotFoundError, got %T", err)
    }

    // Ensure a not found errors is returned when the object doesn't exist
    err = client.Delete("notfound")
    if _, ok := err.(jsondb.NotFoundError); !ok {
        t.Fatalf("expected NotFoundError, got %T", err)
    }
}


