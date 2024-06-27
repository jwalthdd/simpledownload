Simpledownload
==============

Just a package for asynchronous data downloading.

### Usage

```go
downloader := simpledownload.SimpleDownloader{}
downloader.AddDownload("https://download.samplelib.com/mp4/sample-5s.mp4", "/tmp/sample.mp4")
downloader.AddDownload("https://pdfobject.com/pdf/sample.pdf", "/tmp/book.pdf")
downloader.Download()
fmt.Println(downloader.GetStatus())
```


