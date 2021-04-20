// This example demonstrates how to perform a headless file download. Note that for this technique
// to work, the file type must trigger the "Download / Save As" browser dialog. See the download_image
// example for how to save a file which would load inside the browser window without triggering a download.
package main

import (
	"chromedp-study/core"
	"log"
)

func main() {
	log.Println("STARTED")
	core.Screenshot()
}
