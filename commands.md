compositor example
```sh
gst-launch-1.0 videotestsrc pattern=1 ! \
   video/x-raw,format =I420, framerate=\(fraction\)10/1, width=100, height=100 ! \
   compositor name=comp ! videoconvert ! ximagesink \
   videotestsrc ! \
   video/x-raw,format=I420, framerate=\(fraction\)5/1, width=320, height=240 ! comp.
```

Text-render example
```sh
 gst-launch-1.0 -v filesrc location=subtitles.srt ! subparse ! textrender ! videoconvert ! autovideosink
```

a working compositor example when run from ~/Documents
```sh
gst-launch-1.0 compositor name=comp ! videoconvert ! ximagesink videotestsrc ! video/x-raw, width=1920, height=1080, framerate=\(fraction\)60/1 ! comp. filesrc location=subtitles.srt ! subparse ! textrender ! video/x-raw, width=1920, height=1080, framerate=\(fraction\)60/1 ! comp.
```
