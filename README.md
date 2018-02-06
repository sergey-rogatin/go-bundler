# Go-Bundler
[![Build Status](https://travis-ci.org/lvl5hm/go-bundler.svg?branch=master)](https://travis-ci.org/lvl5hm/go-bundler)
[![License: ISC](https://img.shields.io/badge/License-ISC-blue.svg)](https://www.isc.org/downloads/software-support-policy/isc-license/)    
[![NPM](https://nodei.co/npm/go-bundler.png)](https://npmjs.org/package/go-bundler)    
> **STILL IN EARLY DEV AND NOT TESTED, USE AT YOUR OWN RISK**

A bundler for javascript files with minimal config, written in golang.    
Has built-in dev server, auto rebuild on file change and can build html templates along with js files.    
All non-js files are copied to bundle folder and imported as urls.    

*Very fast so far :D*

## How to use
Install bundler
```bash
npm install --save-dev go-bundler
```

Save this in your project folder as `config.json`:
```json
{
  "entry": "index.js",
  "templateHTML": "template.html",
  "bundleDir": "build",
  "watchFiles": true,
  "devServer": {
    "enable": true,
    "port": 8080
  }
}
```

Run npm command 
```bash
go-bundler config.json
```
