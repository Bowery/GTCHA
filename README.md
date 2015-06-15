# GTCHA
`\ˈgä-chə\` **G**if to **T**ell **C**omputers and **H**umans **A**part.

## Development

To start the application listening on port 8080:

```
$ goapp serve
```

The GTCHA is live at `http://localhost:8080/static/gtcha.html` and a demo version (acting as an embed) can be viewed at `http://localhost:8080/static/test.html`.

CSS is processed using Myth:

```
$ myth -w static/gtcha.css static/out.css
```
