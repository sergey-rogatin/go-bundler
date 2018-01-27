# goBundler
<b>STILL IN EARLY DEV AND NOT TESTED, USE AT YOUR OWN RISK</b>

<p>A bundler for javascript files with minimal config, written in golang.
Has built-in dev server, auto rebuild on file change and can build html templates along with js files.
All non-js files are copied to bundle folder and imported as urls.</p>
<p>Very fast so far :D</p>

# How to use

`npm install --save-dev go-bundler`<br/></br>
Save this in your project folder as `config.json`:
<pre>{
  "entry": "index.js",
  "templateHTML": "template.html",
  "bundleDir": "build",
  "watchFiles": true,
  "devServer": {
    "enable": true,
    "port": 8080
  }
}</pre>
Run npm command `go-bundler config.json`
