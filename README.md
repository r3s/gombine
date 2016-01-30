## gombine

`gombine` is a small utility to combine images together.

I created this so that I could combine multiple images from
a comic in to one single image and also learn golang doing it.

Given 5 images as input, `gombine` can place the given
images next to each other or below each other. The default
is to place the images below each other

`gombine` creates an image that can contain the given
images and outputs a jpg or png image based on the
commandline argument. The default format is png.

### Installation
```sh
$ go get github.com/r3s/gombine
$ cd <src folder>
$ go build
```
### Usage

```sh
$ gombine -format=png -side=bottom -out=out.png file1.png file2.png
$ gombine -format=jpg -side=right -out=out.png file1.jpg file2.png
```


![](http://i.imgur.com/ZEgXyau.gif)