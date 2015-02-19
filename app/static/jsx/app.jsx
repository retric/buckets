/** @jsx React.DOM */
var React   = require('react');  // Browserify!
var jquery  = require('jquery');
var Buckets = require('./components/buckets.jsx');

var objects = [ 
  {"rank": 1, "text": "task1"},
  {"rank": 2, "text": "task2"},
  {"rank": 3, "text": "task3"}
];
 
var App = React.createClass({  // Create a component, App.
  render: function() {
    return (
      <Buckets data={objects} />
    );  
  }
});
React.render(  // Render component at #app.
  <App name="Buckets" />,
  document.getElementById('app'));
