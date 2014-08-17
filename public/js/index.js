function filter(element) {
	Home.language = element;
	Home.load('pushevent');
	$('#reset').fadeIn();
}

function reset() {
	Home.language = undefined;
	Home.load('pushevent');
	$('#reset').fadeOut();
}

function open(element) {
	$(element).find('i').removeClass('add sign');
	$(element).find('i').addClass('minus sign');
	$(element).removeClass('green');
	$(element).addClass('red');
}

function close(element) {
	$(element).find('i').removeClass('minus sign');
	$(element).find('i').addClass('add sign');
	$(element).removeClass('red');
	$(element).addClass('green');
}

$(".nav.index").click(function(e) {
	$(".handle-ranking").hide();
	$(".show-ranking").each(function(i, e) {
		close(e);
	});
});


$(".show-ranking").click(function(e) {
	var type = $(this).attr("data-type");
	var icon = $(this).find('i');

	$(".show-ranking").each(function(i, e) {
		if ($(e).attr('data-type') != type) {
			close(e);
		}
	});

	$(".handle-ranking").fadeOut();

	if (icon.hasClass('add sign')) {
		open(this);
		$("#ranking-" + type).fadeIn();
	} else {
		close(this);
	}
});

var Home = {
	init: function() {
		this.compile();
		this.load('pushevent');
	},

	compile: function() {
		this.language = undefined;
		this.data = {};
		this.tpl = Handlebars.compile($("#overview-ranking").html());
	},

	render: function(data) {
		for (var kind in data) {
			for (var type in data[kind]) {
				$("#" + kind + "-" + type).html(this.tpl(data[kind][type]));
				
				if (data[kind][type].length > 0) {
					$("#" + kind + "-top-" + type + "-name").html(data[kind][type][0].key);

					var value = data[kind][type][0].value;
					var value = type == "experience" ? value + " XP" : value + " pushes";
					$("#" + kind + "-top-" + type + "-value").html(value);

					if (type != "language") {
						$("#" + kind + "-top-" + type + "-link").attr("href", "/users/" + data[kind][type][0].key);
					}
				}
			}
		}
	},

	load: function(type) {
		var route = "/api/ranking/home/" + type;

		if (this.language !== undefined) {
			route = route + "/" + this.language;
		}

		$.get(route, function(data) {
			Home.render(data);
		});
	}
};

Home.init();
