$(document).ready(function() {
    $.ajax({
        url: "/api/test"
    }).then(function(data) {
       $('.Username').prepend(data.username);
       $('.Password').prepend(data.password);
    });
});
