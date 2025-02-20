var submit = false;
$('#submitBtn').on('click', function () {
    if (submit) return;
    submit = true;
    $.ajax({
        type: 'post',
        data: $('#form').serializeArray(),
        url: "/auth/register",
        dataType: 'json',
        success: function (data) {
            if (data.status == "OK") {
                alert("Registration successful, please log in");
                location.href = '/auth/login?email=' + $('#inputEmail').val();
                submit = false;
            } else {
                alert(data.message);
                submit = false;
            }
        },
        error: function(xhr, textStatus, errorThrown) {
            alert("Registration failed: " + errorThrown);
            submit = false;
        }
    });
});