buckets
=======

Intro
-----
An experimental web app using [React](http://facebook.github.io/react/)
and a [Go](https://golang.org/) server backend.
[Browserify](http://browserify.org/) and [gulp](http://gulpjs.com/) are used to
bundle all js files and simplify frontend dependency management.

This project assumes that any node binaries installed locally via npm install
(such as bower and gulp) will be accessible via your PATH. This can be
accomplished by using the following bash construct whenever inside in the
project directory: 
    
    PATH=$(npm bin):$PATH

Alternately, you could just install all node packages globally.

Installation
------------ 
Required dependencies: node, go, [gpm](https://github.com/pote/gpm). 

$GOROOT and $GOPATH should be set when running go builds; setup may be
simplified with [gvm](https://github.com/moovweb/gvm).

Steps:

Install package dependencies.

    $ gpm install
    $ npm install
    $ bower install 

Run gulp to start the builds. 
    
    $ gulp
        
Run the buckets executable within app to start the server.
    
    $ cd app
    $ ./buckets

