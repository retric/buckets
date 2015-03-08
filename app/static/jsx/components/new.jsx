var React = require('react');

var NewBucket = React.createClass({

    handleClick: function(event) {
        //TODO: implement server request for new Bucket
    },

    render: function() {
        return (
            <div className="bucketCase">
            <div className="bucket newBucket">
            <a href="" onClick={this.handleClick}>+</a>
            </div></div>
        );
    }

});

module.exports = NewBucket;

