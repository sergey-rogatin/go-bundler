# go-bundler
[![License: ISC](https://img.shields.io/badge/License-ISC-blue.svg)](https://www.isc.org/downloads/software-support-policy/isc-license/)    
[![NPM](https://nodei.co/npm/go-bundler.png)](https://npmjs.org/package/go-bundler)    
> **STILL IN EARLY DEV AND NOT TESTED, USE AT YOUR OWN RISK**

*Source code is here: https://github.com/lvl5hm/go-bundler/tree/src.*

A bundler for javascript files with minimal config, written in golang.    
Has built-in dev server, auto rebuild on file change and can build html templates along with js files.    
All non-js files are copied to bundle folder and imported as urls.    

Very fast so far :D

## How to use
Install bundler
```bash
npm install --save-dev go-bundler
```
Run npm command 
```bash
go-bundler
```
You will be prompted to answer a couple of questions to generate a config file. You will then be able to edit it manually in `go-bundler-config.json`
