package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/arifseft/grpc/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Blog client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "ARIF",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// read blog
	fmt.Println("\nReading the blog")
	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "231qeqw"})
	if err != nil {
		fmt.Printf("Error while reading: %v \n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, err := c.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		fmt.Printf("Error while reading: %v \n", err)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// update blog
	fmt.Println("\nUpdating the blog")
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "ARIF",
		Title:    "My first blog edit",
		Content:  "Content of the first blog",
	}
	updateRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		fmt.Printf("Error while reading: %v \n", err)
	}

	fmt.Printf("Blog was updated: %v\n", updateRes)

	// delete blog
	fmt.Println("\nDeleting the blog")
	deleteRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("Error while reading: %v \n", err)
	}

	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	// list blogs
	fmt.Println("\nListing the blog")
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		fmt.Println(res.GetBlog())
	}

}
