var React = require('react');
var NewBucket = require('./new.jsx');

var Buckets = React.createClass({

  getInitialState: function() {
    return {data: []};
  },

  componentWillReceiveProps: function(nextProps) {
    this.setState({data: nextProps.data.slice(0)});
  },

  addBucket: function(bucket) {
    this.setState({data: this.state.data.concat(bucket)});
  },

  deleteBucket: function(bucketKey) {
    var dataCopy = this.state.data.slice(0);
    dataCopy.splice(bucketKey, 1);
    this.setState({data: dataCopy});
  },

  updateBucket: function(bucketKey, data) {
      var dataCopy = this.state.data.slice(0);
      dataCopy[bucketKey] = data;
      this.setState({data: dataCopy});
  },

  render: function() {
    var bucketlists = this.state.data.map(function(item, index) {
      return (
      <BucketList data={item} key={index} index={index} 
          deleteMe={this.deleteBucket} updateMe={this.updateBucket} />
      );
    }.bind(this));
    return (
      <div className="row">
      {bucketlists}
      <NewBucket addBucket={this.addBucket}/>
      </div>
    );
  }
});

var BucketList = React.createClass({
  getInitialState: function() {
    return {showX: false};
  },

  showX: function() {
    this.setState({showX: true});
  },

  hideX: function() {
    this.setState({showX: false});
  },

  submitDelete: function() {
    var api_url = 'api/buckets/' + this.props.data.ID;  
    $.ajax({
        url: api_url,
        dataType: 'json',
        type: 'DELETE',
        success: function(data) {
            this.props.deleteMe(this.props.index);
        }.bind(this),
        error: function(xhr, status, err) {
            console.error(api_url, status, err.toString());
        }.bind(this)
    });
  },

  onNameSubmit: function(bucketName) {
    var api_url = 'api/buckets/' + this.props.data.ID;  
    $.ajax({
        url: api_url, 
        dataType: 'json',
        type: 'PUT',
        data: bucketName,
        success: function(data) {
            this.props.updateMe(this.props.index, data);
        }.bind(this),
        error: function(xhr, status, err) {
            console.error(api_url, status, err.toString());
        }.bind(this)
    });
  },

  render: function() {
      var partX = this.state.showX ? <X submitDelete={this.submitDelete} /> : null;
      
      var tasks = this.props.data.Tasks;
      var bucketnodes = tasks.map(function(object, index) {
        return (
          <Item key={index} >
          {object.Priority} {object.Name}
          </Item>
        );
      });
      return(
        <div className="bucketCase" onMouseEnter={this.showX} onMouseLeave={this.hideX}>
        <div className="bucket">
        <NameBox name={this.props.data.Name || "null"} onNameSubmit={this.onNameSubmit} />
        {partX}
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

var NameBox = React.createClass({

    getInitialState: function() {
        return {nameClicked: false};
    },

    clickForm: function() {
        this.setState({nameClicked: true});
    },

    blurForm: function() {
        this.setState({nameClicked: false});
    },

    handleSubmit: function(event) {
        event.preventDefault();
        this.setState({nameClicked: false});
        var name = React.findDOMNode(this.refs.newname).value.trim();
        if (!name) {
            return;
        }
        this.props.onNameSubmit(JSON.stringify({ name: name }));
        React.findDOMNode(this.refs.newname).value = '';
    },

    componentWillReceiveProps: function(nextProps) {
        if (nextProps.name != this.state.name) {
            this.setState({nameClicked: false});
        }
        this.setState({name: nextProps.name});
    },

    componentDidUpdate: function() {
        if (this.state.nameClicked) {
            React.findDOMNode(this.refs.newname).focus();
        }
    },

    render: function() {
        var name = this.state.name || this.props.name || "null";
        if (!this.state.nameClicked) {
            return (
                <span className="listName" onClick={this.clickForm}>{name}</span>
                );
        } else {
            return (
                <form className="bucketForm" onSubmit={this.handleSubmit}>
                <input type="text" ref="newname" defaultValue={name} onBlur={this.blurForm}/>
                <input type="submit" className="submit" value="Submit"/></form>
            );
        }
    }
});

var X = React.createClass({

    handleClick: function(event) {
      event.preventDefault();
      this.props.submitDelete();
    },

    render: function() {
      return (
        <a href="" className="X" onClick={this.handleClick}>x</a>
      );
    }
});

module.exports = Buckets;
