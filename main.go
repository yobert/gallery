package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const thumbSize = 128

type File struct {
	Source string
	Dest   string
	Path   string

	SizeNice string

	Image bool
	Video bool

	ViewPath string

	BigPath     string
	BigSizeNice string

	ThumbPath   string
	ThumbWidth  int
	ThumbHeight int
}

type Folder struct {
	Path string
	Name string
}

type Gallery struct {
	Folders []*Folder
	Files   []*File
	Dest    string
	Path    string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("Usage: gallery <source dir> <dest dir>")
	}

	src := path.Clean(os.Args[1])
	dst := path.Clean(os.Args[2])

	//	limit := 500
	count := 0

	files := []*File{}
	galleries := map[string]*Gallery{}

	fmt.Println("Will write output to", dst)
	fmt.Println("Scanning", src, "...")

	err := filepath.Walk(src, func(fpath string, fi os.FileInfo, err error) error {
		//fmt.Println(fpath)

		suf := strings.TrimPrefix(fpath, src)

		if strings.HasPrefix(fi.Name(), ".") {
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if fi.IsDir() {
			gallery := &Gallery{
				Path: suf,
				Dest: dst + suf + "/index.html",
			}
			if _, dup := galleries[gallery.Path]; dup {
				return fmt.Errorf("Duplicate gallery %#v", gallery.Path)
			}
			//if suf != "" {
			//	gallery.Folders = append(gallery.Folders, &Folder{Path: path.Dir(gallery.Path), Name: "<- Back"})
			//}

			if err := os.MkdirAll(dst+suf, 0755); err != nil {
				return err
			}

			galleries[gallery.Path] = gallery
			//fmt.Println("added gallery", gallery.Path, gallery.Dest)

			// now add to parent gallery folder list
			if gallery.Path != "" {
				parent := path.Dir(gallery.Path)
				if parent == "/" {
					parent = ""
				}
				//fmt.Printf("parent would be: %#v\n", parent)
				gp, ok := galleries[parent]
				if !ok {
					return fmt.Errorf("Cannot find parent gallery %#v for gallery %#v", parent, gallery.Path)
				}
				gp.Folders = append(gp.Folders, &Folder{Path: gallery.Path, Name: fi.Name()})
			}
			return nil
		}

		count++
		//		if count > limit {
		//			return nil
		//		}

		newpath := dst + suf

		file := &File{
			Source:   fpath,
			Dest:     newpath,
			Path:     suf,
			ViewPath: suf + "_view/",
			SizeNice: formatSize(fi.Size()),
		}

		if hasImageSuffix(fpath) {
			file.Image = true
		} else if hasVideoSuffix(fpath) {
			file.Video = true
		} else {
			fmt.Println("not image/video:", fpath)
		}

		files = append(files, file)
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("Generating files and thumbnails and stuff ...")

	for _, file := range files {
		//fmt.Println(file.Path)

		//		dir := path.Dir(file.Dest)
		//		if _, err := os.Stat(dir); err != nil {
		//			if err := os.MkdirAll(dir, 0755); err != nil {
		//				return err
		//			}
		//		}

		dp := path.Dir(file.Path)
		if dp == "/" {
			dp = ""
		}

		if _, err := os.Stat(file.Dest); err != nil {
			cmd := exec.Command("cp", "--reflink=auto",
				file.Source,
				file.Dest)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}

		if file.Image {
			big := file.Dest + "_big.jpeg"
			file.BigPath = file.Path + "_big.jpeg"
			if _, err := os.Stat(big); err != nil {
				cmd := exec.Command("gm", "convert",
					file.Source,
					"-auto-orient",
					"-resize", "2048>",
					"-quality", "90",
					"+profile", "*",
					big)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					return err
				}
			}

			fi, err := os.Stat(big)
			if err != nil {
				return err
			}
			file.BigSizeNice = formatSize(fi.Size())

			thumb := file.Dest + "_thumb.jpeg"

			file.ThumbPath = file.Path + "_thumb.jpeg"
			file.ThumbWidth = thumbSize
			file.ThumbHeight = thumbSize

			if _, err := os.Stat(thumb); err != nil {
				cmd := exec.Command("gm", "convert",
					"-size", fmt.Sprintf("%dx%d", thumbSize*2, thumbSize*2),
					file.Source,
					"-auto-orient",
					"-resize", fmt.Sprintf("%dx%d^", thumbSize, thumbSize),
					"-gravity", "center",
					"-extent", fmt.Sprintf("%dx%d", thumbSize, thumbSize),
					"+profile", "*",
					thumb)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					return err
				}
			}
		} else if file.Video {
			big := file.Dest + "_big.mp4"
			file.BigPath = file.Path + "_big.mp4"
			if _, err := os.Stat(big); err != nil {
				cmd := exec.Command("ffmpeg",
					"-i", file.Source,
					"-c:v", "libx264",
					"-preset", "slow",
					"-crf", "23",
					"-c:a", "libmp3lame",
					"-q:a", "6",
					"-strict", "-1",
					big)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					return err
				}
			}

			fi, err := os.Stat(big)
			if err != nil {
				return err
			}
			file.BigSizeNice = formatSize(fi.Size())

			thumb := file.Dest + "_thumb.mp4"

			file.ThumbPath = file.Path + "_thumb.mp4"
			file.ThumbWidth = thumbSize
			file.ThumbHeight = thumbSize

			if _, err := os.Stat(thumb); err != nil {
				cmd := exec.Command("ffmpeg",
					"-i", file.Source,
					"-vf", fmt.Sprintf("crop=min(in_w\\,in_h):min(in_w\\,in_h), scale=-2:%d", thumbSize),
					"-an",
					"-c:v", "libx264",
					"-preset", "veryfast",
					"-crf", "30",
					"-r", "15",
					thumb)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					return err
				}
			}
		} else {
			fmt.Println("MISC FILE?", file.Source)
		}

		gallery, ok := galleries[dp]
		if !ok {
			return fmt.Errorf("Cannot find gallery %#v", dp)
		}

		gallery.Files = append(gallery.Files, file)
	}

	fmt.Println("Generating galleries ...")

	for _, gallery := range galleries {
		if err := renderGallery(gallery); err != nil {
			return err
		}
		for _, file := range gallery.Files {
			if err := renderFile(file); err != nil {
				return err
			}
		}
	}

	fmt.Println("donesies")

	return nil
}
