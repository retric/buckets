buckets
=======

Intro
-----
Experimental boilerplate with Flask and React. Bower is used to manage front-end
dependencies, with browserify run via gulp to compile them.

This project assumes that any node binaries installed locally via npm install
(such as bower and gulp) will be accessible via your PATH. This can be
accomplished by using the following bash construct whenever inside in the
project directory: 
    
    PATH=$(npm bin):$PATH

Alternately, you could just install all node packages globally.

Installation
------------ 
Install python dependencies
    
    $ pip install flask requests 

Install npm dependencies

    $ npm install

Install required frontend libraries using [bower](http://bower.io/#install-bower).
        
    $ bower install 

Run gulp to start builds and flask 
    
    $ gulp
        
Notes
----- 
Boilerplate adapted from
[flaskr](https://github.com/mitsuhiko/flask/tree/master/examples/flaskr/) and
[flask-react](https://github.com/abhiomkar/flask-react)
