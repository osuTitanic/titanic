// Console stub for IE8 and earlier
if (!window.console) {
    window.console = {
        log: function () {},
        info: function () {},
        warn: function () {},
        error: function () {}
    };
}

// FormData stub for IE9 and earlier
if (!window.FormData) {
    window.FormData = function (form) {};
}

// getElementsByClassName polyfill
function getElementsByClassNamePolyfill(className) {
    var allElements = document.getElementsByTagName("*");
    var matchedElements = [];
    var pattern = new RegExp("(^|\\s)" + className + "(\\s|$)");

    for (var i = 0; i < allElements.length; i++) {
        if (pattern.test(allElements[i].className)) {
            matchedElements.push(allElements[i]);
        }
    }
    return matchedElements;
}

if (!document.getElementsByClassName) {
    document.getElementsByClassName = getElementsByClassNamePolyfill;
}

// Element.getElementsByClassName polyfill
if (window.Element && window.Element.prototype && !window.Element.prototype.getElementsByClassName) {
    window.Element.prototype.getElementsByClassName = function (className) {
        var allElements = this.getElementsByTagName("*");
        var matchedElements = [];
        var pattern = new RegExp("(^|\\s)" + className + "(\\s|$)");

        for (var i = 0; i < allElements.length; i++) {
            if (pattern.test(allElements[i].className)) {
                matchedElements.push(allElements[i]);
            }
        }
        return matchedElements;
    };
}

// Array.indexOf polyfill
if (!Array.prototype.indexOf) {
    Array.prototype.indexOf = function (searchElement, fromIndex) {
        var k;
        if (this == null) {
            throw new TypeError('"this" is null or not defined');
        }
        var o = Object(this);
        var len = o.length >>> 0;
        if (len === 0) {
            return -1;
        }
        var n = fromIndex | 0;
        if (n >= len) {
            return -1;
        }
        k = Math.max(n >= 0 ? n : len - Math.abs(n), 0);
        while (k < len) {
            if (k in o && o[k] === searchElement) {
                return k;
            }
            k++;
        }
        return -1;
    };
}

// Array.isArray polyfill
if (!Array.isArray) {
    Array.isArray = function (arg) {
        return Object.prototype.toString.call(arg) === "[object Array]";
    };
}

// Object.keys polyfill
if (!Object.keys) {
    Object.keys = function (obj) {
        var keys = [];
        for (var key in obj) {
            if (Object.prototype.hasOwnProperty.call(obj, key)) {
                keys.push(key);
            }
        }
        return keys;
    };
}

// String.trim polyfill
if (!String.prototype.trim) {
    String.prototype.trim = function () {
        return this.replace(/^[\s\uFEFF\xA0]+|[\s\uFEFF\xA0]+$/g, "");
    };
}

// addEventListener/removeEventListener polyfill
if (window.Element && window.Element.prototype && !window.Element.prototype.addEventListener) {
    window.Element.prototype.addEventListener = function (type, listener) {
        this.attachEvent("on" + type, listener);
    };
    window.Element.prototype.removeEventListener = function (type, listener) {
        this.detachEvent("on" + type, listener);
    };
}

if (!window.addEventListener) {
    window.addEventListener = function (type, listener) {
        window.attachEvent("on" + type, listener);
    };
    window.removeEventListener = function (type, listener) {
        window.detachEvent("on" + type, listener);
    };
}

if (!document.addEventListener) {
    document.addEventListener = function (type, listener) {
        document.attachEvent("on" + type, listener);
    };
    document.removeEventListener = function (type, listener) {
        document.detachEvent("on" + type, listener);
    };
}
