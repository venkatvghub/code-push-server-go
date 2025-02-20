function onLoggedIn() {
    var query = parseQuery();
    if (query.hostname) {
        location.href = '/tokens/' + location.search;
    } else {
        location.href = '/';
    }
}

if (getAccessToken()) {
    onLoggedIn();
}

var submit = false;
$('#submitBtn').on('click', function () {
    if (submit) return;
    submit = true;
    $.ajax({
        type: 'post',
        data: $('#form').serializeArray(),
        url: $('#form').attr('action'),
        dataType: 'json',
        success: function (data) {
            if (data.status == "OK") {
                localStorage.setItem('auth', data.results.tokens);
                submit = false;
                onLoggedIn();
            } else {
                alert(data.message);
                submit = false;
            }
        }
    });
});

// Add Register button if allowed
if ("true" === "true") { // Replace with actual config check if needed
    $('<a id="registerBtn" class="btn btn-lg btn-primary btn-block" href="/auth/register" type="button">Register</a>').insertAfter('#submitBtn');
}