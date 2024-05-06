package main

import (
	"image"
	"testing"
)

// TestLoadImage checks if loadImage correctly handles image loading
func TestLoadImage(t *testing.T) {
	testPaths := []string{"images/image1.jpeg"}
	jobs := loadImage(testPaths)

	for job := range jobs {
		if job.Image == nil {
			t.Error("Failed to load image")
		}
	}
}

// TestSaveImage checks if saveImage correctly handles image saving
func TestSaveImage(t *testing.T) {
	// Mock job with minimal valid image and path setup
	input := make(chan Job, 1)
	input <- Job{
		OutPath: "images/output/image1.jpeg",
		Image:   image.NewRGBA(image.Rect(0, 0, 100, 100)), // Creating a simple image
	}
	close(input)

	results := saveImage(input)
	success := <-results
	if !success {
		t.Error("Failed to save image")
	}
}

// BenchmarkLoadImage tests the performance of the loadImage function
func BenchmarkLoadImage(b *testing.B) {
	imagePaths := []string{"images/image1.jpeg", "images/image2.jpeg", "images/image3.jpeg", "images/image4.jpeg"} //paths to images

	// Benchmark the loadImage function
	for i := 0; i < b.N; i++ {
		out := loadImage(imagePaths)
		for range out {
		}
	}
}

// BenchmarkSaveImage tests the performance of the saveImage function
func BenchmarkSaveImage(b *testing.B) {
	input := make(chan Job, 1)
	input <- Job{
		OutPath: "images/output/image1.jpeg",
		Image:   image.NewRGBA(image.Rect(0, 0, 100, 100)),
	}
	close(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = saveImage(input)
	}
}
