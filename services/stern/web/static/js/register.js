function validateField(element) {
    var type = element.getAttribute("name");
    var value = element.value;
    var descriptionField = $(element).parent()[0].querySelector(".input-description");

    if (!value) return;

    descriptionField.innerText = "Checking...";
    descriptionField.style.fontWeight = "normal";

    var xhr = new XMLHttpRequest();
    xhr.open(
        "GET",
        "/account/register/check?type=" + encodeURIComponent(type) + "&value=" + encodeURIComponent(value),
        true
    );
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status !== 200) {
                descriptionField.innerText = "Could not verify this field. Please try something else!";
                descriptionField.style.fontWeight = "bold";
                return;
            }

            var validationError = xhr.responseText;
            if (validationError.length === 0) {
                descriptionField.innerText = "Looking good!";
                descriptionField.style.fontWeight = "normal";
            } else {
                descriptionField.innerText = validationError;
                descriptionField.style.fontWeight = "bold";
            }
        }
    };
    xhr.send();
}

function isValid(element) {
    var descriptionField = $(element).parent()[0].querySelector(".input-description");
    var type = element.getAttribute("name");
    var value = element.value;

    if (!value) {
        descriptionField.innerText = "This field is required!";
        descriptionField.style.fontWeight = "bold";
        return false;
    }

    var xhr = new XMLHttpRequest();
    xhr.open(
        "GET",
        "/account/register/check?type=" + encodeURIComponent(type) + "&value=" + encodeURIComponent(value),
        true
    );
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status !== 200) {
                descriptionField.innerText = "Could not verify this field. Please try something else!";
                descriptionField.style.fontWeight = "bold";
                return false;
            }

            var validationError = xhr.responseText;
            if (validationError.length > 0) {
                descriptionField.innerText = validationError;
                descriptionField.style.fontWeight = "bold";
                return false;
            }

            descriptionField.innerText = "Looking good!";
            descriptionField.style.fontWeight = "normal";
            return true;
        }
    };
    xhr.send();
    return true; // Assume valid until we receive a response
}

function validateAll(event) {
    event.preventDefault();

    var validationFields = document.querySelectorAll(".validate");
    var promises = [];

    for (var i = 0; i < validationFields.length; i++) {
        (function (field) {
            promises.push(
                new Promise(function (resolve) {
                    isValid(field) && resolve(true);
                })
            );
        })(validationFields[i]);
    }

    Promise.all(promises).then(function (results) {
        var allValid = $.inArray(false, results) === -1;
        if (allValid) {
            event.target.submit();
        }
    });
}

var timeout = null;

var validationFields = document.querySelectorAll(".validate");
for (var j = 0; j < validationFields.length; j++) {
    (function (element) {
        $(element).on("keyup", function (event) {
            clearTimeout(timeout);

            timeout = setTimeout(function () {
                validateField(element);
            }, 500);
        });

        $(element).on("blur", function (event) {
            validateField(element);
        });
    })(validationFields[j]);
}
