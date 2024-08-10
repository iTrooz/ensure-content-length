# Ensure-Content-Length

This webserver will proxy requests add ensure a content-length response header is set, by loading the response body in RAM to calculate its length when there is no such header.

Content-Length may not be set in a number of cases, including in most web servers when a content encoding is set (See https://serverfault.com/questions/529621/force-nginx-to-send-content-length-header-for-static-files-with-gzip and https://serverfault.com/questions/183843/content-length-not-sent-when-gzip-compression-enabled-in-apache/183856#183856)

# How to run
Build:
`go build -o ensure-content-length`

Run:
`./ensure-content-length 8080 https://example.com`

This will make our proxy webserver run on port 8080, and proxy its requests to https://example.com. For example, http://localhost:8080/my/file will become https://example.com/my/file

# Notes
According to https://stackoverflow.com/questions/3819280/content-length-when-using-http-compression, `Content-Length` is the length of the *compressed* content. Not everyone seems to agree though
