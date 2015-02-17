var gulp       = require('gulp');
var browserify = require('browserify'); // Bundles JS.
var del        = require('del'); // Deletes files.
var reactify   = require('reactify'); // Transforms React JSX to JS.
var source     = require('vinyl-source-stream');

var paths = {
  app_js: ['./app/static/jsx/app.jsx'],
  js: ['app/static/js/*.js']
};

// dependency task. clean out existing builds.
gulp.task('clean', function(done) {
  del(['build'], done);
});

// JS task. browserify existing code and compile React JSX files.
gulp.task('js', ['clean'], function() {
  browserify(paths.app_js)
    .transform(reactify)
    .bundle()
    .pipe(source('bundle.js'))
    .pipe(gulp.dest('app/static/js/'));
});

// watch task. rerun tasks when files change.
gulp.task('watch', function() {
  gulp.watch(paths.js, ['js']);
});

gulp.task('default', ['watch', 'js']);
