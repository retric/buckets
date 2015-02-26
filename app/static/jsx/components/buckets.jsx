var React = require('react');

var Buckets = React.createClass({

  render: function() {
    var test = new Array(6);
    for (var i = 0; i < 6; i++) {
        test[i] = 0;
    }

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
        <div className="bucketCase">
        <div className="bucket">
        <ol className="bucketList">
          {bucketnodes}
        </ol>
        </div>
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
