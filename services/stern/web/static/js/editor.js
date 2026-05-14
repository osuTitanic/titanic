function getSelectionRange(textarea) {
    if (typeof textarea.selectionStart === "number" && typeof textarea.selectionEnd === "number") {
        return {
            start: textarea.selectionStart,
            end: textarea.selectionEnd
        };
    }

    if (document.selection && textarea.createTextRange) {
        textarea.focus();
        var selectedRange = document.selection.createRange();
        var duplicateRange = selectedRange.duplicate();
        duplicateRange.moveToElementText(textarea);
        duplicateRange.setEndPoint("EndToEnd", selectedRange);

        var selectedText = selectedRange.text || "";
        var end = duplicateRange.text.length;
        var start = end - selectedText.length;

        if (start < 0) {
            start = 0;
        }
        if (end < start) {
            end = start;
        }

        return {
            start: start,
            end: end
        };
    }

    var length = textarea.value.length;
    return {
        start: length,
        end: length
    };
}

function setSelectionRangeCompat(textarea, start, end) {
    if (textarea.setSelectionRange) {
        textarea.setSelectionRange(start, end);
        return;
    }

    if (textarea.createTextRange) {
        var range = textarea.createTextRange();
        range.collapse(true);
        range.moveStart("character", start);
        range.moveEnd("character", end - start);
        range.select();
    }
}

function insertBBCode(event) {
    event = event || window.event;
    event.preventDefault();
    var element = event.target || event.srcElement;

    if (element.tagName !== "BUTTON") {
        element = $(element).parent()[0]; // whyyy
    }

    var bbcodeTag = element.getAttribute("data-bbcode-tag");
    var property = element.getAttribute("data-property");
    var noClose = element.getAttribute("data-no-close");

    var parent = $(element).parent()[0];
    var textAreas = $(parent).parent()[0].getElementsByTagName("textarea");

    if (textAreas.length === 0) {
        console.warn("No text area found in the parent element.");
        return;
    }

    var editor = textAreas[0];

    if (editor && bbcodeTag) {
        var selection = getSelectionRange(editor);
        var start = selection.start;
        var end = selection.end;
        var selectedText = editor.value.substring(start, end);
        var beforeText = editor.value.substring(0, start);
        var afterText = editor.value.substring(end);
        var bbcodeTagStart = "[" + bbcodeTag + "]";
        var bbcodeTagEnd = "[/" + bbcodeTag + "]";

        if (property !== null) {
            bbcodeTagStart = "[" + bbcodeTag + "=" + property + "]";
        }

        if (noClose !== null) {
            bbcodeTagEnd = "";
        }

        editor.value = beforeText + bbcodeTagStart + selectedText + bbcodeTagEnd + afterText;
        editor.focus();
        setSelectionRangeCompat(editor, start + bbcodeTag.length + 2, end + bbcodeTag.length + 2);
    }
}

function insertTextAtSelection(textarea, content) {
    var selection = getSelectionRange(textarea);
    var start = selection.start;
    var end = selection.end;

    var beforeText = textarea.value.substring(0, start);
    var afterText = textarea.value.substring(end);

    textarea.value = beforeText + content + afterText;

    // Move the caret to just after the inserted text
    var newCaretPos = beforeText.length + content.length;
    textarea.focus();
    setSelectionRangeCompat(textarea, newCaretPos, newCaretPos);
}

function replaceFirstOccurrence(textarea, oldContent, newContent) {
    var value = textarea.value;
    var index = value.indexOf(oldContent);
    var isActiveEditor = document.activeElement === textarea;
    var selection = isActiveEditor ? getSelectionRange(textarea) : null;

    if (index === -1) {
        return false;
    }

    var newValue = value.substring(0, index) + newContent + value.substring(index + oldContent.length);
    textarea.value = newValue;

    if (isActiveEditor && selection) {
        var delta = newContent.length - oldContent.length;
        var replacementEnd = index + oldContent.length;
        var selectionStart = selection.start;
        var selectionEnd = selection.end;

        if (selectionStart > replacementEnd) {
            selectionStart += delta;
        } else if (selectionStart >= index) {
            selectionStart = index + newContent.length;
        }

        if (selectionEnd > replacementEnd) {
            selectionEnd += delta;
        } else if (selectionEnd >= index) {
            selectionEnd = index + newContent.length;
        }

        setSelectionRangeCompat(textarea, selectionStart, selectionEnd);
    }

    return true;
}

function createImageBBCode(content) {
    return "[img]" + content + "[/img]";
}

var imageUploadCounter = 0;

function createImageUploadPlaceholder(blob) {
    imageUploadCounter += 1;
    var imageName = blob && blob.name ? blob.name : "pasted-image-" + imageUploadCounter;
    return "Uploading " + imageName + "...";
}

function getClipboardDataFromEvent(event) {
    var nativeEvent = event && event.originalEvent ? event.originalEvent : event;
    return nativeEvent && nativeEvent.clipboardData ? nativeEvent.clipboardData : null;
}

function getClipboardImageFile(clipboardData) {
    if (!clipboardData) {
        return null;
    }

    var items = clipboardData.items;
    if (items) {
        for (var i = 0; i < items.length; i++) {
            var item = items[i];

            if (item.kind !== "file" || !item.type || item.type.indexOf("image/") !== 0) {
                continue;
            }

            var blob = item.getAsFile();
            if (blob) {
                return blob;
            }
        }
    }

    var files = clipboardData.files;
    if (files) {
        for (var j = 0; j < files.length; j++) {
            if (files[j].type && files[j].type.indexOf("image/") === 0) {
                return files[j];
            }
        }
    }

    return null;
}

var editors = $(".bbcode-editor");

for (var i = 0; i < editors.length; i++) {
    var editor = editors[i];

    $(editor).on("paste", function (event) {
        event = event || window.event;
        var editor = event.currentTarget || event.target || event.srcElement;

        var clipboardData = getClipboardDataFromEvent(event);
        var blob = getClipboardImageFile(clipboardData);

        if (!blob) {
            return;
        }

        event.preventDefault();
        event.stopPropagation();

        var uploadPrompt = createImageUploadPlaceholder(blob);
        insertTextAtSelection(editor, uploadPrompt);

        var formData = new FormData();
        formData.append("input", blob);

        performApiRequest(
            "POST",
            "/forum/images",
            formData,
            function (xhr) {
                var response = JSON.parse(xhr.responseText);
                var imageUrl = response.image.image.url;
                replaceFirstOccurrence(editor, uploadPrompt, createImageBBCode(imageUrl));
            },
            function () {
                replaceFirstOccurrence(editor, uploadPrompt, "");
            }
        );
    });
}

var toolbars = $(".bbcode-toolbar");

for (var i = 0; i < toolbars.length; i++) {
    $(toolbars[i]).on("click", insertBBCode);
}
