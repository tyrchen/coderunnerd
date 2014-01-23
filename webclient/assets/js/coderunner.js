(function ( $ ) {
    var defaults = {
        wsuri: '',
        button: '<button class="coderunner-btn btn btn-default">Run!</button>',
        result: '<hr/><pre class="result"><code class="coderunner-result hljs ">Code execution result</code></pre>',
        buttonSelector: '.coderunner-btn',
        resultSelector: '.coderunner-result',
        contentSelector: '.coderunner-content',
        parentSelector: 'div.coderunner',
        buttonContainer: 'pre.main',
        codeContainer: 'code'
    }

    var attached = false;


    
    $.CodeRunner  = function(options) {
        $.extend(defaults, options);       
    }

    $.fn.CodeRunner = function() {
        if (attached) return this;

        var $parent, $btn;
        var sock = null;

        // append nodes
        this.filter(defaults.parentSelector).append(defaults.result);
        $(defaults.buttonContainer, this).append(defaults.button);
        $(defaults.buttonContainer + ' >' + defaults.codeContainer, this).addClass(defaults.contentSelector.slice(1))

        $btn = $(defaults.buttonSelector);

        $btn.on('click', function() {
            $parent = $(this).closest(defaults.parentSelector); 
            var $content = $(defaults.contentSelector, $parent);
            var suffix = $content.attr('suffix');
            var content = $content.text() || $content.val();
            send({'Suffix': suffix, 'Content': content});
            $(this).blur();
            return false;            
        });

        initWebsocket(defaults.wsuri);

        function send(data) {
            var msg = JSON.stringify(data);
            sock.send(msg);
        };

        function initWebsocket(wsuri) {
            sock = new ReconnectingWebSocket(defaults.wsuri);

            sock.onopen = function() {
                console.log("connected to " + wsuri);
            }

            sock.onclose = function(e) {
                console.log("connection closed (" + e.code + ")");
            }        

            sock.onmessage = function(e) {
                console.log("message received: " + e.data);
                var data;
                try {
                    data = JSON.parse(e.data);
                    $(defaults.resultSelector, $parent).text(data.Content);
                } catch (err) {
                    $(defaults.resultSelector, $parent).text(e.data);
                }
                
            }
        }
        return this;
    };
 
    
}( jQuery ));