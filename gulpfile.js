var gulp       = require('gulp');
var shell      = require('gulp-shell');
var browserify = require('browserify'); // Bundles JS.
var del        = require('del'); // Deletes files.
var reactify   = require('reactify'); // Transforms React JSX to JS.
var source     = require('vinyl-source-stream');
var concat     = require('gulp-concat-sourcemap');
var sass       = require('gulp-sass');

var paths = {
  jsx: ['./app/static/jsx/app.jsx'],
  js: ['app/static/jsx/*.jsx', 'app/static/jsx/components/*.jsx'],
  css: ['app/static/css/scss/*.scss'],
  go: ['app/src/main/*.go'],
};

// dependency task. clean out existing builds.
gulp.task('clean', function(done) {
  del(['build'], done);
});

// JS task. browserify existing code and compile React JSX files.
gulp.task('js', ['clean'], function() {
  browserify(paths.jsx)
    .transform(reactify)
    .bundle()
    .pipe(source('bundle.js'))
    .pipe(gulp.dest('app/static/js/'));
});

var bower = require('wiredep')({
  directory: 'app/static/libs'
});

// lib task. concatenate all bower library files.
gulp.task('libs', function() {
  gulp.src( bower.js )
    .pipe(concat('libs.js'))
    .pipe(gulp.dest('app/static/js'));
});

gulp.task('css', function() {
    gulp.src('app/static/css/scss/*.scss')
      .pipe(sass())
      .pipe(gulp.dest('app/static/css'));
});

// watch task. rerun tasks when files change.
gulp.task('watch', function() {
  gulp.watch(paths.js, ['js']);
  gulp.watch(paths.go,  ['build']);
  gulp.watch(paths.css, ['css']);
});

// main build tasks.
gulp.task('build', ['buildtest'], shell.task(['go install buckets | tee build.log']));
gulp.task('buildtest', shell.task(['go test buckets']));

gulp.task('default', ['watch', 'js', 'libs', 'css', 'build', 'buildtest']);
