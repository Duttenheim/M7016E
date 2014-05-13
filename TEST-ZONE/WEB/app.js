//http://stackoverflow.com/questions/18629327/adding-css-file-to-ejs

var express = require('express');
var app = express();

app.set('views', __dirname + '/views');

app.use(express.static(__dirname + '/public'));

app.get('/', function(req, res){
  res.render('index.ejs', {
        title: 'YACS' 
  });
});

app.get('/images', function(req, res){
  res.render('images.ejs', {
        title: 'YACS' 
  });
});

app.get('/servers', function(req, res){
  res.render('servers.ejs', {
        title: 'YACS' 
  });
});

app.listen(3000);
