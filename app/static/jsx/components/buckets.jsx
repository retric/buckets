var React = require('react');
var NewBucket = require('./new.jsx');

var Buckets = React.createClass({

  render: function() {
    var test = new Array(6);
    for (var i = 0; i < 6; i++) {
        test[i] = 0;
    }

    var bucketlists = this.props.data.map(function(item, index) {
      return (
      <BucketList data={item} key={index} />
      );
    });
    return (
      <div className="row">
      {bucketlists}
      <NewBucket />
      </div>
    );
  }
});

var BucketList = React.createClass({

  render: function() {
      var tasks = this.props.data.Tasks;
      var bucketnodes = tasks.map(function(object, index) {
        return (
          <Item key={index} >
          {object.Priority} {object.Name}
          </Item>
        );
      });
      return(
        <div className="bucketCase">
        <div className="bucket">
        <span className="listName">{this.props.data.Name || "null"}</span>
        <ol className="bucketList">
          {bucketnodes}
        </ol>
        </div>
        </div>
      );
  }

});

var Item = React.createClass({
  
  render: function() {
    return (
      <li className="item">
        {this.props.children}
      </li>
    );
  }

});

module.exports = Buckets;
