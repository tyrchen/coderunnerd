function CodeRunner(wsuri, parent) {
    var sock = null;

    initWebsocket();

    var $coderunner = $('.coderunner');
    var $parent;

    $coderunner.on('click', function(e) {
        $parent = $(this).closest(parent); 
        var $content = $('.coderunner-content', $parent);
        var suffix = $coderunner.attr('suffix');
        var content = $content.text() || $content.val();
        send({'Suffix': suffix, 'Content': content});
        return false;
    });
    
    function send(data) {
        var msg = JSON.stringify(data);
        sock.send(msg);

    };

    function initWebsocket() {
        sock = new ReconnectingWebSocket(wsuri);

        sock.onopen = function() {
            console.log("connected to " + wsuri);
        }

        sock.onclose = function(e) {
            console.log("connection closed (" + e.code + ")");
        }

        sock.onmessage = function(e) {
            console.log("message received: " + e.data);
            var data = JSON.parse(e.data);
            $('.coderunner-result', $parent).text(data.Content);
        }
    }
}