var React = require('react');

var NewBucket = React.createClass({

    getInitialState: function() {
        return {formOn: false};
    },

    handleClick: function(event) {
        this.setState({formOn: !this.state.formOn});
        event.preventDefault();
    },

    handleBucketSubmit: function(bucketPart) {
        this.setState({formOn: false});
        event.preventDefault();
        $.ajax({
            url: 'api/buckets/',
            dataType: 'json',
            type: 'POST',
            data: bucketPart,
            success: function(data) {
                this.props.addBucket(data);
            }.bind(this),
            error: function(xhr, status, err) {
                console.error('api/buckets/', status, err.toString());
            }.bind(this)
        });
    },

    render: function() {
        if (!this.state.formOn) {
            return (
                <div className="bucketCase">
                <div className="bucket newBucket">
                <a href="" onClick={this.handleClick}>+</a>
                </div></div>
                );
        } else {
            return (
                <NewBucketForm onBucketSubmit={this.handleBucketSubmit} xClick={this.handleClick} />
            );
        }
    }

});

var NewBucketForm = React.createClass({
    handleSubmit: function(event) {
        event.preventDefault();
        var name = React.findDOMNode(this.refs.name).value.trim();
        if (!name) {
            return;
        }
        this.props.onBucketSubmit(JSON.stringify({ name: name }));
        React.findDOMNode(this.refs.name).value = '';
    },
    render: function() {
        return (
                <div className="bucketCase">
                <div className="bucket newBucket">
                <form className="bucketForm" onSubmit={this.handleSubmit}>
                <input type="text" placeholder="Bucket name" ref="name" />
                <input type = "submit" className="submit" value="Submit" />
                <a href="" className="X" onClick={this.props.xClick} >x</a>
                </form>
                </div></div>

        );
    }
});

module.exports = NewBucket;

