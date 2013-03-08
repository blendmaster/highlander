/* Highlander: prominent topology for prominent people
 * Author: Steven Ruppert
 * For CSCI 447 Scientific Visualization, Spring 2013
 * Colorado School Of Mines
 */
package main

import (
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

	log.Println("reading image from stdin...")
	img, err := png.Decode(os.Stdin)
	if err != nil {
		log.Fatalln("Couldn't read png image from Stdin!", err)
	}
	log.Println("read image!")

	prominence.PrintProminentFeatures(img, uint16(threshold))
}
