buckets
=======

Intro
-----
An experimental web app using [React](http://facebook.github.io/react/)
and a [Go](https://golang.org/) server backend.
[Browserify](http://browserify.org/) and [gulp](http://gulpjs.com/) are used to
bundle all js files and simplify frontend dependency management.

Installation
------------ 
Required dependencies: node, go, [gpm](https://github.com/pote/gpm). 

$GOROOT and $GOPATH should be set when running go builds. Setup may be
simplified with [gvm](https://github.com/moovweb/gvm) as detailed in the steps
below.

This project assumes that any node binaries installed locally via npm install
(such as bower and gulp) will also be accessible via your $PATH. This can be
accomplished by using the following bash construct whenever inside the
project directory: 
    
    PATH=$(npm bin):$PATH

Alternately, you could just install all node packages globally.


### Steps:

Install package dependencies.

    $ npm install
    $ bower install 

Setup Go environment.

    $ gvm install go1.4.2 
    $ gvm use go1.4.2
    $ gpm install
    $ cd app
    $ gvm pkgset create --local
    $ gvm pkgset use --local

After initial install, only gvm use/pkgset use are required for future setup.

Run gulp to start the builds. 
    
    $ gulp
        
Run the executable within app/bin to start the server.
    
    $ cd app/bin
    $ ./main

