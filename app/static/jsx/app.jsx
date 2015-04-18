/** @jsx React.DOM */
var React   = require('react');  // Browserify!
var jquery  = require('jquery');
var Buckets = require('./components/buckets.jsx');

var App = React.createClass({  // Create a component, App.
  getInitialState: function() {
    return {data: []};
  },
  loadBucketsFromServer: function() {
    $.ajax({
        url: 'api/buckets/',
        dataType: 'json',
        success: function(data) {
            this.setState({data: data});
        }.bind(this),
        error: function(xhr, status, err) {
            console.error('api/buckets/', status, err.toString());
        }.bind(this)
        });
  },
  componentDidMount: function() {
    this.loadBucketsFromServer();
  },
  render: function() {
    return (
      <Buckets data={this.state.data} />
    );  
  }
});

React.render(  // Render component at #app.
  <App name="Buckets" />,
  document.getElementById('app'));

