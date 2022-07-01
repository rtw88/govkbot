// Copyright 2013 Herbert G. Fischer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imagick

/*
#include <MagickWand/MagickWand.h>
#include <MagickCore/MagickCore.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Image struct {
	img *C.Image
}

func NewMagickImage(info *ImageInfo, width, height uint,
	background *PixelInfo) (*Image, *ExceptionInfo) {

	var exc C.ExceptionInfo

	img := C.NewMagickImage(
		info.info,
		C.size_t(width), C.size_t(height),
		background.pi,
		&exc)

	runtime.KeepAlive(info)

	if err := checkExceptionInfo(&exc); err != nil {
		return nil, err
	}

	return &Image{img: img}, nil
}

// ImageCommandResult is returned by a call to
// ConvertImageCommand. It contains the ImageInfo
// and Metadata string generated by the convert command.
type ImageCommandResult struct {
	Info *ImageInfo
	Meta string
}

/*
ConvertImageCommand reads one or more images, applies one or more image
processing operations, and writes out the image in the same or differing
format.
The first item in the args list is expected to be the program name, ie "convert".
*/
func ConvertImageCommand(args []string) (*ImageCommandResult, error) {
	size := len(args)

	cmdArr := make([]*C.char, size)
	for i, s := range args {
		cmdArr[i] = C.CString(s)
	}

	empty := C.CString("")
	metaStr := C.AcquireString(empty)
	C.free(unsafe.Pointer(empty))

	defer func() {
		for i := range cmdArr {
			C.free(unsafe.Pointer(cmdArr[i]))
		}

		C.DestroyString(metaStr)
	}()

	imageInfo := newImageInfo()

	var exc *C.ExceptionInfo = C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exc)

	ok := C.ConvertImageCommand(
		imageInfo.info,
		C.int(size), // argc
		&cmdArr[0],  // argv
		&metaStr,    // metadata
		exc,         // exception
	)
	if C.int(ok) == 0 {
		imageInfo.Destroy()
		return nil, newExceptionInfo(exc)
	}

	ret := &ImageCommandResult{
		Info: imageInfo,
		Meta: C.GoString(metaStr),
	}
	return ret, nil
}
