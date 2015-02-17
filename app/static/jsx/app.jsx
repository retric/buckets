/** @jsx React.DOM */
var React = require('react');  // Browserify!
 
var HelloMessage = React.createClass({  // Create a component, HelloMessage.
  render: function() {
    return <div>Hello {this.props.name}</div>;  // Display a property.
  }
});
React.renderComponent(  // Render HelloMessage component at #name.
  <HelloMessage name="John" />,
  document.getElementById('name'));
