package main

import (
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing" // ensure import is correct
	"image"
	"strings"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
	Error     error
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		for _, p := range paths {
			img, err := imageprocessing.ReadImage(p)
			job := Job{
				InputPath: p,
				Image:     img,
				OutPath:   strings.Replace(p, "images/", "images/output/", 1),
				Error:     err,
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			if job.Error == nil { //Only resize if there was no error loading the image
				job.Image = imageprocessing.Resize(job.Image)
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			if job.Error == nil { //Only convert to grayscale if previous steps were successful
				job.Image = imageprocessing.Grayscale(job.Image)
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan bool {
	out := make(chan bool)
	go func() {
		for job := range input {
			if job.Error == nil {
				err := imageprocessing.WriteImage(job.OutPath, job.Image)
				if err != nil {
					job.Error = err
					out <- false
					continue
				}
				out <- true
			} else {
				out <- false
			}
		}
		close(out)
	}()
	return out
}

func main() {
	imagePaths := []string{"images/image1.jpeg", "images/image2.jpeg", "images/image3.jpeg", "images/image4.jpeg"}
	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	for success := range writeResults {
		if success {
			fmt.Println("Success!")
		} else {
			fmt.Println("Failed!")
		}
	}
}
