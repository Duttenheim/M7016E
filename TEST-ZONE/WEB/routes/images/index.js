/*
 * GET images page.
 */

var request = require("request");

exports.index = function(req, res){
	//curl -XGET 130.240.134.116:5000/v1/search?
	//ssh -i yacs.pem.mdlp ubuntu@130.240.134.116 sudo docker ps
	
	//var images = exec("curl -XGET 130.240.134.116:5000/v1/search, {silent:true}").output;
	
	var IP = '130.240.134.118';	
	if(req.query.search != null){
		var url = 'http://'+IP+':5000/v1/search?q='+req.query.search;
		var text = "Showing results for \'"+req.query.search + "\'";
	}else{
		var url = 'http://'+IP+':5000/v1/search?';
		var text = " "
	}
	//console.log(url)
	request({
		url: url,
		json: true
	}, function (error, response, body) {
		var zeImages;
		if (!error && response.statusCode == 200) {
			zeImages = body;
			//~ console.log(body)
			if (body.num_results==0){
				console.log("No images found");
				request({
				url: 'http://'+IP+':5000/v1/search?',
				json: true
				}, function (error, response, body) {
					//~ console.log(body);
					res.render('images', { title: 'Private repository images' , images: body, server_addr: IP , text:"no images found for "+req.query.search});
				})
			};
		}
		else {
			console.log(error)
			zeImages = error;
		}
		if(!error && body.num_results!=0){
			res.render('images', { title: 'Private repository images' , images: zeImages, server_addr: IP, text :text });
		}
	})
    
};
