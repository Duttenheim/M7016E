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

app.get('/containers', function(req, res){
  res.render('containers.ejs', {
        title: 'YACS' 
  });
});

app.get('/nodes', function(req, res){
  res.render('nodes.ejs', {
        title: 'YACS' 
  });
});


app.listen(3000);
