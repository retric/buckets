var React = require('react');

var Buckets = React.createClass({

  render: function() {
    var test = [ 0, 2, 4 ];
    var data = this.props.data;
    var bucketlists = test.map(function(item, index) {
      return (
      <BucketList data={data} />
      );
    });
    return (
      <div className="row">
      {bucketlists}
      </div>
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
        <div className="col-md-4">
        <ol className="bucketList">
          {bucketnodes}
        </ol>
        </div>
      );
  }

});

var Item = React.createClass( {
  
  render: function() {
    return (
      <li className="item">
        {this.props.children}
      </li>
    );
  }

});

module.exports = Buckets;
