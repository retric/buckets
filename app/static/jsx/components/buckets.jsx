var React = require('react');

var Buckets = React.createClass({

  render: function() {
    return (
      <BucketList data={this.props.data} />
    );
  }
});

var BucketList = React.createClass({

  render: function() {
      var bucketnodes = this.props.data.map(function(object, index) {
        return (
          <Item key={index} >
          {object.rank} {object.text}
          </Item>
        );
      });
      return(
        <div className="bucketList">
          {bucketnodes}
        </div>
      );
  }

});

var Item = React.createClass( {
  
  render: function() {
    return (
      <div className="item">
        {this.props.children}
      </div>
    );
  }

});

module.exports = Buckets;
