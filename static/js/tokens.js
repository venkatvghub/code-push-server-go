ensureLogin();
var submit = false;

$('#submitBtn').on('click', function () {
    if (submit) return;
    submit = true;
    var query = parseQuery();
    var createdBy = query.hostname || 'Login-' + new Date().getTime();
    var time = new Date().getTime();

    var postParams = {
        createdBy: createdBy,
        friendlyName: "Login-" + time,
        ttl: 60*60*24*30*1000,
        description: "Login-" + time,
        isSession: true
    };
    var accessToken = getAccessToken();
    $.ajax({
        type: 'post',
        data: JSON.stringify(postParams),
        contentType: 'application/json',
        headers: { Authorization: 'Bearer ' + accessToken },
        url: '/accessKeys',
        dataType: 'json',
        success: function (data) {
            submit = false;
            $('#tipsSuccess').show();
            $('#key').val(data.accessKey.name);
            $('#key').show();
            $('#tipsClose').show();
        },
        error: function(xhr, textStatus, errorThrown) {
            submit = false;
            if (errorThrown == 'Unauthorized') {
                alert("please login again!");
                location.href = '/auth/login';
            } else {
                alert(errorThrown);
            }
        }
    });
});