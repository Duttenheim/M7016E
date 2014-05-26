/*
 * GET images page.
 */

var request = require("request")

exports.index = function(req, res){
	//curl -XGET 130.240.134.116:5000/v1/search
	//ssh -i yacs.pem.mdlp ubuntu@130.240.134.116 sudo docker ps
	
	//var images = exec("curl -XGET 130.240.134.116:5000/v1/search, {silent:true}").output;
	
	var url = 'http://130.240.134.116:5000/v1/search';
	request({
		url: url,
		json: true
	}, function (error, response, body) {

		if (!error && response.statusCode == 200) {
			console.log(body)
			res.render('images', { title: 'Private repository images' , images: body })
		}
		else {
			console.log(error)
			res.end()
		}
	})
    
};
