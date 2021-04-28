# pdf-extract

Extracts images from PDFs in a hopefully better way than `pdfimages`.

# How it is Better

1. Doesn't have a 1000 image limit
2. Extracts both the mask and image as one so no more combining layers.

# Limitations and Warnings

1. Currently exports everything as a PNG, in order to maintain image alpha 
2. No guarantees this runs on anyone else computer
3. No guarantees this doesn't leak memory like a beast

# Use

```
$ pdfextract -f my_tokens.pdf -d /tmp/tokens/
$ ls /tmp/tokens | head -n 3
page_000_id_000.png
page_000_id_001.png
page_000_id_002.png
```


# Installation on MacOS

```
brew install poppler cairo golang
go install
```

