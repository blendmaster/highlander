/* Highlander: prominent topology for prominent people
 * Author: Steven Ruppert
 * For CSCI 447 Scientific Visualization, Spring 2013
 * Colorado School Of Mines
 */
package main

import (
	"fmt"
	"log"
	//"image"
	"image/png"
	//"io"
	"flag"
	"highlander/prominence"
	"os"
)

func main() {

	var threshold uint
	flag.UintVar(&threshold, "threshold", 8000, "threshold of prominence for output features")
	flag.Parse()

	img, err := png.Decode(os.Stdin)
	if err != nil {
		log.Fatalln("Couldn't read png image from Stdin!", err)
	}

	// TODO pass threshold and image as command line args
	features := prominence.ProminentFeatures(img, uint16(threshold))
	for e := features.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
