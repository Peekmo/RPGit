jQuery(function($) {
	if ($('.nav-menu .item').length) {
		$('.nav-menu .item')
		.tab()
		.tab('activate tab', 'first')
		.tab('activate navigation', 'first')
	;
	}

	// Init sidebar
	$('.ui.sidebar').sidebar();

	if (window.innerWidth >= 1000) {
		$('#show-menu').fadeOut();
		$('.sidebar').sidebar('show');

		if (window.innerWidth < 1000 || window.innerWidth > 1500) {
			$('.sidebar').sidebar('pull page');
		}
	}

	$('#close').click(function(e) {
		$('.sidebar').sidebar('hide', function() {
			$('#show-menu').fadeIn();
		});
	});

	$('#home').click(function(e) {
		document.location.href= "/";
	});

	$('.header-rpg').click(function(e) {
		document.location.href = "/";
	});

	$('.filter-language').click(function(e) {
		$('.filter-language').removeClass('filter-language-selected');
		$(this).addClass('filter-language-selected');
	});

	$('#reset').click(function(e) {
		$('.filter-language').removeClass('filter-language-selected');
	});


	$('#show-menu').click(function(e) {
		$('#show-menu').fadeOut();
		$('.sidebar').sidebar('show');


		if (window.innerWidth < 1000 || window.innerWidth > 1500) {
			$('.sidebar').sidebar('pull page');
		}
	});

	$('#search-btn').click(function(e) {
		var name = $('#search-field').val();
		if (name != "") {
			document.location.href = '/users/' + name;
		}
	});	

	$('#search-field').keyup(function(e) {
		if (e.keyCode == 13) {
			var name = $(this).val();
			if (name != "") {
				document.location.href = '/users/' + name;
			}
		}
	});
});
