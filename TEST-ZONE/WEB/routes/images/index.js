/*
 * GET images page.
 */

var request = require("request")

exports.index = function(req, res){
	//curl -XGET 130.240.134.116:5000/v1/search
	//ssh -i yacs.pem.mdlp ubuntu@130.240.134.116 sudo docker ps
	
	//var images = exec("curl -XGET 130.240.134.116:5000/v1/search, {silent:true}").output;
	
	var url = 'http://130.240.134.118:5000/v1/search';
	request({
		url: url,
		json: true
	}, function (error, response, body) {
		var zeImages;
		var zeError;
		if (!error && response.statusCode == 200) {
			zeImages = body;
			console.log(body)
		}
		else {
			console.log(error)
			zeError = error;
		}
		res.render('images', { title: 'Private repository images' , images: zeImages, errorMsg: zeError  })
	})
    
};
