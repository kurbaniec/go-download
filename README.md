<h1 align="center">
  <br>
  <img src="images/Logo/Logo.png" alt="go download" width="250"></a>
  <br>
  Go Download!
  <br>
</h1>


<h4 align="center">A YouTube Downloader for audio files written in Go.</h4>

## 📝Requirements

This package uses [FFmpeg](https://www.ffmpeg.org/)  to convert YouTube's audio streams to audio files.

Please install it locally and add it to your system path.

## 🛠️ Build

```bash
go build -o ./bin/godownload[.exe|.sh|...] downloader/src
```

## 🚴 Run

```
godownload <URL> [Options]
```

### 📎 Options

```
-q, --quality       Select stream quality. Use "high", "medium" or "low".
                    This option will also define the codec and file type 
                    of the donwload. This option defaults to "high".
                    "high"   => codec flac      .flac
                    "medium" => codec libvorbis .ogg
                    "low"    => codec mp3       .mp3
-m, --manual		Choose codec and file type independently from the 
                    selected stream quality. This flac needs to be used
                    with the following options.
-a, --audio-format  Chose file extension of the file (flac, ogg, mp3, ...).
-c, --codec			Specify the codec to use (flac, libvorbis, mp3, ...).
```



* [Reverse-Engineering Youtube - Alexey Golub | 07.03.2020](https://tyrrrz.me/blog/reverse-engineering-youtube)
* [YoutubeExplode | 07.03.2020](https://github.com/Tyrrrz/YoutubeExplode)
* [Pytube3 | 12.04.202](https://github.com/nficano/pytube)
* [Grabbing Youtube URL | 07.03.2020](https://stackoverflow.com/questions/8317199/grabbing-youtube-video-url-from-curl-or-get-video-info)
* [Deprecating `url_encoded_ftm_stream_map` | 07.03.2020](https://github.com/nficano/pytube/issues/467)
* [Command Line Flags in Go  with package flags | 12.04.2020](https://godoc.org/github.com/jessevdk/go-flags)
* [Convert files with FFmep | 12.04.2020](https://superuser.com/questions/339023/convert-audio-file-to-flac-with-ffmpeg)
* [The Go Gopher Logo - Renee French | 13.04.2020](https://commons.wikimedia.org/wiki/File:Gogophercolor.png)

