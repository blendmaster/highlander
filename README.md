# Highlander

Prominent Topology for prominent people. This is Steven Ruppert's Project
2 for **CSCI 447 Scientific Visualization** of the Spring 2013 semester at
the Colorado School of Mines.

## Usage

Highlander is written in Go, Google's wacky new language with opinions. Thus,
to make the go compiler pick up the import statements, you'll have to add the
highlander directory to `GOPATH`:

    export GOPATH=wherever/you/download/this

Then, you can run the main method in `highlander.go`. The main method reads
a 16bit grayscale PNG from `stdin`, and outputs a CSV formatted list of
prominent saddles and peaks:

    go run highlander.go -threshold 5000 < image.png > prominences.csv

You'll get logging statements on `stderr`, if you want to see progress.

To generate images like the screenshots below, use the included GNU Octave
function `display_prominences.m`, in an octave shell:

    display_prominences('image.png', 'prominences.csv')

Then you can `write` the figure to an image.

### Performance

Not great. My implementation of the algorithm uses memory in a way that Go's
garbage collector doesn't like. I had to resize some of the inputs in the
screenshots to run in a reasonable amount of time.

## Screenshots

![earth, 3000](http://i.imgur.com/KxknAKD.png)
![moon, 3000](http://i.imgur.com/jkgDoPU.png)
![Puget sound, 8000](http://i.imgur.com/yEOAYkM.png)
![Puget sound, 6000](http://i.imgur.com/T8wIGpv.png)
