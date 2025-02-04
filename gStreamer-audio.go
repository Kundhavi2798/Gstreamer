package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/notedit/gst"
)

func main() {
	// Initialize GStreamer
	//gst.Init(nil)
	// Create a GStreamer pipeline
	pipeline, err := gst.ParseLaunch(`
		audiotestsrc name=testtone ! volume name=testtone_volume ! autoaudiosink
		autoaudiosrc name=mic ! level name=mic_level ! fakesink
	`)
	if err != nil {
		log.Fatalf("Failed to create pipeline: %v", err)
	}

	// Get elements
	testToneVolume := pipeline.GetByName("testtone_volume")
	if testToneVolume == nil {
		log.Fatalf("Failed to get volume element")
	}

	micLevel := pipeline.GetByName("mic_level")
	if micLevel == nil {
		log.Fatalf("Failed to get level element")
	}

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	fmt.Println("Running... Press Ctrl+C to stop")

	// Monitor microphone volume in a goroutine
	go func() {
		bus := pipeline.GetBus()
		for {
			msg := bus.Pull(gst.MessageElement)
			fmt.Println("The message pulled from bus", msg)
			//if msg != nil && msg.GetStructure() != nil {
			//	structure := msg.GetStructure()
			//	if structure.Name == "level" {
			//		var rms []float64
			//		structure.GetArray("rms", &rms)
			//		if len(rms) > 0 && rms[0] > -50.0 { // Adjust threshold
			//			testToneVolume.SetProperty("volume", 0.0) // Mute test tone when speaking
			//		} else {
			//			testToneVolume.SetProperty("volume", 1.0) // Play test tone when silent
			//		}
			//	}
			//}
			time.Sleep(100 * time.Millisecond) // Prevent CPU overload
		}
	}()

	// Keep the program running
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
