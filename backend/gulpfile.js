var childProcess = require('child_process');
var gulp = require('gulp');
var gutil = require('gulp-util');
var spawn = childProcess.spawn;
var spawnSync = childProcess.spawnSync;

var isAPI = true;

function log(text) {
  gutil.log(gutil.colors.yellow(text));
}

function run(command, args, callback, ignoreStatus) {
  log(command + ' ' + args.join(' '));
  var status = spawnSync(command, args, {
    stdio: ['ignore', process.stdout, process.stderr],
  }).status;
  callback(ignoreStatus ? 0 : status);
}

gulp.task('gpm', function(callback) {
  run('gpm', ['install'], callback);
});

gulp.task('build', function(callback) {
  run('go', ['install'], callback);
});

var app = null;

gulp.task('run', ['build'], function() {
  if (app) {
    log('app is running, sending SIGTERM');
    // without this, an onExit might be fired after a new app instance has been
    // launched, and we'll end up with app == null, while it's still running
    app.removeAllListeners();
    app.kill();
  }

  log('launching app');
  app = spawn(
    'app', [isAPI ? 'serve' : 'work'],
    { stdio: ['ignore', process.stdout, process.stderr] }
  );
  var onExit = function() { app = null; };
  app.on('exit', onExit);
  app.on('error', onExit);
});

gulp.task('fmt', function(callback) {
  run('go', ['fmt', './...'], callback);
});

gulp.task('vet', function(callback) {
  run('go', ['tool', 'vet', '-composites=false', '.'], callback);
});

gulp.task('watch', function() {
  gulp.watch(['**/*.go'], ['fmt', 'vet', 'run']);
  gulp.watch(['Godeps'], ['gpm', 'run']);
});

gulp.task('dev', ['watch', 'gpm', 'fmt', 'vet', 'run']);
gulp.task('_worker', function() {
  isAPI = false;
});

gulp.task('default', ['dev']);
gulp.task('worker', ['_worker', 'dev']);
