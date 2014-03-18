var express = require('express');
var app = express();
var routes = require('./routes');

app.set('port', process.env.PORT || 8080);

app.get('/', redirectHome);
app.get('/images', redirectImages);
app.get('/containers', redirectContainers);
app.get('/nodes', redirectNodes);


function redirectHome(req,res){
	routes.home(req,res);
}

function redirectImages(req,res){
	routes.images(req,res);
}

function redirectContainers(req,res){
	routes.containers(req,res);
}

function redirectNodes(req,res){
	routes.nodes(req,res);
}

console.log('Express server listening on port ' + app.get('port'));
app.listen(app.get('port'));
