var url = require('url');
var querystring = require('querystring');

exports.home = function(req, res){
	res.sendfile('views/index.html');
};

exports.images = function(req, res){
	res.sendfile('views/images.html');
};

exports.containers = function(req, res){
	res.sendfile('views/containers.html');
};

exports.nodes = function(req, res){
	res.sendfile('views/nodes.html');
};
