function filter(element) {
	$('.rp-floating-item').each(function() {
		if ($(this).text() != element) {
			$(this).parent().fadeOut();
			$('#reset').fadeIn();
		} else {
			$(this).parent().fadeIn();
		}	
	});
}

function reset() {
	$('.rp-floating-item').each(function() {
		$(this).parent().fadeIn();
		$('#reset').fadeOut();
	});
}
