/** @jsx React.DOM */
var React   = require('react');  // Browserify!
var jquery  = require('jquery');
var Buckets = require('./components/buckets.jsx');
 
var App = React.createClass({  // Create a component, App.
  render: function() {
    return (
      <Buckets />
    );  
  }
});
React.render(  // Render component at #app.
  <App name="Buckets" />,
  document.getElementById('app'));
