Simple video info server
=======

A simple server that allows you to download video and get brief information about it

How run app
======
- First, install ffmpeg 3.x libraries on your system.

  (For example for Mac OS:**brew install sy1vain/ffmpeg/ffmpeg@3.4**)

- For run: go run -tags ffmpeg33 *.go

- For build: go build -tags ffmpeg33

- Request example: curl -X POST \
                     http://localhost:4000/api/uploadVideo \
                     -F file=@/path_to_video_file
- Response example:

```json
{
    "video": {
        "name": "h263",
        "width": 176,
        "height": 144,
        "bitRate": 124428,
        "duration": "36.133984375s"
    },
    "audio": {
        "name": "amrnb",
        "bitRate": 12804,
        "duration": "36.021s"
    }
}
```
