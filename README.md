## gombine

`gombine` is a small utility to combine images together.

Given 5 images as input, `gombine` can place the given
images next to each other or below each other. The default
is to place the images below each other

`gombine` creates an image that can contain the given
images and outputs a jpg or png image based on the
commandline argument. The default format is png.

### Installation
```sh
$ go get github.com/r3s/Go-Crypter
$ cd $GOPATH/bin
$ ./Go-Crypter testfile.txt
```
### Usage

```sh
$ gombine -format=png -side=bottom -out=out.png file1.png file2.png
$ gombine -format=jpg -side=right -out=out.png file1.jpg file2.png
```


