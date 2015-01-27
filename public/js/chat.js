$(function() {
	var socket = new SockJS(window.location.origin + '/websocket/sockjs/room');
	socket.onopen = function() {
		console.log('open');
	};
	socket.onmessage = function(e) {
		console.log('message', e.data);
		$("#thread").append(e.data).append("<br />")
	};
	socket.onclose = function() {
		console.log('close');
	};

	$('#send').click(function(e) {
		var message = $('#message').val()
		$('#message').val('')
		socket.send(message)
	});

	$('#message').keypress(function(e) {
		if (e.charCode == 13 || e.keyCode == 13) {
			$('#send').click()
			e.preventDefault()
		}
	})
});