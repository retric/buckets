buckets
=======

Intro
-----
Boilerplate with Flask and React

Installation
------------
Install python dependencies
    
    pip install flask requests 

Install required frontend libraries using [bower](http://bower.io/#install-bower).
        
    bower install 

Transform JSX to JS using [React tool](http://facebook.github.io/react/docs/tooling-integration.html#productionizing-precompiled-jsx) for development purpose
        
    jsx --watch app/static/jsx app/static/js
        
Run Flask server
        
    python buckets.py
